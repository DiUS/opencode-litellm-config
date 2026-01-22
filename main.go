package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const usage = `Usage: %s [flags] <output-file>

Generates OpenCode config from LiteLLM models endpoint.

Flags:
  -h, --help                 Show help
  --base-url string          LiteLLM base URL (default "https://litellm.dius.network/v1")
  --provider-name string     Provider display name (default "LiteLLM Dius")
  --provider-key string      Provider key in config (default "litellm-dius")

Environment:
  DIUS_LITELLM_SK            Required. API key for LiteLLM
`

type cliConfig struct {
	baseURL      string
	providerName string
	providerKey  string
	outputFile   string
}

type Model struct {
	ID string `json:"id"`
}

type ModelsResponse struct {
	Data  []Model `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

type ModelConfig struct {
	Name string `json:"name"`
}

type ProviderOptions struct {
	BaseURL string `json:"baseURL"`
}

type Provider struct {
	NPM     string                 `json:"npm"`
	Name    string                 `json:"name"`
	Options ProviderOptions        `json:"options"`
	Models  map[string]ModelConfig `json:"models"`
}

type Config struct {
	Schema   string              `json:"$schema"`
	Provider map[string]Provider `json:"provider"`
}

func main() {
	cfg := parseFlags()
	apiKey := requireEnvVar("DIUS_LITELLM_SK")
	models := fetchModels(cfg.baseURL, apiKey)
	config := buildConfig(cfg, models)
	createParentDirs(cfg.outputFile)
	writeConfig(cfg.outputFile, config)
	fmt.Printf("Config written to %s\n", cfg.outputFile)
}

func parseFlags() cliConfig {
	baseURL := flag.String("base-url", "https://litellm.dius.network/v1", "LiteLLM base URL")
	providerName := flag.String("provider-name", "LiteLLM Dius", "Provider display name")
	providerKey := flag.String("provider-key", "litellm-dius", "Provider key in config")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	return cliConfig{
		baseURL:      *baseURL,
		providerName: *providerName,
		providerKey:  *providerKey,
		outputFile:   flag.Arg(0),
	}
}

func requireEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fatal("%s env var not set", name)
	}
	return value
}

func fetchModels(baseURL, apiKey string) []Model {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", baseURL+"/models", nil)
	if err != nil {
		fatal("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fatal("failed to fetch models: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fatal("API returned status %d", resp.StatusCode)
	}

	var modelsResp ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		fatal("failed to parse response: %v", err)
	}

	if modelsResp.Error != nil {
		fatal("API returned error: %s", modelsResp.Error.Message)
	}

	return modelsResp.Data
}

func buildConfig(cfg cliConfig, models []Model) Config {
	modelConfigs := make(map[string]ModelConfig)
	for _, m := range models {
		modelConfigs[m.ID] = ModelConfig{Name: m.ID}
	}

	return Config{
		Schema: "https://opencode.ai/config.json",
		Provider: map[string]Provider{
			cfg.providerKey: {
				NPM:  "@ai-sdk/openai-compatible",
				Name: cfg.providerName,
				Options: ProviderOptions{
					BaseURL: cfg.baseURL,
				},
				Models: modelConfigs,
			},
		},
	}
}

func createParentDirs(path string) {
	dir := filepath.Dir(path)
	if dir == "." {
		return
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		fatal("failed to create directory: %v", err)
	}
}

func writeConfig(path string, config Config) {
	f, err := os.Create(path)
	if err != nil {
		fatal("failed to create file: %v", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config); err != nil {
		fatal("failed to write config: %v", err)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

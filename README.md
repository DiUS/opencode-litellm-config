# opencode-litellm-config

Generate an [OpenCode](https://opencode.ai) JSON config for connecting to the DiUS LiteLLM proxy (OpenAI-compatible endpoint). Queries the LiteLLM `/models` endpoint to define available models.

You can run the latest binary, or run the project with Go

## Installation

Download the latest binary from [Releases](https://github.com/DiUS/opencode-litellm-config/releases):

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/DiUS/opencode-litellm-config/releases/latest/download/opencode-litellm-config-darwin-arm64
chmod +x opencode-litellm-config-darwin-arm64

# macOS (Intel)
curl -LO https://github.com/DiUS/opencode-litellm-config/releases/latest/download/opencode-litellm-config-darwin-amd64
chmod +x opencode-litellm-config-darwin-amd64

# Linux
curl -LO https://github.com/DiUS/opencode-litellm-config/releases/latest/download/opencode-litellm-config-linux-amd64
chmod +x opencode-litellm-config-linux-amd64
```

## Usage

```bash
export DIUS_LITELLM_SK="your-api-key"
./opencode-litellm-config-darwin-arm64 output.json
```

Or run from source:

```bash
go run main.go output.json
```

## Output
The output json file can be placed in an appropriate [OpenCode config location](https://opencode.ai/docs/config/) such as `~/.config/opencode/opencode.json` for global user config.


### Flags

- `--base-url` - LiteLLM base URL (default: `https://litellm.dius.network/v1`)
- `--provider-name` - Provider display name (default: `LiteLLM Dius`)
- `--provider-key` - Provider key in config (default: `litellm-dius`)
- `-v, --version` - Print version

## Contributing + Releases
Change for your needs and push to main, I don't  mind -- then make a release

```bash
make release VERSION=1.0.0  # explicit version
make release-major          # bump major (1.0.0 → 2.0.0)
make release-minor          # bump minor (1.0.0 → 1.1.0)
make release-patch          # bump patch (1.0.0 → 1.0.1)
```

Builds binaries for linux/darwin/windows, generates checksums, tags, and creates GitHub Release.

## License

Apache 2.0

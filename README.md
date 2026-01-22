# opencode-litellm-config

Generate OpenCode provider config from a LiteLLM `/models` endpoint.

## Usage

```bash
export DIUS_LITELLM_SK="your-api-key"
go run main.go output.json
```

### Flags

- `--base-url` - LiteLLM base URL (default: `https://litellm.dius.network/v1`)
- `--provider-name` - Provider display name (default: `LiteLLM Dius`)
- `--provider-key` - Provider key in config (default: `litellm-dius`)

## License

Apache 2.0

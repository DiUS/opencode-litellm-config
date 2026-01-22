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
- `-v, --version` - Print version

## Releases

Prerequisites: `gh` CLI installed and authenticated, changes committed and pushed.

```bash
make release VERSION=1.0.0  # explicit version
make release-major          # bump major (1.0.0 → 2.0.0)
make release-minor          # bump minor (1.0.0 → 1.1.0)
make release-patch          # bump patch (1.0.0 → 1.0.1)
```

Builds binaries for linux/darwin/windows, generates checksums, tags, and creates GitHub Release.

## License

Apache 2.0

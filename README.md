# confmerge

A command-line tool for deep-merging layered YAML/TOML config files with override precedence rules.

---

## Installation

```bash
go install github.com/yourusername/confmerge@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/confmerge.git
cd confmerge
go build -o confmerge .
```

---

## Usage

Merge multiple config files where later files take precedence over earlier ones:

```bash
confmerge base.yaml staging.yaml overrides.yaml -o merged.yaml
```

**Example:**

`base.yaml`
```yaml
server:
  host: localhost
  port: 8080
  timeout: 30
```

`overrides.yaml`
```yaml
server:
  host: production.example.com
  port: 443
```

**Result:**
```yaml
server:
  host: production.example.com
  port: 443
  timeout: 30
```

TOML files are also supported and can be mixed with YAML inputs:

```bash
confmerge defaults.toml local.toml -o config.toml
```

### Flags

| Flag | Description |
|------|-------------|
| `-o, --output` | Output file path (defaults to stdout) |
| `--format` | Output format: `yaml` or `toml` (auto-detected by default) |
| `--strict` | Exit with error on conflicting scalar types |

---

## License

MIT © 2024 yourusername
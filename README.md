# reviewrules-security

A [reviewrules](https://github.com/popoffvg/reviewrules) rule set that checks a diff for common security mistakes. Each rule is one fork, and each rule has tight `scope` globs, so a diff starts only the rules that apply to its files.

## What it checks

| Rule id | What it flags | Scope | Severity |
|---|---|---|---|
| `secrets` | a password, API key, token, or private key written as a literal | all files | high |
| `injection` | untrusted data in SQL, shell commands, eval, template source, unsafe deserialization, XML with external entities | source files and shell scripts | high |
| `file-and-url` | path traversal, SSRF, open redirect from request values | source files | high |
| `access-control` | a handler with no auth or ownership check that its siblings have, token checks turned off | web source files (no C/C++, no shell) | high |
| `memory-safety` | unbounded copies, non-literal format strings, use after free, unchecked Rust `unsafe` | C, C++, Rust | high |
| `crypto-tls` | TLS verification off, weak password hashing, DES/RC4/ECB, fixed IV, weak random for secrets | source files and shell scripts | medium |
| `data-exposure` | secrets in logs, stack traces to the client, debug on, open CORS with credentials, cookies without Secure/HttpOnly | web source files | medium |
| `workflow-injection` | `${{ }}` with PR or issue data in `run`, `pull_request_target` that runs PR code, secrets or write permissions given to PR code | `.github/workflows/**`, `.github/actions/**`, `action.yml` | high |
| `container-config` | privileged pods, host mounts, root containers, `curl \| sh`, open firewall rules, public buckets | Dockerfile, compose, k8s/helm folders, Terraform | medium |
| `dependency-source` | http package URLs, unpinned git dependencies, new registries, look-alike names, changed lockfile hashes | package manifests and lockfiles | medium |

Source files are Go, Python, JS/TS, Java, Kotlin, Ruby, PHP, C/C++, Rust, and C#. The exact globs are in [.reviewrules/config.json](.reviewrules/config.json).

Each rule is written as sentences, not as a `lead: a; b; c` list, so reviewrules does not split it into one fork per item. One fork reads one rule with all its patterns.

## Add it to a repo

Add this entry to `.reviewrules/settings.toml` of the repo under review:

```toml
[[rule_sets]]
name = "security"
url = "https://github.com/popoffvg/reviewrules-security.git"
ref = "v0.1.0"
```

The rules show as `security/<rule id>` in the summary. The repo under review decides which files the diff ignores, so add generated files and large lockfiles to its `ignorePatterns` when you do not want them checked.

## Cost

Cost is the number of forks that a diff starts: one fork for each rule whose scope matches a changed file.

| Diff | Rules that run | Forks |
|---|---|---|
| Go code only | `secrets`, `injection`, `file-and-url`, `access-control`, `crypto-tls`, `data-exposure` | 6 |
| Go code and `go.mod` | the six above and `dependency-source` | 7 |
| docs or other text only | `secrets` | 1 |
| one workflow file | `secrets`, `workflow-injection` | 2 |

A run on the fixture below (10 rule forks, `correctness = false`, `agent = "codex"`, model `gpt-6-luna`) costs about $0.08.

To run fewer forks, set `correctness = false` in the repo settings, or keep a repo-local copy of the set with fewer rules.

## Fixture

`examples/files/` holds one planted mistake and one safe look-alike for each rule. Copy the files into a new git repo as one commit on top of a clean base and run reviewrules on it. Every rule must give FAIL on the planted line and no finding on the look-alike.

| Rule id | File | Planted mistake | Safe look-alike |
|---|---|---|---|
| `secrets` | `config/settings.py` | `SMTP_PASSWORD` literal | `API_KEY` from env, `SAMPLE_TOKEN` placeholder |
| `injection` | `internal/store/users.go` | `FindByName` joins a query value into SQL | `FindByEmail` uses `$1` and a constant table name |
| `file-and-url` | `internal/files/download.go` | `ServeReport` joins a query value into a path | `FetchAvatar` uses a fixed host and `url.PathEscape` |
| `access-control` | `app/routes/users.js` | `DELETE /users/:id` with no `requireAuth` or `requireAdmin` | `/health` is public, `DELETE /sessions/:id` filters by `req.user.id` |
| `memory-safety` | `src/parse.c` | `strcpy` into `label[32]` | `snprintf` with `sizeof title` |
| `crypto-tls` | `internal/client/http.go` | `InsecureSkipVerify: true` | MD5 for an ETag |
| `data-exposure` | `services/auth.py` | `password` in a log line | log of `user.id` |
| `workflow-injection` | `.github/workflows/pr-title.yml` | PR title in `run` | `head_ref` through `env`, `github.sha` in `run` |
| `container-config` | `deploy/k8s/worker.yaml` | `privileged: true` | sidecar with `runAsNonRoot` |
| `dependency-source` | `package.json` | `csv-fast-parse` from an http URL | `acme-ui` pinned to a full commit sha |

## License

MIT, see [LICENSE](LICENSE).

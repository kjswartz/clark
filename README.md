# clark

Automate seeding GitHub repositories with shared workflow assets (issue templates, slash commands, etc.) via an automated pull request.

`clark` clones one or more target repositories, copies the contents of this project's `resources/` directory into a new `.github/` directory inside each target repo (preserving sub‐folders), commits those files on a new branch (`clark-install-branch`), pushes the branch, and then opens a pull request back to the repo's `main` branch.

## Features

- Interactive prompt for owner and a comma‑separated list of repository names.
- Copies everything under `resources/` into `.github/` of the target repo (e.g. `resources/commands/report.yml` -> `.github/commands/report.yml`).
- Safeguard: skips creating a PR if branch `clark-install-branch` already exists remotely.
- Single, minimal dependency surface (uses `google/go-github`).

## Current Limitations

- Base branch is hardcoded to `main`.
- Branch name is fixed to `clark-install-branch` (will fail/skip if it already exists).
- No non‑interactive / flag based CLI yet (fully prompt driven).
- No dry‑run mode.
- No per‑file templating beyond the static contents in `resources/`.

## Requirements

- Go 1.22+ (to build from source; see `go.mod` for exact version expectations).
- A GitHub personal access token with `repo` scope exported as environment variable `TOKEN`.
- `git` available in your PATH (used for cloning, branching, committing, pushing).

### GitHub Token

Create a classic PAT (or fine‑grained token with appropriate repository content permissions) and export it:

```bash
export TOKEN="ghp_your_token_here"
```

Validate it's set:

```bash
test -n "$TOKEN" || echo "TOKEN not set"
```

## Install / Build

### Option 1: Go install (from your clone or directly)

```bash
git clone https://github.com/kjswartz/clark.git
cd clark
go build -o clark .
```

This creates a local `clark` binary in the project root.

### Option 2: Add to your GOPATH bin (if repo is public and you prefer `go install`)

```bash
go install github.com/kjswartz/clark@latest
```

Ensure `GOBIN` or `GOPATH/bin` is on your PATH so you can run `clark` directly.

## Usage

Run the binary and answer the prompts:

```bash
./clark
```

Example session:

```
What is the name for the owner? my-org
What are the names of the repositories? service-one, service-two
Owner: my-org, Repositories: [service-one service-two]
Submitting pull request for my-org/service-one...
Cloning repository and adding files...
Pull request submitted successfully for my-org/service-one...
Submitting pull request for my-org/service-two...
Cloning repository and adding files...
Pull request submitted successfully for my-org/service-two...
```

After completion, each target repo will have:

```
.github/
	commands/
		report.yml
	ISSUE_TEMPLATE/
		simple-issue.yml
	... (any other files you add under resources/)
```

And an open pull request titled:

> Automated Pull Request

With body:

> This is an automated pull request to add files to your repo from clark.

## How It Works (Internals)

1. Ensure `TOKEN` is present; construct an authenticated GitHub client.
2. For each repo name provided:
	 - Check if branch `clark-install-branch` already exists (via API). If yes, skip.
	 - `git clone` the repository into a temporary directory.
	 - Copy every file from local `resources/` into `.github/` inside the clone (creating directories as needed).
	 - Create/switch to the branch `clark-install-branch` locally.
	 - Stage, commit, and push.
	 - Call GitHub API to open a pull request (`head=clark-install-branch` -> `base=main`).
3. Temporary clone directory is deleted (uses `os.MkdirTemp` + `defer os.RemoveAll`).

## Customizing the Payload

Add or modify files under `resources/`. Subdirectories are preserved when copied into `.github/`. Example locations you might add:

- `resources/workflows/*.yml` → `.github/workflows/*.yml`
- `resources/ISSUE_TEMPLATE/*.yml` → `.github/ISSUE_TEMPLATE/*.yml`
- `resources/commands/*.yml` → `.github/commands/*.yml`

Rebuild / rerun to generate new pull requests (rename or delete the previous branch first if it already exists).

## Roadmap Ideas

- CLI flags (e.g. `--owner`, `--repos`, `--base`, `--branch`, `--non-interactive`).
- Dry‑run / diff output.
- Support for choosing default branch dynamically (detect from repo metadata).
- Optional templating (e.g. Go text/template or Liquid) fed by flags / environment.
- Idempotent update mode (force‑push same branch instead of skipping).
- Logging verbosity levels and structured output (JSON).

## Safety Notes

- The tool creates and pushes a new branch; it doesn't modify default branch directly.
- Skips repositories where the branch already exists to reduce accidental overwrites.
- Review the PR before merging to confirm the added files are expected.

## Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| `TOKEN environment variable not set` | Missing export | `export TOKEN=...` |
| `failed to clone repository` | Bad repo name / permissions / network | Verify name & token scopes |
| Branch skipped message | Existing branch from earlier run | Delete branch or merge/close prior PR |

## Development

Run vet & build:

```bash
go vet ./...
go build ./...
```

## License

See `LICENSE` file.

---

Contributions & feedback welcome. Open an issue or PR to propose enhancements.

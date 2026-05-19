# Repository Guidelines

## Project Structure & Module Organization

This repository contains a small Go HTTP service packaged for Docker.

- `cmd/aip_food_lookup/` contains the Go module and service entry point.
- `docker/` contains the runtime `Dockerfile`, `docker-compose.yml`, and `.env.example`.
- `scripts/build.sh` builds the Linux binary, builds the Docker image, and exports an image tarball.
- There is currently no dedicated test directory; add Go tests beside the package they cover using `*_test.go`.

## Build, Test, and Development Commands

Run commands from the repository root unless noted.

- `cd cmd/aip_food_lookup && go run .` starts the service locally on port `8080`.
- `cd cmd/aip_food_lookup && go test ./...` runs all Go tests in the module.
- `gofmt -w cmd/aip_food_lookup/*.go` formats Go source files.
- `bash scripts/build.sh` builds a static Linux binary at `docker/aip_food_lookup`, builds the `aip_food_lookup` Docker image, and saves `docker/aip_food_lookup.tar`.

The build script assumes a Unix-like shell, Docker, and `sudo`. On Windows, run it from WSL or adapt the commands manually.

## Coding Style & Naming Conventions

Use standard Go conventions:

- Format Go code with `gofmt`.
- Use tabs for Go indentation, as produced by `gofmt`.
- Keep package names short, lowercase, and underscore-free.
- Prefer clear handler/function names such as `handler`, `healthHandler`, or `lookupHandler`.

Use comments sparingly: comment complex code and all non-trivial methods, but do not add noise. Never remove existing comments unless explicitly requested. Keep Docker and compose changes minimal and deployment-focused. Do not commit generated binaries or image archives.

## Testing Guidelines

Use Go's built-in `testing` package. Name test files `*_test.go` and test functions `TestName`. Place tests near the code under test, for example `cmd/aip_food_lookup/main_test.go`.

Always add unit tests for new or modified code. Run `go test ./...` from `cmd/aip_food_lookup` before opening a pull request. Add tests for new handlers, parsing logic, configuration behavior, and error paths.

## Commit & Pull Request Guidelines

Recent commit history uses short, imperative or descriptive messages, for example `Sanitize deployment configuration` and `start project`. Keep commit subjects concise and specific.

Pull requests should include:

- A brief summary of the change.
- Verification steps, such as `go test ./...` or Docker build output.
- Notes for config or deployment changes.
- Linked issues when applicable.

## Security & Configuration Tips

Do not commit real hosts, emails, credentials, tokens, private keys, or generated `.env` files. Use `docker/.env.example` for variable names only, and keep real values in an untracked `docker/.env`.

Make sure code follows security best practices for storing and handling PII, passwords, secrets, and credentials. Prefer environment variables or managed secret stores over hard-coded values. If sensitive data is committed, remove it from files and rewrite history before pushing.

## Agent-Specific Instructions

Always ask qualifying questions before proceeding with code changes. Do not refactor unrelated code while applying a change, including stylistic rewrites such as converting unrelated loops or reshaping unaffected logic.

After changes are complete, ask whether the maintainer is ready to produce post-change deliverables. Only if they confirm, provide any requested combination of commit, PR, and ticket materials.

Post-change deliverables should include:

- Commit message format: `<type>: <concise summary>`.
- PR description template with summary, verification, config/security notes, and linked issues.
- Jira ticket description template with context, scope, acceptance criteria, and validation steps.

# Copilot Coding Agent Instructions for gehoer

## Overview
This repository uses Go and is organized into several packages (e.g., `engraver`, `game`, `music`, etc.). The main entry point is `main.go`. The project implements SMuFL in Go and uses raylib for rendering. The goal is to be able to display scores, so that they can be encapsulated for use in ear training and note reading exercises. 

Furthermore the app wishes to teach the circle of fifts. The app will implement a memorisation algorithm, and finally be released as an iOS app.

## Best Practices for Copilot Coding Agent

### 1. Code Style
- Follow idiomatic Go conventions (gofmt, goimports).
- Use descriptive variable and function names.
- Organize code into logical packages as per the current structure.
- Write concise, clear comments for exported functions and types.

### 2. Testing
- Place tests in the same package as the code, using `_test.go` suffix.
- Use Go's standard `testing` package.
- Prefer table-driven tests for functions with multiple cases.

### 3. Dependencies
- Use Go modules (`go.mod`, `go.sum`) for dependency management.
- Run `go mod tidy` after adding/removing dependencies.

### 4. Commits & Pull Requests
- Write clear, descriptive commit messages.
- Reference related issues or features in PR descriptions.
- Ensure all code passes tests and builds before merging.

### 5. Documentation
- Update `README.md` with any major changes to usage or setup.
- Document new packages and exported APIs with Go doc comments.

### 6. Security & Secrets
- Do not commit secrets or sensitive data.
- Do not expose private keys, credentials, or API tokens.

### 7. Copilot Coding Agent Specific
- When making changes, prefer minimal, targeted diffs.
- Use the existing package structure; do not create new top-level folders unless necessary.
- When adding new features, update or add tests as appropriate.
- If unsure about a design decision, prefer to match the style and patterns of existing code.

---

For more details, see [Best practices for Copilot coding agent in your repository](https://gh.io/copilot-coding-agent-tips)

# Go (Golang) Learning Plan

> A practical, project-driven 10–12 week plan to become proficient in Go, with resources, exercises, milestones, and success criteria.

## Overview
This plan is targeted at a developer who knows at least one other programming language (for example JavaScript, Python, Java, or C#) and wants to learn Go for backend services, tooling, and systems programming. It balances language fundamentals, concurrency, testing, tooling, and small-to-medium projects.

## Assumptions
- You have programming experience in at least one other language.
- You can spend 5–10 hours per week (adjust timeline proportionally).
- Goal: become productive building real Go services and libraries, understand idiomatic Go, and write well-tested concurrent code.

## Requirements checklist (mapped to this deliverable)
- [x] Create a plan to learn Go (Golang).
- [x] Create an `.md` file with the plan in the workspace at `c:\source\Tryouts\golang\golang-learning-plan.md`.

## Contract (what success looks like)
- Inputs: Your current familiarity with programming (assumed), time availability.
- Outputs: Progress through the weekly plan, completed exercises & projects, a ready-to-run sample service.
- Error modes: If you get stuck on tooling or CI, use the resource links and community forums.
- Success criteria: pass unit tests for projects, implement at least 2 small services, and demonstrate idiomatic error handling and concurrency.

## Timeline (10–12 weeks, flexible)
- Weeks 1–2: Setup & Language Basics
- Weeks 3–4: Tooling, Modules, Testing
- Weeks 5–6: Concurrency and Patterns
- Weeks 7–8: Web Services, Databases, APIs
- Weeks 9–10: System Design, Performance & Deployment
- Weeks 11–12 (optional): Advanced topics and a capstone project

## Weekly plan details

### Week 1 — Setup & Basics
- Install Go (use latest stable, currently a Go 1.x release). Verify with `go version`.
- Tools: `gofmt`, `go vet`, `golint` (optional), `gopls` (language server), and an editor (VS Code + Go extension recommended).
- Learn basics: variables, types, control flow, functions, packages.
- Exercises:
  - Complete the official Tour of Go: https://tour.golang.org/
  - Write small programs: FizzBuzz, palindrome check, file line counter.
- Milestone: can compile (`go build`) and run a simple program.

### Week 2 — Types, Structs, Interfaces
- Deeper: slices, maps, arrays, structs, methods, interfaces, pointer semantics.
- Exercises:
  - Implement a stack and queue using slices.
  - Implement `String()` methods and custom types.
- Milestone: understand zero values, pointer receivers vs value receivers, and when to use interfaces.

### Week 3 — Modules, Packaging, and Project Layout
- Learn `go mod` (modules), how to create a module, add dependencies, semantic import paths.
- Project layout: simple module structure, `cmd/`, `pkg/`, `internal/` conventions.
- Exercises:
  - Initialize a module, add a third-party package, build and vendor if needed.
- Milestone: create a small reusable package and consume it.

### Week 4 — Testing, Benchmarking, and Documentation
- Learn `testing` package, table-driven tests, benchmarks (`testing.B`), and writing effective tests.
- Learn `godoc` conventions and comment-based documentation.
- Exercises:
  - Add unit tests for earlier code; add a few benchmarks.
- Milestone: tests run with `go test ./...` and have >80% coverage for small modules.

### Week 5 — Concurrency Basics
- Learn goroutines, channels, buffered channels, select, and basic synchronization.
- Patterns: worker pools, fan-in/fan-out, pipelines.
- Exercises:
  - Implement a parallel worker pool that processes jobs from a channel.
  - Compare concurrent vs sequential execution timings (benchmarks).
- Milestone: implement a safe concurrent pipeline and understand data races.

### Week 6 — Advanced Concurrency & `context`
- Learn `sync` (Mutex, WaitGroup), `atomic` operations, `context.Context` for cancellation and timeouts.
- Debugging data races with `-race` flag.
- Exercises:
  - Build an HTTP client with cancellation and timeouts using `context`.
- Milestone: run `go test -race` and fix any races.

### Week 7 — Building Web Services
- Use `net/http`, `http.ServeMux`, middleware patterns, request routing (or minimal router library), JSON handling.
- Learn about request lifecycle, timeouts, graceful shutdown.
- Exercises:
  - Build a small REST API to store items in memory with CRUD endpoints.
- Milestone: service responds to JSON requests and has basic integration tests.

### Week 8 — Persistence, ORM/SQL, and Migrations
- Learn `database/sql`, prepared statements, connection pooling, and simple ORMs (e.g., `sqlx`, `gorm` — prefer learning `database/sql` first).
- Learn migration tools (`golang-migrate/migrate`) and connection config.
- Exercises:
  - Add a SQLite or Postgres backend to the REST API and implement migrations.
- Milestone: API persists objects to the DB and supports basic queries.

### Week 9 — Observability, Logging, and Metrics
- Structured logging (`log`, `zap`, or `logrus`), tracing basics (OpenTelemetry), and metrics (`prometheus` client).
- Learn graceful shutdown, health checks, and readiness probes.
- Exercises:
  - Add structured logs and Prometheus metrics to your service.
- Milestone: service exposes `/metrics` and logs structured JSON.

### Week 10 — Testing at Scale, CI, and Packaging
- Integration tests, testcontainers (optional), mocking patterns and interfaces.
- CI: GitHub Actions or other; build, test, lint, and publish a Docker image.
- Exercises:
  - Add a simple GitHub Actions workflow that runs `go test`, `golint`, and builds a Docker image.
- Milestone: passing CI for the repo; a runnable Docker image produced.

### Weeks 11–12 — Advanced Topics & Capstone
- Pick one or two advanced areas: generics (Go 1.18+), compiler tools, plugins, gRPC, WebAssembly, or advanced performance tuning.
- Capstone project suggestions:
  - A small, documented microservice with DB, tests, CI, and Docker deployment.
  - A CLI tool using `cobra` or `urfave/cli` with useful functionality.
- Milestone: deliver a small production-ready service or CLI with documentation.

## Projects & Exercises (progressive)
- Starter: CLI text tools (grep-like, word count), and small puzzles.
- Intermediate: REST API with DB and metrics.
- Advanced: gRPC microservice, distributed job runner, or CLI tool for developers.

## Resources (official & recommended)
- Official docs: https://golang.org/doc/
- Tour: https://tour.golang.org/
- Effective Go: https://golang.org/doc/effective_go.html
- Go By Example: https://gobyexample.com/
- Book: *The Go Programming Language* by Donovan & Kernighan
- Community: `r/golang` on Reddit, Go Slack, StackOverflow
- Tools: `gofmt`, `go vet`, `staticcheck`, `gopls`, `golangci-lint` (aggregator of linters)

## Quality gates & checklist for each project
- Build: `go build ./...` (PASS)
- Lint/Static analysis: `golangci-lint run` (PASS/Address warnings)
- Tests: `go test ./...` and `go test -race ./...` (PASS)
- Benchmarks: optional `go test -bench .` (review results)
- Docker: `docker build` (if relevant)

## Good practices and idioms
- Keep code small and focused; prefer composition over inheritance.
- Return errors and check them — prefer wrapping with `%w`.
- Use contexts for cancellation and propagation of request-scoped values.
- Make tests deterministic and use table-driven tests.

## Suggested daily/weekly routine
- Daily: 30–60 minutes learning + 30–60 minutes coding exercises.
- Weekly: one focused mini-project or feature, review and refactor past code.

## Next steps (after following this plan)
- Contribute a small PR to an open-source Go project.
- Start using Go at work in small components or utilities.
- Learn related ecosystem tools (Docker, Kubernetes, CI/CD) as needed.

## Notes & assumptions
- Timeline is flexible; if you can dedicate more time, shorten the weeks. If you have less time, extend.
- If you prefer a different focus (CLI, systems, or frontend via WASM), I can produce a tailored plan.

---

Happy learning — open `c:\source\Tryouts\golang\golang-learning-plan.md` in your editor and start with the Tour of Go.

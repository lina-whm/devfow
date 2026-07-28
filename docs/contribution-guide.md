# Contribution Guide

## Branch Strategy

- `main` — production-ready code
- `develop` — integration branch
- `feature/*` — new features (merge to develop)
- `fix/*` — bug fixes (merge to develop)
- `release/*` — release candidates (merge to main)

## PR Process

1. Create a feature/fix branch from `develop`
2. Implement your changes
3. Run tests: `make test-backend && make test-frontend`
4. Run linting: `make lint-backend && make lint-frontend`
5. Create a pull request to `develop`
6. Request review from at least one team member
7. Address review feedback
8. Merge after approval

## Coding Standards

### Go
- Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `gofmt` before committing
- Avoid global state; use dependency injection
- Errors must be handled or explicitly ignored (`_ =`)
- Interface segregation — small, focused interfaces
- No cyclic imports

### TypeScript/React
- Use TypeScript strict mode
- Prefer functional components with hooks
- Use TanStack Query for server state
- Use React Hook Form for forms with Zod schemas
- Follow the existing component patterns (shadcn/ui style)

### General
- Write tests for all new code
- Keep functions small (under 40 lines)
- Document public APIs with Go doc comments or JSDoc
- Write meaningful commit messages

## Commit Message Format

```
type(scope): short description

[optional body]

[optional footer]
```

Types: feat, fix, refactor, test, docs, chore, perf, security

Examples:
- `feat(auth): implement JWT refresh token rotation`
- `fix(board): resolve optimistic update race condition`
- `docs(api): add OpenAPI spec for task endpoints`

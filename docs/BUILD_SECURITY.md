# Build security baseline

## Go runtime

- Backend builds, Docker images, CI, Release and Security Scan must use Go `1.26.5` or newer within the same supported release line.
- Go `1.26.4` is not allowed because `govulncheck` reports `GO-2026-5856` in `crypto/tls`, fixed by Go `1.26.5`.
- When changing `backend/go.mod`, keep `.github/workflows/*.yml`, root `Dockerfile`, `backend/Dockerfile` and `deploy/Dockerfile` aligned.

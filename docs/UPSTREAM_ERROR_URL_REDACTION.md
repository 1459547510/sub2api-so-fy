# Upstream Error URL Redaction

## Behavior

Client-facing upstream error messages that contain complete `http://` or
`https://` URLs, or common network transport diagnostics, are replaced with:

```text
Service temporarily unavailable, please retry later
```

This applies to:

- OpenAI WebSocket fallback errors.
- Error passthrough rules, including custom passthrough messages.
- Failover-exhausted error passthrough for OpenAI, Anthropic-compatible, and
  Gemini-compatible gateway responses, including streaming error frames.

HTTP status codes, error types, failover scheduling, and account selection are
unchanged.

Ordinary upstream business errors without infrastructure details continue to be
returned according to the configured passthrough rules.

## Operations

The client response is redacted independently from the internal upstream error
message. Existing Ops diagnostics continue to retain the upstream message while
sensitive query values such as `key`, `access_token`, and `refresh_token` remain
masked by the existing sanitizer.

## Verification

Run from `backend/`:

```bash
go test ./internal/service -run 'Test(ResolveOpenAIWSFallbackErrorResponse|ApplyErrorPassthroughRule)' -count=1
go test ./internal/handler -run 'Test(OpenAIGatewayHandler|GatewayHandler|.*Failover.*|.*Stream.*)' -count=1
```

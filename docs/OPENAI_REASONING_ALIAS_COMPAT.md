# OpenAI Reasoning Alias Compatibility

## Scope

OpenAI-compatible API-key accounts can serve a Responses API request through
the Chat Completions fallback path. Some upstream pools return reasoning in the
`reasoning` field instead of `reasoning_content`.

Sub2 accepts both field names for streaming and non-streaming responses. When
both are present, `reasoning_content` takes precedence. The fallback converter
emits reasoning as standard Responses reasoning items and
`response.reasoning_summary_text.delta` events, while final answer text remains
in message output events.

Pool mode only controls same-account retry and scheduling state. It does not
intentionally rewrite successful response text.

## Deployment

- Included from upstream release `v0.1.175`.
- No database migration or configuration change is required.
- Existing account pool settings remain valid.

## Verification

Run the compatibility regression tests:

```bash
cd backend
go test ./internal/pkg/apicompat -run 'TestChatReasoningAlias' -count=1
```

The streaming test must produce `response.reasoning_summary_text.delta` for an
upstream `delta.reasoning` value.

## Rollback

Revert the fork merge commit for upstream `v0.1.175` and publish a follow-up
release. No data rollback is required.

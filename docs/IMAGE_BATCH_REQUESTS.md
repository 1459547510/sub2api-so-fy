# Image `n` Requests

The synchronous OpenAI-compatible image endpoints accept an integer `n` from
`1` through `10` in both JSON and multipart requests:

```json
{
  "model": "gpt-image-2",
  "prompt": "a lighthouse at dusk",
  "n": 4
}
```

The response keeps the standard shape and returns every successful image in
`data[]`. For predictable compatibility, the gateway runs `n` single-image
child requests concurrently and combines their `data[]` entries. Clients do
not need to know which account handled each child request.

Each child request is routed independently and consumes the normal account
concurrency slot. This means account pools, account failover, and pool-mode
retry rules continue to apply. Only images that are actually returned are
recorded as successful usage. A partial upstream outage may therefore produce
fewer entries than requested; the gateway records the successful entries and
returns them as one standard response.

`n` is intended for non-streaming image responses. Streaming requests retain
the existing single upstream stream behavior. Multipart image-edit uploads
can use the same `n` field; the uploaded media is reused for each child
request.

Values outside `1-10` or non-integer values return `400 invalid_request_error`. The gateway
does not expose account identifiers, upstream credentials, or provider-internal
request details in the public response.

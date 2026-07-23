# Leo Video Model Specifications

Sub2 currently exposes three LeoStudio video models. The request contract is
aligned with LeoStudio commit `2fd5c21b01a049817962812cf4675ade7727cc12`.

| Model | Resolutions | Duration | Aspect ratios | Prompt limit |
|---|---|---:|---|---:|
| `seedance-2.0` | 480p, 720p, 1080p | 4-15 seconds, default 8 | 480p/1080p: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`, `21:9`, `9:21`; 720p excludes `9:21` | 5000 characters |
| `seedance-2.0-fast` | 480p, 720p | 4-15 seconds, default 8 | 480p: all seven Seedance ratios; 720p excludes `9:21` | 5000 characters |
| `seedance-2.0-mini` | 720p | 4-15 seconds, default 8 | `16:9` only | 5000 characters |

All three models support generated audio. Each request may contain at most one
start frame, one end frame, four image references, three video references, and
one audio reference. `image_url` is the legacy start-frame field and cannot be
combined with `start_frame_url` or `guidances.start_frame`.

Clients send the friendly `resolution` field to Sub2. LeoStudio sends
`parameters.mode` for Seedance 2.0 and Seedance 2.0 Fast, but does not send that
field for Seedance 2.0 Mini.

`happy-horse-1.1` and `grok-imagine-1.5` are intentionally not exposed by Sub2.

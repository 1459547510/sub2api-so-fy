# Leo Video Model Specifications

Sub2 currently exposes three LeoStudio video models. The request contract is
aligned with LeoStudio `feat/web-admin` commit `3ed1f43438325e56635f4435ff23b4c91c4b2db9` for the currently exposed models.

| Model | Resolutions | Duration | Aspect ratios | Prompt limit |
|---|---|---:|---|---:|
| `seedance-2.0` | 480p, 720p, 1080p | 480p/720p: 4-15 seconds; 1080p: 4-12 seconds; default 8 | 480p/1080p: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`, `21:9`, `9:21`; 720p excludes `9:21` | 5000 characters |
| `seedance-2.0-fast` | 480p, 720p | 4-15 seconds, default 8 | 480p: all seven Seedance ratios; 720p excludes `9:21` | 5000 characters |
| `seedance-2.0-mini` | 720p | 4-15 seconds, default 8 | `16:9` only | 5000 characters |

All three models support generated audio. Each request may contain at most one
start frame, one end frame, four image references, three video references, and
one audio reference. `image_url` is the legacy start-frame field and cannot be
combined with `start_frame_url` or `guidances.start_frame`.

Start and end frames are independent and may be submitted together. Image
references are a separate mode and cannot be combined with either frame.

`guidances.video_reference_base` accepts a Leonardo video UUID or an absolute
MP4/MOV URL reachable by LeoStudio; URL entries use `type: "UPLOADED"`. The
`guidances.audio_reference` field accepts a Leonardo audio UUID or an absolute
MP3/WAV URL reachable by LeoStudio. UUID audio entries may include a 2-30
second `duration`; URL audio entries omit it, and every audio reference must be
paired with an image or video reference.

Clients send the friendly `resolution` field to Sub2. LeoStudio sends
`parameters.mode` for Seedance 2.0 and Seedance 2.0 Fast, but does not send that
field for Seedance 2.0 Mini.

`happy-horse-1.1` and `grok-imagine-1.5` are intentionally not exposed by Sub2.

## Uploaded media references

The customer-facing upload endpoint accepts MP4/MOV video references and MP3
or PCM 16/24-bit WAV audio references. Videos are limited to 100 MiB and three
items per generation. Audio is limited to 15 MiB, one item per generation, and
2-30 seconds. The workbench uploads files and uses the returned opaque
`media_url` with `type: "UPLOADED"`; UUID fields remain an advanced API
compatibility path and are not required for customers.

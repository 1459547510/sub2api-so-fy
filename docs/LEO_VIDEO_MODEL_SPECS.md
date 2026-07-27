# Leo Video Model Specifications

This document describes the public video contract exposed by Sub2API after
normalizing the current video capability registry. Model IDs are public request
values; internal account and task identifiers are not part of the customer
contract.

## Model matrix

| Model | Resolutions | Duration | Aspect ratios | Prompt limit |
| --- | --- | --- | --- | ---: |
| `seedance-2.0` | `480p`, `720p`, `1080p` | `4-15s`; `1080p` is `4-12s`; default `8s` | `480p`/`1080p`: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`, `21:9`, `9:21`; `720p` excludes `9:21` | 5000 |
| `seedance-2.0-fast` | `480p`, `720p` | `4-15s`; default `8s` | `480p`: all seven Seedance ratios; `720p` excludes `9:21` | 5000 |
| `seedance-2.0-mini` | `480p`, `720p` | `4-15s`; default `8s` | Both resolutions: `16:9`, `1:1`, `9:16` | 5000 |
| `happy-horse-1.1` | `720p`, `1080p` | `3-15s`; default `5s` | Both resolutions: `16:9`, `4:3`, `1:1`, `3:4`, `9:16` | 2500 |
| `grok-imagine-1.5` | `400p`, `544p`, `720p`, `960p` | `3-15s`; default `6s` | `400p`/`720p`: `16:9`, `9:16`; `544p`/`960p`: `1:1` | 5000 |

The platform rejects unsupported model, resolution, duration, and aspect-ratio
combinations before dispatch. `seedance-2.0` is the compatibility alias for
`seedance` when an account mapping uses that name.

### Model-specific guidance

- Every model accepts at most one start frame. Start and end frames are
  independent, so both may be supplied in the same request when the model
  supports an end frame.
- Reference-image mode and frame mode are mutually exclusive. A request may
  use either the start/end frame fields or reference images, but not both.
- Seedance models allow up to four reference images, three reference videos,
  and one reference audio item.
- `happy-horse-1.1` allows a start frame and up to nine reference images. It
  does not accept an end frame, reference video, or reference audio. Its
  optional `prompt_enhance` value is `AUTO`, `ON`, or `OFF`; `ON` cannot be
  combined with a start frame.
- `grok-imagine-1.5` requires exactly one start frame and does not accept an
  end frame, reference image, reference video, or reference audio. It does not
  accept `prompt_enhance` or other unsupported guidance options.

## Media upload contract

The customer upload endpoint is `POST /v1/videos/uploads`. Uploaded media is
returned as a public-contract `media_url`; use that value in the generation
request.

### Images

- The form field is `image` (the legacy `file` field remains accepted).
- The file must be a readable PNG, JPEG, or WebP image.
- The size limit is 10 MiB per file.
- A generation may contain at most four reference images for Seedance or nine
  for Happy Horse. Grok does not accept reference-image guidance.

### Reference videos

- The form field is `video`.
- The filename must end in `.mp4` or `.mov`.
- The bytes must be an ISO Base Media container with an `ftyp` header and must
  decode as a readable video; renaming an unsupported file is not sufficient.
- The size limit is 100 MiB per file, with at most three video references per
  generation.
- The public contract does not publish one universal duration, frame-rate,
  codec, audio-track, or dimension whitelist. A file that passes the container
  check can still be rejected by a selected model during processing.
- Reference video guidance is available to Seedance models only.

### Reference audio

- The form field is `audio`.
- MP3 files must contain readable MP3 frames.
- WAV files must be RIFF/WAVE PCM with 16-bit or 24-bit samples. Compressed or
  floating-point WAV files are rejected.
- The size limit is 15 MiB per file, the duration must be 2-30 seconds, and a
  generation may contain at most one reference audio item.
- A reference audio item must be paired with at least one reference image or
  reference video. Generated output audio (`audio: true`) is a separate option.
- Reference audio guidance is available to Seedance models only.

## URL and request rules

- Direct media URLs must be absolute HTTP(S) URLs reachable by the video
  service while the job is queued. Data URLs, Base64 values, and multipart
  media embedded in the JSON generation request are not accepted.
- The `image_url` field is the single-start-frame compatibility field. Use
  `start_frame_url` and `end_frame_url` for explicit frame pairs.
- Use `guidances.image_reference`, `guidances.video_reference_base`, and
  `guidances.audio_reference` for uploaded reference media. Audio URL entries
  omit a client-supplied duration because the upload service validates it.
- The asynchronous endpoint returns only the Sub2API `job_id`, `status`, and
  `status_url` in its acceptance response. Poll the status URL and download a
  completed MP4 from the Sub2API content endpoint.

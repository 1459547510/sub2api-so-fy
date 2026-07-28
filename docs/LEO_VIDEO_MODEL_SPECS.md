# Leo Video Model Specifications

This document describes the public video contract exposed by Sub2API after
normalizing the current video capability registry. Model IDs are public request
values; internal account and task identifiers are not part of the customer
contract.

## Model matrix

The following table is the server-side validation contract. A request must
use a supported resolution/aspect-ratio pair from the same row.

| Model | Resolutions | Duration / default | Aspect ratios | Prompt limit | Reference inputs |
| --- | --- | --- | --- | ---: | --- |
| `seedance-2.0` | `480p`, `720p`, `1080p` | `4-15s`; default `8s`; `1080p` is `4-12s` | `480p`/`1080p`: `16:9`, `9:16`, `1:1`, `4:3`, `3:4`, `21:9`, `9:21`; `720p` excludes `9:21` | 5000 | 1 start frame, 1 end frame, 4 images, 3 videos, 1 audio |
| `seedance-2.0-fast` | `480p`, `720p` | `4-15s`; default `8s` | `480p`: all seven Seedance ratios; `720p` excludes `9:21` | 5000 | 1 start frame, 1 end frame, 4 images, 3 videos, 1 audio |
| `seedance-2.0-mini` | `480p`, `720p` | `4-15s`; default `8s` | Both resolutions: `16:9`, `1:1`, `9:16` | 5000 | 1 start frame, 1 end frame, 4 images, 3 videos, 1 audio |
| `happy-horse-1.1` | `720p`, `1080p` | `3-15s`; default `5s` | Both resolutions: `16:9`, `4:3`, `1:1`, `3:4`, `9:16` | 2500 | 1 start frame or 9 images; no end frame/video/audio; `prompt_enhance` |
| `grok-imagine-1.5` | `400p`, `544p`, `720p`, `960p` | `3-15s`; default `6s` | `400p`/`720p`: `16:9`, `9:16`; `544p`/`960p`: `1:1` | 5000 | exactly 1 start frame; no end frame/image/video/audio |

The platform rejects unsupported model, resolution, duration, and aspect-ratio
combinations before dispatch. `seedance-2.0` is the compatibility alias for
`seedance` when an account mapping uses that name.

## Common request fields

| Field | Type | Required | Rules |
| --- | --- | --- | --- |
| `model` | string | yes | One of the five public model IDs above. `seedance` is accepted as an alias for `seedance-2.0`. |
| `prompt` | string | yes | Scene, action, camera, and style description. Maximum is 5000 characters except `happy-horse-1.1`, which is 2500. |
| `resolution` | string | no | Defaults to `720p` when omitted. The selected model must support the value. |
| `duration` | integer | no | Whole seconds. Defaults to the model default in the matrix. |
| `aspect_ratio` | string | no | Defaults to the first supported ratio for the selected resolution. |
| `audio` | boolean | no | Whether generated output should include audio. Defaults to `false`; this is separate from a reference-audio input. |
| `prompt_enhance` | string | no | `happy-horse-1.1` only: `AUTO`, `ON`, or `OFF`. `ON` cannot be combined with a start frame. |
| `image_url` | string | no | One start-frame absolute HTTP(S) URL. Use `start_frame_url` and `end_frame_url` for an explicit frame pair. |
| `start_frame_url` / `end_frame_url` | string | no | Absolute HTTP(S) frame URLs. Frame mode cannot be combined with reference images. |
| `guidances` | object | no | Nested `image_reference`, `video_reference_base`, and `audio_reference` arrays. Use the media object formats below. |

## Reference video request format

Upload a local reference video first. The upload response's `media_url` is the
value used in the generation request:

```bash
curl -X POST "$SUB2_BASE_URL/v1/videos/uploads" \
  -H "Authorization: Bearer $SUB2_API_KEY" \
  -F "video=@./reference.mp4"
```

```json
{
  "media_url": "https://media.example/uploaded/reference.mp4",
  "media_type": "video",
  "content_type": "video/mp4",
  "size": 428516
}
```

Use the returned URL under `guidances.video_reference_base[].video`:

```json
{
  "model": "seedance-2.0",
  "prompt": "Preserve the motion and timing of the reference video",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "guidances": {
    "video_reference_base": [
      {
        "video": {
          "url": "https://media.example/uploaded/reference.mp4",
          "type": "UPLOADED"
        }
      }
    ]
  }
}
```

The video object must contain either `url` or a UUID `id`, never both. For a
URL, `type` must be `UPLOADED`; do not add a `duration` property. A generation
may contain at most three reference videos, and only Seedance models accept
this guidance. The file itself must be a readable MP4/MOV ISO Base Media
container and no larger than 100 MiB. Data URLs, Base64, and multipart media
inside the generation JSON are not accepted.

## Reference audio format

Use `guidances.audio_reference[].audio` with the same `url` and
`type: "UPLOADED"` shape. Audio URL entries must omit `duration`; an audio
reference must be paired with at least one image or reference video. Only
Seedance models accept reference audio. The upload endpoint accepts MP3 with
readable frames or RIFF/WAVE PCM 16/24-bit WAV, up to 15 MiB and 2-30 seconds.

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

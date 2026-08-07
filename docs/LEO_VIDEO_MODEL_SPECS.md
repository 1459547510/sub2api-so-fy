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
| `ltx-2.3-pro` | `1080p`, `1440p`, `2160p` | `6s`, `8s`, `10s`; default `6s` | `16:9` only | 5000 | 1 start frame, 1 end frame; no image/video/audio references; generated audio and `prompt_enhance` |
| `ltx-2.3-fast` | `1080p`, `1440p`, `2160p` | even values from `6s` to `20s`; default `6s` | `16:9` only | 5000 | 1 start frame, 1 end frame; no image/video/audio references; generated audio and `prompt_enhance` |
| `hailuo-03` | `1440p` | `5-15s`; default `5s` | `16:9`, `1:1`, `9:16` | 2000 | up to 5 images or 3 audio references totaling 15s; 1 start frame and 1 end frame; frames exclude image references; no video reference |
| `gemini-omni-flash` | `720p` | `3-10s`; default `5s` | `16:9`, `9:16` | 2500 | up to 5 images; no generated audio, frames, video reference, or audio reference |
| `kling-2.1` | `1080p` | `5s` or `10s`; default `5s` | `16:9`, `1:1`, `9:16` | 2500 | required start frame, optional end frame; `prompt_enhance` |
| `kling-2.5` | `720p`, `1080p` | `5s` or `10s`; default `5s` | `16:9`, `1:1`, `9:16` | 2500 | start and end frames; `prompt_enhance` |
| `kling-2.5-turbo-standard` | `720p` | `5s` or `10s`; default `5s` | `16:9`, `1:1`, `9:16` | 2500 | required start frame; no end frame; `prompt_enhance` |
| `kling-2.6` | `auto`, `1080p` | `5s` or `10s`; default `5s` | `auto` with `auto`, otherwise `16:9`, `1:1`, `9:16` | 2500 | one start frame; generated audio |
| `kling-video-o-1` | `1080p` | `3-10s`; default `5s` | `16:9`, `1:1`, `9:16` | 2500 | frames or images; video reference must be an existing `GENERATED` asset, not an uploaded URL |
| `kling-3.0` | `auto`, `720p`, `1080p`, `2160p` | `3-15s`; default `5s` | `auto` with `auto`, otherwise `16:9`, `1:1`, `9:16` | 2500 | 1 start frame, 1 end frame; generated audio |
| `kling-3.0-turbo` | `auto`, `720p`, `1080p` | `3-15s`; default `5s` | `auto` with `auto`, otherwise `16:9`, `1:1`, `9:16` | 2500 | 1 start frame; generated audio |
| `kling-video-o-3` | `720p`, `1080p`, `2160p` | `3-15s`; default `5s` | `16:9`, `1:1`, `9:16` | 2500 | up to 7 images or 1 `GENERATED` video; with video, max 4 images and 10s; frames exclude other references |
| `veo-3.1-generate-001` | `720p`, `1080p`, `2160p` | `4s`, `6s`, `8s`; default `8s` | `16:9`, `9:16` | 9999 | 1 start frame, 1 end frame, up to 3 images; generated audio |
| `veo-3.1-fast-generate-001` | `720p`, `1080p`, `2160p` | `4s`, `6s`, `8s`; default `8s` | `16:9`, `9:16` | 9999 | 1 start frame, 1 end frame; no image reference; generated audio |
| `veo-3.1-lite` | `720p`, `1080p` | `4s`, `6s`, `8s`; default `8s` | `16:9`, `9:16` | 9999 | 1 start frame, 1 end frame; generated audio |

The platform rejects unsupported model, resolution, duration, and aspect-ratio
combinations before dispatch. `auto` is a real resolution value only for Grok
Imagine 1.5 and Kling 2.6/3.0/3.0 Turbo, and it must be paired with
`aspect_ratio: "auto"`. `seedance-2.0` is the compatibility alias for
`seedance` when an account mapping uses that name.

## Common request fields

| Field | Type | Required | Rules |
| --- | --- | --- | --- |
| `model` | string | yes | One of the public model IDs in the matrix. `seedance` is accepted as an alias for `seedance-2.0`. |
| `prompt` | string | yes | Scene, action, camera, and style description. The maximum is model-specific and is shown in the matrix. |
| `resolution` | string | no | Defaults to the selected model's default. The selected model must support the value; `auto` is only available on the models listed above. |
| `duration` | integer | no | Whole seconds or a model-specific enumerated value. Defaults to the model default in the matrix. |
| `aspect_ratio` | string | no | Defaults to the first supported ratio for the selected resolution. `auto` is valid only with `auto` resolution. |
| `audio` | boolean | no | Whether generated output should include audio. Defaults to `false`; this is separate from a reference-audio input. |
| `prompt_enhance` | string | no | `happy-horse-1.1`, Kling 2.1/2.5/Turbo Standard, and both LTX models: `AUTO`, `ON`, or `OFF`. Only Happy Horse forbids `ON` with a start frame. |
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

The video object uses the upload response's `media_url` as `url`, with
`type: "UPLOADED"`; do not add a `duration` property. Seedance models accept
uploaded video references. Kling O-1 and O-3 video guidance is a different
contract: it accepts only an existing `GENERATED` video asset, so a customer
cannot turn a local upload into that type. The file itself must be a readable
MP4/MOV ISO Base Media container and no larger than 100 MiB. Data URLs, Base64,
and multipart media inside the generation JSON are not accepted.

## Reference audio format

Use `guidances.audio_reference[].audio` with the same `url` and
`type: "UPLOADED"` shape. Audio URL entries must omit `duration`; an audio
reference must be paired with at least one image or reference video when the
selected model requires media. Seedance and Hailuo 03 accept reference audio;
Hailuo allows up to three items with a combined duration of 15 seconds. The
upload endpoint accepts MP3 with readable frames or RIFF/WAVE PCM 16/24-bit WAV,
up to 15 MiB and 2-30 seconds per file.

### Model-specific guidance

- Every model accepts at most one start frame. Start and end frames are
  independent, so both may be supplied in the same request when the model
  supports an end frame.
- Reference-image mode and frame mode are mutually exclusive. A request may
  use either the start/end frame fields or reference images, but not both.
- Seedance models allow up to four reference images, three reference videos,
  and one reference audio item; reference audio must be paired with an image
  or video reference.
- `hailuo-03` allows up to five images or three audio references totaling 15
  seconds. Its end frame requires a start frame, and frame mode excludes all
  image references.
- `gemini-omni-flash` allows up to five images, but no generated audio, frames,
  video reference, or audio reference.
- `kling-2.1` requires a start frame and optionally accepts an end frame. Its
  supported durations are 5 and 10 seconds. `kling-2.5` accepts both frames,
  while `kling-2.5-turbo-standard` requires a start frame and rejects an end
  frame.
- `kling-2.6`, `kling-3.0`, and `kling-3.0-turbo` support `auto` resolution;
  when used, the request must also set `aspect_ratio` to `auto`.
- `kling-video-o-1` accepts frames or image references. Its video reference
  type is `GENERATED` only; an uploaded URL is rejected. `kling-video-o-3`
  accepts up to seven images or one `GENERATED` video, and when a video is
  supplied it allows at most four images and a maximum duration of 10 seconds.
  Frame mode excludes other references for both O-series models.
- `kling-3.0` and `kling-3.0-turbo` support generated audio. `kling-3.0`
  accepts a start/end frame pair; Turbo accepts a start frame only.
- `veo-3.1-generate-001` supports up to three images, while
  `veo-3.1-fast-generate-001` and `veo-3.1-lite` do not accept image
  references. All three Veo models support start/end frames, generated audio,
  16:9 or 9:16, and only 4, 6, or 8 seconds.
- `happy-horse-1.1` allows a start frame and up to nine reference images. It
  does not accept an end frame, reference video, or reference audio. Its
  optional `prompt_enhance` value is `AUTO`, `ON`, or `OFF`; `ON` cannot be
  combined with a start frame.
- `grok-imagine-1.5` requires exactly one start frame and does not accept an
  end frame, reference image, reference video, or reference audio. It does not
  accept `prompt_enhance` or other unsupported guidance options.
- Both LTX 2.3 models accept one start frame and one end frame together, but no
  image, video, or audio references. They support generated audio and
  `prompt_enhance`; `ON` may be combined with a start frame. Do not send
  `seed` or `mode`.

## Media upload contract

The customer upload endpoint is `POST /v1/videos/uploads`. Uploaded media is
returned as a public-contract `media_url`; use that value in the generation
request.

### Images

- The form field is `image` (the legacy `file` field remains accepted).
- The file must be a readable PNG, JPEG, or WebP image.
- The size limit is 10 MiB per file.
- A generation may contain at most four reference images for Seedance or nine
  for Happy Horse. Grok and LTX do not accept reference-image guidance.

### Reference videos

- The form field is `video`.
- The filename must end in `.mp4` or `.mov`.
- The bytes must be an ISO Base Media container with an `ftyp` header and must
  decode as a readable video; renaming an unsupported file is not sufficient.
- The size limit is 100 MiB per file. The number and type of video references
  are model-specific; uploaded video is accepted only where the model allows
  `UPLOADED` video.
- The public contract does not publish one universal duration, frame-rate,
  codec, audio-track, or dimension whitelist. A file that passes the container
  check can still be rejected by a selected model during processing.
- Uploaded video guidance is available to Seedance models. Kling O-1/O-3 use
  a separate `GENERATED`-asset-only contract.

### Reference audio

- The form field is `audio`.
- MP3 files must contain readable MP3 frames.
- WAV files must be RIFF/WAVE PCM with 16-bit or 24-bit samples. Compressed or
  floating-point WAV files are rejected.
- The size limit is 15 MiB per file, the duration must be 2-30 seconds. A
  generation may contain at most one reference audio item for Seedance and up
  to three for Hailuo 03; Hailuo's combined duration limit is 15 seconds.
- A reference audio item must be paired with at least one reference image or
  reference video. Generated output audio (`audio: true`) is a separate option.
- Reference audio guidance is available to Seedance and Hailuo 03. Other
  models reject audio-reference guidance.

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

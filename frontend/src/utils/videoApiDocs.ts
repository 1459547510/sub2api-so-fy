import { buildGatewayUrl } from '@/api/client'
import { activeSeedanceV2Docs } from './seedanceV2DocsCatalog'

// Seedance V2 matrix/examples follow frontend/src/utils/seedanceV2DocsSource.ts.

export type VideoModelMatrixRow = {
  model: string
  resolution: string
  duration: string
  aspectRatio: string
  promptLimit: string
  references: string
}

export type VideoModelExample = {
  model: string
  description: string
  code: string
}

const matrixKeys = [
  ['seedance-2.0', 'seedance20'],
  ['seedance-2.0-fast', 'seedance20Fast'],
  ['bytedance/seedance-2.5', 'seedance25'],
  ['seedance-2.5', 'seedance25'],
  ['happy-horse-1.1', 'happyHorse'],
  ['grok-imagine-1.5', 'grokImagine'],
  ['ltx-2.3-pro', 'ltx23Pro'],
  ['ltx-2.3-fast', 'ltx23Fast'],
  ['hailuo-03', 'hailuo03'],
  ['gemini-omni-flash', 'geminiOmniFlash'],
  ['kling-2.1', 'kling21'],
  ['kling-2.5', 'kling25'],
  ['kling-2.5-turbo-standard', 'kling25TurboStandard'],
  ['kling-2.6', 'kling26'],
  ['kling-video-o-1', 'klingVideoO1'],
  ['kling-3.0', 'kling30'],
  ['kling-3.0-turbo', 'kling30Turbo'],
  ['kling-video-o-3', 'klingVideoO3'],
  ['veo-3.1-generate-001', 'veo31Generate'],
  ['veo-3.1-fast-generate-001', 'veo31Fast'],
  ['veo-3.1-lite', 'veo31Lite'],
] as const

const v2SeedanceMatrixKeys = new Set(activeSeedanceV2Docs.v2MatrixKeys)

export const videoModelMatrixRows: VideoModelMatrixRow[] = matrixKeys.map(([model, key]) => ({
  model,
  resolution: `video.apiDocs.matrix.${key}.resolution`,
  duration: `video.apiDocs.matrix.${key}.duration`,
  aspectRatio: `video.apiDocs.matrix.${key}.aspectRatio`,
  promptLimit: `video.apiDocs.matrix.${key}.promptLimit`,
  references: `video.apiDocs.matrix.${key}.references`,
}))

export const v2VideoModelMatrixRows: VideoModelMatrixRow[] = matrixKeys.map(([model, key]) => {
  const prefix = v2SeedanceMatrixKeys.has(key) ? 'video.apiDocs.v2.matrix' : 'video.apiDocs.matrix'
  return {
    model,
    resolution: `${prefix}.${key}.resolution`,
    duration: `${prefix}.${key}.duration`,
    aspectRatio: `${prefix}.${key}.aspectRatio`,
    promptLimit: `${prefix}.${key}.promptLimit`,
    references: `${prefix}.${key}.references`,
  }
})

export function resolveVideoApiDocsBaseUrl() {
  try {
    const fallback = typeof window === 'undefined' ? 'https://your-sub2-domain.example' : window.location.origin
    return new URL(buildGatewayUrl('/v1/videos/generations'), fallback).origin
  } catch {
    return 'https://your-sub2-domain.example'
  }
}

export function videoApiAuthHeaders() {
  return `-H "Authorization: Bearer $SUB2_API_KEY" \\
  -H "Content-Type: application/json"`
}

export function videoModelRequestExample(baseUrl: string, body: string) {
  return `curl -X POST "${baseUrl}/v1/videos/generations" \\
  ${videoApiAuthHeaders()} \\
  -H "Prefer: respond-async" \\
  -d '${body}'`
}

const modelExampleBodies: Array<{ model: string; description: string; body: string }> = [
  {
    model: 'seedance-2.0',
    description: 'video.apiDocs.models.seedance20',
    body: `{
  "model": "seedance-2.0",
  "prompt": "A slow aerial shot over a coastal city at sunrise",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "audio": false
}`,
  },
  {
    model: 'seedance-2.0-fast',
    description: 'video.apiDocs.models.seedance20Fast',
    body: `{
  "model": "seedance-2.0-fast",
  "prompt": "A fast tracking shot through a neon city street",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "audio": false
}`,
  },
  {
    model: 'bytedance/seedance-2.5',
    description: 'video.apiDocs.models.seedance25',
    body: `{
  "model": "bytedance/seedance-2.5",
  "prompt": "A wide tracking shot across a futuristic coastal city",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "21:9",
  "audio": true
}`,
  },
  {
    model: 'seedance-2.5',
    description: 'video.apiDocs.models.seedance25',
    body: `{
  "model": "seedance-2.5",
  "prompt": "A wide tracking shot across a futuristic coastal city",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "21:9",
  "audio": true
}`,
  },
  {
    model: 'happy-horse-1.1',
    description: 'video.apiDocs.models.happyHorse',
    body: `{
  "model": "happy-horse-1.1",
  "prompt": "A horse runs through a sunlit meadow",
  "resolution": "720p",
  "duration": 3,
  "aspect_ratio": "16:9",
  "audio": false,
  "prompt_enhance": "OFF",
  "start_frame_url": "https://media.example/start-frame.png"
}`,
  },
  {
    model: 'grok-imagine-1.5',
    description: 'video.apiDocs.models.grokImagine',
    body: `{
  "model": "grok-imagine-1.5",
  "prompt": "A cinematic camera move through a neon city",
  "resolution": "400p",
  "duration": 3,
  "aspect_ratio": "16:9",
  "audio": false,
  "start_frame_url": "https://media.example/start-frame.png"
}`,
  },
  {
    model: 'ltx-2.3-pro',
    description: 'video.apiDocs.models.ltx23Pro',
    body: `{
  "model": "ltx-2.3-pro",
  "prompt": "A cinematic mountain landscape with synchronized ambience",
  "resolution": "2160p",
  "duration": 10,
  "aspect_ratio": "16:9",
  "audio": true,
  "prompt_enhance": "ON",
  "start_frame_url": "https://media.example/start-frame.png",
  "end_frame_url": "https://media.example/end-frame.png"
}`,
  },
  {
    model: 'ltx-2.3-fast',
    description: 'video.apiDocs.models.ltx23Fast',
    body: `{
  "model": "ltx-2.3-fast",
  "prompt": "A continuous tracking shot through a futuristic city",
  "resolution": "1440p",
  "duration": 20,
  "aspect_ratio": "16:9",
  "audio": true,
  "prompt_enhance": "AUTO"
}`,
  },
  {
    model: 'hailuo-03',
    description: 'video.apiDocs.models.hailuo03',
    body: `{
  "model": "hailuo-03",
  "prompt": "A paper lantern floats above a quiet river at night",
  "resolution": "1440p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "audio": true,
  "guidances": {
    "image_reference": [{
      "image": { "url": "https://media.example/lantern.png", "type": "UPLOADED" },
      "strength": "MID"
    }],
    "audio_reference": [{
      "audio": { "url": "https://media.example/river.mp3", "type": "UPLOADED" }
    }]
  }
}`,
  },
  {
    model: 'gemini-omni-flash',
    description: 'video.apiDocs.models.geminiOmniFlash',
    body: `{
  "model": "gemini-omni-flash",
  "prompt": "A watercolor train crosses a mountain bridge",
  "resolution": "720p",
  "duration": 5,
  "aspect_ratio": "16:9",
  "audio": false,
  "guidances": {
    "image_reference": [{
      "image": { "url": "https://media.example/train.png", "type": "UPLOADED" }
    }]
  }
}`,
  },
  {
    model: 'kling-2.1',
    description: 'video.apiDocs.models.kling21',
    body: `{
  "model": "kling-2.1",
  "prompt": "The camera slowly circles the dancer",
  "resolution": "1080p",
  "duration": 5,
  "aspect_ratio": "16:9",
  "audio": false,
  "start_frame_url": "https://media.example/dancer.png",
  "prompt_enhance": "ON"
}`,
  },
  {
    model: 'kling-2.5',
    description: 'video.apiDocs.models.kling25',
    body: `{
  "model": "kling-2.5",
  "prompt": "Transition from a city street into a star field",
  "resolution": "1080p",
  "duration": 10,
  "aspect_ratio": "16:9",
  "audio": false,
  "start_frame_url": "https://media.example/city.png",
  "end_frame_url": "https://media.example/stars.png"
}`,
  },
  {
    model: 'kling-2.5-turbo-standard',
    description: 'video.apiDocs.models.kling25TurboStandard',
    body: `{
  "model": "kling-2.5-turbo-standard",
  "prompt": "A close-up of a flower opening in morning light",
  "resolution": "720p",
  "duration": 5,
  "aspect_ratio": "9:16",
  "audio": false,
  "start_frame_url": "https://media.example/flower.png"
}`,
  },
  {
    model: 'kling-2.6',
    description: 'video.apiDocs.models.kling26',
    body: `{
  "model": "kling-2.6",
  "prompt": "A chef plates a colorful dish in one continuous shot",
  "resolution": "auto",
  "duration": 5,
  "aspect_ratio": "auto",
  "audio": true,
  "start_frame_url": "https://media.example/chef.png"
}`,
  },
  {
    model: 'kling-video-o-1',
    description: 'video.apiDocs.models.klingVideoO1',
    body: `{
  "model": "kling-video-o-1",
  "prompt": "A robot turns toward the camera",
  "resolution": "1080p",
  "duration": 6,
  "aspect_ratio": "16:9",
  "audio": false,
  "start_frame_url": "https://media.example/robot.png"
}`,
  },
  {
    model: 'kling-3.0',
    description: 'video.apiDocs.models.kling30',
    body: `{
  "model": "kling-3.0",
  "prompt": "A sailboat leaves the harbor at golden hour",
  "resolution": "1080p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "audio": true,
  "start_frame_url": "https://media.example/harbor.png",
  "end_frame_url": "https://media.example/open-water.png"
}`,
  },
  {
    model: 'kling-3.0-turbo',
    description: 'video.apiDocs.models.kling30Turbo',
    body: `{
  "model": "kling-3.0-turbo",
  "prompt": "A cyclist rides through a forest trail",
  "resolution": "auto",
  "duration": 8,
  "aspect_ratio": "auto",
  "audio": true,
  "start_frame_url": "https://media.example/cyclist.png"
}`,
  },
  {
    model: 'kling-video-o-3',
    description: 'video.apiDocs.models.klingVideoO3',
    body: `{
  "model": "kling-video-o-3",
  "prompt": "A glass sculpture catches moving sunlight",
  "resolution": "2160p",
  "duration": 8,
  "aspect_ratio": "1:1",
  "audio": true,
  "start_frame_url": "https://media.example/sculpture.png",
  "end_frame_url": "https://media.example/sculpture-end.png"
}`,
  },
  {
    model: 'veo-3.1-generate-001',
    description: 'video.apiDocs.models.veo31Generate',
    body: `{
  "model": "veo-3.1-generate-001",
  "prompt": "A cinematic drone shot over snowy peaks",
  "resolution": "2160p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "audio": true,
  "image_url": "https://media.example/peaks.png"
}`,
  },
  {
    model: 'veo-3.1-fast-generate-001',
    description: 'video.apiDocs.models.veo31Fast',
    body: `{
  "model": "veo-3.1-fast-generate-001",
  "prompt": "Rain falls on a quiet city street",
  "resolution": "1080p",
  "duration": 6,
  "aspect_ratio": "9:16",
  "audio": true,
  "start_frame_url": "https://media.example/street.png",
  "end_frame_url": "https://media.example/street-end.png"
}`,
  },
  {
    model: 'veo-3.1-lite',
    description: 'video.apiDocs.models.veo31Lite',
    body: `{
  "model": "veo-3.1-lite",
  "prompt": "A warm breeze moves through a field of wheat",
  "resolution": "720p",
  "duration": 4,
  "aspect_ratio": "16:9",
  "audio": true,
  "start_frame_url": "https://media.example/wheat.png"
}`,
  },
]

const v2SeedanceExampleOverrides = activeSeedanceV2Docs.exampleOverrides

export function buildVideoModelExamples(baseUrl: string): VideoModelExample[] {
  return modelExampleBodies.map((example) => ({
    model: example.model,
    description: example.description,
    code: videoModelRequestExample(baseUrl, example.body),
  }))
}

export function buildV2VideoModelExamples(baseUrl: string): VideoModelExample[] {
  return modelExampleBodies.map((example) => {
    const override = v2SeedanceExampleOverrides[example.model]
    return {
      model: example.model,
      description: override?.description ?? example.description,
      code: videoModelRequestExample(baseUrl, override?.body ?? example.body),
    }
  })
}

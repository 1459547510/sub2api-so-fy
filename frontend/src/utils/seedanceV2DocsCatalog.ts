import { SEEDANCE_V2_DOCS_SOURCE, type SeedanceV2DocsSource } from './seedanceV2DocsSource'

type SeedanceMatrixCopy = {
  resolution: string
  duration: string
  aspectRatio: string
  promptLimit: string
  references: string
}

type SeedanceLocaleCopy = {
  matrixDescription: string
  seedance20: SeedanceMatrixCopy
  seedance20Fast: SeedanceMatrixCopy
  seedance20Mini: SeedanceMatrixCopy
  seedance25?: SeedanceMatrixCopy
  models: {
    seedance20: string
    seedance20Fast: string
    seedance20Mini: string
    seedance25?: string
  }
}

export type SeedanceV2DocsCatalog = {
  v2MatrixKeys: string[]
  exampleOverrides: Record<string, { description: string; body: string }>
  zh: SeedanceLocaleCopy
  en: SeedanceLocaleCopy
}

const promptLimitZh = '本网关按 5000 字符校验'
const promptLimitEn = 'This gateway enforces a 5000-character limit'

const seedance20Example = (duration: number) => `{
  "model": "seedance-2.0",
  "prompt": "A slow aerial shot over a coastal city at sunrise",
  "resolution": "4k",
  "duration": ${duration},
  "aspect_ratio": "16:9",
  "audio": false
}`

const seedance20FastExample = (duration: number) => `{
  "model": "seedance-2.0-fast",
  "prompt": "A fast tracking shot through a neon city street",
  "resolution": "720p",
  "duration": ${duration},
  "aspect_ratio": "21:9",
  "audio": false
}`

const seedance20MiniExample = (resolution: string, duration: number) => `{
  "model": "seedance-2.0-mini",
  "prompt": "A paper boat drifting across a calm lake",
  "resolution": "${resolution}",
  "duration": ${duration},
  "aspect_ratio": "21:9",
  "audio": false
}`

const seedance25Example = `{
  "model": "seedance-2.5",
  "prompt": "A wide tracking shot across a futuristic coastal city",
  "resolution": "720p",
  "duration": 4,
  "aspect_ratio": "21:9",
  "audio": true
}`

const bytedanceSeedance25Example = `{
  "model": "bytedance/seedance-2.5",
  "prompt": "A wide tracking shot across a futuristic coastal city",
  "resolution": "720p",
  "duration": 4,
  "aspect_ratio": "21:9",
  "audio": true
}`

/** 上一份 Seedance V2 文档。 */
const previousCatalog: SeedanceV2DocsCatalog = {
  v2MatrixKeys: ['seedance20', 'seedance20Fast', 'seedance20Mini'],
  exampleOverrides: {
    'seedance-2.0': { description: 'video.apiDocs.v2.models.seedance20', body: seedance20Example(5) },
    'seedance-2.0-fast': { description: 'video.apiDocs.v2.models.seedance20Fast', body: seedance20FastExample(5) },
    'seedance-2.0-mini': { description: 'video.apiDocs.v2.models.seedance20Mini', body: seedance20MiniExample('720p', 5) },
  },
  zh: {
    matrixDescription: '先通过 /v1/models 确认当前 Key 可用的模型，再按下表选择该模型允许的分辨率、时长、画面比例和参考输入。Seedance 2.0 / Fast / Mini 按时长 4-15 秒（默认 5 秒）校验，各分辨率画面比例相同。不支持的组合会返回 400 或 422。',
    seedance20: {
      resolution: '480p、720p、1080p、4k',
      duration: '4-15 秒，默认 5 秒；各分辨率相同',
      aspectRatio: '各分辨率均为 16:9、4:3、1:1、3:4、9:16、21:9',
      promptLimit: promptLimitZh,
      references: '首帧 1、尾帧 1、参考图 9、参考视频 3、参考音频 3；可与首尾帧同时使用',
    },
    seedance20Fast: {
      resolution: '480p、720p',
      duration: '4-15 秒，默认 5 秒',
      aspectRatio: '各分辨率均为 16:9、4:3、1:1、3:4、9:16、21:9',
      promptLimit: promptLimitZh,
      references: '首帧 1、尾帧 1、参考图 9、参考视频 3、参考音频 3；可与首尾帧同时使用',
    },
    seedance20Mini: {
      resolution: '480p、720p',
      duration: '4-15 秒，默认 5 秒',
      aspectRatio: '各分辨率均为 16:9、4:3、1:1、3:4、9:16、21:9',
      promptLimit: promptLimitZh,
      references: '首帧 1、尾帧 1、参考图 9、参考视频 3、参考音频 3；可与首尾帧同时使用',
    },
    models: {
      seedance20: 'Seedance 2.0：支持 480p/720p/1080p/4k，时长 4-15 秒，默认 5 秒；画面比例为 16:9、4:3、1:1、3:4、9:16、21:9。最多 9 张参考图、3 个参考视频和 3 个参考音频。本示例为 4k 文生视频。',
      seedance20Fast: 'Seedance 2.0 Fast：仅 480p/720p，时长 4-15 秒，默认 5 秒；画面比例与标准版相同。最多 9 张参考图、3 个参考视频和 3 个参考音频。',
      seedance20Mini: 'Seedance 2.0 Mini：仅 480p/720p，时长 4-15 秒，默认 5 秒；画面比例与标准版相同，不再限制为 16:9/1:1/9:16。最多 9 张参考图、3 个参考视频和 3 个参考音频。',
    },
  },
  en: {
    matrixDescription: 'Read /v1/models first to confirm the models available to the current Key, then choose a supported resolution, duration, aspect ratio, and reference input from the matching row. Seedance 2.0 / Fast / Mini accept 4-15 seconds (default 5 seconds) with the same aspect ratios at every listed resolution. Invalid combinations return 400 or 422.',
    seedance20: {
      resolution: '480p, 720p, 1080p, 4k',
      duration: '4-15s, default 5s; the same range at every resolution',
      aspectRatio: '16:9, 4:3, 1:1, 3:4, 9:16, 21:9 at every resolution',
      promptLimit: promptLimitEn,
      references: '1 start frame, 1 end frame, 9 images, 3 videos, 3 audio references; frames may be combined with other references',
    },
    seedance20Fast: {
      resolution: '480p, 720p',
      duration: '4-15s, default 5s',
      aspectRatio: '16:9, 4:3, 1:1, 3:4, 9:16, 21:9 at every resolution',
      promptLimit: promptLimitEn,
      references: '1 start frame, 1 end frame, 9 images, 3 videos, 3 audio references; frames may be combined with other references',
    },
    seedance20Mini: {
      resolution: '480p, 720p',
      duration: '4-15s, default 5s',
      aspectRatio: '16:9, 4:3, 1:1, 3:4, 9:16, 21:9 at every resolution',
      promptLimit: promptLimitEn,
      references: '1 start frame, 1 end frame, 9 images, 3 videos, 3 audio references; frames may be combined with other references',
    },
    models: {
      seedance20: 'Seedance 2.0 supports 480p/720p/1080p/4k for 4-15 seconds, default 5 seconds, at 16:9, 4:3, 1:1, 3:4, 9:16, or 21:9. Up to 9 image, 3 video, and 3 audio references are supported. This example is 4k text-to-video.',
      seedance20Fast: 'Seedance 2.0 Fast supports 480p/720p for 4-15 seconds, default 5 seconds, with the same aspect ratios as standard. Up to 9 image, 3 video, and 3 audio references are supported.',
      seedance20Mini: 'Seedance 2.0 Mini supports 480p/720p for 4-15 seconds, default 5 seconds, with the same aspect ratios as standard. It is no longer limited to 16:9/1:1/9:16. Up to 9 image, 3 video, and 3 audio references are supported.',
    },
  },
}

/** 当前 Seedance V2 文档。 */
const currentCatalog: SeedanceV2DocsCatalog = {
  v2MatrixKeys: ['seedance20', 'seedance20Fast', 'seedance20Mini', 'seedance25'],
  exampleOverrides: {
    'seedance-2.0': { description: 'video.apiDocs.v2.models.seedance20', body: seedance20Example(4) },
    'seedance-2.0-fast': { description: 'video.apiDocs.v2.models.seedance20Fast', body: seedance20FastExample(4) },
    'seedance-2.0-mini': { description: 'video.apiDocs.v2.models.seedance20Mini', body: seedance20MiniExample('1080p', 4) },
    'bytedance/seedance-2.5': { description: 'video.apiDocs.v2.models.seedance25', body: bytedanceSeedance25Example },
    'seedance-2.5': { description: 'video.apiDocs.v2.models.seedance25', body: seedance25Example },
  },
  zh: {
    matrixDescription: '先通过 /v1/models 确认当前 Key 可用的模型，再按下表选择该模型允许的分辨率、时长、画面比例和参考输入。Seedance 2.0 / Fast 时长 4-15 秒（默认 4 秒）；Mini 时长 4-12 秒（默认 4 秒）并支持 1080p；2.5 仅支持离散时长 4、5、6、8、10、12、15、20、25、30 秒（默认 4 秒）。各分辨率画面比例相同。不支持的组合会返回 400 或 422。',
    seedance20: {
      resolution: '480p、720p、1080p、4k，默认 480p',
      duration: '4-15 秒，默认 4 秒；各分辨率相同',
      aspectRatio: '各分辨率均为 16:9、9:16、1:1、4:3、3:4、21:9',
      promptLimit: promptLimitZh,
      references: '首帧 1、尾帧 1、参考图 12、参考视频 3、参考音频 3；可与首尾帧同时使用',
    },
    seedance20Fast: {
      resolution: '480p、720p，默认 480p',
      duration: '4-15 秒，默认 4 秒',
      aspectRatio: '各分辨率均为 16:9、9:16、1:1、4:3、3:4、21:9',
      promptLimit: promptLimitZh,
      references: '首帧 1、尾帧 1、参考图 12、参考视频 3、参考音频 3；可与首尾帧同时使用',
    },
    seedance20Mini: {
      resolution: '480p、720p、1080p，默认 480p',
      duration: '4-12 秒，默认 4 秒',
      aspectRatio: '各分辨率均为 16:9、9:16、1:1、4:3、3:4、21:9',
      promptLimit: promptLimitZh,
      references: '首帧 1、尾帧 1，图片合计最多 2 张；不支持参考视频和参考音频',
    },
    seedance25: {
      resolution: '480p、720p，默认 480p',
      duration: '4、5、6、8、10、12、15、20、25、30 秒，默认 4 秒',
      aspectRatio: '各分辨率均为 16:9、9:16、1:1、4:3、3:4、21:9',
      promptLimit: promptLimitZh,
      references: '首帧 1、尾帧 1、参考图 30、参考视频 3、参考音频 3；可与首尾帧同时使用',
    },
    models: {
      seedance20: 'Seedance 2.0：支持 480p/720p/1080p/4k，默认 480p；时长 4-15 秒，默认 4 秒；画面比例为 16:9、9:16、1:1、4:3、3:4、21:9。最多 12 张参考图、3 个参考视频和 3 个参考音频。本示例为 4k 文生视频。',
      seedance20Fast: 'Seedance 2.0 Fast：仅 480p/720p，默认 480p；时长 4-15 秒，默认 4 秒；画面比例与标准版相同。最多 12 张参考图、3 个参考视频和 3 个参考音频。',
      seedance20Mini: 'Seedance 2.0 Mini：支持 480p/720p/1080p，默认 480p；时长 4-12 秒，默认 4 秒；画面比例与标准版相同。图片合计最多 2 张，不支持参考视频和参考音频。本示例为 1080p 文生视频。',
      seedance25: 'Seedance 2.5：仅 480p/720p，默认 480p；时长只能为 4、5、6、8、10、12、15、20、25、30 秒，默认 4 秒；画面比例与 2.0 相同。最多 30 张参考图、3 个参考视频和 3 个参考音频。',
    },
  },
  en: {
    matrixDescription: 'Read /v1/models first to confirm the models available to the current Key, then choose a supported resolution, duration, aspect ratio, and reference input from the matching row. Seedance 2.0 / Fast accept 4-15 seconds (default 4 seconds). Mini accepts 4-12 seconds (default 4 seconds) including 1080p. 2.5 accepts only 4, 5, 6, 8, 10, 12, 15, 20, 25, or 30 seconds (default 4 seconds). Aspect ratios are the same at every listed resolution. Invalid combinations return 400 or 422.',
    seedance20: {
      resolution: '480p, 720p, 1080p, 4k, default 480p',
      duration: '4-15s, default 4s; the same range at every resolution',
      aspectRatio: '16:9, 9:16, 1:1, 4:3, 3:4, 21:9 at every resolution',
      promptLimit: promptLimitEn,
      references: '1 start frame, 1 end frame, 12 images, 3 videos, 3 audio references; frames may be combined with other references',
    },
    seedance20Fast: {
      resolution: '480p, 720p, default 480p',
      duration: '4-15s, default 4s',
      aspectRatio: '16:9, 9:16, 1:1, 4:3, 3:4, 21:9 at every resolution',
      promptLimit: promptLimitEn,
      references: '1 start frame, 1 end frame, 12 images, 3 videos, 3 audio references; frames may be combined with other references',
    },
    seedance20Mini: {
      resolution: '480p, 720p, 1080p, default 480p',
      duration: '4-12s, default 4s',
      aspectRatio: '16:9, 9:16, 1:1, 4:3, 3:4, 21:9 at every resolution',
      promptLimit: promptLimitEn,
      references: '1 start frame, 1 end frame, at most 2 images in total; video and audio references are not supported',
    },
    seedance25: {
      resolution: '480p, 720p, default 480p',
      duration: '4, 5, 6, 8, 10, 12, 15, 20, 25, or 30s, default 4s',
      aspectRatio: '16:9, 9:16, 1:1, 4:3, 3:4, 21:9 at every resolution',
      promptLimit: promptLimitEn,
      references: '1 start frame, 1 end frame, 30 images, 3 videos, 3 audio references; frames may be combined with other references',
    },
    models: {
      seedance20: 'Seedance 2.0 supports 480p/720p/1080p/4k, default 480p, for 4-15 seconds, default 4 seconds, at 16:9, 9:16, 1:1, 4:3, 3:4, or 21:9. Up to 12 image, 3 video, and 3 audio references are supported. This example is 4k text-to-video.',
      seedance20Fast: 'Seedance 2.0 Fast supports 480p/720p, default 480p, for 4-15 seconds, default 4 seconds, with the same aspect ratios as standard. Up to 12 image, 3 video, and 3 audio references are supported.',
      seedance20Mini: 'Seedance 2.0 Mini supports 480p/720p/1080p, default 480p, for 4-12 seconds, default 4 seconds, with the same aspect ratios as standard. At most 2 images in total; video and audio references are not supported. This example is 1080p text-to-video.',
      seedance25: 'Seedance 2.5 supports 480p/720p, default 480p, and only 4, 5, 6, 8, 10, 12, 15, 20, 25, or 30 seconds, default 4 seconds, with the same aspect ratios as 2.0. Up to 30 image, 3 video, and 3 audio references are supported.',
    },
  },
}

export const seedanceV2DocsCatalog: Record<SeedanceV2DocsSource, SeedanceV2DocsCatalog> = {
  current: currentCatalog,
  previous: previousCatalog,
}

export const activeSeedanceV2Docs = seedanceV2DocsCatalog[SEEDANCE_V2_DOCS_SOURCE]

export function applySeedanceV2DocsToDashboard<T>(messages: T, locale: 'zh' | 'en', source: SeedanceV2DocsSource = SEEDANCE_V2_DOCS_SOURCE): T {
  const root = messages as {
    video?: { apiDocs?: { v2?: { matrix?: Record<string, unknown>; models?: Record<string, unknown> } } }
  }
  const apiDocs = root.video?.apiDocs?.v2
  if (!apiDocs?.matrix || !apiDocs.models) {
    return messages
  }

  const copy = seedanceV2DocsCatalog[source][locale]
  apiDocs.matrix.description = copy.matrixDescription
  apiDocs.matrix.seedance20 = copy.seedance20
  apiDocs.matrix.seedance20Fast = copy.seedance20Fast
  apiDocs.matrix.seedance20Mini = copy.seedance20Mini
  apiDocs.models.seedance20 = copy.models.seedance20
  apiDocs.models.seedance20Fast = copy.models.seedance20Fast
  apiDocs.models.seedance20Mini = copy.models.seedance20Mini

  if (copy.seedance25 && copy.models.seedance25) {
    apiDocs.matrix.seedance25 = copy.seedance25
    apiDocs.models.seedance25 = copy.models.seedance25
  } else {
    delete apiDocs.matrix.seedance25
    delete apiDocs.models.seedance25
  }

  return messages
}

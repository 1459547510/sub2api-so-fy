/** Channel status values (must match service.Status* constants in Go). */
export const CHANNEL_STATUS_ACTIVE = 'active' as const
export const CHANNEL_STATUS_DISABLED = 'disabled' as const
export type ChannelStatus = typeof CHANNEL_STATUS_ACTIVE | typeof CHANNEL_STATUS_DISABLED

/** Billing mode values (must match service.BillingMode* constants in Go). */
export const BILLING_MODE_TOKEN = 'token' as const
export const BILLING_MODE_PER_REQUEST = 'per_request' as const
export const BILLING_MODE_IMAGE = 'image' as const
export const BILLING_MODE_VIDEO = 'video' as const
export type BillingMode =
  | typeof BILLING_MODE_TOKEN
  | typeof BILLING_MODE_PER_REQUEST
  | typeof BILLING_MODE_IMAGE
  | typeof BILLING_MODE_VIDEO

export const LEO_VIDEO_MODELS = [
  'seedance-2.0',
  'seedance-2.0-fast',
  'seedance-2.0-mini',
  'bytedance/seedance-2.5',
  'seedance-2.5',
  'happy-horse-1.1',
  'grok-imagine-1.5',
  'ltx-2.3-pro',
  'ltx-2.3-fast',
  'hailuo-03',
  'gemini-omni-flash',
  'kling-2.1',
  'kling-2.5',
  'kling-2.5-turbo-standard',
  'kling-2.6',
  'kling-video-o-1',
  'kling-3.0',
  'kling-3.0-turbo',
  'kling-video-o-3',
  'veo-3.1-generate-001',
  'veo-3.1-fast-generate-001',
  'veo-3.1-lite',
] as const
export const VIDEO_PRICING_RESOLUTIONS = ['480p', '720p', '1080p'] as const

/** Billing-model-source values (must match service.BillingModelSource* constants in Go). */
export const BILLING_MODEL_SOURCE_REQUESTED = 'requested' as const
export const BILLING_MODEL_SOURCE_UPSTREAM = 'upstream' as const
export const BILLING_MODEL_SOURCE_CHANNEL_MAPPED = 'channel_mapped' as const
export const BILLING_MODEL_SOURCE_RESPONSE = 'response_model' as const
export type BillingModelSource =
  | typeof BILLING_MODEL_SOURCE_REQUESTED
  | typeof BILLING_MODEL_SOURCE_UPSTREAM
  | typeof BILLING_MODEL_SOURCE_CHANNEL_MAPPED
  | typeof BILLING_MODEL_SOURCE_RESPONSE

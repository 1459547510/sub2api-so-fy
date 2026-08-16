export function configuredPrice(value: unknown): number | null {
  if (value === null || value === undefined || value === '') return null
  const price = Number(value)
  return Number.isFinite(price) && price >= 0 ? price : null
}

export type MediaPriceTier = {
  label: string | null
  value: number
}

const IMAGE_TIER_FIELDS = [
  { label: '1K', field: 'image_price_1k' },
  { label: '2K', field: 'image_price_2k' },
  { label: '4K', field: 'image_price_4k' },
] as const

export type ImagePriceFields = {
  image_price_1k: unknown
  image_price_2k: unknown
  image_price_4k: unknown
}

export function imageTiers(group: ImagePriceFields): MediaPriceTier[] {
  const configured = IMAGE_TIER_FIELDS
    .map((tier) => ({ label: tier.label, value: configuredPrice(group[tier.field]) }))
    .filter((tier): tier is { label: string; value: number } => tier.value !== null)

  if (configured.length >= 2 && configured.every((tier) => tier.value === configured[0].value)) {
    return [{ label: null, value: configured[0].value }]
  }
  return configured
}

export const MEDIA_PRICING_PLATFORMS = new Set([
  'leo',
  'openai_media',
  'video',
  'composite',
  'grok',
])

export const VIDEO_RESOLUTIONS = [
  { label: '480p', field: 'video_price_480p' },
  { label: '720p', field: 'video_price_720p' },
  { label: '1080p', field: 'video_price_1080p' },
] as const

export type VideoPriceFields = ImagePriceFields & {
  video_price_480p: unknown
  video_price_720p: unknown
  video_price_1080p: unknown
  video_model_prices?: Record<string, Record<string, unknown>> | null
}

export type MediaPriceCard = {
  key: string
  title: string
  kind: 'image' | 'video'
  unit: 'per_image' | 'per_second'
  tiers: MediaPriceTier[]
}

export function hasConfiguredImagePrice(group: ImagePriceFields): boolean {
  return imageTiers(group).length > 0
}

export function flatVideoTiers(group: VideoPriceFields): MediaPriceTier[] {
  return VIDEO_RESOLUTIONS
    .map((tier) => ({ label: tier.label, value: configuredPrice(group[tier.field]) }))
    .filter((tier): tier is { label: string; value: number } => tier.value !== null)
}

function modelOverridePrice(
  overrides: Record<string, unknown> | undefined,
  resolution: string,
): number | null {
  return configuredPrice(overrides?.[resolution])
}

export function hasConfiguredVideoModelOverride(group: VideoPriceFields): boolean {
  for (const overrides of Object.values(group.video_model_prices ?? {})) {
    if (VIDEO_RESOLUTIONS.some((tier) => modelOverridePrice(overrides, tier.label) !== null)) {
      return true
    }
  }
  return false
}

export function keepMediaPricingGroup(group: VideoPriceFields & { platform: string }): boolean {
  if (!MEDIA_PRICING_PLATFORMS.has(group.platform)) return false
  return hasConfiguredImagePrice(group)
    || flatVideoTiers(group).length > 0
    || hasConfiguredVideoModelOverride(group)
}

export function buildVideoCards(group: VideoPriceFields & { name: string }): MediaPriceCard[] {
  const keys = Object.keys(group.video_model_prices ?? {})
    .filter((model) => VIDEO_RESOLUTIONS.some((tier) => (
      modelOverridePrice(group.video_model_prices?.[model], tier.label) !== null
    )))
    .sort((a, b) => a.localeCompare(b))

  if (keys.length > 0) {
    return keys.map((model) => ({
      key: `video:${model}`,
      title: model,
      kind: 'video',
      unit: 'per_second',
      tiers: VIDEO_RESOLUTIONS
        .map((tier) => ({
          label: tier.label,
          value: modelOverridePrice(group.video_model_prices?.[model], tier.label)
            ?? configuredPrice(group[tier.field]),
        }))
        .filter((tier): tier is { label: string; value: number } => tier.value !== null),
    }))
  }

  const flats = flatVideoTiers(group)
  if (flats.length === 0) return []
  return [{
    key: 'video-fallback',
    title: group.name,
    kind: 'video',
    unit: 'per_second',
    tiers: flats,
  }]
}

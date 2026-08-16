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

function collectConfiguredTiers<T extends { label: string }>(
  items: ReadonlyArray<T>,
  getValue: (item: T) => number | null,
): MediaPriceTier[] {
  const tiers: MediaPriceTier[] = []
  for (const item of items) {
    const value = getValue(item)
    if (value !== null) {
      tiers.push({ label: item.label, value })
    }
  }
  return tiers
}

export function imageTiers(group: ImagePriceFields): MediaPriceTier[] {
  const configured = collectConfiguredTiers(
    IMAGE_TIER_FIELDS,
    (tier) => configuredPrice(group[tier.field]),
  )

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
  return collectConfiguredTiers(
    VIDEO_RESOLUTIONS,
    (tier) => configuredPrice(group[tier.field]),
  )
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

export type MediaModelRef = {
  name: string
  platform?: string
  pricing?: { billing_mode?: string } | null
}

export type PlazaImageSource = {
  id: number
  models: MediaModelRef[]
}

export type ChannelModelSource = {
  platforms?: Array<{
    groups?: Array<{ id: number }>
    supported_models?: MediaModelRef[]
  }>
}

export function isMediaImageModel(model: MediaModelRef): boolean {
  return model.pricing?.billing_mode === 'image'
}

export function isMediaVideoModel(model: MediaModelRef): boolean {
  if (model.pricing?.billing_mode === 'image') return false
  if (model.pricing?.billing_mode === 'video') return true
  return MEDIA_PRICING_PLATFORMS.has(model.platform ?? '')
}

export function collectGroupModels(
  groupId: number,
  plazaGroups: PlazaImageSource[] = [],
  channels: ChannelModelSource[] = [],
): MediaModelRef[] {
  const collected: MediaModelRef[] = [
    ...(plazaGroups.find((item) => item.id === groupId)?.models ?? []),
  ]
  for (const channel of channels) {
    for (const section of channel.platforms ?? []) {
      if (!(section.groups ?? []).some((item) => item.id === groupId)) continue
      collected.push(...(section.supported_models ?? []))
    }
  }

  const byName = new Map<string, MediaModelRef>()
  for (const model of collected) {
    const name = model.name?.trim()
    if (!name || byName.has(name)) continue
    byName.set(name, model)
  }
  return [...byName.values()]
}

function overrideVideoModels(group: VideoPriceFields): string[] {
  return Object.keys(group.video_model_prices ?? {}).filter((model) => (
    VIDEO_RESOLUTIONS.some((tier) => (
      modelOverridePrice(group.video_model_prices?.[model], tier.label) !== null
    ))
  ))
}

function uniqueSortedNames(names: Iterable<string>): string[] {
  return [...new Set([...names].map((name) => name.trim()).filter(Boolean))]
    .sort((a, b) => a.localeCompare(b))
}

function videoCardForModel(
  group: VideoPriceFields,
  model: string,
): MediaPriceCard | null {
  const tiers = collectConfiguredTiers(
    VIDEO_RESOLUTIONS,
    (tier) => modelOverridePrice(group.video_model_prices?.[model], tier.label)
      ?? configuredPrice(group[tier.field]),
  )
  if (tiers.length === 0) return null
  return {
    key: `video:${model}`,
    title: model,
    kind: 'video',
    unit: 'per_second',
    tiers,
  }
}

export function buildVideoCards(
  group: VideoPriceFields & { name: string },
  models: MediaModelRef[] = [],
): MediaPriceCard[] {
  const names = uniqueSortedNames([
    ...overrideVideoModels(group),
    ...models.filter(isMediaVideoModel).map((model) => model.name),
  ])
  const cards = names
    .map((model) => videoCardForModel(group, model))
    .filter((card): card is MediaPriceCard => card !== null)

  if (cards.length > 0) return cards

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

export type MediaPricingGroup = VideoPriceFields & {
  id: number
  name: string
  platform: string
}

export type MediaPriceGroupSection = {
  id: number
  name: string
  cards: MediaPriceCard[]
}

export function buildImageCards(
  group: MediaPricingGroup,
  models: MediaModelRef[] = [],
): MediaPriceCard[] {
  if (!hasConfiguredImagePrice(group)) return []
  const tiers = imageTiers(group)
  const names = uniqueSortedNames(models.filter(isMediaImageModel).map((model) => model.name))

  if (names.length > 0) {
    return names.map((name) => ({
      key: `image:${name}`,
      title: name,
      kind: 'image',
      unit: 'per_image',
      tiers,
    }))
  }

  return [{
    key: 'image-fallback',
    title: group.name,
    kind: 'image',
    unit: 'per_image',
    tiers,
  }]
}

export function buildMediaPricingSections(
  groups: MediaPricingGroup[],
  plazaGroups: PlazaImageSource[] = [],
  channels: ChannelModelSource[] = [],
): MediaPriceGroupSection[] {
  return groups
    .filter(keepMediaPricingGroup)
    .map((group) => {
      const models = collectGroupModels(group.id, plazaGroups, channels)
      return {
        id: group.id,
        name: group.name,
        cards: [
          ...buildImageCards(group, models),
          ...buildVideoCards(group, models),
        ],
      }
    })
    .filter((section) => section.cards.length > 0)
}

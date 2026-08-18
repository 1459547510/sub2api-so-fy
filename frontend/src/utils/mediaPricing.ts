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

const VIDEO_TIER_ORDER = ['400p', '480p', '544p', '720p', '960p', '1080p', '1440p', '2160p'] as const
const IMAGE_TIER_ORDER = ['1K', '2K', '4K'] as const

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

export type MediaPriceInterval = {
  tier_label?: string
  per_request_price?: number | null
}

export type MediaModelPricing = {
  billing_mode?: string
  per_request_price?: number | null
  intervals?: MediaPriceInterval[]
}

export type MediaModelRef = {
  name: string
  platform?: string
  pricing?: MediaModelPricing | null
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

function pricingIntervalCount(model: MediaModelRef): number {
  return model.pricing?.intervals?.filter((interval) => (
    configuredPrice(interval.per_request_price) !== null && String(interval.tier_label || '').trim()
  )).length ?? 0
}

function preferRicherModel(existing: MediaModelRef, incoming: MediaModelRef): MediaModelRef {
  // Available-channel copies are merged after plaza. Live channel interval
  // prices must win so 渠道定价 edits show up on the public price page.
  if (pricingIntervalCount(incoming) > 0) return incoming
  if (pricingIntervalCount(existing) > 0) return existing
  if (!existing.pricing && incoming.pricing) return incoming
  return existing
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
    if (!name) continue
    const current = byName.get(name)
    byName.set(name, current ? preferRicherModel(current, { ...model, name }) : { ...model, name })
  }
  return [...byName.values()]
}

export function videoPriceFamilyKeys(model: string): string[] {
  const normalized = model.trim().toLowerCase().replace(/^(xai|x-ai|grok)\//, '')
  if (!normalized) return []
  const keys = [normalized]
  if (
    normalized === 'grok-imagine-1.5' ||
    normalized === 'grok-imagine-video-1.5' ||
    normalized === 'grok-imagine-video-1.5-preview' ||
    normalized === 'grok-video-1.5' ||
    normalized.includes('video-1.5')
  ) {
    keys.push('grok-imagine-video-1.5')
  } else if (
    normalized === 'grok-imagine-video' ||
    normalized === 'grok-imagine-video-preview' ||
    normalized === 'grok-video' ||
    normalized === 'grok-video-latest'
  ) {
    keys.push('grok-imagine-video')
  }
  return [...new Set(keys)]
}

function sortVideoTiers(tiers: MediaPriceTier[]): MediaPriceTier[] {
  return [...tiers].sort((left, right) => {
    const leftOrder = VIDEO_TIER_ORDER.indexOf((left.label || '') as typeof VIDEO_TIER_ORDER[number])
    const rightOrder = VIDEO_TIER_ORDER.indexOf((right.label || '') as typeof VIDEO_TIER_ORDER[number])
    if (leftOrder === -1 && rightOrder === -1) return (left.label || '').localeCompare(right.label || '')
    if (leftOrder === -1) return 1
    if (rightOrder === -1) return -1
    return leftOrder - rightOrder
  })
}

function videoTiersFromPricing(pricing: MediaModelPricing | null | undefined): MediaPriceTier[] {
  if (!pricing || pricing.billing_mode === 'image') return []
  const seen = new Set<string>()
  const tiers: MediaPriceTier[] = []
  for (const interval of pricing.intervals ?? []) {
    const label = String(interval.tier_label || '').trim().toLowerCase()
    const value = configuredPrice(interval.per_request_price)
    if (!label || value === null || seen.has(label)) continue
    seen.add(label)
    tiers.push({ label, value })
  }
  return sortVideoTiers(tiers)
}

function imageTiersFromPricing(pricing: MediaModelPricing | null | undefined): MediaPriceTier[] {
  if (pricing?.billing_mode !== 'image') return []
  const byLabel = new Map<string, number>()
  for (const interval of pricing.intervals ?? []) {
    const label = String(interval.tier_label || '').trim()
    const value = configuredPrice(interval.per_request_price)
    if (!label || value === null) continue
    byLabel.set(label, value)
  }
  const fromIntervals = IMAGE_TIER_ORDER
    .filter((label) => byLabel.has(label))
    .map((label) => ({ label, value: byLabel.get(label)! }))
  if (fromIntervals.length > 0) {
    if (fromIntervals.length >= 2 && fromIntervals.every((tier) => tier.value === fromIntervals[0].value)) {
      return [{ label: null, value: fromIntervals[0].value }]
    }
    return fromIntervals
  }
  const single = configuredPrice(pricing.per_request_price)
  return single === null ? [] : [{ label: null, value: single }]
}

function overrideTiersForModel(group: VideoPriceFields, model: string): MediaPriceTier[] {
  const prices = group.video_model_prices ?? {}
  for (const family of videoPriceFamilyKeys(model)) {
    const overrides = prices[family]
    if (!overrides) continue
    const tiers = collectConfiguredTiers(
      VIDEO_RESOLUTIONS,
      (tier) => configuredPrice(overrides[tier.label]),
    )
    if (tiers.length > 0) return tiers
  }
  return []
}

function mergeVideoTiers(
  group: VideoPriceFields,
  channelTiers: MediaPriceTier[],
  overrideTiers: MediaPriceTier[],
): MediaPriceTier[] {
  if (channelTiers.length > 0) {
    return sortVideoTiers(channelTiers)
  }
  if (overrideTiers.length === 0) return []
  const byLabel = new Map(overrideTiers.map((tier) => [tier.label, tier]))
  return collectConfiguredTiers(
    VIDEO_RESOLUTIONS,
    (tier) => byLabel.get(tier.label)?.value ?? configuredPrice(group[tier.field]),
  )
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
  model: MediaModelRef | string,
): MediaPriceCard | null {
  const name = typeof model === 'string' ? model : model.name
  const pricing = typeof model === 'string' ? undefined : model.pricing
  const tiers = mergeVideoTiers(group, videoTiersFromPricing(pricing), overrideTiersForModel(group, name))
  if (tiers.length > 0) {
    return {
      key: `video:${name}`,
      title: name,
      kind: 'video',
      unit: 'per_second',
      tiers,
    }
  }

  return null
}

export function buildVideoCards(
  group: VideoPriceFields & { name: string },
  models: MediaModelRef[] = [],
): MediaPriceCard[] {
  const byName = new Map(
    models.filter(isMediaVideoModel).map((model) => [model.name.trim(), model]),
  )
  const modelNames = [...byName.keys()]
  const extraOverrides = overrideVideoModels(group).filter((family) => (
    !modelNames.some((name) => name === family || videoPriceFamilyKeys(name).includes(family))
  ))
  const names = uniqueSortedNames([
    ...extraOverrides,
    ...modelNames,
  ])
  const cards = names
    .map((name) => videoCardForModel(group, byName.get(name) ?? name))
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

function imageTiersForModel(group: MediaPricingGroup, model: MediaModelRef): MediaPriceTier[] {
  const fromPricing = imageTiersFromPricing(model.pricing)
  if (fromPricing.length > 0) return fromPricing
  return hasConfiguredImagePrice(group) ? imageTiers(group) : []
}

export function buildImageCards(
  group: MediaPricingGroup,
  models: MediaModelRef[] = [],
): MediaPriceCard[] {
  const imageModels = models.filter(isMediaImageModel)
  const names = uniqueSortedNames(imageModels.map((model) => model.name))
  const byName = new Map(imageModels.map((model) => [model.name.trim(), model]))

  if (names.length > 0) {
    return names.flatMap((name) => {
      const model = byName.get(name)
      const tiers = model ? imageTiersForModel(group, model) : imageTiers(group)
      if (tiers.length === 0) return []
      return [{
        key: `image:${name}`,
        title: name,
        kind: 'image' as const,
        unit: 'per_image' as const,
        tiers,
      }]
    })
  }

  if (!hasConfiguredImagePrice(group)) return []
  return [{
    key: 'image-fallback',
    title: group.name,
    kind: 'image',
    unit: 'per_image',
    tiers: imageTiers(group),
  }]
}

export function buildMediaPricingSections(
  groups: MediaPricingGroup[],
  plazaGroups: PlazaImageSource[] = [],
  channels: ChannelModelSource[] = [],
): MediaPriceGroupSection[] {
  return groups
    .filter((group) => MEDIA_PRICING_PLATFORMS.has(group.platform))
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

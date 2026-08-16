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

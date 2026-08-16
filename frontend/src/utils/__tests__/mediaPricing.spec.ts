import { describe, expect, it } from 'vitest'
import { configuredPrice, imageTiers } from '../mediaPricing'

describe('configuredPrice', () => {
  it('treats null, undefined, and empty string as unconfigured before Number()', () => {
    expect(configuredPrice(null)).toBeNull()
    expect(configuredPrice(undefined)).toBeNull()
    expect(configuredPrice('')).toBeNull()
  })

  it('accepts 0 and numeric strings', () => {
    expect(configuredPrice(0)).toBe(0)
    expect(configuredPrice('0.08')).toBe(0.08)
  })

  it('rejects non-numeric and negative values', () => {
    expect(configuredPrice('x')).toBeNull()
    expect(configuredPrice(-1)).toBeNull()
  })
})

describe('imageTiers', () => {
  it('collapses two or three matching configured tiers to one unlabeled price', () => {
    expect(imageTiers({ image_price_1k: 0.02, image_price_2k: 0.02, image_price_4k: 0.02 }))
      .toEqual([{ label: null, value: 0.02 }])
    expect(imageTiers({ image_price_1k: 0.02, image_price_2k: 0.02, image_price_4k: null }))
      .toEqual([{ label: null, value: 0.02 }])
  })

  it('keeps a single configured tier labeled', () => {
    expect(imageTiers({ image_price_1k: null, image_price_2k: null, image_price_4k: 0.16 }))
      .toEqual([{ label: '4K', value: 0.16 }])
  })

  it('splits different values and omits unconfigured tiers', () => {
    expect(imageTiers({ image_price_1k: 0.08, image_price_2k: 0.12, image_price_4k: null }))
      .toEqual([
        { label: '1K', value: 0.08 },
        { label: '2K', value: 0.12 },
      ])
  })
})

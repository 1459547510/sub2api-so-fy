import { describe, expect, it } from 'vitest'
import {
  configuredPrice,
  imageTiers,
  keepMediaPricingGroup,
  buildVideoCards,
} from '../mediaPricing'

const emptyPrices = {
  image_price_1k: null,
  image_price_2k: null,
  image_price_4k: null,
  video_price_480p: null,
  video_price_720p: null,
  video_price_1080p: null,
  video_model_prices: undefined as Record<string, Record<string, unknown>> | undefined,
}

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

describe('keepMediaPricingGroup', () => {
  it('keeps media platforms that have nested video model prices only', () => {
    expect(keepMediaPricingGroup({
      platform: 'grok',
      ...emptyPrices,
      video_model_prices: { 'grok-imagine-video': { '720p': 0.07 } },
    })).toBe(true)
  })

  it('drops non-media platforms and groups with no configured media prices', () => {
    expect(keepMediaPricingGroup({ platform: 'openai', ...emptyPrices, image_price_1k: 0.08 })).toBe(false)
    expect(keepMediaPricingGroup({ platform: 'leo', ...emptyPrices })).toBe(false)
    expect(keepMediaPricingGroup({
      platform: 'leo',
      ...emptyPrices,
      video_model_prices: { 'seedance-2.0': { '4k': 1, '720p': null } },
    })).toBe(false)
  })
})

describe('buildVideoCards', () => {
  it('keeps a model key only when it has a 480p/720p/1080p override', () => {
    const cards = buildVideoCards({
      name: 'Grok media',
      ...emptyPrices,
      video_price_720p: 0.1,
      video_model_prices: {
        'grok-imagine-video': { '720p': 0.07, '4k': 9 },
        empty: { '720p': null },
      },
    })
    expect(cards.map((card) => card.title)).toEqual(['grok-imagine-video'])
    expect(cards[0].tiers).toEqual([{ label: '720p', value: 0.07 }])
  })

  it('fills missing override resolutions from flat prices and does not collapse matches', () => {
    const cards = buildVideoCards({
      name: 'Grok media',
      ...emptyPrices,
      video_price_480p: 0.08,
      video_price_720p: 0.08,
      video_model_prices: { 'grok-imagine-video-1.5': { '1080p': 0.25 } },
    })
    expect(cards[0].tiers).toEqual([
      { label: '480p', value: 0.08 },
      { label: '720p', value: 0.08 },
      { label: '1080p', value: 0.25 },
    ])
  })

  it('uses a group-name fallback only when no override key is kept', () => {
    const fallback = buildVideoCards({
      name: 'Leo group',
      ...emptyPrices,
      video_price_720p: 0.14,
      video_model_prices: { leftover: { '4k': 1 } },
    })
    expect(fallback).toEqual([{
      key: 'video-fallback',
      title: 'Leo group',
      kind: 'video',
      unit: 'per_second',
      tiers: [{ label: '720p', value: 0.14 }],
    }])

    const noFallback = buildVideoCards({
      name: 'Leo group',
      ...emptyPrices,
      video_price_480p: 0.08,
      video_model_prices: { 'seedance-2.0': { '720p': 0.14 } },
    })
    expect(noFallback.map((card) => card.title)).toEqual(['seedance-2.0'])
  })
})

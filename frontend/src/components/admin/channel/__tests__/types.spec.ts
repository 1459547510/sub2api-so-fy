import { describe, expect, it } from 'vitest'
import {
  createPricingFormEntry,
  createSyncedPricingEntries,
  formIntervalsToAPI,
  getNextLeoVideoModel,
  normalizeVideoIntervals,
  validateIntervals,
  validateVideoPricing,
  type IntervalFormEntry,
} from '../types'

function makeInterval(over: Partial<IntervalFormEntry>): IntervalFormEntry {
  return {
    min_tokens: 0,
    max_tokens: null,
    tier_label: '',
    input_price: null,
    output_price: null,
    cache_write_price: null,
    cache_read_price: null,
    per_request_price: null,
    sort_order: 0,
    ...over,
  }
}

function t(key: string, params?: Record<string, unknown>): string {
  return `${key}${params ? ` ${JSON.stringify(params)}` : ''}`
}

describe('validateIntervals', () => {
  describe('token mode', () => {
    it('rejects unbounded interval that is not last', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: null, input_price: 1, output_price: 1 }),
        makeInterval({ min_tokens: 200000, max_tokens: 500000, input_price: 2, output_price: 2 }),
      ]
      expect(validateIntervals(intervals, 'token', t)).toContain('unboundedLast')
    })

    it('accepts unbounded interval at the end', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: 200000, input_price: 1, output_price: 1 }),
        makeInterval({ min_tokens: 200000, max_tokens: null, input_price: 2, output_price: 2 }),
      ]
      expect(validateIntervals(intervals, 'token', t)).toBeNull()
    })

    it('rejects overlapping intervals', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: 250000, input_price: 1, output_price: 1 }),
        makeInterval({ min_tokens: 200000, max_tokens: 500000, input_price: 2, output_price: 2 }),
      ]
      expect(validateIntervals(intervals, 'token', t)).toContain('overlap')
    })

    it('rejects unbounded interval in token mode', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: null, input_price: 1, output_price: 1 }),
        makeInterval({ min_tokens: 100, max_tokens: 200, input_price: 2, output_price: 2 }),
      ]
      expect(validateIntervals(intervals, 'token', t)).toContain('unboundedLast')
    })
  })

  describe('image / per_request mode', () => {
    it('allows multiple unbounded tiers identified by label', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', per_request_price: 0.04 }),
        makeInterval({ tier_label: '2K', per_request_price: 0.06 }),
        makeInterval({ tier_label: '4K', per_request_price: 0.08 }),
      ]
      expect(validateIntervals(intervals, 'image', t)).toBeNull()
      expect(validateIntervals(intervals, 'per_request', t)).toBeNull()
    })

    it('still rejects negative prices', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', per_request_price: -1 }),
      ]
      expect(validateIntervals(intervals, 'image', t)).toContain('negativePrice')
    })

    it('still rejects max <= min on a single tier', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', min_tokens: 100, max_tokens: 50, per_request_price: 0.04 }),
      ]
      expect(validateIntervals(intervals, 'image', t)).toContain('maxGreaterThanMin')
    })
  })
})

describe('video pricing', () => {
  it('creates fixed video resolution intervals serialized as per-request prices', () => {
    const entry = createPricingFormEntry(['seedance-2.0'], 'video')
    entry.intervals[0].per_request_price = '0.01'
    entry.intervals[1].per_request_price = '0.02'
    entry.intervals[2].per_request_price = '0.03'

    expect(entry.intervals.map(iv => iv.tier_label)).toEqual(['480p', '720p', '1080p'])
    expect(formIntervalsToAPI(entry.intervals).map(iv => iv.per_request_price)).toEqual([0.01, 0.02, 0.03])
  })

  it('normalizes existing video intervals to the supported resolutions and order', () => {
    const intervals = normalizeVideoIntervals([
      makeInterval({ tier_label: '1080P', per_request_price: 0.3 }),
      makeInterval({ tier_label: 'custom', per_request_price: 9 }),
      makeInterval({ tier_label: '480p', per_request_price: 0.1 }),
    ])

    expect(intervals.map(iv => [iv.tier_label, iv.per_request_price])).toEqual([
      ['480p', 0.1],
      ['720p', null],
      ['1080p', 0.3],
    ])
  })

  it('prefills Leo video models in order without duplicating existing models', () => {
    expect(getNextLeoVideoModel([])).toBe('seedance-2.0')
    expect(getNextLeoVideoModel([
      createPricingFormEntry(['SEEDANCE-2.0'], 'video'),
    ])).toBe('seedance-2.0-fast')
    expect(getNextLeoVideoModel([
      createPricingFormEntry(['seedance-2.0'], 'video'),
      createPricingFormEntry(['seedance-2.0-fast'], 'video'),
    ])).toBe('seedance-2.0-mini')
    expect(getNextLeoVideoModel([
      createPricingFormEntry(['seedance-2.0'], 'video'),
      createPricingFormEntry(['seedance-2.0-fast'], 'video'),
      createPricingFormEntry(['seedance-2.0-mini'], 'video'),
    ])).toBeUndefined()
  })

  it('creates a Mini pricing entry with only its supported 720p tier', () => {
    const entry = createPricingFormEntry(['seedance-2.0-mini'], 'video')
    expect(entry.intervals.map(iv => iv.tier_label)).toEqual(['720p'])
    entry.intervals[0].per_request_price = 0.06
    expect(validateVideoPricing(entry, t)).toBeNull()
  })

  it('splits synchronized Leo models into independent video entries', () => {
    const entries = createSyncedPricingEntries('leo', ['seedance-2.0', 'seedance-2.0-fast'])

    expect(entries.map(entry => [entry.models, entry.billing_mode])).toEqual([
      [['seedance-2.0'], 'video'],
      [['seedance-2.0-fast'], 'video'],
    ])
    expect(createSyncedPricingEntries('openai', ['gpt-a', 'gpt-b'])).toHaveLength(1)
  })

  it('requires one model and all three non-negative prices', () => {
    const entry = createPricingFormEntry(['seedance-2.0'], 'video')
    expect(validateVideoPricing(entry, t)).toContain('videoPricesRequired')

    entry.intervals.forEach((interval, index) => {
      interval.per_request_price = index / 100
    })
    expect(validateVideoPricing(entry, t)).toBeNull()

    entry.models.push('seedance-2.0-fast')
    expect(validateVideoPricing(entry, t)).toContain('videoSingleModelRequired')

    entry.models.pop()
    entry.intervals[2].per_request_price = -1
    expect(validateVideoPricing(entry, t)).toContain('videoPricesRequired')
  })
})

import { describe, expect, it } from 'vitest'
import {
  apiIntervalsToForm,
  apiTimePricingToForm,
  createDefaultTimePricingForm,
  createPricingFormEntry,
  createSyncedPricingEntries,
  formIntervalsToAPI,
  formPricedVideoIntervalsToAPI,
  formTimePricingToAPI,
  getNextLeoVideoModel,
  isUnsoldVideoPricingEntry,
  isValidPositiveMultiplier,
  normalizeVideoIntervals,
  validateIntervals,
  validateVideoPricing,
  validateTimePricing,
  type IntervalFormEntry,
  type TimePricingFormEntry,
  type TimePricingPeriodFormEntry,
} from '../types'

describe('interval multiplier conversion', () => {
  it('preserves component multipliers without MTok conversion', () => {
    const form = apiIntervalsToForm([{
      min_tokens: 272000,
      max_tokens: null,
      tier_label: '',
      input_price: null,
      output_price: null,
      cache_write_price: null,
      cache_read_price: null,
      input_multiplier: 2,
      output_multiplier: 1.5,
      cache_write_multiplier: 2,
      cache_read_multiplier: 2,
      per_request_price: null,
      sort_order: 0,
    }])

    expect(form[0].input_multiplier).toBe(2)
    expect(form[0].output_multiplier).toBe(1.5)
    expect(formIntervalsToAPI(form)[0]).toMatchObject({
      input_multiplier: 2,
      output_multiplier: 1.5,
      cache_write_multiplier: 2,
      cache_read_multiplier: 2,
    })
  })
})

describe('positive multiplier validation', () => {
  it('accepts empty and positive values but rejects zero and negative values', () => {
    expect(isValidPositiveMultiplier(null)).toBe(true)
    expect(isValidPositiveMultiplier('')).toBe(true)
    expect(isValidPositiveMultiplier('0.5')).toBe(true)
    expect(isValidPositiveMultiplier(0)).toBe(false)
    expect(isValidPositiveMultiplier(-1)).toBe(false)
  })

  it('rejects a zero interval multiplier', () => {
    expect(validateIntervals([
      makeInterval({ min_tokens: 100, input_multiplier: 0 }),
    ], 'token', t)).toContain('multiplierPositive')
  })
})

function makeInterval(over: Partial<IntervalFormEntry>): IntervalFormEntry {
  return {
    min_tokens: 0,
    max_tokens: null,
    tier_label: '',
    input_price: null,
    output_price: null,
    cache_write_price: null,
    cache_read_price: null,
    input_multiplier: null,
    output_multiplier: null,
    cache_write_multiplier: null,
    cache_read_multiplier: null,
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
    entry.intervals[3].per_request_price = '0.04'

    expect(entry.intervals.map(iv => iv.tier_label)).toEqual(['480p', '720p', '1080p', '2160p'])
    expect(formIntervalsToAPI(entry.intervals).map(iv => iv.per_request_price)).toEqual([0.01, 0.02, 0.03, 0.04])
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
    ])).toBe('bytedance/seedance-2.5')
  })

  it('creates a Mini pricing entry with only its supported 720p tier', () => {
    const entry = createPricingFormEntry(['seedance-2.0-mini'], 'video')
    expect(entry.intervals.map(iv => iv.tier_label)).toEqual(['720p'])
    entry.intervals[0].per_request_price = 0.06
    expect(validateVideoPricing(entry, t)).toBeNull()
  })

  it('creates Happy Horse pricing with only native 720p and 1080p tiers', () => {
    const entry = createPricingFormEntry(['happy-horse-1.1'], 'video')
    expect(entry.intervals.map(iv => iv.tier_label)).toEqual(['720p', '1080p'])
    entry.intervals[0].per_request_price = 0.15
    entry.intervals[1].per_request_price = 0.19
    expect(validateVideoPricing(entry, t)).toBeNull()
  })

  it('creates Grok pricing with only native 400p, 544p, 720p, and 960p tiers', () => {
    const entry = createPricingFormEntry(['grok-imagine-1.5'], 'video')
    expect(entry.intervals.map(iv => iv.tier_label)).toEqual(['400p', '544p', '720p', '960p'])
    entry.intervals.forEach((interval, index) => {
      interval.per_request_price = [0.10, 0.10, 0.17, 0.17][index]
    })
    expect(validateVideoPricing(entry, t)).toBeNull()
  })

  it('creates LTX pricing entries with native 1080p, 1440p, and 2160p tiers', () => {
    for (const model of ['ltx-2.3-pro', 'ltx-2.3-fast']) {
      const entry = createPricingFormEntry([model], 'video')
      expect(entry.intervals.map(iv => iv.tier_label)).toEqual(['1080p', '1440p', '2160p'])
      entry.intervals.forEach((interval, index) => {
        interval.per_request_price = [0.06, 0.21, 0.24][index]
      })
      expect(validateVideoPricing(entry, t)).toBeNull()
    }
  })

  it('creates native pricing tiers for the latest Leo models', () => {
    const expected: Record<string, string[]> = {
      'bytedance/seedance-2.5': ['480p', '720p'],
      'seedance-2.5': ['480p', '720p'],
      'hailuo-03': ['1440p'],
      'gemini-omni-flash': ['720p'],
      'kling-2.1': ['1080p'],
      'kling-2.5': ['720p', '1080p'],
      'kling-2.5-turbo-standard': ['720p'],
      'kling-2.6': ['1080p'],
      'kling-video-o-1': ['1080p'],
      'kling-3.0': ['720p', '1080p', '2160p'],
      'kling-3.0-turbo': ['720p', '1080p'],
      'kling-video-o-3': ['720p', '1080p', '2160p'],
      'veo-3.1-generate-001': ['720p', '1080p', '2160p'],
      'veo-3.1-fast-generate-001': ['720p', '1080p', '2160p'],
      'veo-3.1-lite': ['720p', '1080p'],
    }
    for (const [model, resolutions] of Object.entries(expected)) {
      expect(normalizeVideoIntervals([], model).map(iv => iv.tier_label)).toEqual(resolutions)
    }
  })

  it('splits synchronized Leo models into independent video entries', () => {
    const entries = createSyncedPricingEntries('leo', ['seedance-2.0', 'seedance-2.0-fast'])

    expect(entries.map(entry => [entry.models, entry.billing_mode])).toEqual([
      [['seedance-2.0'], 'video'],
      [['seedance-2.0-fast'], 'video'],
    ])
    expect(createSyncedPricingEntries('openai', ['gpt-a', 'gpt-b'])).toHaveLength(1)
  })

  it('uses synchronized Leo model kinds for image and video pricing', () => {
    const entries = createSyncedPricingEntries(
      'leo',
      ['seedance-2.0', 'recraft-v4'],
      [
        { id: 'seedance-2.0', kind: 'video' },
        { id: 'recraft-v4', kind: 'image' },
      ],
    )

    expect(entries.map(entry => [entry.models[0], entry.billing_mode])).toEqual([
      ['seedance-2.0', 'video'],
      ['recraft-v4', 'image'],
    ])
  })

  it('treats an all-empty video card as unsold and still requires one model plus non-negative prices when selling', () => {
    const entry = createPricingFormEntry(['seedance-2.5'], 'video')
    expect(isUnsoldVideoPricingEntry(entry)).toBe(true)
    expect(validateVideoPricing(entry, t)).toBeNull()
    expect(formPricedVideoIntervalsToAPI(entry.intervals, entry.models[0])).toEqual([])

    entry.intervals[0].per_request_price = 0.01
    expect(isUnsoldVideoPricingEntry(entry)).toBe(false)
    expect(validateVideoPricing(entry, t)).toBeNull()
    expect(formPricedVideoIntervalsToAPI(entry.intervals, entry.models[0]).map(iv => [iv.tier_label, iv.per_request_price])).toEqual([
      ['480p', 0.01],
    ])

    entry.intervals.forEach((interval, index) => {
      interval.per_request_price = index / 100
    })
    expect(validateVideoPricing(entry, t)).toBeNull()

    entry.models.push('seedance-2.0-fast')
    expect(validateVideoPricing(entry, t)).toContain('videoSingleModelRequired')

    entry.models.pop()
    entry.intervals[1].per_request_price = -1
    expect(validateVideoPricing(entry, t)).toContain('videoPricesRequired')
  })
})

describe('time pricing', () => {
  it('uses a disabled Shanghai default', () => {
    const form = createDefaultTimePricingForm()
    expect(form).toEqual({ timezone: 'Asia/Shanghai', periods: [], weekdays_only: false })
    expect(formTimePricingToAPI(form)).toBeNull()
  })

  it('defaults missing API day scope to every day', () => {
    expect(apiTimePricingToForm({
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00', end_time: '12:00', multiplier: 2 }],
    }).weekdays_only).toBe(false)
  })

  it('round-trips day scope and formats multiplier', () => {
    const form = apiTimePricingToForm({
      timezone: 'Asia/Shanghai',
      weekdays_only: true,
      periods: [{ start_time: '09:00', end_time: '12:00', multiplier: 2 }],
    })
    expect(form.weekdays_only).toBe(true)
    expect(form.periods[0]).toEqual({
      start_time: '09:00:00',
      end_time: '12:00:00',
      multiplier: '2.00',
    })
    expect(formTimePricingToAPI(form)).toEqual({
      timezone: 'Asia/Shanghai',
      weekdays_only: true,
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: 2 }],
    })
  })

  it.each([
    ['separated', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' }, { start_time: '14:00:00', end_time: '18:00:00', multiplier: '2.00' }], null],
    ['adjacent', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' }, { start_time: '12:00:00', end_time: '14:00:00', multiplier: '1.50' }], null],
    ['midnight split', [{ start_time: '22:00:00', end_time: '00:00:00', multiplier: '2.00' }, { start_time: '00:00:00', end_time: '02:00:00', multiplier: '2.00' }], null],
    ['overlap by one second', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' }, { start_time: '11:59:59', end_time: '14:00:00', multiplier: '2.00' }], 'overlap'],
    ['cross midnight', [{ start_time: '22:00:00', end_time: '02:00:00', multiplier: '2.00' }], 'range'],
    ['equal midnight', [{ start_time: '00:00:00', end_time: '00:00:00', multiplier: '2.00' }], 'range'],
    ['missing seconds', [{ start_time: '09:00', end_time: '12:00', multiplier: '2.00' }], 'format'],
    ['zero', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '0.00' }], 'multiplier'],
    ['three decimals', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '1.001' }], 'multiplier'],
  ])('%s', (_name, periods, errorKey) => {
    const result = validateTimePricing({
      timezone: 'Asia/Shanghai',
      periods: periods as TimePricingPeriodFormEntry[],
    }, t)
    if (errorKey === null) expect(result).toBeNull()
    else expect(result).toContain(String(errorKey))
  })

  it('rejects non-IANA timezone', () => {
    expect(validateTimePricing({
      timezone: 'UTC+8',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' }],
    }, t)).toContain('timezone')
  })

  it.each([
    ['missing', undefined],
    ['blank', '   '],
  ])('rejects a %s timezone without throwing during conversion', (_name, timezone) => {
    const form = {
      timezone,
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' }],
    } as unknown as TimePricingFormEntry

    expect(validateTimePricing(form, t)).toContain('timezone')
    expect(() => formTimePricingToAPI(form)).not.toThrow()
    expect(formTimePricingToAPI(form)?.timezone).toBe('')
  })
})

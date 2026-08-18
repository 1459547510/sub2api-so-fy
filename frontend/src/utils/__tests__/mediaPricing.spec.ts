import { describe, expect, it } from 'vitest'
import {
  configuredPrice,
  imageTiers,
  keepMediaPricingGroup,
  buildVideoCards,
  buildMediaPricingSections,
  collectGroupModels,
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

  it('uses a group-name fallback only when no model name can be resolved', () => {
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

  it('uses each model’s channel interval prices instead of the group flat rates', () => {
    const cards = buildVideoCards({
      name: '视频图片分组',
      ...emptyPrices,
      video_price_480p: 0.12,
      video_price_720p: 0.17,
      video_price_1080p: 0.6,
    }, [
      {
        name: 'seedance-2.0',
        platform: 'leo',
        pricing: {
          billing_mode: 'video',
          intervals: [
            { tier_label: '480p', per_request_price: 0.11 },
            { tier_label: '720p', per_request_price: 0.19 },
            { tier_label: '1080p', per_request_price: 0.41 },
          ],
        },
      },
      {
        name: 'kling-3.0',
        platform: 'leo',
        pricing: {
          billing_mode: 'video',
          intervals: [
            { tier_label: '720p', per_request_price: 0.33 },
            { tier_label: '1080p', per_request_price: 0.88 },
            { tier_label: '2160p', per_request_price: 1.6 },
          ],
        },
      },
      { name: 'gpt-image-2', platform: 'openai_media', pricing: { billing_mode: 'image' } },
    ])
    expect(cards.map((card) => card.title)).toEqual(['kling-3.0', 'seedance-2.0'])
    expect(cards[0].tiers).toEqual([
      { label: '720p', value: 0.33 },
      { label: '1080p', value: 0.88 },
      { label: '2160p', value: 1.6 },
    ])
    expect(cards[1].tiers).toEqual([
      { label: '480p', value: 0.11 },
      { label: '720p', value: 0.19 },
      { label: '1080p', value: 0.41 },
    ])
  })

  it('keeps live channel interval prices when group video_model_prices also exist', () => {
    const cards = buildVideoCards({
      name: '视频图片分组',
      ...emptyPrices,
      video_price_720p: 0.17,
      video_model_prices: { 'seedance-2.0': { '720p': 0.31, '1080p': 0.55 } },
    }, [{
      name: 'seedance-2.0',
      platform: 'leo',
      pricing: {
        billing_mode: 'video',
        intervals: [
          { tier_label: '480p', per_request_price: 0.11 },
          { tier_label: '720p', per_request_price: 0.19 },
          { tier_label: '1080p', per_request_price: 0.41 },
          { tier_label: '2160p', per_request_price: 0.9 },
        ],
      },
    }])
    expect(cards).toEqual([{
      key: 'video:seedance-2.0',
      title: 'seedance-2.0',
      kind: 'video',
      unit: 'per_second',
      tiers: [
        { label: '480p', value: 0.11 },
        { label: '720p', value: 0.19 },
        { label: '1080p', value: 0.41 },
        { label: '2160p', value: 0.9 },
      ],
    }])
  })

  it('maps Grok Imagine workbench IDs onto group video family overrides', () => {
    const cards = buildVideoCards({
      name: 'Grok media',
      ...emptyPrices,
      video_price_720p: 0.07,
      video_model_prices: {
        'grok-imagine-video-1.5': { '480p': 0.08, '720p': 0.14, '1080p': 0.25 },
      },
    }, [{ name: 'grok-imagine-1.5', platform: 'grok', pricing: { billing_mode: 'video' } }])
    expect(cards).toEqual([{
      key: 'video:grok-imagine-1.5',
      title: 'grok-imagine-1.5',
      kind: 'video',
      unit: 'per_second',
      tiers: [
        { label: '480p', value: 0.08 },
        { label: '720p', value: 0.14 },
        { label: '1080p', value: 0.25 },
      ],
    }])
  })

  it('does not stamp group flat prices onto catalog models that have no per-model price', () => {
    const cards = buildVideoCards({
      name: 'Leo group',
      ...emptyPrices,
      video_price_720p: 0.17,
    }, [
      { name: 'seedance-2.0', platform: 'leo', pricing: { billing_mode: 'video' } },
      { name: 'kling-3.0', platform: 'leo' },
    ])
    expect(cards).toEqual([{
      key: 'video-fallback',
      title: 'Leo group',
      kind: 'video',
      unit: 'per_second',
      tiers: [{ label: '720p', value: 0.17 }],
    }])
  })
})

describe('collectGroupModels', () => {
  it('prefers the copy that still has per-model interval prices', () => {
    const models = collectGroupModels(
      7,
      [{ id: 7, models: [{ name: 'seedance-2.0', pricing: { billing_mode: 'video' } }] }],
      [{
        platforms: [{
          groups: [{ id: 7 }],
          supported_models: [{
            name: 'seedance-2.0',
            platform: 'leo',
            pricing: {
              billing_mode: 'video',
              intervals: [{ tier_label: '720p', per_request_price: 0.19 }],
            },
          }],
        }],
      }],
    )
    expect(models[0].pricing?.intervals).toEqual([{ tier_label: '720p', per_request_price: 0.19 }])
  })

  it('uses live channel interval prices over plaza copies', () => {
    const models = collectGroupModels(
      7,
      [{
        id: 7,
        models: [{
          name: 'seedance-2.0',
          pricing: {
            billing_mode: 'video',
            intervals: [
              { tier_label: '480p', per_request_price: 0.22 },
              { tier_label: '720p', per_request_price: 0.31 },
              { tier_label: '1080p', per_request_price: 0.55 },
            ],
          },
        }],
      }],
      [{
        platforms: [{
          groups: [{ id: 7 }],
          supported_models: [{
            name: 'seedance-2.0',
            platform: 'leo',
            pricing: {
              billing_mode: 'video',
              intervals: [
                { tier_label: '480p', per_request_price: 0.11 },
                { tier_label: '720p', per_request_price: 0.19 },
                { tier_label: '1080p', per_request_price: 0.41 },
                { tier_label: '2160p', per_request_price: 0.9 },
              ],
            },
          }],
        }],
      }],
    )
    expect(models[0].pricing?.intervals).toEqual([
      { tier_label: '480p', per_request_price: 0.11 },
      { tier_label: '720p', per_request_price: 0.19 },
      { tier_label: '1080p', per_request_price: 0.41 },
      { tier_label: '2160p', per_request_price: 0.9 },
    ])
  })

  it('unions plaza and available-channel models for the group', () => {
    const models = collectGroupModels(
      7,
      [{ id: 7, models: [{ name: 'seedance-2.0', pricing: { billing_mode: 'video' } }] }],
      [{
        platforms: [{
          groups: [{ id: 7 }],
          supported_models: [
            { name: 'seedance-2.0', platform: 'leo', pricing: { billing_mode: 'video' } },
            { name: 'kling-3.0', platform: 'leo' },
          ],
        }],
      }],
    )
    expect(models.map((model) => model.name)).toEqual(['seedance-2.0', 'kling-3.0'])
  })
})

describe('buildMediaPricingSections', () => {
  it('keeps API group order, puts image cards before video cards, and sorts titles', () => {
    const sections = buildMediaPricingSections(
      [
        {
          id: 2,
          name: 'Second',
          platform: 'openai_media',
          ...emptyPrices,
          image_price_1k: 0.08,
        },
        {
          id: 1,
          name: 'First',
          platform: 'leo',
          ...emptyPrices,
          video_price_720p: 0.14,
          video_model_prices: {
            'seedance-2.0-fast': { '720p': 0.1 },
            'seedance-2.0': { '720p': 0.14 },
          },
        },
      ],
      [{
        id: 2,
        models: [
          { name: 'gpt-image-2', pricing: { billing_mode: 'image' } },
          { name: 'alpha-image', pricing: { billing_mode: 'image' } },
          { name: 'token-model', pricing: { billing_mode: 'token' } },
        ],
      }],
    )

    expect(sections.map((section) => section.name)).toEqual(['Second', 'First'])
    expect(sections[0].cards.map((card) => card.title)).toEqual(['alpha-image', 'gpt-image-2'])
    expect(sections[1].cards.map((card) => card.title)).toEqual(['seedance-2.0', 'seedance-2.0-fast'])
  })

  it('falls back to the group name for image cards when plaza has no image models', () => {
    const sections = buildMediaPricingSections(
      [{ id: 3, name: 'Grok images', platform: 'grok', ...emptyPrices, image_price_1k: 0.02, image_price_2k: 0.02, image_price_4k: 0.02 }],
      [],
    )
    expect(sections[0].cards).toEqual([{
      key: 'image-fallback',
      title: 'Grok images',
      kind: 'image',
      unit: 'per_image',
      tiers: [{ label: null, value: 0.02 }],
    }])
  })

  it('does not create image cards when the group has no image price', () => {
    const sections = buildMediaPricingSections(
      [{
        id: 4,
        name: 'Video only',
        platform: 'leo',
        ...emptyPrices,
        video_model_prices: { 'seedance-2.0': { '720p': 0.14 } },
      }],
      [{ id: 4, models: [{ name: 'gpt-image-2', pricing: { billing_mode: 'image' } }] }],
    )
    expect(sections[0].cards.every((card) => card.kind === 'video')).toBe(true)
  })

  it('keeps distinct per-model prices for a realistic Leo catalog', () => {
    const sections = buildMediaPricingSections(
      [{
        id: 7,
        name: '视频图片分组',
        platform: 'composite',
        ...emptyPrices,
        image_price_1k: 0.08,
        image_price_2k: 0.12,
        video_price_480p: 0.12,
        video_price_720p: 0.17,
        video_price_1080p: 0.6,
      }],
      [{
        id: 7,
        models: [{ name: 'gpt-image-2', pricing: { billing_mode: 'image' } }],
      }],
      [{
        platforms: [{
          groups: [{ id: 7 }],
          supported_models: [
            {
              name: 'seedance-2.0',
              platform: 'leo',
              pricing: {
                billing_mode: 'video',
                intervals: [
                  { tier_label: '480p', per_request_price: 0.11 },
                  { tier_label: '720p', per_request_price: 0.19 },
                  { tier_label: '1080p', per_request_price: 0.41 },
                ],
              },
            },
            {
              name: 'seedance-2.0-mini',
              platform: 'leo',
              pricing: {
                billing_mode: 'video',
                intervals: [{ tier_label: '720p', per_request_price: 0.06 }],
              },
            },
            {
              name: 'kling-3.0',
              platform: 'leo',
              pricing: {
                billing_mode: 'video',
                intervals: [
                  { tier_label: '720p', per_request_price: 0.33 },
                  { tier_label: '1080p', per_request_price: 0.88 },
                  { tier_label: '2160p', per_request_price: 1.6 },
                ],
              },
            },
            {
              name: 'hailuo-03',
              platform: 'leo',
              pricing: {
                billing_mode: 'video',
                intervals: [{ tier_label: '1440p', per_request_price: 0.45 }],
              },
            },
            {
              name: 'ltx-2.3-fast',
              platform: 'leo',
              pricing: {
                billing_mode: 'video',
                intervals: [
                  { tier_label: '1080p', per_request_price: 0.06 },
                  { tier_label: '1440p', per_request_price: 0.21 },
                  { tier_label: '2160p', per_request_price: 0.24 },
                ],
              },
            },
            {
              name: 'grok-imagine-1.5',
              platform: 'leo',
              pricing: {
                billing_mode: 'video',
                intervals: [
                  { tier_label: '400p', per_request_price: 0.10 },
                  { tier_label: '544p', per_request_price: 0.10 },
                  { tier_label: '720p', per_request_price: 0.17 },
                  { tier_label: '960p', per_request_price: 0.17 },
                ],
              },
            },
          ],
        }],
      }],
    )

    const lines = sections[0].cards.map((card) => (
      `${card.kind} ${card.title} ${card.tiers.map((tier) => `${tier.label ?? 'flat'}=${tier.value}`).join(' ')}`
    ))
    expect(lines).toEqual([
      'image gpt-image-2 1K=0.08 2K=0.12',
      'video grok-imagine-1.5 400p=0.1 544p=0.1 720p=0.17 960p=0.17',
      'video hailuo-03 1440p=0.45',
      'video kling-3.0 720p=0.33 1080p=0.88 2160p=1.6',
      'video ltx-2.3-fast 1080p=0.06 1440p=0.21 2160p=0.24',
      'video seedance-2.0 480p=0.11 720p=0.19 1080p=0.41',
      'video seedance-2.0-mini 720p=0.06',
    ])
    expect(lines.every((line) => !line.includes('0.12') || line.startsWith('image'))).toBe(true)
    expect(new Set(sections[0].cards.filter((card) => card.kind === 'video').map((card) => (
      card.tiers.map((tier) => tier.value).join(',')
    ))).size).toBe(6)
  })

  it('keeps a media group that only has channel interval prices', () => {
    const sections = buildMediaPricingSections(
      [{
        id: 7,
        name: '视频图片分组',
        platform: 'composite',
        ...emptyPrices,
      }],
      [],
      [{
        platforms: [{
          groups: [{ id: 7 }],
          supported_models: [
            {
              name: 'seedance-2.0',
              platform: 'leo',
              pricing: {
                billing_mode: 'video',
                intervals: [{ tier_label: '720p', per_request_price: 0.21 }],
              },
            },
            { name: 'gpt-4o', platform: 'openai', pricing: { billing_mode: 'token' } },
          ],
        }],
      }],
    )
    expect(sections[0].cards).toEqual([{
      key: 'video:seedance-2.0',
      title: 'seedance-2.0',
      kind: 'video',
      unit: 'per_second',
      tiers: [{ label: '720p', value: 0.21 }],
    }])
  })
})

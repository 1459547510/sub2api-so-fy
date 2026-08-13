import { shallowMount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Select from '@/components/common/Select.vue'
import PricingEntryCard from '../PricingEntryCard.vue'
import { createPricingFormEntry } from '../types'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (key: string) => key }),
}))

vi.mock('@/api/admin/channels', () => ({
  default: {
    getModelDefaultPricing: vi.fn(),
  },
}))

describe('PricingEntryCard video mode', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('offers video billing only when the parent explicitly allows it', async () => {
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry: createPricingFormEntry(['seedance-2.0'], 'video'), platform: 'leo' },
    })

    expect(wrapper.findComponent(Select).props('options')).not.toContainEqual(
      expect.objectContaining({ value: 'video' })
    )
    expect(wrapper.findAll('[data-testid^="video-price-"]')).toHaveLength(0)

    await wrapper.setProps({ allowVideo: true })
    expect(wrapper.findComponent(Select).props('options')).toContainEqual(
      expect.objectContaining({ value: 'video' })
    )
    expect(wrapper.findAll('[data-testid^="video-price-"]')).toHaveLength(3)
  })

  it('renders fixed resolution prices and emits interval updates', async () => {
    const entry = createPricingFormEntry(['seedance-2.0'], 'video')
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry, platform: 'leo', allowVideo: true },
    })

    expect(wrapper.findAll('[data-testid^="video-price-"]')).toHaveLength(3)
    await wrapper.get('[data-testid="video-price-720p"]').setValue('0.025')

    const updated = wrapper.emitted('update')?.at(-1)?.[0]
    expect(updated).toMatchObject({ billing_mode: 'video' })
    expect((updated as typeof entry).intervals.map(iv => [iv.tier_label, iv.per_request_price])).toEqual([
      ['480p', null],
      ['720p', '0.025'],
      ['1080p', null],
    ])
  })

  it('initializes the three video intervals when mode changes', () => {
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry: createPricingFormEntry(['seedance-2.0']), platform: 'leo', allowVideo: true },
    })

    wrapper.findComponent(Select).vm.$emit('update:modelValue', 'video')
    const updated = wrapper.emitted('update')?.at(-1)?.[0] as ReturnType<typeof createPricingFormEntry>
    expect(updated.billing_mode).toBe('video')
    expect(updated.intervals.map(iv => iv.tier_label)).toEqual(['480p', '720p', '1080p'])
  })

  it('only renders the Mini model 720p price', () => {
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry: createPricingFormEntry(['seedance-2.0-mini'], 'video'), platform: 'leo', allowVideo: true },
    })

    expect(wrapper.findAll('[data-testid^="video-price-"]')).toHaveLength(1)
    expect(wrapper.find('[data-testid="video-price-720p"]').exists()).toBe(true)
  })

  it('renders generic default pricing for non-Leo group video entries', () => {
    const entry = createPricingFormEntry(['video-model'], 'video')
    entry.per_request_price = '0.08'
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry, platform: 'openai', hideTokenIntervals: true },
    })

    expect(wrapper.text()).toContain('admin.channels.form.defaultVideoPrice')
    expect(wrapper.findAll('[data-testid^="video-price-"]')).toHaveLength(0)
    expect(wrapper.get('input[type="number"]').element.value).toBe('0.08')
  })
})

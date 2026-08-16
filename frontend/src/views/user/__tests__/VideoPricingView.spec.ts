import { flushPromises, mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import VideoPricingView from '@/views/user/VideoPricingView.vue'

const { groupsApi, modelPlazaApi } = vi.hoisted(() => ({
  groupsApi: {
    getAvailable: vi.fn(),
  },
  modelPlazaApi: {
    getModelPlaza: vi.fn(),
  },
}))

vi.mock('@/api/groups', () => ({ default: groupsApi }))
vi.mock('@/api/modelPlaza', () => ({ default: modelPlazaApi }))
vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return { ...actual, useI18n: () => ({ t: (key: string) => key }) }
})

function mountPage() {
  return mount(VideoPricingView, {
    global: {
      stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        Icon: { template: '<span />' },
        VideoSectionTabs: { template: '<nav data-testid="video-section-tabs" />' },
      },
    },
  })
}

describe('VideoPricingView', () => {
  it('renders raw group prices as cards and ignores multipliers', async () => {
    groupsApi.getAvailable.mockResolvedValueOnce([{
      id: 7,
      name: 'Video image group',
      platform: 'composite',
      rate_multiplier: 9,
      image_rate_independent: true,
      image_rate_multiplier: 9,
      image_price_1k: 0.08,
      image_price_2k: 0.12,
      image_price_4k: null,
      video_rate_independent: true,
      video_rate_multiplier: 9,
      video_price_480p: null,
      video_price_720p: 0.14,
      video_price_1080p: null,
      video_model_prices: { 'seedance-2.0': { '720p': 0.14 } },
    }])
    modelPlazaApi.getModelPlaza.mockResolvedValueOnce({
      groups: [{
        id: 7,
        models: [{ name: 'gpt-image-2', pricing: { billing_mode: 'image' } }],
      }],
    })

    const wrapper = mountPage()
    await flushPromises()

    expect(wrapper.text()).toContain('gpt-image-2')
    expect(wrapper.text()).toContain('seedance-2.0')
    expect(wrapper.text()).toContain('$0.08')
    expect(wrapper.text()).toContain('$0.12')
    expect(wrapper.text()).toContain('$0.14')
    expect(wrapper.text()).not.toContain('0.72')
    expect(wrapper.text()).not.toContain('1.08')
    expect(groupsApi.getUserGroupRates).toBeUndefined()
  })

  it('shows unavailable when groups fail even if plaza succeeds', async () => {
    groupsApi.getAvailable.mockRejectedValueOnce(new Error('groups down'))
    modelPlazaApi.getModelPlaza.mockResolvedValueOnce({ groups: [] })

    const wrapper = mountPage()
    await flushPromises()

    expect(wrapper.text()).toContain('video.pricing.unavailable')
  })

  it('shows empty when no media prices exist', async () => {
    groupsApi.getAvailable.mockResolvedValueOnce([{
      id: 1,
      name: 'Chat',
      platform: 'openai',
      image_price_1k: 0.08,
    }])
    modelPlazaApi.getModelPlaza.mockRejectedValueOnce(new Error('plaza off'))

    const wrapper = mountPage()
    await flushPromises()

    expect(wrapper.text()).toContain('video.pricing.empty')
  })

  it('keeps previous price cards when refresh fails', async () => {
    groupsApi.getAvailable.mockResolvedValueOnce([{
      id: 7,
      name: 'Video image group',
      platform: 'composite',
      video_price_720p: 0.14,
      video_model_prices: { 'seedance-2.0': { '720p': 0.14 } },
    }])
    groupsApi.getAvailable.mockRejectedValueOnce(new Error('groups down'))
    modelPlazaApi.getModelPlaza.mockResolvedValue({ groups: [] })

    const wrapper = mountPage()
    await flushPromises()

    expect(wrapper.text()).toContain('seedance-2.0')

    await wrapper.get('button').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('seedance-2.0')
  })
})

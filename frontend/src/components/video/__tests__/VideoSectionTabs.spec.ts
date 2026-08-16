import { mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { describe, expect, it, vi } from 'vitest'
import VideoSectionTabs from '@/components/video/VideoSectionTabs.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (key: string) => key }),
}))

describe('VideoSectionTabs', () => {
  it('links the workbench, API docs, and pricing tabs and marks the active tab', async () => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/video-generation', name: 'VideoGeneration', component: { template: '<div />' } },
        { path: '/video-generation/api-docs', name: 'VideoApiDocs', component: { template: '<div />' } },
        { path: '/video-generation/pricing', name: 'VideoPricing', component: { template: '<div />' } },
      ],
    })
    await router.push('/video-generation/pricing')
    await router.isReady()

    const wrapper = mount(VideoSectionTabs, {
      global: {
        plugins: [router],
        stubs: { Icon: { template: '<span />' } },
      },
    })

    const links = wrapper.findAll('a')
    expect(links.map((link) => link.attributes('href'))).toEqual([
      '/video-generation',
      '/video-generation/api-docs',
      '/video-generation/pricing',
    ])
    expect(links.map((link) => link.attributes('aria-selected'))).toEqual(['false', 'false', 'true'])
  })
})

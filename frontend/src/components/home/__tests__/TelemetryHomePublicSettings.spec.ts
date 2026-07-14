import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import TelemetryHome from '@/components/home/TelemetryHome.vue'
import type { CustomMenuItem } from '@/types'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    }),
  }
})

const RouterLinkStub = {
  props: ['to'],
  template: '<a :data-to="to"><slot /></a>',
}

function menuItem(index: number): CustomMenuItem {
  return {
    id: `menu-${index}`,
    label: `Menu ${index}`,
    icon_svg: '',
    url: '',
    visibility: 'user',
    sort_order: index,
  }
}

function mountHome(overrides: Record<string, unknown> = {}) {
  return mount(TelemetryHome, {
    props: {
      siteName: 'Sub2API',
      siteLogo: '',
      siteSubtitle: 'Unified AI Gateway',
      apiBaseUrl: 'https://api.example.com/v1',
      contactInfo: 'support@example.com',
      customMenuItems: [1, 2, 3, 4, 5].map(menuItem),
      docUrl: '',
      isAuthenticated: false,
      dashboardPath: '/dashboard',
      ...overrides,
    },
    global: {
      stubs: {
        RouterLink: RouterLinkStub,
        LocaleSwitcher: true,
        Icon: true,
      },
    },
  })
}

describe('TelemetryHome public settings', () => {
  beforeEach(() => {
    vi.spyOn(HTMLCanvasElement.prototype, 'getContext').mockReturnValue(null)
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders visitor settings without replacing the telemetry metrics', async () => {
    const wrapper = mountHome()

    expect(wrapper.get('[data-testid="site-subtitle"]').text()).toBe('Unified AI Gateway')
    expect(wrapper.get('[data-testid="api-base-url"]').text()).toContain('https://api.example.com/v1')
    expect(wrapper.get('[data-testid="contact-info"]').text()).toContain('support@example.com')
    expect(wrapper.findAll('[data-testid="direct-custom-menu"]')).toHaveLength(2)

    await wrapper.get('[data-testid="custom-menu-toggle"]').trigger('click')
    expect(wrapper.findAll('[data-testid="overflow-custom-menu"]')).toHaveLength(3)
    expect(wrapper.find('.telemetry-metrics').exists()).toBe(true)

    wrapper.unmount()
  })

  it('does not reserve space for empty optional settings', () => {
    const wrapper = mountHome({
      siteSubtitle: '',
      apiBaseUrl: '',
      contactInfo: '',
      customMenuItems: [],
    })

    expect(wrapper.find('[data-testid="site-subtitle"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="site-meta"]').exists()).toBe(false)
    expect(wrapper.find('[data-testid="custom-menu-toggle"]').exists()).toBe(false)

    wrapper.unmount()
  })

  it('closes the custom menu with Escape', async () => {
    const wrapper = mountHome()

    await wrapper.get('[data-testid="custom-menu-toggle"]').trigger('click')
    expect(wrapper.find('[data-testid="custom-menu-panel"]').exists()).toBe(true)

    document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    await nextTick()

    expect(wrapper.find('[data-testid="custom-menu-panel"]').exists()).toBe(false)

    wrapper.unmount()
  })
})

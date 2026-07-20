import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import ApiCodeBlock from '@/components/video/ApiCodeBlock.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (key: string) => key }),
}))

describe('ApiCodeBlock', () => {
  it('copies the displayed code', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: { writeText } })
    const wrapper = mount(ApiCodeBlock, {
      props: { label: 'curl', code: 'curl https://example.com' },
      global: { stubs: { Icon: { template: '<span />' } } },
    })

    await wrapper.get('button').trigger('click')

    expect(writeText).toHaveBeenCalledWith('curl https://example.com')
    expect(wrapper.get('button').attributes('title')).toBe('common.copied')
  })
})

import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import VideoApiDocsView from '@/views/user/VideoApiDocsView.vue'
import ApiCodeBlock from '@/components/video/ApiCodeBlock.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return { ...actual, useI18n: () => ({ t: (key: string) => key }) }
})

describe('VideoApiDocsView', () => {
  it('documents the complete public asynchronous video workflow', () => {
    const wrapper = mount(VideoApiDocsView, {
      global: {
        stubs: {
          AppLayout: { template: '<div><slot /></div>' },
          Icon: { template: '<span />' },
          VideoSectionTabs: { template: '<nav data-testid="video-section-tabs" />' },
        },
      },
    })

    const text = wrapper.text()
    for (const endpoint of [
      '/v1/videos/uploads',
      '/v1/videos/generations',
      '/v1/videos/jobs?limit=50&status=running',
      '/v1/videos/jobs/{job_id}',
      '/v1/videos/jobs/{job_id}/content',
    ]) {
      expect(text).toContain(endpoint)
    }

    const examples = wrapper.findAllComponents(ApiCodeBlock).map((component) => component.props('code'))
    expect(examples.some((code) => code.includes('Prefer: respond-async'))).toBe(true)
    expect(examples.some((code) => code.includes('"image_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"start_frame_url"') && code.includes('"end_frame_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"image_reference"') && code.includes('"order": 1'))).toBe(true)
    expect(examples.some((code) => code.includes('"video_reference_base"') && code.includes('reference.mp4'))).toBe(true)
    expect(examples.some((code) => code.includes('"audio_reference"') && code.includes('reference.mp3'))).toBe(true)
    expect(examples.some((code) => code.includes('-F "video=@./reference.mp4"'))).toBe(true)
    expect(examples.some((code) => code.includes('-F "audio=@./reference.mp3"'))).toBe(true)
    expect(examples.some((code) => code.includes('"media_url"'))).toBe(true)
    expect(examples.every((code) => !(code.includes('"start_frame_url"') && code.includes('"image_reference"')))).toBe(true)
    expect(examples.filter((code) => code.startsWith('curl')).every((code) => !/^\+\s/m.test(code))).toBe(true)
    expect(wrapper.html()).not.toMatch(/Leonardo|LeoStudio|upstream_job_id/i)
  })
})

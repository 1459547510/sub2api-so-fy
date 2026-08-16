import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import VideoApiDocsView from '@/views/user/VideoApiDocsView.vue'
import ApiCodeBlock from '@/components/video/ApiCodeBlock.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return { ...actual, useI18n: () => ({ t: (key: string) => key }) }
})

describe('VideoApiDocsView', () => {
  it('documents the V2 public image and asynchronous video workflow', () => {
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
      '/v1/models',
      '/v1/images/generations',
      '/v1/images/edits',
      '/v1/videos/uploads',
      '/v1/videos/generations',
      '/v1/videos/jobs?limit=50&status=running',
      '/v1/videos/jobs/{job_id}',
      '/v1/videos/jobs/{job_id}/content',
    ]) {
      expect(text).toContain(endpoint)
    }
    expect(text).toContain('video.apiDocs.v2.compatibility.title')
    expect(text).toContain('video.apiDocs.v2.matrix.title')
    expect(wrapper.findAll('#model-matrix tbody tr')).toHaveLength(2)
    const modelExampleHeadings = wrapper.findAll('#model-examples h3').map((heading) => heading.text())
    expect(modelExampleHeadings).toHaveLength(2)
    expect(new Set(modelExampleHeadings).size).toBe(2)

    const examples = wrapper.findAllComponents(ApiCodeBlock).map((component) => component.props('code'))
    expect(examples.some((code) => code.includes('Prefer: respond-async'))).toBe(true)
    expect(examples.some((code) => code.includes('"image_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"start_frame_url"') && code.includes('"end_frame_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"image_reference"') && code.includes('"order": 0'))).toBe(true)
    expect(examples.some((code) => code.includes('"video_reference_base"') && code.includes('reference.mp4') && code.includes('"type": "UPLOADED"'))).toBe(true)
    expect(examples.some((code) => code.includes('"audio_reference"') && code.includes('reference.mp3'))).toBe(true)
    expect(examples.some((code) => code.includes('/v1/images/generations') && code.includes('response_format'))).toBe(true)
    expect(examples.some((code) => code.includes('/v1/images/edits') && code.includes('image[]=@./input.png'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "<image-model-id>"') && code.includes('images'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "<video-model-id>"') && code.includes('/v1/videos/generations'))).toBe(true)
    expect(examples.some((code) => code.includes('-F "video=@./reference.mp4"'))).toBe(true)
    expect(examples.some((code) => code.includes('-F "audio=@./reference.mp3"'))).toBe(true)
    expect(examples.some((code) => code.includes('"media_url"'))).toBe(true)
    expect(examples.every((code) => !(code.includes('"start_frame_url"') && code.includes('"image_reference"')))).toBe(true)
    expect(examples.filter((code) => code.startsWith('curl')).every((code) => !/^\+\s/m.test(code))).toBe(true)
    expect(wrapper.html()).not.toMatch(/Leonardo|Leo\s*Studio|Grok|upstream|provider|account_id|target_platform|composite|internal\/video-inputs/i)
    expect(wrapper.find('#pricing').exists()).toBe(false)
    expect(wrapper.text()).not.toContain('video.apiDocs.v2.nav.pricing')
    expect(wrapper.text()).not.toContain('video.apiDocs.v2.pricing.title')
  })
})

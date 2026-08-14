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
    expect(text).toContain('video.apiDocs.matrix.title')
    expect(text).toContain('video.apiDocs.uploadMediaFormat')
    expect(wrapper.findAll('#model-matrix tbody tr')).toHaveLength(22)
    const modelExampleHeadings = wrapper.findAll('#model-examples h3').map((heading) => heading.text())
    expect(modelExampleHeadings).toHaveLength(22)
    expect(new Set(modelExampleHeadings).size).toBe(22)

    const examples = wrapper.findAllComponents(ApiCodeBlock).map((component) => component.props('code'))
    expect(examples.some((code) => code.includes('Prefer: respond-async'))).toBe(true)
    expect(examples.some((code) => code.includes('"image_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"start_frame_url"') && code.includes('"end_frame_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"image_reference"') && code.includes('"order": 1'))).toBe(true)
    expect(examples.some((code) => code.includes('"video_reference_base"') && code.includes('reference.mp4') && code.includes('"type": "UPLOADED"'))).toBe(true)
    expect(examples.some((code) => code.includes('"audio_reference"') && code.includes('reference.mp3'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "gemini-omni-flash"') && code.includes('"image_reference"') && !code.includes('"image_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "happy-horse-1.1"') && code.includes('prompt_enhance'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "grok-imagine-1.5"') && code.includes('start_frame_url'))).toBe(true)
    for (const model of [
      'seedance-2.0', 'seedance-2.0-fast', 'seedance-2.0-mini', 'bytedance/seedance-2.5', 'seedance-2.5', 'happy-horse-1.1', 'grok-imagine-1.5', 'ltx-2.3-pro', 'ltx-2.3-fast',
      'hailuo-03', 'gemini-omni-flash', 'kling-2.1', 'kling-2.5', 'kling-2.5-turbo-standard', 'kling-2.6', 'kling-video-o-1',
      'kling-3.0', 'kling-3.0-turbo', 'kling-video-o-3', 'veo-3.1-generate-001', 'veo-3.1-fast-generate-001', 'veo-3.1-lite',
    ]) {
      expect(examples.some((code) => code.includes(`"model": "${model}"`) && code.includes('/v1/videos/generations'))).toBe(true)
    }
    expect(examples.some((code) => code.includes('-F "video=@./reference.mp4"'))).toBe(true)
    expect(examples.some((code) => code.includes('-F "audio=@./reference.mp3"'))).toBe(true)
    expect(examples.some((code) => code.includes('"media_url"'))).toBe(true)
    expect(examples.every((code) => !(code.includes('"start_frame_url"') && code.includes('"image_reference"')))).toBe(true)
    expect(examples.filter((code) => code.startsWith('curl')).every((code) => !/^\+\s/m.test(code))).toBe(true)
    expect(wrapper.html()).not.toMatch(/Leonardo|Leo\s*Studio|UUID|upstream|provider|internal\/video-inputs|upstream_job_id/i)
  })
})

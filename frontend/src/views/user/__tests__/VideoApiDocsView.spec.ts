import { mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { describe, expect, it, vi } from 'vitest'
import VideoApiDocsView from '@/views/user/VideoApiDocsView.vue'
import ApiCodeBlock from '@/components/video/ApiCodeBlock.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return { ...actual, useI18n: () => ({ t: (key: string) => key }) }
})

const documentedVideoModels = [
  'seedance-2.0', 'seedance-2.0-fast', 'seedance-2.0-mini', 'bytedance/seedance-2.5', 'seedance-2.5', 'happy-horse-1.1', 'grok-imagine-1.5', 'ltx-2.3-pro', 'ltx-2.3-fast',
  'hailuo-03', 'gemini-omni-flash', 'kling-2.1', 'kling-2.5', 'kling-2.5-turbo-standard', 'kling-2.6', 'kling-video-o-1',
  'kling-3.0', 'kling-3.0-turbo', 'kling-video-o-3', 'veo-3.1-generate-001', 'veo-3.1-fast-generate-001', 'veo-3.1-lite',
]

async function mountDocs(path: string) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/video-generation/api-docs', name: 'VideoApiDocs', component: VideoApiDocsView },
      { path: '/video-generation/api-docs/v1', name: 'VideoApiDocsV1', component: VideoApiDocsView },
    ],
  })
  await router.push(path)
  await router.isReady()

  const wrapper = mount(VideoApiDocsView, {
    global: {
      plugins: [router],
      stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        Icon: { template: '<span />' },
        VideoSectionTabs: { template: '<nav data-testid="video-section-tabs" />' },
      },
    },
  })

  return { wrapper, router }
}

describe('VideoApiDocsView', () => {
  it('documents the V2 public image and asynchronous video workflow with per-model parameters', async () => {
    const { wrapper } = await mountDocs('/video-generation/api-docs')

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
    expect(text).toContain('video.apiDocs.version.v1')
    expect(wrapper.findAll('#model-matrix tbody tr')).toHaveLength(22)
    expect(wrapper.text()).toContain('video.apiDocs.v2.matrix.seedance20.resolution')
    expect(wrapper.text()).not.toContain('video.apiDocs.matrix.seedance20.resolution')
    const modelExampleHeadings = wrapper.findAll('#model-examples h3').map((heading) => heading.text())
    expect(modelExampleHeadings).toEqual(documentedVideoModels)
    expect(new Set(modelExampleHeadings).size).toBe(22)

    const examples = wrapper.findAllComponents(ApiCodeBlock).map((component) => component.props('code'))
    expect(examples.some((code) => code.includes('Prefer: respond-async'))).toBe(true)
    expect(examples.some((code) => code.includes('"image_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"start_frame_url"') && code.includes('"end_frame_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"image_reference"') && code.includes('"order": 0'))).toBe(true)
    expect(examples.some((code) => code.includes('"video_reference_base"') && code.includes('reference.mp4') && code.includes('"type": "UPLOADED"'))).toBe(true)
    expect(examples.some((code) => code.includes('"audio_reference"') && code.includes('reference.mp3'))).toBe(true)
    expect(examples.some((code) => code.includes('/v1/images/generations') && code.includes('response_format'))).toBe(true)
    expect(examples.some((code) => code.includes('/v1/images/generations') && code.includes('"image_urls"'))).toBe(true)
    expect(examples.some((code) => code.includes('/v1/images/edits') && code.includes('image[]=@./input.png'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "gpt-image-2"') && code.includes('images'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "seedance-2.0"') && code.includes('"resolution": "4k"'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "seedance-2.0"') && code.includes('/v1/videos/generations'))).toBe(true)
    expect(examples.some((code) => code.includes('-F "video=@./reference.mp4"'))).toBe(true)
    expect(examples.some((code) => code.includes('-F "audio=@./reference.mp3"'))).toBe(true)
    expect(examples.some((code) => code.includes('"media_url"'))).toBe(true)
    expect(examples.every((code) => !(code.includes('"start_frame_url"') && code.includes('"image_reference"')))).toBe(true)
    expect(examples.filter((code) => code.startsWith('curl')).every((code) => !/^\+\s/m.test(code))).toBe(true)
    for (const model of documentedVideoModels) {
      expect(examples.some((code) => code.includes(`"model": "${model}"`) && code.includes('/v1/videos/generations'))).toBe(true)
    }
    expect(wrapper.html()).not.toMatch(/Leonardo|LeoStudio|Leo\s*Studio|\bLeo\b|\bKrea\b|上游|upstream|provider|account_id|target_platform|composite|internal\/video-inputs/i)
    expect(wrapper.find('#pricing').exists()).toBe(false)
    expect(wrapper.text()).not.toContain('video.apiDocs.v2.nav.pricing')
    expect(wrapper.text()).not.toContain('video.apiDocs.v2.pricing.title')
  })

  it('publishes the original V1 video documentation with the same per-model matrix', async () => {
    const { wrapper } = await mountDocs('/video-generation/api-docs/v1')

    const text = wrapper.text()
    expect(text).toContain('video.apiDocs.title')
    expect(text).toContain('video.apiDocs.badge')
    expect(text).toContain('video.apiDocs.matrix.title')
    expect(text).toContain('video.apiDocs.matrix.seedance20.resolution')
    expect(text).not.toContain('video.apiDocs.v2.matrix.seedance20.resolution')
    expect(text).toContain('video.apiDocs.uploadMediaFormat')
    expect(text).not.toContain('video.apiDocs.v2.compatibility.title')
    expect(text).not.toContain('/v1/images/generations')
    expect(wrapper.findAll('#model-matrix tbody tr')).toHaveLength(22)
    const modelExampleHeadings = wrapper.findAll('#model-examples h3').map((heading) => heading.text())
    expect(modelExampleHeadings).toEqual(documentedVideoModels)

    const examples = wrapper.findAllComponents(ApiCodeBlock).map((component) => component.props('code'))
    expect(examples.some((code) => code.includes('"model": "gemini-omni-flash"') && code.includes('"image_reference"') && !code.includes('"image_url"'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "happy-horse-1.1"') && code.includes('prompt_enhance'))).toBe(true)
    expect(examples.some((code) => code.includes('"model": "grok-imagine-1.5"') && code.includes('start_frame_url'))).toBe(true)
    for (const model of documentedVideoModels) {
      expect(examples.some((code) => code.includes(`"model": "${model}"`) && code.includes('/v1/videos/generations'))).toBe(true)
    }
    expect(wrapper.html()).not.toMatch(/Leonardo|Leo\s*Studio|UUID|upstream|provider|internal\/video-inputs|upstream_job_id/i)
  })
})

import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import VideoGenerationView from '@/views/user/VideoGenerationView.vue'
import { keysAPI } from '@/api'
import {
  cancelVideoJob,
  createVideoJob,
  downloadVideoOutput,
  listVideoJobs,
  uploadVideoInput,
} from '@/api/videoGeneration'

vi.mock('@/api', () => ({
  keysAPI: { list: vi.fn() },
}))

vi.mock('@/api/videoGeneration', () => ({
  cancelVideoJob: vi.fn(),
  createVideoJob: vi.fn(),
  downloadVideoOutput: vi.fn(),
  listVideoJobs: vi.fn(),
  uploadVideoInput: vi.fn(),
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showError: vi.fn(), showSuccess: vi.fn() }),
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return { ...actual, useI18n: () => ({ t: (key: string) => key }) }
})

const leoKey = {
  id: 1,
  name: 'Leo production',
  key: 'sub2-leo-key',
  status: 'active',
  group: { platform: 'leo', allow_image_generation: true },
}

const grokKey = {
  id: 2,
  name: 'Grok only',
  key: 'sub2-grok-key',
  status: 'active',
  group: { platform: 'grok', allow_image_generation: true },
}

function mountView() {
  return mount(VideoGenerationView, {
    global: {
      stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        Icon: { template: '<span />' },
      },
    },
  })
}

describe('VideoGenerationView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    let imageIndex = 0
    Object.defineProperty(URL, 'createObjectURL', { configurable: true, value: vi.fn((blob: Blob) => blob.type === 'video/mp4' ? 'blob:local-video' : `blob:local-image-${++imageIndex}`) })
    Object.defineProperty(URL, 'revokeObjectURL', { configurable: true, value: vi.fn() })
    vi.mocked(keysAPI.list).mockResolvedValue({ items: [leoKey, grokKey] } as any)
    vi.mocked(listVideoJobs).mockResolvedValue({ data: [] })
  })

  it('only shows active Leo keys in the workbench selector', async () => {
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.text()).toContain('Leo production')
    expect(wrapper.text()).not.toContain('Grok only')
  })

  it('submits a text prompt and shows the accepted job', async () => {
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-text', status: 'pending', status_url: '/v1/videos/jobs/vidjob-text',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="video-prompt"]').setValue('A kite over the ocean')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    expect(createVideoJob).toHaveBeenCalledWith('sub2-leo-key', expect.objectContaining({
      model: 'seedance-2.0', prompt: 'A kite over the ocean',
    }))
    expect(wrapper.text()).toContain('vidjob-text')
  })

  it('submits multiple remote reference image URLs', async () => {
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-remote', status: 'pending', status_url: '/v1/videos/jobs/vidjob-remote',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="mode-url"]').trigger('click')
    await wrapper.get('[data-testid="video-image-url"]').setValue('https://example.com/ref-1.png\nhttps://example.com/ref-2.png')
    await wrapper.get('[data-testid="video-prompt"]').setValue('Animate this frame')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    expect(createVideoJob).toHaveBeenCalledWith('sub2-leo-key', expect.objectContaining({
      image_urls: ['https://example.com/ref-1.png', 'https://example.com/ref-2.png'],
    }))
  })

  it('uploads multiple local reference images before creating the job', async () => {
    vi.mocked(uploadVideoInput).mockImplementation(async (_key, file) => ({
      upload_id: file.name,
      image_url: `http://127.0.0.1/internal/video-inputs/${file.name}`,
      content_type: file.type,
      size: file.size,
    }))
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-local', status: 'pending', status_url: '/v1/videos/jobs/vidjob-local',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="mode-local"]').trigger('click')
    const input = wrapper.get('[data-testid="video-image-file"]')
    const files = [
      new File(['one'], 'ref-1.png', { type: 'image/png' }),
      new File(['two'], 'ref-2.jpg', { type: 'image/jpeg' }),
    ]
    Object.defineProperty(input.element, 'files', { configurable: true, value: files })
    await input.trigger('change')
    await wrapper.get('[data-testid="video-prompt"]').setValue('Keep both references')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    expect(uploadVideoInput).toHaveBeenCalledTimes(2)
    expect(createVideoJob).toHaveBeenCalledWith('sub2-leo-key', expect.objectContaining({
      image_urls: [
        'http://127.0.0.1/internal/video-inputs/ref-1.png',
        'http://127.0.0.1/internal/video-inputs/ref-2.jpg',
      ],
    }))
  })

  it('cancels pending jobs and renders completed video output', async () => {
    vi.mocked(downloadVideoOutput).mockResolvedValue(new Blob(['mp4'], { type: 'video/mp4' }))
    vi.mocked(listVideoJobs).mockResolvedValue({ data: [
      { job_id: 'vidjob-done', status: 'completed', status_url: '', model: 'seedance-2.0', prompt: 'done', created_at: new Date().toISOString(), updated_at: new Date().toISOString(), result: { data: [{ mp4_url: 'https://cdn.example/done.mp4' }] } },
      { job_id: 'vidjob-pending', status: 'pending', status_url: '', model: 'seedance-2.0', prompt: 'wait', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
    ] })
    vi.mocked(cancelVideoJob).mockResolvedValue({ job_id: 'vidjob-pending', status: 'canceled', status_url: '', created_at: new Date().toISOString(), updated_at: new Date().toISOString() })
    const wrapper = mountView()
    await flushPromises()

    expect(downloadVideoOutput).toHaveBeenCalledWith('sub2-leo-key', 'vidjob-done')
    expect(wrapper.find('video').attributes('src')).toBe('blob:local-video')
    await wrapper.get('[data-testid="cancel-vidjob-pending"]').trigger('click')
    await flushPromises()
    expect(cancelVideoJob).toHaveBeenCalledWith('sub2-leo-key', 'vidjob-pending')
    wrapper.unmount()
    expect(URL.revokeObjectURL).toHaveBeenCalledWith('blob:local-video')
  })

  it('shows a useful empty state when no Leo key is available', async () => {
    vi.mocked(keysAPI.list).mockResolvedValue({ items: [grokKey] } as any)
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.text()).toContain('video.noKey')
    expect(wrapper.get('[data-testid="submit-video"]').attributes('disabled')).toBeDefined()
  })
})

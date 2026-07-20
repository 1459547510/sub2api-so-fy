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
        VideoSectionTabs: { template: '<nav data-testid="video-section-tabs" />' },
      },
    },
  })
}

describe('VideoGenerationView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    window.localStorage.clear()
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

  it('submits one remote reference image as a start frame', async () => {
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-remote-single', status: 'pending', status_url: '/v1/videos/jobs/vidjob-remote-single',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="mode-url"]').trigger('click')
    await wrapper.get('[data-testid="video-image-url"]').setValue('https://example.com/start.png')
    await wrapper.get('[data-testid="video-prompt"]').setValue('Animate this frame')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    expect(createVideoJob).toHaveBeenCalledWith('sub2-leo-key', expect.objectContaining({
      image_url: 'https://example.com/start.png',
    }))
    const payload = vi.mocked(createVideoJob).mock.calls[0][1]
    expect(payload.image_urls).toBeUndefined()
    expect(payload.guidances).toBeUndefined()
  })

  it('submits with an in-memory custom API Key when no saved Leo key exists', async () => {
    vi.mocked(keysAPI.list).mockResolvedValue({ items: [grokKey] } as any)
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-custom', status: 'pending', status_url: '/v1/videos/jobs/vidjob-custom',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="key-mode-custom"]').trigger('click')
    await wrapper.get('[data-testid="video-custom-api-key"]').setValue(' custom-sub2-key ')
    await wrapper.get('[data-testid="video-prompt"]').setValue('A train crossing a snowy valley')
    expect(wrapper.get('[data-testid="submit-video"]').attributes('disabled')).toBeUndefined()
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    expect(createVideoJob).toHaveBeenCalledWith('custom-sub2-key', expect.objectContaining({
      prompt: 'A train crossing a snowy valley',
    }))
    expect(window.localStorage.length).toBe(0)
  })

  it('submits multiple remote reference image URLs with explicit order', async () => {
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

    const payload = vi.mocked(createVideoJob).mock.calls[0][1]
    expect(payload.image_url).toBeUndefined()
    expect(payload.image_urls).toBeUndefined()
    expect(payload.guidances).toEqual({ image_reference: [
      { image: { url: 'https://example.com/ref-1.png', type: 'UPLOADED' }, strength: 'MID', order: 0 },
      { image: { url: 'https://example.com/ref-2.png', type: 'UPLOADED' }, strength: 'MID', order: 1 },
    ] })
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
    await wrapper.get('[data-testid="key-mode-custom"]').trigger('click')
    await wrapper.get('[data-testid="video-custom-api-key"]').setValue('custom-upload-key')
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
    expect(uploadVideoInput).toHaveBeenNthCalledWith(1, 'custom-upload-key', files[0])
    expect(uploadVideoInput).toHaveBeenNthCalledWith(2, 'custom-upload-key', files[1])
    expect(createVideoJob).toHaveBeenCalledWith('custom-upload-key', expect.any(Object))
    const payload = vi.mocked(createVideoJob).mock.calls[0][1]
    expect(payload.image_urls).toBeUndefined()
    expect(payload.guidances).toEqual({ image_reference: [
      { image: { url: 'http://127.0.0.1/internal/video-inputs/ref-1.png', type: 'UPLOADED' }, strength: 'MID', order: 0 },
      { image: { url: 'http://127.0.0.1/internal/video-inputs/ref-2.jpg', type: 'UPLOADED' }, strength: 'MID', order: 1 },
    ] })
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

  it('uses the custom API Key to list, cancel, and download jobs', async () => {
    vi.mocked(keysAPI.list).mockResolvedValue({ items: [grokKey] } as any)
    vi.mocked(downloadVideoOutput).mockResolvedValue(new Blob(['mp4'], { type: 'video/mp4' }))
    vi.mocked(listVideoJobs).mockResolvedValue({ data: [
      { job_id: 'vidjob-custom-done', status: 'completed', status_url: '', model: 'seedance-2.0', prompt: 'done', created_at: new Date().toISOString(), updated_at: new Date().toISOString(), result: { data: [{ mp4_url: '/v1/videos/jobs/vidjob-custom-done/content' }] } },
      { job_id: 'vidjob-custom-pending', status: 'pending', status_url: '', model: 'seedance-2.0', prompt: 'wait', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
    ] })
    vi.mocked(cancelVideoJob).mockResolvedValue({ job_id: 'vidjob-custom-pending', status: 'canceled', status_url: '', created_at: new Date().toISOString(), updated_at: new Date().toISOString() })
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="key-mode-custom"]').trigger('click')
    await wrapper.get('[data-testid="video-custom-api-key"]').setValue('custom-history-key')
    await wrapper.get('[data-testid="refresh-video-jobs"]').trigger('click')
    await flushPromises()

    expect(listVideoJobs).toHaveBeenCalledWith('custom-history-key', { limit: 50 })
    expect(downloadVideoOutput).toHaveBeenCalledWith('custom-history-key', 'vidjob-custom-done')
    await wrapper.get('[data-testid="cancel-vidjob-custom-pending"]').trigger('click')
    await flushPromises()
    expect(cancelVideoJob).toHaveBeenCalledWith('custom-history-key', 'vidjob-custom-pending')
  })

  it('clears jobs from the previous Key when switching modes', async () => {
    vi.mocked(listVideoJobs).mockResolvedValue({ data: [
      { job_id: 'vidjob-saved-only', status: 'pending', status_url: '', model: 'seedance-2.0', prompt: 'saved key job', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
    ] })
    const wrapper = mountView()
    await flushPromises()
    expect(wrapper.text()).toContain('vidjob-saved-only')

    await wrapper.get('[data-testid="key-mode-custom"]').trigger('click')
    await flushPromises()

    expect(wrapper.text()).not.toContain('vidjob-saved-only')
    expect(listVideoJobs).toHaveBeenCalledTimes(1)
  })

  it('shows a useful empty state when no Leo key is available', async () => {
    vi.mocked(keysAPI.list).mockResolvedValue({ items: [grokKey] } as any)
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.text()).toContain('video.noKey')
    expect(wrapper.get('[data-testid="submit-video"]').attributes('disabled')).toBeDefined()
  })
})

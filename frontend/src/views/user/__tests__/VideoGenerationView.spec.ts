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

const appStoreMocks = vi.hoisted(() => ({ showError: vi.fn(), showSuccess: vi.fn() }))

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
  useAppStore: () => appStoreMocks,
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

  it('only shows resolutions supported by the selected model', async () => {
    const wrapper = mountView()
    await flushPromises()

    const resolution = wrapper.get('[data-testid="video-resolution"]')
    expect(resolution.findAll('option').map((option) => option.attributes('value'))).toEqual(['480p', '720p', '1080p'])
    expect((resolution.element as HTMLSelectElement).value).toBe('720p')

    await wrapper.get('[data-testid="video-model"]').setValue('seedance-2.0-fast')
    expect(resolution.findAll('option').map((option) => option.attributes('value'))).toEqual(['480p', '720p'])
    expect((resolution.element as HTMLSelectElement).value).toBe('720p')

    await resolution.setValue('480p')
    await wrapper.get('[data-testid="video-model"]').setValue('seedance-2.0')
    expect(resolution.findAll('option').map((option) => option.attributes('value'))).toEqual(['480p', '720p', '1080p'])
    expect((resolution.element as HTMLSelectElement).value).toBe('720p')
  })

  it('limits Seedance Mini to 720p and 16:9', async () => {
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="video-model"]').setValue('seedance-2.0-mini')

    expect(wrapper.get('[data-testid="video-resolution"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['720p'])
    expect(wrapper.get('[data-testid="video-aspect-ratio"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['16:9'])
  })

  it('uses LeoStudio discrete duration options instead of free numeric input', async () => {
    const wrapper = mountView()
    await flushPromises()

    const duration = wrapper.get('[data-testid="video-duration"]')
    expect(duration.element.tagName).toBe('SELECT')
    expect(duration.findAll('option').map((option) => Number(option.attributes('value')))).toEqual([4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15])
    expect((duration.element as HTMLSelectElement).value).toBe('8')
    expect(wrapper.get('[data-testid="video-prompt"]').attributes('maxlength')).toBe('5000')
  })

  it('removes unsupported aspect ratios and falls back when resolution changes', async () => {
    const wrapper = mountView()
    await flushPromises()

    const aspectRatio = wrapper.get('[data-testid="video-aspect-ratio"]')
    await wrapper.get('[data-testid="video-resolution"]').setValue('480p')
    await aspectRatio.setValue('9:21')
    await wrapper.get('[data-testid="video-resolution"]').setValue('720p')

    expect(aspectRatio.findAll('option').map((option) => option.attributes('value'))).not.toContain('9:21')
    expect((aspectRatio.element as HTMLSelectElement).value).toBe('16:9')
  })

  it('rejects programmatically injected unsupported model parameters before uploading', async () => {
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="mode-local"]').trigger('click')
    const input = wrapper.get('[data-testid="video-image-file"]')
    Object.defineProperty(input.element, 'files', {
      configurable: true,
      value: [new File(['ref'], 'reference.png', { type: 'image/png' })],
    })
    await input.trigger('change')
    await wrapper.get('[data-testid="video-prompt"]').setValue('This invalid request must stay in the browser')

    const setupState = (wrapper.vm as any).$?.setupState
    setupState.duration = 16
    await wrapper.vm.$nextTick()
    expect(wrapper.get('[data-testid="submit-video"]').attributes('disabled')).toBeDefined()
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    expect(appStoreMocks.showError).toHaveBeenCalledWith('video.invalidModelParameters')
    expect(uploadVideoInput).not.toHaveBeenCalled()
    expect(createVideoJob).not.toHaveBeenCalled()
  })

  it('submits the exact supported model parameter combination', async () => {
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-fast', status: 'pending', status_url: '/v1/videos/jobs/vidjob-fast',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="video-model"]').setValue('seedance-2.0-fast')
    await wrapper.get('[data-testid="video-resolution"]').setValue('720p')
    await wrapper.get('[data-testid="video-aspect-ratio"]').setValue('21:9')
    await wrapper.get('[data-testid="video-duration"]').setValue('15')
    await wrapper.get('[data-testid="video-audio"]').setValue(true)
    await wrapper.get('[data-testid="video-prompt"]').setValue('A vertical city reveal')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    expect(createVideoJob).toHaveBeenCalledWith('sub2-leo-key', {
      model: 'seedance-2.0-fast',
      prompt: 'A vertical city reveal',
      resolution: '720p',
      duration: 15,
      aspect_ratio: '21:9',
      audio: true,
    })
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
    expect(wrapper.get('[data-testid="video-start-frame-url"]').attributes('disabled')).toBeDefined()
    expect(wrapper.get('[data-testid="video-end-frame-url"]').attributes('disabled')).toBeDefined()
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
    Object.defineProperty(input.element, 'files', { configurable: true, value: [files[0]] })
    await input.trigger('change')
    Object.defineProperty(input.element, 'files', { configurable: true, value: [files[1]] })
    await input.trigger('change')
    expect(wrapper.findAll('[data-testid="video-settings"] img')).toHaveLength(2)
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

  it('disables local frame inputs and rejects frame selection after four reference images', async () => {
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="mode-local"]').trigger('click')

    const referenceInput = wrapper.get('[data-testid="video-image-file"]')
    for (let index = 1; index <= 4; index += 1) {
      const file = new File([String(index)], `reference-${index}.png`, { type: 'image/png' })
      Object.defineProperty(referenceInput.element, 'files', { configurable: true, value: [file] })
      await referenceInput.trigger('change')
    }

    const startInput = wrapper.get('[data-testid="video-start-frame-file"]')
    expect(startInput.attributes('disabled')).toBeDefined()
    expect(wrapper.get('[data-testid="video-end-frame-file"]').attributes('disabled')).toBeDefined()
    Object.defineProperty(startInput.element, 'files', {
      configurable: true,
      value: [new File(['start'], 'start.png', { type: 'image/png' })],
    })
    startInput.element.dispatchEvent(new Event('change'))
    await flushPromises()

    expect(wrapper.find('[data-testid="remove-start-frame"]').exists()).toBe(false)
    expect(appStoreMocks.showError).toHaveBeenCalledWith('video.imageInputConflict')
  })

  it('disables and rejects local reference images after a frame is selected', async () => {
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="mode-local"]').trigger('click')

    const startInput = wrapper.get('[data-testid="video-start-frame-file"]')
    Object.defineProperty(startInput.element, 'files', {
      configurable: true,
      value: [new File(['start'], 'start.png', { type: 'image/png' })],
    })
    await startInput.trigger('change')

    const referenceInput = wrapper.get('[data-testid="video-image-file"]')
    expect(referenceInput.attributes('disabled')).toBeDefined()
    Object.defineProperty(referenceInput.element, 'files', {
      configurable: true,
      value: [new File(['ref'], 'reference.png', { type: 'image/png' })],
    })
    referenceInput.element.dispatchEvent(new Event('change'))
    await flushPromises()

    expect(wrapper.findAll('[data-testid="video-settings"] img')).toHaveLength(1)
    expect(appStoreMocks.showError).toHaveBeenCalledWith('video.imageInputConflict')
  })

  it('blocks a programmatic mixed URL payload before upload or job creation', async () => {
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="mode-url"]').trigger('click')
    await wrapper.get('[data-testid="video-start-frame-url"]').setValue('https://example.com/start.png')

    const referenceInput = wrapper.get('[data-testid="video-image-url"]')
    expect(referenceInput.attributes('disabled')).toBeDefined()
    const textarea = referenceInput.element as HTMLTextAreaElement
    textarea.value = 'https://example.com/reference.png'
    textarea.dispatchEvent(new Event('input'))
    await wrapper.get('[data-testid="video-prompt"]').setValue('This mixed request must not leave the page')
    await flushPromises()

    expect(wrapper.get('[data-testid="submit-video"]').attributes('disabled')).toBeDefined()
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()
    expect(appStoreMocks.showError).toHaveBeenCalledWith('video.imageInputConflict')
    expect(uploadVideoInput).not.toHaveBeenCalled()
    expect(createVideoJob).not.toHaveBeenCalled()
  })

  it('uploads start-frame and end-frame files concurrently without reference guidance', async () => {
    const pendingUploads = new Map<string, (value: any) => void>()
    vi.mocked(uploadVideoInput).mockImplementation((_key, file) => new Promise((resolve) => {
      pendingUploads.set(file.name, resolve)
    }))
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-frames', status: 'pending', status_url: '/v1/videos/jobs/vidjob-frames',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="mode-local"]').trigger('click')

    const start = new File(['start'], 'start.png', { type: 'image/png' })
    const end = new File(['end'], 'end.jpg', { type: 'image/jpeg' })
    const setFile = async (testId: string, file: File) => {
      const input = wrapper.get(`[data-testid="${testId}"]`)
      Object.defineProperty(input.element, 'files', { configurable: true, value: [file] })
      await input.trigger('change')
    }
    await setFile('video-start-frame-file', start)
    await setFile('video-end-frame-file', end)
    await wrapper.get('[data-testid="video-prompt"]').setValue('Interpolate between both frames')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await Promise.resolve()

    expect(uploadVideoInput).toHaveBeenCalledTimes(2)
    expect(uploadVideoInput).toHaveBeenNthCalledWith(1, 'sub2-leo-key', start)
    expect(uploadVideoInput).toHaveBeenNthCalledWith(2, 'sub2-leo-key', end)
    expect(createVideoJob).not.toHaveBeenCalled()

    for (const [name, resolve] of pendingUploads) {
      resolve({ upload_id: name, image_url: `http://127.0.0.1/internal/video-inputs/${name}`, content_type: 'image/png', size: 3 })
    }
    await flushPromises()

    const payload = vi.mocked(createVideoJob).mock.calls[0][1]
    expect(payload.start_frame_url).toBe('http://127.0.0.1/internal/video-inputs/start.png')
    expect(payload.end_frame_url).toBe('http://127.0.0.1/internal/video-inputs/end.jpg')
    expect(payload.image_url).toBeUndefined()
    expect(payload.guidances).toBeUndefined()
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

  it('keeps the completed preview stable when polling refreshes the same job', async () => {
    const completedJob = { job_id: 'vidjob-stable', status: 'completed', status_url: '', model: 'seedance-2.0', prompt: 'done', created_at: new Date().toISOString(), updated_at: new Date().toISOString() }
    const pendingJob = { job_id: 'vidjob-active', status: 'running', status_url: '', model: 'seedance-2.0', prompt: 'active', created_at: new Date().toISOString(), updated_at: new Date().toISOString() }
    vi.mocked(listVideoJobs).mockResolvedValue({ data: [completedJob, pendingJob] })
    vi.mocked(downloadVideoOutput).mockResolvedValue(new Blob(['mp4'], { type: 'video/mp4' }))
    const wrapper = mountView()
    await flushPromises()

    expect(downloadVideoOutput).toHaveBeenCalledTimes(1)
    expect(wrapper.get('[data-testid="video-preview"]').attributes('src')).toBe('blob:local-video')
    await wrapper.get('[data-testid="refresh-video-jobs"]').trigger('click')
    await flushPromises()

    expect(downloadVideoOutput).toHaveBeenCalledTimes(1)
    expect(URL.revokeObjectURL).not.toHaveBeenCalledWith('blob:local-video')
    expect(wrapper.get('[data-testid="video-preview"]').attributes('src')).toBe('blob:local-video')
    wrapper.unmount()
  })

  it('keeps the download available after playback fails and refetches on retry', async () => {
    const completedJob = { job_id: 'vidjob-retry', status: 'completed', status_url: '', model: 'seedance-2.0', prompt: 'retry', created_at: new Date().toISOString(), updated_at: new Date().toISOString() }
    vi.mocked(listVideoJobs).mockResolvedValue({ data: [completedJob] })
    vi.mocked(downloadVideoOutput).mockResolvedValue(new Blob(['mp4'], { type: 'video/mp4' }))
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="video-preview"]').trigger('error')
    expect(wrapper.get('[data-testid="video-preview-error"]').text()).toContain('video.previewError')
    expect(wrapper.get('[data-testid="download-video-output"]').attributes('download')).toBe('vidjob-retry.mp4')

    await wrapper.get('[data-testid="retry-video-preview"]').trigger('click')
    await flushPromises()
    expect(URL.revokeObjectURL).toHaveBeenCalledWith('blob:local-video')
    expect(downloadVideoOutput).toHaveBeenCalledTimes(2)
    expect(wrapper.get('[data-testid="download-video-output"]').attributes('download')).toBe('vidjob-retry.mp4')
    wrapper.unmount()
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

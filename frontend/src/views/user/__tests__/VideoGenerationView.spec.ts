import { flushPromises, mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import VideoGenerationView from '@/views/user/VideoGenerationView.vue'
import { keysAPI } from '@/api'
import {
  cancelVideoJob,
  createVideoJob,
  downloadVideoOutput,
  listGatewayModels,
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
  listGatewayModels: vi.fn(),
  listVideoJobs: vi.fn(),
  uploadVideoInput: vi.fn(),
}))

const allWorkbenchModels = [
  'seedance-2.0',
  'seedance-2.0-fast',
  'seedance-2.0-mini',
  'seedance-2.5',
  'happy-horse-1.1',
  'grok-imagine-1.5',
  'ltx-2.3-pro',
  'ltx-2.3-fast',
  'hailuo-03',
  'gemini-omni-flash',
  'kling-2.1',
  'kling-2.5',
  'kling-2.5-turbo-standard',
  'kling-2.6',
  'kling-video-o-1',
  'kling-3.0',
  'kling-3.0-turbo',
  'kling-video-o-3',
  'veo-3.1-generate-001',
  'veo-3.1-fast-generate-001',
  'veo-3.1-lite',
]

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
    vi.mocked(listGatewayModels).mockResolvedValue(allWorkbenchModels)
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllGlobals()
  })

  it('hides Grok Imagine image IDs and MiniMax chat IDs from the workbench list', async () => {
    vi.mocked(listGatewayModels).mockResolvedValue([
      'grok-imagine-1.5',
      'grok-imagine-image',
      'grok-imagine-image-quality',
      'grok-imagine',
      'minimax-m3',
    ])
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.get('[data-testid="video-model"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['grok-imagine-1.5'])
  })

  it('loads workbench models from the selected key /v1/models list', async () => {
    vi.mocked(listGatewayModels).mockResolvedValue(['seedance-2.0', 'kling-3.0', 'gpt-image-2', 'gpt-4o'])
    const wrapper = mountView()
    await flushPromises()

    expect(listGatewayModels).toHaveBeenCalledWith('sub2-leo-key')
    const options = wrapper.get('[data-testid="video-model"]').findAll('option').map((option) => option.attributes('value'))
    expect(options).toEqual(['kling-3.0', 'seedance-2.0'])
  })

  it('reloads workbench models when the saved key changes', async () => {
    const otherKey = {
      id: 3,
      name: 'Leo other',
      key: 'sub2-other-key',
      status: 'active',
      group: { platform: 'leo', allow_image_generation: true },
    }
    vi.mocked(keysAPI.list).mockResolvedValue({ items: [leoKey, otherKey] } as any)
    vi.mocked(listGatewayModels).mockImplementation(async (apiKey: string) => (
      apiKey === 'sub2-other-key' ? ['kling-3.0'] : ['seedance-2.0', 'gpt-image-2']
    ))
    const wrapper = mountView()
    await flushPromises()
    expect(wrapper.get('[data-testid="video-model"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['seedance-2.0'])

    await wrapper.get('[data-testid="video-api-key"]').setValue(3)
    await flushPromises()

    expect(listGatewayModels).toHaveBeenCalledWith('sub2-other-key')
    expect(wrapper.get('[data-testid="video-model"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['kling-3.0'])
    expect((wrapper.get('[data-testid="video-model"]').element as HTMLSelectElement).value).toBe('kling-3.0')
  })

  it('loads workbench models from a custom key after debounce', async () => {
    vi.useFakeTimers()
    const wrapper = mountView()
    await flushPromises()
    expect(listGatewayModels).toHaveBeenCalledWith('sub2-leo-key')

    vi.mocked(listGatewayModels).mockResolvedValue(['hailuo-03', 'gpt-image-2'])
    await wrapper.get('[data-testid="key-mode-custom"]').trigger('click')
    await wrapper.get('[data-testid="video-custom-api-key"]').setValue('custom-sub2-key')
    expect(listGatewayModels).toHaveBeenCalledTimes(1)

    await vi.advanceTimersByTimeAsync(400)
    await flushPromises()

    expect(listGatewayModels).toHaveBeenCalledWith('custom-sub2-key')
    expect(wrapper.get('[data-testid="video-model"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['hailuo-03'])
    vi.useRealTimers()
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
    expect(resolution.findAll('option').map((option) => option.attributes('value'))).toEqual(['480p', '720p', '1080p', '4k'])
    expect((resolution.element as HTMLSelectElement).value).toBe('720p')

    await wrapper.get('[data-testid="video-model"]').setValue('seedance-2.0-fast')
    expect(resolution.findAll('option').map((option) => option.attributes('value'))).toEqual(['480p', '720p'])
    expect((resolution.element as HTMLSelectElement).value).toBe('720p')

    await resolution.setValue('480p')
    await wrapper.get('[data-testid="video-model"]').setValue('seedance-2.0')
    expect(resolution.findAll('option').map((option) => option.attributes('value'))).toEqual(['480p', '720p', '1080p', '4k'])
    expect((resolution.element as HTMLSelectElement).value).toBe('720p')
  })

  it('uses Krea Seedance Mini ratios instead of the old Leo portrait-only set', async () => {
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="video-model"]').setValue('seedance-2.0-mini')

    expect(wrapper.get('[data-testid="video-resolution"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['480p', '720p'])
    expect(wrapper.get('[data-testid="video-aspect-ratio"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['16:9', '4:3', '1:1', '3:4', '9:16', '21:9'])
  })

  it('limits Seedance 2.5 to its published resolution, duration, and aspect options', async () => {
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="video-model"]').setValue('seedance-2.5')

    expect(wrapper.get('[data-testid="video-resolution"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['480p', '720p'])
    expect(wrapper.get('[data-testid="video-duration"]').findAll('option').map((option) => Number(option.attributes('value')))).toEqual(Array.from({ length: 27 }, (_, index) => index + 4))
    expect(wrapper.get('[data-testid="video-aspect-ratio"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['16:9', '9:16', '1:1', '4:3', '3:4', '21:9'])
  })

  it('exposes Happy Horse and Grok capability-specific controls', async () => {
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="video-model"]').setValue('happy-horse-1.1')
    await wrapper.get('[data-testid="mode-local"]').trigger('click')
    expect(wrapper.get('[data-testid="video-resolution"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['720p', '1080p'])
    expect(wrapper.get('[data-testid="video-duration"]').findAll('option').map((option) => Number(option.attributes('value')))[0]).toBe(3)
    expect(wrapper.get('[data-testid="video-prompt"]').attributes('maxlength')).toBe('2500')
    expect(wrapper.get('[data-testid="video-prompt-enhance"]').exists()).toBe(true)
    expect(wrapper.get('[data-testid="video-end-frame-file"]').attributes('disabled')).toBeDefined()
    expect(wrapper.get('[data-testid="video-reference-file"]').attributes('disabled')).toBeDefined()
    expect(wrapper.get('[data-testid="audio-reference-file"]').attributes('disabled')).toBeDefined()

    await wrapper.get('[data-testid="video-model"]').setValue('grok-imagine-1.5')
    expect(wrapper.get('[data-testid="video-resolution"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['auto', '400p', '544p', '720p', '960p'])
    await wrapper.get('[data-testid="video-resolution"]').setValue('544p')
    expect(wrapper.get('[data-testid="video-aspect-ratio"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['1:1'])
    expect(wrapper.get('[data-testid="video-start-frame-file"]').attributes('disabled')).toBeUndefined()
    expect(wrapper.get('[data-testid="video-image-file"]').attributes('disabled')).toBeDefined()
  })

  it('filters the new Leo model parameters and disables unsupported generated audio', async () => {
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="video-model"]').setValue('hailuo-03')
    expect(wrapper.get('[data-testid="video-resolution"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['1440p'])
    expect(wrapper.get('[data-testid="video-duration"]').findAll('option').map((option) => Number(option.attributes('value')))[0]).toBe(5)
    expect(wrapper.get('[data-testid="video-duration"]').findAll('option').map((option) => Number(option.attributes('value'))).at(-1)).toBe(15)

    await wrapper.get('[data-testid="video-model"]').setValue('kling-2.6')
    expect(wrapper.get('[data-testid="video-resolution"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['auto', '1080p'])
    await wrapper.get('[data-testid="video-resolution"]').setValue('auto')
    expect(wrapper.get('[data-testid="video-aspect-ratio"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['auto'])

    await wrapper.get('[data-testid="video-model"]').setValue('gemini-omni-flash')
    expect(wrapper.get('[data-testid="video-audio"]').attributes('disabled')).toBeDefined()
    expect(wrapper.get('[data-testid="video-reference-file"]').attributes('disabled')).toBeDefined()
  })

  it('submits a single Gemini image as a reference image instead of a start frame', async () => {
    vi.mocked(uploadVideoInput).mockResolvedValue({
      upload_id: 'gemini-reference',
      media_url: 'http://127.0.0.1/internal/video-inputs/gemini.png',
      media_type: 'image',
      content_type: 'image/png',
      size: 3,
    })
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-gemini-image', status: 'pending', status_url: '/v1/videos/jobs/vidjob-gemini-image',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="video-model"]').setValue('gemini-omni-flash')
    await wrapper.get('[data-testid="mode-local"]').trigger('click')
    const input = wrapper.get('[data-testid="video-image-file"]')
    Object.defineProperty(input.element, 'files', { configurable: true, value: [new File(['png'], 'gemini.png', { type: 'image/png' })] })
    await input.trigger('change')
    await wrapper.get('[data-testid="video-prompt"]').setValue('Animate the reference image')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    const payload = vi.mocked(createVideoJob).mock.calls[0][1]
    expect(payload.image_url).toBeUndefined()
    expect(payload.guidances).toEqual({ image_reference: [
      { image: { url: 'http://127.0.0.1/internal/video-inputs/gemini.png', type: 'UPLOADED' }, strength: 'MID', order: 0 },
    ] })
  })

  it('rejects local video upload for generated-only Kling video references', async () => {
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="video-model"]').setValue('kling-video-o-3')

    const input = wrapper.get('[data-testid="video-reference-file"]')
    Object.defineProperty(input.element, 'files', { configurable: true, value: [new File(['mp4'], 'reference.mp4', { type: 'video/mp4' })] })
    await input.trigger('change')

    expect(appStoreMocks.showError).toHaveBeenCalledWith('video.modelGuidanceUnsupported')
    expect(wrapper.findAll('[data-testid="audio-reference-preview"]')).toHaveLength(0)
    expect(uploadVideoInput).not.toHaveBeenCalled()
  })

  it('uploads multiple Hailuo audio references and keeps the public URL shape', async () => {
    class MockAudio {
      duration = 5
      onloadedmetadata: (() => void) | null = null
      onerror: (() => void) | null = null
      set src(_value: string) {
        queueMicrotask(() => this.onloadedmetadata?.())
      }
    }
    vi.stubGlobal('Audio', MockAudio)
    vi.mocked(uploadVideoInput).mockImplementation(async (_key, file, kind = 'image') => ({
      upload_id: file.name,
      media_url: `http://127.0.0.1/internal/video-inputs/${file.name}`,
      media_type: kind,
      content_type: file.type,
      size: file.size,
    }))
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-hailuo-audio', status: 'pending', status_url: '/v1/videos/jobs/vidjob-hailuo-audio',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()
    await wrapper.get('[data-testid="video-model"]').setValue('hailuo-03')
    await wrapper.get('[data-testid="mode-local"]').trigger('click')

    const imageInput = wrapper.get('[data-testid="video-image-file"]')
    Object.defineProperty(imageInput.element, 'files', { configurable: true, value: [new File(['png'], 'reference.png', { type: 'image/png' })] })
    await imageInput.trigger('change')
    const audioFiles = [
      new File(['one'], 'reference-1.mp3', { type: 'audio/mpeg' }),
      new File(['two'], 'reference-2.mp3', { type: 'audio/mpeg' }),
      new File(['three'], 'reference-3.mp3', { type: 'audio/mpeg' }),
    ]
    const audioInput = wrapper.get('[data-testid="audio-reference-file"]')
    Object.defineProperty(audioInput.element, 'files', { configurable: true, value: audioFiles })
    await audioInput.trigger('change')
    await flushPromises()
    await wrapper.get('[data-testid="video-prompt"]').setValue('Match the reference audio')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    const payload = vi.mocked(createVideoJob).mock.calls[0][1]
    expect(payload.guidances?.audio_reference).toEqual(audioFiles.map((file) => ({ audio: { url: `http://127.0.0.1/internal/video-inputs/${file.name}`, type: 'UPLOADED' } })))
    expect(JSON.stringify(payload)).not.toMatch(/"id"\s*:/)
  })

  it('exposes exact LTX 2.3 parameters and allows prompt enhancement with a start frame', async () => {
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-ltx', status: 'pending', status_url: '/v1/videos/jobs/vidjob-ltx',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()

    await wrapper.get('[data-testid="video-model"]').setValue('ltx-2.3-pro')
    expect(wrapper.get('[data-testid="video-resolution"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['1080p', '1440p', '2160p'])
    expect(wrapper.get('[data-testid="video-duration"]').findAll('option').map((option) => Number(option.attributes('value')))).toEqual([6, 8, 10])
    expect(wrapper.get('[data-testid="video-aspect-ratio"]').findAll('option').map((option) => option.attributes('value'))).toEqual(['16:9'])
    expect(wrapper.get('[data-testid="video-prompt-enhance"]').exists()).toBe(true)

    await wrapper.get('[data-testid="video-model"]').setValue('ltx-2.3-fast')
    expect(wrapper.get('[data-testid="video-duration"]').findAll('option').map((option) => Number(option.attributes('value')))).toEqual([6, 8, 10, 12, 14, 16, 18, 20])
    await wrapper.get('[data-testid="video-resolution"]').setValue('2160p')
    await wrapper.get('[data-testid="video-duration"]').setValue('20')
    await wrapper.get('[data-testid="video-prompt-enhance"]').setValue('ON')
    await wrapper.get('[data-testid="mode-url"]').trigger('click')
    await wrapper.get('[data-testid="video-start-frame-url"]').setValue('https://example.com/start.png')
    await wrapper.get('[data-testid="video-prompt"]').setValue('A cinematic coastal flyover')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()

    expect(createVideoJob).toHaveBeenCalledWith('sub2-leo-key', {
      model: 'ltx-2.3-fast',
      prompt: 'A cinematic coastal flyover',
      resolution: '2160p',
      duration: 20,
      aspect_ratio: '16:9',
      audio: false,
      prompt_enhance: 'ON',
      start_frame_url: 'https://example.com/start.png',
    })
  })

  it('keeps the same 4-15 second range at every Krea Seedance resolution including 4k', async () => {
    const wrapper = mountView()
    await flushPromises()

    const duration = wrapper.get('[data-testid="video-duration"]')
    const resolution = wrapper.get('[data-testid="video-resolution"]')
    expect(duration.element.tagName).toBe('SELECT')
    expect(duration.findAll('option').map((option) => Number(option.attributes('value')))).toEqual([4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15])
    expect((duration.element as HTMLSelectElement).value).toBe('5')

    await duration.setValue('15')
    await resolution.setValue('4k')
    expect(duration.findAll('option').map((option) => Number(option.attributes('value')))).toEqual([4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15])
    expect((duration.element as HTMLSelectElement).value).toBe('15')

    await resolution.setValue('720p')
    expect(duration.findAll('option').map((option) => Number(option.attributes('value')))).toEqual([4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15])
    expect(wrapper.get('[data-testid="video-prompt"]').attributes('maxlength')).toBe('5000')
  })

  it('keeps Krea Seedance aspect ratios when resolution changes', async () => {
    const wrapper = mountView()
    await flushPromises()

    const aspectRatio = wrapper.get('[data-testid="video-aspect-ratio"]')
    await wrapper.get('[data-testid="video-resolution"]').setValue('480p')
    await aspectRatio.setValue('21:9')
    await wrapper.get('[data-testid="video-resolution"]').setValue('4k')

    expect(aspectRatio.findAll('option').map((option) => option.attributes('value'))).toEqual(['16:9', '4:3', '1:1', '3:4', '9:16', '21:9'])
    expect((aspectRatio.element as HTMLSelectElement).value).toBe('21:9')
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
    setupState.aspectRatio = '9:21'
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

  it('uploads reference video and audio files and builds media guidance without UUIDs', async () => {
    const originalLoad = HTMLMediaElement.prototype.load
    HTMLMediaElement.prototype.load = function () {
      Object.defineProperty(this, 'duration', { configurable: true, value: 5 })
      this.onloadedmetadata?.(new Event('loadedmetadata'))
    }
    class MockAudio {
      duration = 5
      onloadedmetadata: (() => void) | null = null
      onerror: (() => void) | null = null
      set src(_value: string) {
        queueMicrotask(() => this.onloadedmetadata?.())
      }
    }
    vi.stubGlobal('Audio', MockAudio)
    vi.mocked(uploadVideoInput).mockImplementation(async (_key, file, kind = 'image') => ({
      upload_id: file.name,
      media_url: `http://127.0.0.1/internal/video-inputs/${file.name}`,
      media_type: kind,
      content_type: file.type,
      size: file.size,
    }))
    vi.mocked(createVideoJob).mockResolvedValue({
      job_id: 'vidjob-media', status: 'pending', status_url: '/v1/videos/jobs/vidjob-media',
      created_at: new Date().toISOString(), updated_at: new Date().toISOString(),
    })
    const wrapper = mountView()
    await flushPromises()

    const videoInput = wrapper.get('[data-testid="video-reference-file"]')
    Object.defineProperty(videoInput.element, 'files', { configurable: true, value: [new File(['mp4'], 'reference.mp4', { type: 'video/mp4' })] })
    await videoInput.trigger('change')
    await flushPromises()

    const audioInput = wrapper.get('[data-testid="audio-reference-file"]')
    Object.defineProperty(audioInput.element, 'files', { configurable: true, value: [new File(['ID3'], 'reference.mp3', { type: 'audio/mpeg' })] })
    await audioInput.trigger('change')
    await flushPromises()
    await wrapper.get('[data-testid="video-prompt"]').setValue('Match the reference media')
    await wrapper.get('[data-testid="video-settings"] form').trigger('submit')
    await flushPromises()
    HTMLMediaElement.prototype.load = originalLoad

    expect(uploadVideoInput).toHaveBeenCalledWith('sub2-leo-key', expect.any(File), 'video')
    expect(uploadVideoInput).toHaveBeenCalledWith('sub2-leo-key', expect.any(File), 'audio')
    const payload = vi.mocked(createVideoJob).mock.calls[0][1]
    expect(payload.guidances).toEqual({
      video_reference_base: [{ video: { url: 'http://127.0.0.1/internal/video-inputs/reference.mp4', type: 'UPLOADED' } }],
      audio_reference: [{ audio: { url: 'http://127.0.0.1/internal/video-inputs/reference.mp3', type: 'UPLOADED' } }],
    })
    expect(JSON.stringify(payload)).not.toMatch(/"id"\s*:/)
  })

  it('warns immediately for unsupported reference media and blocks audio without a visual reference', async () => {
    class MockAudio {
      duration = 5
      onloadedmetadata: (() => void) | null = null
      onerror: (() => void) | null = null
      set src(_value: string) {
        queueMicrotask(() => this.onloadedmetadata?.())
      }
    }
    vi.stubGlobal('Audio', MockAudio)
    const wrapper = mountView()
    await flushPromises()

    const videoInput = wrapper.get('[data-testid="video-reference-file"]')
    Object.defineProperty(videoInput.element, 'files', { configurable: true, value: [new File(['webm'], 'reference.webm', { type: 'video/webm' })] })
    await videoInput.trigger('change')
    expect(appStoreMocks.showError).toHaveBeenCalledWith('video.invalidVideo')

    const oversizedVideo = new File(['mp4'], 'large-reference.mp4', { type: 'video/mp4' })
    Object.defineProperty(oversizedVideo, 'size', { configurable: true, value: 100 * 1024 * 1024 + 1 })
    Object.defineProperty(videoInput.element, 'files', { configurable: true, value: [oversizedVideo] })
    await videoInput.trigger('change')
    expect(appStoreMocks.showError).toHaveBeenCalledWith('video.videoTooLarge')

    const audioInput = wrapper.get('[data-testid="audio-reference-file"]')
    Object.defineProperty(audioInput.element, 'files', { configurable: true, value: [new File(['ID3'], 'reference.mp3', { type: 'audio/mpeg' })] })
    await audioInput.trigger('change')
    await flushPromises()
    expect(appStoreMocks.showError).toHaveBeenCalledWith('video.audioNeedsVisualReference')
    expect(createVideoJob).not.toHaveBeenCalled()
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
    expect(wrapper.get('[data-testid="video-end-frame-file"]').attributes('disabled')).toBeUndefined()

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
    expect(wrapper.get('[data-testid="submit-video"]').attributes('disabled')).toBeUndefined()
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

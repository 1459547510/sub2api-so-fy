import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  cancelVideoJob,
  createVideoJob,
  downloadVideoOutput,
  listVideoJobs,
  uploadVideoInput,
} from '@/api/videoGeneration'

describe('video generation API', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('uploads a local image without setting a manual multipart boundary', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(new Response(JSON.stringify({ upload_id: 'upload-1' }), { status: 200 })))
    vi.stubGlobal('fetch', fetchMock)
    const file = new File(['png'], 'frame.png', { type: 'image/png' })

    await uploadVideoInput('sub2-key', file)

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(init.method).toBe('POST')
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer sub2-key')
    expect((init.headers as Record<string, string>)['Content-Type']).toBeUndefined()
    expect(init.body).toBeInstanceOf(FormData)
  })

  it('uses dedicated multipart fields for reference video and audio uploads', async () => {
    const fetchMock = vi.fn().mockImplementation(() => Promise.resolve(new Response(JSON.stringify({ upload_id: 'upload-1' }), { status: 200 })))
    vi.stubGlobal('fetch', fetchMock)
    const video = new File(['mp4'], 'reference.mp4', { type: 'video/mp4' })
    const audio = new File(['mp3'], 'reference.mp3', { type: 'audio/mpeg' })

    await uploadVideoInput('sub2-key', video, 'video')
    await uploadVideoInput('sub2-key', audio, 'audio')

    const videoForm = fetchMock.mock.calls[0][1].body as FormData
    const audioForm = fetchMock.mock.calls[1][1].body as FormData
    expect(videoForm.get('video')).toBe(video)
    expect(videoForm.get('image')).toBeNull()
    expect(audioForm.get('audio')).toBe(audio)
    expect(audioForm.get('image')).toBeNull()
  })

  it('submits an async video request with the selected API key', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify({ job_id: 'vidjob-1', status: 'pending' }), { status: 202 }))
    vi.stubGlobal('fetch', fetchMock)

    await createVideoJob('sub2-key', {
      model: 'seedance-2.0',
      prompt: 'waves',
      start_frame_url: 'https://example.com/start.png',
      end_frame_url: 'https://example.com/end.png',
      image_urls: ['https://example.com/ref-1.png', 'https://example.com/ref-2.png'],
      guidances: { image_reference: [{ image: { id: 'asset-1', type: 'GENERATED' }, strength: 'HIGH' }] },
    })

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    const headers = init.headers as Record<string, string>
    expect(headers.Authorization).toBe('Bearer sub2-key')
    expect(headers.Prefer).toBe('respond-async')
    expect(JSON.parse(String(init.body))).toMatchObject({
      model: 'seedance-2.0',
      prompt: 'waves',
      start_frame_url: 'https://example.com/start.png',
      end_frame_url: 'https://example.com/end.png',
      image_urls: ['https://example.com/ref-1.png', 'https://example.com/ref-2.png'],
      guidances: { image_reference: [{ strength: 'HIGH' }] },
    })
  })

  it('lists jobs with the selected API key and cancels with DELETE', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(new Response(JSON.stringify({ data: [] }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ job_id: 'vidjob-1', status: 'canceled' }), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await listVideoJobs('sub2-key', { limit: 50, status: 'pending' })
    await cancelVideoJob('sub2-key', 'vidjob-1')

    const [, listInit] = fetchMock.mock.calls[0] as [string, RequestInit]
    const [cancelUrl, cancelInit] = fetchMock.mock.calls[1] as [string, RequestInit]
    expect(listInit.headers).toMatchObject({ Authorization: 'Bearer sub2-key' })
    expect(fetchMock.mock.calls[0][0]).toContain('limit=50')
    expect(fetchMock.mock.calls[0][0]).toContain('status=pending')
    expect(cancelUrl).toContain('/v1/videos/jobs/vidjob-1')
    expect(cancelInit.method).toBe('DELETE')
    expect(cancelInit.headers).toMatchObject({ Authorization: 'Bearer sub2-key' })
  })

  it('downloads saved video output with the selected API key', async () => {
    const video = new Blob(['mp4'], { type: 'video/mp4' })
    const fetchMock = vi.fn().mockResolvedValue(new Response(video, { status: 200, headers: { 'Content-Type': 'video/mp4' } }))
    vi.stubGlobal('fetch', fetchMock)

    const result = await downloadVideoOutput('sub2-key', 'vidjob-1')

    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toContain('/v1/videos/jobs/vidjob-1/content')
    expect(init.headers).toMatchObject({ Authorization: 'Bearer sub2-key' })
    expect(result.type).toBe('video/mp4')
  })
})

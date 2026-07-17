import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  cancelVideoJob,
  createVideoJob,
  listVideoJobs,
  uploadVideoInput,
} from '@/api/videoGeneration'

describe('video generation API', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('uploads a local image without setting a manual multipart boundary', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify({ upload_id: 'upload-1' }), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)
    const file = new File(['png'], 'frame.png', { type: 'image/png' })

    await uploadVideoInput('sub2-key', file)

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(init.method).toBe('POST')
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer sub2-key')
    expect((init.headers as Record<string, string>)['Content-Type']).toBeUndefined()
    expect(init.body).toBeInstanceOf(FormData)
  })

  it('submits an async video request with the selected API key', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify({ job_id: 'vidjob-1', status: 'pending' }), { status: 202 }))
    vi.stubGlobal('fetch', fetchMock)

    await createVideoJob('sub2-key', { model: 'seedance-2.0', prompt: 'waves' })

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    const headers = init.headers as Record<string, string>
    expect(headers.Authorization).toBe('Bearer sub2-key')
    expect(headers.Prefer).toBe('respond-async')
    expect(JSON.parse(String(init.body))).toMatchObject({ model: 'seedance-2.0', prompt: 'waves' })
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
})

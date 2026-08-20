import { buildGatewayUrl } from './client'

export type VideoJobStatus = 'pending' | 'running' | 'settling' | 'completed' | 'failed' | 'canceled' | string

export interface VideoGuidanceAsset {
  url?: string
  type?: string
  width?: number
  height?: number
  duration?: number
  motion_has_audio?: boolean
}

export interface VideoGuidances {
  start_frame?: Array<{ image: VideoGuidanceAsset }>
  end_frame?: Array<{ image: VideoGuidanceAsset }>
  image_reference?: Array<{ image: VideoGuidanceAsset; strength?: 'LOW' | 'MID' | 'HIGH' | string; order?: number }>
  video_reference_base?: Array<{ video: VideoGuidanceAsset }>
  audio_reference?: Array<{ audio: VideoGuidanceAsset }>
}

export interface VideoGenerationRequest {
  model: string
  prompt: string
  resolution?: '400p' | '480p' | '544p' | '720p' | '960p' | '1080p' | '1440p' | '2160p' | string
  duration?: number
  aspect_ratio?: string
  audio?: boolean
  prompt_enhance?: 'AUTO' | 'ON' | 'OFF' | string
  image_url?: string
  start_frame_url?: string
  end_frame_url?: string
  image_urls?: string[]
  guidances?: VideoGuidances
}

export interface VideoUploadResponse {
  upload_id: string
  media_url?: string
  media_type: 'image' | 'video' | 'audio' | string
  image_url?: string
  video_url?: string
  audio_url?: string
  content_type: string
  size: number
}

export interface VideoJobResult {
  data?: Array<Record<string, unknown>>
  [key: string]: unknown
}

export interface VideoJob {
  job_id: string
  status: VideoJobStatus
  status_url: string
  model?: string
  prompt?: string
  created_at: string
  updated_at: string
  result?: VideoJobResult
  error?: { message: string }
}

export interface VideoJobsResponse {
  data: VideoJob[]
  total?: number
  limit?: number
  offset?: number
}

export interface VideoJobListOptions {
  limit?: number
  offset?: number
  status?: string
}

function authHeaders(apiKey: string, extra?: HeadersInit): HeadersInit {
  return {
    Authorization: `Bearer ${apiKey}`,
    ...extra,
  }
}

async function parseVideoError(response: Response): Promise<Error> {
  try {
    const body = await response.json()
    const message = body?.error?.message || body?.message || response.statusText
    const error = new Error(message)
    ;(error as any).status = response.status
    ;(error as any).code = body?.error?.code || response.status
    return error
  } catch {
    const error = new Error(response.statusText || `HTTP ${response.status}`)
    ;(error as any).status = response.status
    ;(error as any).code = response.status
    return error
  }
}

export async function uploadVideoInput(apiKey: string, file: File, kind: 'image' | 'video' | 'audio' = 'image'): Promise<VideoUploadResponse> {
  const form = new FormData()
  form.append(kind, file)
  const response = await fetch(buildGatewayUrl('/v1/videos/uploads'), {
    method: 'POST',
    headers: authHeaders(apiKey),
    body: form,
  })
  if (!response.ok) throw await parseVideoError(response)
  return response.json()
}

export async function createVideoJob(apiKey: string, payload: VideoGenerationRequest): Promise<VideoJob> {
  const response = await fetch(buildGatewayUrl('/v1/videos/generations'), {
    method: 'POST',
    headers: authHeaders(apiKey, {
      'Content-Type': 'application/json',
      Prefer: 'respond-async',
    }),
    body: JSON.stringify(payload),
  })
  if (!response.ok) throw await parseVideoError(response)
  return response.json()
}

export async function listVideoJobs(apiKey: string, options: VideoJobListOptions = {}): Promise<VideoJobsResponse> {
  const params = new URLSearchParams({ limit: String(options.limit || 50) })
  if (options.offset != null) params.set('offset', String(options.offset))
  if (options.status) params.set('status', options.status)
  const response = await fetch(buildGatewayUrl(`/v1/videos/jobs?${params.toString()}`), {
    headers: authHeaders(apiKey),
  })
  if (!response.ok) throw await parseVideoError(response)
  return response.json()
}

export async function getVideoJob(apiKey: string, jobId: string): Promise<VideoJob> {
  const response = await fetch(buildGatewayUrl(`/v1/videos/jobs/${encodeURIComponent(jobId)}`), {
    headers: authHeaders(apiKey),
  })
  if (!response.ok) throw await parseVideoError(response)
  return response.json()
}

export async function downloadVideoOutput(apiKey: string, jobId: string): Promise<Blob> {
  const response = await fetch(buildGatewayUrl(`/v1/videos/jobs/${encodeURIComponent(jobId)}/content`), {
    headers: authHeaders(apiKey),
  })
  if (!response.ok) throw await parseVideoError(response)
  return response.blob()
}

export async function listGatewayModels(apiKey: string): Promise<string[]> {
  const response = await fetch(buildGatewayUrl('/v1/models'), {
    headers: authHeaders(apiKey),
  })
  if (!response.ok) throw await parseVideoError(response)
  const body = await response.json()
  const data = Array.isArray(body?.data) ? body.data : []
  return data
    .map((item: { id?: string; ID?: string }) => String(item?.id || item?.ID || '').trim())
    .filter(Boolean)
}

export async function cancelVideoJob(apiKey: string, jobId: string): Promise<VideoJob> {
  const response = await fetch(buildGatewayUrl(`/v1/videos/jobs/${encodeURIComponent(jobId)}`), {
    method: 'DELETE',
    headers: authHeaders(apiKey),
  })
  if (!response.ok) throw await parseVideoError(response)
  return response.json()
}

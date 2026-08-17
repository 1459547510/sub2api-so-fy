import { describe, expect, it } from 'vitest'
import { buildV2VideoModelExamples, buildVideoModelExamples, v2VideoModelMatrixRows, videoModelMatrixRows } from '@/utils/videoApiDocs'

describe('videoApiDocs', () => {
  it('keeps one matrix row and one request example for each documented video model', () => {
    const models = videoModelMatrixRows.map((row) => row.model)
    const examples = buildVideoModelExamples('https://docs.example')

    expect(models).toHaveLength(22)
    expect(examples.map((example) => example.model)).toEqual(models)
    expect(new Set(models).size).toBe(22)
    expect(examples.every((example) => example.code.includes('/v1/videos/generations'))).toBe(true)
    expect(examples.every((example) => example.code.includes(`"model": "${example.model}"`))).toBe(true)
  })

  it('documents Krea Seedance limits on the V2 matrix and examples', () => {
    const seedance = v2VideoModelMatrixRows.find((row) => row.model === 'seedance-2.0')
    const examples = buildV2VideoModelExamples('https://docs.example')

    expect(seedance?.resolution).toBe('video.apiDocs.v2.matrix.seedance20.resolution')
    expect(videoModelMatrixRows.find((row) => row.model === 'seedance-2.0')?.resolution).toBe('video.apiDocs.matrix.seedance20.resolution')
    expect(examples.find((example) => example.model === 'seedance-2.0')?.code).toContain('"resolution": "4k"')
    expect(examples.find((example) => example.model === 'seedance-2.0')?.code).toContain('"duration": 5')
    expect(examples.find((example) => example.model === 'seedance-2.0-mini')?.code).toContain('"aspect_ratio": "21:9"')
  })
})

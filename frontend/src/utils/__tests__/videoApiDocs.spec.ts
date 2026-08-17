import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'
import enDashboard from '@/i18n/locales/en/dashboard'
import zhDashboard from '@/i18n/locales/zh/dashboard'
import { buildV2VideoModelExamples, buildVideoModelExamples, v2VideoModelMatrixRows, videoModelMatrixRows } from '@/utils/videoApiDocs'

const publicDocsVendorName = /Leonardo|LeoStudio|Leo\s*Studio|\bLeo\b|\bKrea\b|上游|provider|upstream/i

function collectStrings(value: unknown, out: string[] = []): string[] {
  if (typeof value === 'string') {
    out.push(value)
    return out
  }
  if (Array.isArray(value)) {
    for (const item of value) collectStrings(item, out)
    return out
  }
  if (value && typeof value === 'object') {
    for (const item of Object.values(value)) collectStrings(item, out)
  }
  return out
}

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

  it('keeps V2 public copy free of upstream vendor names', () => {
    const markdown = readFileSync(resolve(dirname(fileURLToPath(import.meta.url)), '../../../../docs/WEB_API_INTEGRATION_V2_CN.md'), 'utf8')
    for (const text of [
      ...collectStrings(zhDashboard.video.apiDocs.v2),
      ...collectStrings(enDashboard.video.apiDocs.v2),
      markdown,
    ]) {
      expect(text).not.toMatch(publicDocsVendorName)
    }
  })

  it('documents Seedance 2.0 limits on the V2 matrix and examples', () => {
    const seedance = v2VideoModelMatrixRows.find((row) => row.model === 'seedance-2.0')
    const examples = buildV2VideoModelExamples('https://docs.example')

    expect(seedance?.resolution).toBe('video.apiDocs.v2.matrix.seedance20.resolution')
    expect(videoModelMatrixRows.find((row) => row.model === 'seedance-2.0')?.resolution).toBe('video.apiDocs.matrix.seedance20.resolution')
    expect(examples.find((example) => example.model === 'seedance-2.0')?.code).toContain('"resolution": "4k"')
    expect(examples.find((example) => example.model === 'seedance-2.0')?.code).toContain('"duration": 5')
    expect(examples.find((example) => example.model === 'seedance-2.0-mini')?.code).toContain('"aspect_ratio": "21:9"')
  })
})

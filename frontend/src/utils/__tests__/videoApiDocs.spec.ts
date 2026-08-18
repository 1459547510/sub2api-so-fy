import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'
import enDashboard from '@/i18n/locales/en/dashboard'
import zhDashboard from '@/i18n/locales/zh/dashboard'
import { applySeedanceV2DocsToDashboard, seedanceV2DocsCatalog } from '@/utils/seedanceV2DocsCatalog'
import { SEEDANCE_V2_DOCS_SOURCE } from '@/utils/seedanceV2DocsSource'
import { buildV2VideoModelExamples, buildVideoModelExamples, v2VideoModelMatrixRows, videoModelMatrixRows } from '@/utils/videoApiDocs'

const publicDocsVendorName = /Leonardo|LeoStudio|Leo\s*Studio|\bLeo\b|\bKrea\b|\bTrioma\b|byteplus|上游|provider|upstream/i

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

  it('keeps public API docs free of upstream vendor names', () => {
    const markdown = readFileSync(resolve(dirname(fileURLToPath(import.meta.url)), '../../../../docs/WEB_API_INTEGRATION_V2_CN.md'), 'utf8')
    for (const text of [
      ...collectStrings(zhDashboard.video.apiDocs),
      ...collectStrings(enDashboard.video.apiDocs),
      ...collectStrings(seedanceV2DocsCatalog.current.zh),
      ...collectStrings(seedanceV2DocsCatalog.current.en),
      ...collectStrings(seedanceV2DocsCatalog.previous.zh),
      ...collectStrings(seedanceV2DocsCatalog.previous.en),
      markdown,
    ]) {
      expect(text).not.toMatch(publicDocsVendorName)
    }
  })

  it('keeps the previous Seedance V2 catalog available for an internal switch', () => {
    expect(SEEDANCE_V2_DOCS_SOURCE).toBe('current')
    expect(SEEDANCE_V2_DOCS_SOURCE).not.toMatch(/trioma|krea/i)
    expect(Object.keys(seedanceV2DocsCatalog).sort()).toEqual(['current', 'previous'])
    expect(seedanceV2DocsCatalog.previous.zh.seedance20Mini.duration).toContain('默认 5 秒')
    expect(seedanceV2DocsCatalog.previous.zh.seedance20Mini.resolution).toBe('480p、720p')
    expect(seedanceV2DocsCatalog.previous.v2MatrixKeys).not.toContain('seedance25')

    const previousDocs = applySeedanceV2DocsToDashboard(structuredClone(zhDashboard), 'zh', 'previous')
    expect(previousDocs.video.apiDocs.v2.matrix.seedance20Mini.duration).toContain('默认 5 秒')
    expect(previousDocs.video.apiDocs.v2.matrix.seedance25).toBeUndefined()
    expect(zhDashboard.video.apiDocs.v2.matrix.seedance20Mini.duration).toContain('4-12')
  })

  it('documents Seedance limits on the V2 matrix and examples', () => {
    const seedance = v2VideoModelMatrixRows.find((row) => row.model === 'seedance-2.0')
    const examples = buildV2VideoModelExamples('https://docs.example')
    const v2Matrix = zhDashboard.video.apiDocs.v2.matrix
    const v2Models = zhDashboard.video.apiDocs.v2.models

    expect(seedance?.resolution).toBe('video.apiDocs.v2.matrix.seedance20.resolution')
    expect(v2VideoModelMatrixRows.find((row) => row.model === 'seedance-2.5')?.duration).toBe('video.apiDocs.v2.matrix.seedance25.duration')
    expect(videoModelMatrixRows.find((row) => row.model === 'seedance-2.0')?.resolution).toBe('video.apiDocs.matrix.seedance20.resolution')
    expect(examples.find((example) => example.model === 'seedance-2.0')?.code).toContain('"resolution": "4k"')
    expect(examples.find((example) => example.model === 'seedance-2.0')?.code).toContain('"duration": 4')
    expect(examples.find((example) => example.model === 'seedance-2.0-mini')?.code).toContain('"resolution": "1080p"')
    expect(examples.find((example) => example.model === 'seedance-2.0-mini')?.code).toContain('"aspect_ratio": "21:9"')
    expect(examples.find((example) => example.model === 'seedance-2.5')?.code).toContain('"duration": 4')
    expect(v2Matrix.seedance20.duration).toContain('默认 4 秒')
    expect(v2Matrix.seedance20Mini.resolution).toContain('1080p')
    expect(v2Matrix.seedance20Mini.duration).toContain('4-12')
    expect(v2Matrix.seedance20Mini.references).toContain('最多 2 张')
    expect(v2Matrix.seedance25.duration).toContain('4、5、6、8、10、12、15、20、25、30')
    expect(v2Models.seedance20).toContain('默认 4 秒')
    expect(enDashboard.video.apiDocs.v2.matrix.seedance20Mini.duration).toContain('4-12s')
    expect(enDashboard.video.apiDocs.v2.matrix.seedance25.duration).toContain('4, 5, 6, 8, 10, 12, 15, 20, 25, or 30s')
  })
})

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

    expect(models).toHaveLength(21)
    expect(examples.map((example) => example.model)).toEqual(models)
    expect(new Set(models).size).toBe(21)
    expect(examples.every((example) => example.code.includes('/v1/videos/generations'))).toBe(true)
    expect(examples.every((example) => example.code.includes(`"model": "${example.model}"`))).toBe(true)
  })

  it('keeps public API docs free of upstream vendor names', () => {
    const docsDir = resolve(dirname(fileURLToPath(import.meta.url)), '../../../../docs')
    const markdown = readFileSync(resolve(docsDir, 'WEB_API_INTEGRATION_V2_CN.md'), 'utf8')
    const seedanceNative = readFileSync(resolve(docsDir, 'seedance-api.md'), 'utf8')
    for (const text of [
      ...collectStrings(zhDashboard.video.apiDocs),
      ...collectStrings(enDashboard.video.apiDocs),
      ...collectStrings(seedanceV2DocsCatalog.current.zh),
      ...collectStrings(seedanceV2DocsCatalog.current.en),
      ...collectStrings(seedanceV2DocsCatalog.previous.zh),
      ...collectStrings(seedanceV2DocsCatalog.previous.en),
      ...collectStrings(seedanceV2DocsCatalog.trioma.zh),
      ...collectStrings(seedanceV2DocsCatalog.trioma.en),
      markdown,
      seedanceNative,
    ]) {
      expect(text).not.toMatch(publicDocsVendorName)
    }
  })

  it('keeps unused Seedance V2 catalogs available for an internal switch', () => {
    expect(SEEDANCE_V2_DOCS_SOURCE).toBe('previous')
    expect(SEEDANCE_V2_DOCS_SOURCE).not.toMatch(/trioma|krea/i)
    expect(Object.keys(seedanceV2DocsCatalog).sort()).toEqual(['current', 'previous', 'trioma'])
    expect(seedanceV2DocsCatalog.previous.v2MatrixKeys).not.toContain('seedance25')
    for (const catalog of Object.values(seedanceV2DocsCatalog)) {
      expect(catalog.v2MatrixKeys).not.toContain('seedance20Mini')
      expect(JSON.stringify(catalog)).not.toContain('seedance-2.0-mini')
    }

    const currentDocs = applySeedanceV2DocsToDashboard(structuredClone(zhDashboard), 'zh', 'current')
    expect(currentDocs.video.apiDocs.v2.matrix.seedance20.references).toContain('参考图 12')
    expect(currentDocs.video.apiDocs.v2.matrix.seedance25?.duration).toContain('4、5、6、8、10、12、15、20、25、30')
    const triomaDocs = applySeedanceV2DocsToDashboard(structuredClone(zhDashboard), 'zh', 'trioma')
    expect(triomaDocs.video.apiDocs.v2.matrix.seedance20.references).toContain('参考图 9')
    expect(triomaDocs.video.apiDocs.v2.matrix.seedance25.references).toContain('参考视频 10')
    expect(zhDashboard.video.apiDocs.v2.matrix.seedance20.duration).toContain('默认 5 秒')
    expect(zhDashboard.video.apiDocs.v2.matrix.seedance25).toBeUndefined()
  })

  it('documents Seedance limits on the V2 matrix and examples', () => {
    const seedance = v2VideoModelMatrixRows.find((row) => row.model === 'seedance-2.0')
    const examples = buildV2VideoModelExamples('https://docs.example')
    const v2Matrix = zhDashboard.video.apiDocs.v2.matrix
    const v2Models = zhDashboard.video.apiDocs.v2.models

    expect(seedance?.resolution).toBe('video.apiDocs.v2.matrix.seedance20.resolution')
    expect(v2VideoModelMatrixRows.find((row) => row.model === 'seedance-2.5')?.duration).toBe('video.apiDocs.matrix.seedance25.duration')
    expect(videoModelMatrixRows.find((row) => row.model === 'seedance-2.0')?.resolution).toBe('video.apiDocs.matrix.seedance20.resolution')
    expect(examples.find((example) => example.model === 'seedance-2.0')?.code).toContain('"resolution": "4k"')
    expect(examples.find((example) => example.model === 'seedance-2.0')?.code).toContain('"duration": 5')
    expect(examples.find((example) => example.model === 'seedance-2.5')?.code).toContain('"duration": 8')
    expect(v2Matrix.seedance20.duration).toContain('默认 5 秒')
    expect(v2Matrix.seedance20.references).toContain('参考图 9')
    expect(v2Matrix.seedance25).toBeUndefined()
    expect(v2Models.seedance20).toContain('默认 5 秒')
    expect(v2Models.seedance20).toContain('最多 9 张参考图')
    expect(v2Models.seedance25).toBeUndefined()
    expect(enDashboard.video.apiDocs.v2.matrix.seedance20.duration).toContain('default 5s')
    expect(enDashboard.video.apiDocs.v2.matrix.seedance20.references).toContain('9 images')
    expect(enDashboard.video.apiDocs.v2.matrix.seedance25).toBeUndefined()
  })
})

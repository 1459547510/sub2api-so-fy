import { describe, expect, it } from 'vitest'
import { COMPOSITE_ROUTE_PLATFORMS } from '@/components/admin/channel/compositeGroups'

describe('Composite channel platform options', () => {
  it('includes CN and media providers for pricing and model mapping', () => {
    expect(COMPOSITE_ROUTE_PLATFORMS).toEqual(expect.arrayContaining([
      'leo',
      'openai_media',
      'kimi',
      'zhipu',
      'deepseek',
    ]))
  })
})

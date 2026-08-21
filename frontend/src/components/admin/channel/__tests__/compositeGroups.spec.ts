import { describe, expect, it } from 'vitest'
import {
  COMPOSITE_ROUTE_PLATFORMS,
  compositeGroupAppliesToPlatform,
  groupMatchesChannelPlatform,
} from '../compositeGroups'

describe('compositeGroups', () => {
  it('enables all concrete platforms from a composite group', () => {
    expect([...COMPOSITE_ROUTE_PLATFORMS]).toEqual([
      'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
      'leo', 'openai_media', 'kimi', 'zhipu', 'deepseek',
    ])
  })

  it('lets a composite group satisfy media and CN pricing tabs', () => {
    expect(compositeGroupAppliesToPlatform('leo')).toBe(true)
    expect(compositeGroupAppliesToPlatform('openai_media')).toBe(true)
    expect(compositeGroupAppliesToPlatform('kimi')).toBe(true)
    expect(groupMatchesChannelPlatform('composite', 'leo')).toBe(true)
    expect(groupMatchesChannelPlatform('composite', 'openai_media')).toBe(true)
    expect(groupMatchesChannelPlatform('composite', 'openai')).toBe(true)
    expect(groupMatchesChannelPlatform('composite', 'kimi')).toBe(true)
    expect(groupMatchesChannelPlatform('leo', 'openai_media')).toBe(false)
    expect(groupMatchesChannelPlatform(undefined, 'leo')).toBe(false)
  })
})

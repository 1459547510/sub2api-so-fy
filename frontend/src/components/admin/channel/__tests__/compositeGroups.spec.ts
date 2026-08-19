import { describe, expect, it } from 'vitest'
import {
  COMPOSITE_ROUTE_PLATFORMS,
  compositeGroupAppliesToPlatform,
  groupMatchesChannelPlatform,
} from '../compositeGroups'

describe('compositeGroups', () => {
  it('does not auto-enable media or CN platforms from a composite group', () => {
    expect([...COMPOSITE_ROUTE_PLATFORMS]).toEqual([
      'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
    ])
    expect(COMPOSITE_ROUTE_PLATFORMS).not.toContain('leo')
    expect(COMPOSITE_ROUTE_PLATFORMS).not.toContain('openai_media')
    expect(COMPOSITE_ROUTE_PLATFORMS).not.toContain('kimi')
  })

  it('lets a composite group satisfy leo and openai_media pricing tabs', () => {
    expect(compositeGroupAppliesToPlatform('leo')).toBe(true)
    expect(compositeGroupAppliesToPlatform('openai_media')).toBe(true)
    expect(compositeGroupAppliesToPlatform('kimi')).toBe(false)
    expect(groupMatchesChannelPlatform('composite', 'leo')).toBe(true)
    expect(groupMatchesChannelPlatform('composite', 'openai_media')).toBe(true)
    expect(groupMatchesChannelPlatform('composite', 'openai')).toBe(true)
    expect(groupMatchesChannelPlatform('composite', 'kimi')).toBe(false)
    expect(groupMatchesChannelPlatform('leo', 'openai_media')).toBe(false)
    expect(groupMatchesChannelPlatform(undefined, 'leo')).toBe(false)
  })
})

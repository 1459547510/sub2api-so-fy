import type { GroupPlatform } from '@/types'

/** Platforms a composite group enables as routing tabs. Media and CN providers stay off this list. */
export const COMPOSITE_ROUTE_PLATFORMS: readonly GroupPlatform[] = [
  'anthropic', 'openai', 'gemini', 'antigravity', 'grok',
]

/** Composite groups also attach to media pricing tabs so video channels can save. */
export function compositeGroupAppliesToPlatform(platform: GroupPlatform): boolean {
  return COMPOSITE_ROUTE_PLATFORMS.includes(platform)
    || platform === 'leo'
    || platform === 'openai_media'
}

export function groupMatchesChannelPlatform(
  groupPlatform: GroupPlatform | undefined,
  platform: GroupPlatform,
): boolean {
  if (!groupPlatform) return false
  return groupPlatform === platform
    || (groupPlatform === 'composite' && compositeGroupAppliesToPlatform(platform))
}

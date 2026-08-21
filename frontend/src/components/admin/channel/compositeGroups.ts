import type { GroupPlatform } from '@/types'
import { CONCRETE_PLATFORM_OPTIONS } from '@/constants/platforms'

/** Platforms a composite group enables as routing and pricing tabs. */
export const COMPOSITE_ROUTE_PLATFORMS: readonly GroupPlatform[] = CONCRETE_PLATFORM_OPTIONS.map(option => option.value)

export function compositeGroupAppliesToPlatform(platform: GroupPlatform): boolean {
  return COMPOSITE_ROUTE_PLATFORMS.includes(platform)
}

export function groupMatchesChannelPlatform(
  groupPlatform: GroupPlatform | undefined,
  platform: GroupPlatform,
): boolean {
  if (!groupPlatform) return false
  return groupPlatform === platform
    || (groupPlatform === 'composite' && compositeGroupAppliesToPlatform(platform))
}

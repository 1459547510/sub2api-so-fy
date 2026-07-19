import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

import { useAppStore } from '@/stores/app'
import type { PublicSettings } from '@/types'
import { FeatureFlags, makeSidebarFlag } from '@/utils/featureFlags'

const componentPath = resolve(dirname(fileURLToPath(import.meta.url)), '../AppSidebar.vue')
const componentSource = readFileSync(componentPath, 'utf8')
const stylePath = resolve(dirname(fileURLToPath(import.meta.url)), '../../../style.css')
const styleSource = readFileSync(stylePath, 'utf8')

describe('AppSidebar custom SVG styles', () => {
  it('does not override uploaded SVG fill or stroke colors', () => {
    expect(componentSource).toContain('.sidebar-svg-icon {')
    expect(componentSource).toContain('color: currentColor;')
    expect(componentSource).toContain('display: block;')
    expect(componentSource).not.toContain('stroke: currentColor;')
    expect(componentSource).not.toContain('fill: none;')
  })
})

describe('AppSidebar scroll position persistence', () => {
  it('binds a template ref to the sidebar nav element', () => {
    expect(componentSource).toContain('ref="sidebarNavRef"')
    expect(componentSource).toContain('sidebar-nav')
  })

  it('declares sidebarNavRef in script setup', () => {
    expect(componentSource).toContain("const sidebarNavRef = ref<HTMLElement | null>(null)")
  })

  it('saves scroll position on beforeUnmount', () => {
    expect(componentSource).toContain('onBeforeUnmount')
    expect(componentSource).toContain('appStore.sidebarScrollTop')
    expect(componentSource).toContain('sidebarNavRef.value.scrollTop')
  })

  it('restores scroll position on mount', () => {
    expect(componentSource).toContain('onMounted')
    expect(componentSource).toContain('appStore.sidebarScrollTop')
    expect(componentSource).toContain('nextTick')
  })
})

describe('AppSidebar video generation feature switch', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('uses the shared opt-out feature flag for the video generation menu', () => {
    expect(FeatureFlags.videoGeneration).toMatchObject({
      key: 'video_generation_enabled',
      mode: 'opt-out',
    })
    expect(componentSource).toContain(
      'const flagVideoGeneration = makeSidebarFlag(FeatureFlags.videoGeneration)',
    )
    expect(componentSource).toMatch(
      /path: '\/video-generation',[^\n]*featureFlag: flagVideoGeneration/,
    )
  })

  it.each([
    ['settings not loaded', null, true],
    ['field missing', {}, true],
    ['explicitly disabled', { video_generation_enabled: false }, false],
    ['explicitly enabled', { video_generation_enabled: true }, true],
  ])('resolves %s to %s', (_label, settings, expected) => {
    const appStore = useAppStore()
    appStore.cachedPublicSettings = settings as PublicSettings | null

    expect(makeSidebarFlag(FeatureFlags.videoGeneration)()).toBe(expected)
  })
})

describe('AppSidebar header styles', () => {
  it('does not clip the version badge dropdown', () => {
    const sidebarHeaderBlockMatch = styleSource.match(/\.sidebar-header\s*\{[\s\S]*?\n {2}\}/)
    const sidebarBrandBlockMatch = componentSource.match(/\.sidebar-brand\s*\{[\s\S]*?\n\}/)

    expect(sidebarHeaderBlockMatch).not.toBeNull()
    expect(sidebarBrandBlockMatch).not.toBeNull()
    expect(sidebarHeaderBlockMatch?.[0]).not.toContain('@apply overflow-hidden;')
    expect(sidebarBrandBlockMatch?.[0]).not.toContain('overflow: hidden;')
  })
})

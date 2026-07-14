import { existsSync, readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

const dir = dirname(fileURLToPath(import.meta.url))
const homeViewPath = resolve(dir, '../HomeView.vue')
const telemetryHomePath = resolve(dir, '../../components/home/TelemetryHome.vue')

describe('HomeView D telemetry default', () => {
  it('uses the D telemetry experience only for the default home branch', () => {
    const homeViewSource = readFileSync(homeViewPath, 'utf8')

    expect(homeViewSource).toContain('v-if="homeContent"')
    expect(homeViewSource).toContain('<TelemetryHome')
    expect(existsSync(telemetryHomePath)).toBe(true)
  })

  it('keeps the four D phases and cleans up browser animation listeners', () => {
    const telemetrySource = existsSync(telemetryHomePath)
      ? readFileSync(telemetryHomePath, 'utf8')
      : ''

    expect(telemetrySource).toContain("const phaseOrder = ['core', 'route', 'code', 'meter']")
    expect(telemetrySource).toContain('ref="telemetryCanvas"')
    expect(telemetrySource).toContain('onBeforeUnmount')
  })

  it('provides a compact layout for wide mobile landscape viewports', () => {
    const telemetrySource = existsSync(telemetryHomePath)
      ? readFileSync(telemetryHomePath, 'utf8')
      : ''

    expect(telemetrySource).toContain('@media (min-width: 821px) and (max-height: 680px)')
  })

  it('keeps the locale menu above the lifecycle controls', () => {
    const telemetrySource = existsSync(telemetryHomePath)
      ? readFileSync(telemetryHomePath, 'utf8')
      : ''
    const navZIndex = Number(telemetrySource.match(/\.telemetry-nav\s*\{[\s\S]*?z-index:\s*(\d+)/)?.[1])
    const phasesZIndex = Number(telemetrySource.match(/\.telemetry-phases\s*\{[\s\S]*?z-index:\s*(\d+)/)?.[1])

    expect(navZIndex).toBeGreaterThan(phasesZIndex)
  })

  it('maps visitor-facing public settings into the D homepage', () => {
    const homeViewSource = readFileSync(homeViewPath, 'utf8')

    expect(homeViewSource).toContain(':site-subtitle="siteSubtitle"')
    expect(homeViewSource).toContain(':api-base-url="apiBaseUrl"')
    expect(homeViewSource).toContain(':contact-info="contactInfo"')
    expect(homeViewSource).toContain(':custom-menu-items="customMenuItems"')
    expect(homeViewSource).toContain("item.visibility === 'user'")
    expect(homeViewSource).toContain('a.sort_order - b.sort_order')
  })

  it('keeps homepage actions on one line when public metadata is present', () => {
    const telemetrySource = readFileSync(telemetryHomePath, 'utf8')
    const primaryRule = telemetrySource.match(/\.telemetry-primary\s*\{[\s\S]*?\}/)?.[0] ?? ''
    const secondaryRules = [...telemetrySource.matchAll(/\.telemetry-secondary\s*\{[\s\S]*?\}/g)]
      .map((match) => match[0])

    expect(primaryRule).toContain('white-space: nowrap')
    expect(secondaryRules.some((rule) => rule.includes('white-space: nowrap'))).toBe(true)
  })
})

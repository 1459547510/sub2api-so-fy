<template>
  <main
    ref="hero"
    class="telemetry-home"
    :data-phase="currentPhase"
    :aria-label="siteName"
  >
    <div class="telemetry-photos" aria-hidden="true">
      <img
        v-for="phase in phaseOrder"
        :key="phase"
        :src="phaseImages[phase]"
        alt=""
        class="telemetry-photo"
        :class="`telemetry-photo-${phase}`"
      />
    </div>

    <canvas
      ref="telemetryCanvas"
      class="telemetry-canvas"
      role="img"
      :aria-label="t('home.telemetry.aria.canvas')"
    ></canvas>

    <div class="telemetry-cinema" aria-hidden="true"></div>
    <div class="telemetry-haze" aria-hidden="true"></div>
    <div class="telemetry-grain" aria-hidden="true"></div>
    <div class="telemetry-flash" aria-hidden="true"></div>

    <nav class="telemetry-nav" :aria-label="t('home.telemetry.aria.mainNav')">
      <RouterLink class="telemetry-brand" to="/home" :aria-label="t('home.telemetry.aria.home')">
        <span class="telemetry-brand-mark">
          <img v-if="siteLogo" :src="siteLogo" alt="" />
          <span v-else aria-hidden="true">S</span>
        </span>
        <span class="telemetry-brand-copy">
          <span class="telemetry-brand-name">{{ siteName }}</span>
          <span
            v-if="siteSubtitle"
            class="telemetry-brand-subtitle"
            data-testid="site-subtitle"
            :title="siteSubtitle"
          >
            {{ siteSubtitle }}
          </span>
        </span>
      </RouterLink>

      <div class="telemetry-nav-actions">
        <RouterLink class="telemetry-nav-link telemetry-nav-optional" to="/models">
          {{ t('home.telemetry.nav.models') }}
        </RouterLink>
        <a
          v-if="docUrl"
          class="telemetry-nav-link telemetry-nav-optional"
          :href="docUrl"
          target="_blank"
          rel="noopener noreferrer"
        >
          {{ t('home.telemetry.nav.docs') }}
        </a>
        <RouterLink
          v-for="item in directCustomMenuItems"
          :key="item.id"
          class="telemetry-nav-link telemetry-custom-menu-direct"
          :to="`/custom/${item.id}`"
          data-testid="direct-custom-menu"
          :title="item.label"
        >
          {{ item.label }}
        </RouterLink>
        <div
          v-if="customMenuItems.length"
          ref="customMenuRef"
          class="telemetry-custom-menu"
          :class="{ 'has-desktop-overflow': overflowCustomMenuItems.length > 0 }"
        >
          <button
            type="button"
            class="telemetry-custom-menu-toggle"
            data-testid="custom-menu-toggle"
            :aria-label="t('home.telemetry.nav.more')"
            :title="t('home.telemetry.nav.more')"
            :aria-expanded="customMenuOpen"
            aria-controls="telemetry-custom-menu-panel"
            @click="toggleCustomMenu"
          >
            <Icon name="menu" size="sm" />
          </button>
          <div
            v-if="customMenuOpen"
            id="telemetry-custom-menu-panel"
            class="telemetry-custom-menu-panel"
            data-testid="custom-menu-panel"
            role="menu"
          >
            <div class="telemetry-custom-menu-desktop">
              <RouterLink
                v-for="item in overflowCustomMenuItems"
                :key="item.id"
                class="telemetry-custom-menu-item"
                :to="`/custom/${item.id}`"
                data-testid="overflow-custom-menu"
                role="menuitem"
                @click="closeCustomMenu"
              >
                {{ item.label }}
              </RouterLink>
            </div>
            <div class="telemetry-custom-menu-compact">
              <RouterLink
                v-for="item in customMenuItems"
                :key="item.id"
                class="telemetry-custom-menu-item"
                :to="`/custom/${item.id}`"
                data-testid="compact-custom-menu"
                role="menuitem"
                @click="closeCustomMenu"
              >
                {{ item.label }}
              </RouterLink>
            </div>
          </div>
        </div>
        <LocaleSwitcher class="telemetry-locale" />
        <RouterLink class="telemetry-nav-key" :to="primaryPath">
          {{
            isAuthenticated
              ? t('home.telemetry.nav.dashboard')
              : t('home.telemetry.nav.login')
          }}
        </RouterLink>
      </div>
    </nav>

    <div class="telemetry-phases" role="tablist" :aria-label="t('home.telemetry.aria.lifecycle')">
      <button
        v-for="(phase, index) in phaseOrder"
        :key="phase"
        type="button"
        role="tab"
        class="telemetry-phase"
        :class="{ active: currentPhase === phase }"
        :aria-selected="currentPhase === phase"
        @click="setPhase(phase)"
      >
        <span class="telemetry-phase-number">{{ String(index + 1).padStart(2, '0') }}</span>
        {{ t(`home.telemetry.phase.${phase}`) }}
      </button>
    </div>

    <span class="telemetry-wheel-hint" aria-hidden="true">
      {{ t('home.telemetry.wheelHint') }}
    </span>

    <section class="telemetry-copy" aria-live="polite">
      <p class="telemetry-eyebrow">
        <span class="telemetry-live-dot" aria-hidden="true"></span>
        <span>{{ siteName }} / {{ currentContent.eyebrow }}</span>
      </p>
      <h1 class="telemetry-title">
        <span>{{ currentContent.titleLine1 }}</span>
        <span>{{ currentContent.titleLine2 }}</span>
      </h1>
      <p class="telemetry-lead">{{ currentContent.lead }}</p>
      <div class="telemetry-actions">
        <RouterLink class="telemetry-primary" :to="primaryPath">
          {{ t('home.telemetry.actions.primary') }}
          <span aria-hidden="true">&rarr;</span>
        </RouterLink>
        <a
          v-if="docUrl"
          class="telemetry-secondary"
          :href="docUrl"
          target="_blank"
          rel="noopener noreferrer"
        >
          {{ t('home.telemetry.actions.docs') }}
        </a>
        <div
          v-if="apiBaseUrl || contactInfo"
          class="telemetry-site-meta"
          data-testid="site-meta"
        >
          <span
            v-if="apiBaseUrl"
            class="telemetry-site-meta-item"
            data-testid="api-base-url"
            :title="apiBaseUrl"
          >
            <span class="telemetry-site-meta-label">{{ t('home.telemetry.meta.apiBase') }}</span>
            <span class="telemetry-site-meta-value">{{ apiBaseUrl }}</span>
          </span>
          <span
            v-if="contactInfo"
            class="telemetry-site-meta-item"
            data-testid="contact-info"
            :title="contactInfo"
          >
            <span class="telemetry-site-meta-label">{{ t('home.telemetry.meta.contact') }}</span>
            <span class="telemetry-site-meta-value">{{ contactInfo }}</span>
          </span>
        </div>
      </div>
    </section>

    <div class="telemetry-signal" aria-live="polite">
      <strong>{{ currentContent.signal }}</strong>
      <span>{{ currentContent.caption }}</span>
    </div>

    <div class="telemetry-metrics" aria-live="polite">
      <div class="telemetry-metric">
        <span class="telemetry-metric-label">{{ t('home.telemetry.metrics.requestId') }}</span>
        <span class="telemetry-metric-value">{{ currentContent.request }}</span>
      </div>
      <div class="telemetry-metric">
        <span class="telemetry-metric-label">{{ t('home.telemetry.metrics.route') }}</span>
        <span class="telemetry-metric-value">{{ currentContent.route }}</span>
      </div>
      <div class="telemetry-metric">
        <span class="telemetry-metric-label">{{ t('home.telemetry.metrics.latency') }}</span>
        <span class="telemetry-metric-value">{{ currentContent.latency }}</span>
      </div>
      <div class="telemetry-metric">
        <span class="telemetry-metric-label">{{ t('home.telemetry.metrics.tokens') }}</span>
        <span class="telemetry-metric-value">{{ currentContent.tokens }}</span>
      </div>
      <div class="telemetry-metric">
        <span class="telemetry-metric-label">{{ t('home.telemetry.metrics.cost') }}</span>
        <span class="telemetry-metric-value">{{ currentContent.cost }}</span>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import coreImage from '@/assets/home/telemetry/core.webp'
import routeImage from '@/assets/home/telemetry/route.webp'
import codeImage from '@/assets/home/telemetry/code.webp'
import meterImage from '@/assets/home/telemetry/meter.webp'
import type { CustomMenuItem } from '@/types'

const props = defineProps<{
  siteName: string
  siteLogo: string
  siteSubtitle: string
  apiBaseUrl: string
  contactInfo: string
  customMenuItems: CustomMenuItem[]
  docUrl: string
  isAuthenticated: boolean
  dashboardPath: string
}>()

const { t } = useI18n()

const phaseOrder = ['core', 'route', 'code', 'meter'] as const
type Phase = (typeof phaseOrder)[number]
type Rgb = readonly [number, number, number]

interface PhaseContent {
  eyebrow: string
  titleLine1: string
  titleLine2: string
  lead: string
  signal: string
  caption: string
  request: string
  route: string
  latency: string
  tokens: string
  cost: string
}

interface Particle {
  z: number
  x: number
  y: number
  speed: number
  size: number
  alpha: number
  phase: number
  branch: number
  glyph: string
}

interface RoutePath {
  branch: number
  phase: number
}

interface ProjectedPoint {
  x: number
  y: number
  scale: number
}

const phaseImages: Record<Phase, string> = {
  core: coreImage,
  route: routeImage,
  code: codeImage,
  meter: meterImage
}

const phaseColors: Record<Phase, Rgb> = {
  core: [82, 185, 255],
  route: [50, 227, 181],
  code: [171, 140, 255],
  meter: [255, 197, 102]
}

const currentPhase = ref<Phase>('core')
const hero = ref<HTMLElement | null>(null)
const telemetryCanvas = ref<HTMLCanvasElement | null>(null)
const customMenuRef = ref<HTMLElement | null>(null)
const customMenuOpen = ref(false)
const primaryPath = computed(() => (props.isAuthenticated ? props.dashboardPath : '/login'))
const directCustomMenuItems = computed(() => props.customMenuItems.slice(0, 2))
const overflowCustomMenuItems = computed(() => props.customMenuItems.slice(2))

function getPhaseContent(phase: Phase): PhaseContent {
  const key = `home.telemetry.phases.${phase}`
  return {
    eyebrow: t(`${key}.eyebrow`),
    titleLine1: t(`${key}.titleLine1`),
    titleLine2: t(`${key}.titleLine2`),
    lead: t(`${key}.lead`),
    signal: t(`${key}.signal`),
    caption: t(`${key}.caption`),
    request: t(`${key}.request`),
    route: t(`${key}.route`),
    latency: t(`${key}.latency`),
    tokens: t(`${key}.tokens`),
    cost: t(`${key}.cost`)
  }
}

const currentContent = computed(() => getPhaseContent(currentPhase.value))

let context: CanvasRenderingContext2D
let width = 0
let height = 0
let pixelRatio = 1
let particles: Particle[] = []
let routes: RoutePath[] = []
let wheelLocked = false
let phaseTimer: number | undefined
let wheelTimer: number | undefined
let animationFrame = 0
let pulseDepth = 22
let burst = 0
let lastTime = 0
let reducedMotion = false
let motionQuery: MediaQueryList | null = null
let mounted = false

const pointer = { x: 0, y: 0, tx: 0, ty: 0 }
const glyphs = ['POST', '/v1', '{', '}', 'data', '200', 'token', 'stream', 'delta', '01']

function randomParticle(index: number, mode: Phase): Particle {
  const particle: Particle = {
    z: 1.4 + Math.random() * 22,
    x: (Math.random() - 0.5) * 5.6,
    y: (Math.random() - 0.5) * 4.2,
    speed: 0.9 + Math.random() * 1.9,
    size: 0.5 + Math.random() * 1.5,
    alpha: 0.2 + Math.random() * 0.7,
    phase: Math.random() * Math.PI * 2,
    branch: (index % 5) - 2,
    glyph: glyphs[index % glyphs.length]
  }

  if (mode === 'route') particle.z = 2 + Math.random() * 20
  if (mode === 'code') particle.speed = 1.4 + Math.random() * 2.4
  if (mode === 'meter') {
    particle.speed = 1 + Math.random() * 1.4
    particle.phase = Math.random() * Math.PI * 2
  }

  return particle
}

function seedField(mode: Phase) {
  const mobile = width < 821
  const count = mode === 'core' ? (mobile ? 38 : 64) : mode === 'code' ? (mobile ? 48 : 82) : mobile ? 44 : 74
  particles = Array.from({ length: count }, (_, index) => randomParticle(index, mode))
  routes = [-2, -1, 0, 1, 2].map((branch) => ({
    branch,
    phase: Math.random() * Math.PI * 2
  }))
  pulseDepth = 22
}

function resizeCanvas() {
  const canvas = telemetryCanvas.value
  const root = hero.value
  if (!canvas || !root) return

  const bounds = root.getBoundingClientRect()
  width = Math.max(1, Math.round(bounds.width))
  height = Math.max(1, Math.round(bounds.height))
  pixelRatio = Math.min(window.devicePixelRatio || 1, width < 821 ? 1.25 : 1.55)
  canvas.width = Math.round(width * pixelRatio)
  canvas.height = Math.round(height * pixelRatio)
  canvas.style.width = `${width}px`
  canvas.style.height = `${height}px`
  context.setTransform(pixelRatio, 0, 0, pixelRatio, 0, 0)
  seedField(currentPhase.value)

  if (reducedMotion) scheduleFrame()
}

function projection(x: number, y: number, z: number): ProjectedPoint {
  const focal = Math.min(width, height) * (width < 821 ? 0.88 : 1.02)
  const centerX = width * (width < 821 ? 0.58 : 0.71) + pointer.x * 21
  const centerY = height * (width < 821 ? 0.63 : 0.48) + pointer.y * 15
  const depth = Math.max(0.7, z)
  return {
    x: centerX + (x * focal) / depth,
    y: centerY + (y * focal) / depth,
    scale: focal / depth
  }
}

function rgba(alpha: number, values: Rgb = phaseColors[currentPhase.value]) {
  return `rgba(${values[0]},${values[1]},${values[2]},${alpha})`
}

function drawGlow(x: number, y: number, radius: number, values: Rgb, opacity: number) {
  const gradient = context.createRadialGradient(x, y, 0, x, y, radius)
  gradient.addColorStop(0, rgba(opacity, values))
  gradient.addColorStop(0.2, rgba(opacity * 0.72, values))
  gradient.addColorStop(1, rgba(0, values))
  context.fillStyle = gradient
  context.beginPath()
  context.arc(x, y, radius, 0, Math.PI * 2)
  context.fill()
}

function drawDust(particle: Particle, delta: number, speedScale = 1) {
  if (!reducedMotion) particle.z -= delta * particle.speed * speedScale * (1 + burst * 0.9)
  if (particle.z < 0.85) {
    Object.assign(particle, randomParticle(Math.floor(Math.random() * 40), currentPhase.value), {
      z: 21 + Math.random() * 3
    })
  }

  const current = projection(particle.x, particle.y, particle.z)
  const previous = projection(particle.x, particle.y, particle.z + 0.6 + particle.speed * 0.3)
  const alpha = Math.min(0.8, particle.alpha * (1.3 - particle.z / 28))
  context.strokeStyle = rgba(alpha)
  context.lineWidth = Math.max(0.5, Math.min(2.4, current.scale * 0.018 * particle.size))
  context.beginPath()
  context.moveTo(previous.x, previous.y)
  context.lineTo(current.x, current.y)
  context.stroke()
  context.fillStyle = rgba(Math.min(1, alpha + 0.18))
  context.beginPath()
  context.arc(
    current.x,
    current.y,
    Math.max(0.5, Math.min(2.8, current.scale * 0.014 * particle.size)),
    0,
    Math.PI * 2
  )
  context.fill()
}

function drawCore(delta: number, elapsed: number) {
  particles.forEach((particle) => drawDust(particle, delta, 0.58))

  if (!reducedMotion) pulseDepth -= delta * (5.2 + burst * 7)
  if (pulseDepth < 0.82) pulseDepth = 22

  const sway = Math.sin(elapsed * 0.7) * 0.16
  const pulse = projection(sway, Math.cos(elapsed * 0.55) * 0.1, pulseDepth)
  const trail = projection(sway, 0, Math.min(24, pulseDepth + 5.6))
  context.strokeStyle = rgba(0.78)
  context.lineWidth = Math.max(1.2, Math.min(7, pulse.scale * 0.035))
  context.beginPath()
  context.moveTo(trail.x, trail.y)
  context.lineTo(pulse.x, pulse.y)
  context.stroke()

  const radius = Math.max(10, Math.min(88, pulse.scale * 0.24 * (1 + burst * 0.5)))
  drawGlow(pulse.x, pulse.y, radius, phaseColors.core, 0.92)
  context.strokeStyle = rgba(0.48)
  context.lineWidth = 1
  for (let index = 0; index < 3; index += 1) {
    context.beginPath()
    context.arc(pulse.x, pulse.y, radius * (0.3 + index * 0.28), 0, Math.PI * 2)
    context.stroke()
  }
}

function branchWorld(branch: number, z: number, phase = 0) {
  const progress = 1 - z / 22
  const spread = Math.pow(Math.max(0, progress), 1.25)
  return {
    x: branch * 1.25 * spread + Math.sin(z * 0.32 + phase) * 0.16,
    y: branch * 0.52 * spread + Math.cos(z * 0.26 + phase) * 0.18
  }
}

function drawRoute(delta: number, elapsed: number) {
  routes.forEach((route, routeIndex) => {
    context.strokeStyle = rgba(route.branch === 0 ? 0.68 : 0.34 + routeIndex * 0.035)
    context.lineWidth = route.branch === 0 ? 1.5 : 1
    context.beginPath()
    for (let step = 0; step <= 26; step += 1) {
      const z = 22 - step * 0.8
      const lane = branchWorld(route.branch, z, route.phase + Math.sin(elapsed * 0.4) * 0.1)
      const point = projection(lane.x, lane.y, z)
      if (step === 0) context.moveTo(point.x, point.y)
      else context.lineTo(point.x, point.y)
    }
    context.stroke()
  })

  ;[7.5, 3.8].forEach((z, index) => {
    const left = branchWorld(-2 + index, z, routes[index].phase)
    const right = branchWorld(2 - index, z, routes[4 - index].phase)
    const start = projection(left.x, left.y, z)
    const end = projection(right.x, right.y, z)
    context.strokeStyle = rgba(0.22 + index * 0.08)
    context.lineWidth = 1
    context.beginPath()
    context.moveTo(start.x, start.y)
    context.quadraticCurveTo((start.x + end.x) / 2, (start.y + end.y) / 2 - 18, end.x, end.y)
    context.stroke()
  })

  particles.forEach((particle, index) => {
    if (!reducedMotion) particle.z -= delta * particle.speed * (1 + burst * 0.7)
    if (particle.z < 0.9) {
      Object.assign(particle, randomParticle(Math.floor(Math.random() * 50), 'route'), {
        z: 21 + Math.random() * 2
      })
    }
    const lane = branchWorld(particle.branch, particle.z, particle.phase)
    const current = projection(lane.x, lane.y, particle.z)
    const previousLane = branchWorld(particle.branch, particle.z + 0.65, particle.phase)
    const previous = projection(previousLane.x, previousLane.y, particle.z + 0.65)
    context.strokeStyle = rgba(0.5 + particle.alpha * 0.35)
    context.lineWidth = Math.max(0.7, Math.min(3, current.scale * 0.018))
    context.beginPath()
    context.moveTo(previous.x, previous.y)
    context.lineTo(current.x, current.y)
    context.stroke()
    if (index % 4 === 0) {
      drawGlow(
        current.x,
        current.y,
        Math.max(3, Math.min(14, current.scale * 0.06)),
        phaseColors.route,
        0.35
      )
    }
  })
}

function drawCode(delta: number) {
  particles.forEach((particle, index) => {
    if (!reducedMotion) particle.z -= delta * particle.speed * (1 + burst * 0.8)
    if (particle.z < 0.8) {
      Object.assign(particle, randomParticle(index, 'code'), { z: 20 + Math.random() * 4 })
    }
    if (!reducedMotion) particle.x += Math.sin(particle.phase + particle.z * 0.14) * delta * 0.06

    const current = projection(particle.x, particle.y, particle.z)
    const previous = projection(particle.x, particle.y, particle.z + 0.8)
    const fontSize = Math.max(8, Math.min(27, current.scale * 0.085))
    context.strokeStyle = rgba(0.2 + particle.alpha * 0.42)
    context.lineWidth = Math.max(0.6, Math.min(2.4, fontSize * 0.07))
    context.beginPath()
    context.moveTo(previous.x, previous.y)
    context.lineTo(current.x, current.y)
    context.stroke()
    context.fillStyle = rgba(0.35 + particle.alpha * 0.58)
    context.font = `700 ${fontSize}px ui-monospace, SFMono-Regular, Consolas, monospace`
    context.fillText(particle.glyph, current.x, current.y)
  })
}

function meterWorld(particle: Pick<Particle, 'z' | 'phase'>) {
  const normalized = Math.max(0, particle.z / 22)
  const radius = Math.pow(normalized, 1.65) * 3.2
  const angle = particle.phase + (22 - particle.z) * 0.68
  return {
    x: Math.cos(angle) * radius + 0.72 * (1 - normalized),
    y: Math.sin(angle) * radius * 0.55 + 0.32 * (1 - normalized)
  }
}

function drawMeter(delta: number) {
  context.strokeStyle = rgba(0.2)
  context.lineWidth = 1
  for (let lane = 0; lane < 4; lane += 1) {
    context.beginPath()
    for (let step = 0; step <= 28; step += 1) {
      const z = 22 - step * 0.76
      const world = meterWorld({ z, phase: lane * Math.PI * 0.5 })
      const point = projection(world.x, world.y, z)
      if (step === 0) context.moveTo(point.x, point.y)
      else context.lineTo(point.x, point.y)
    }
    context.stroke()
  }

  particles.forEach((particle, index) => {
    if (!reducedMotion) particle.z -= delta * particle.speed * (1 + burst * 0.75)
    if (particle.z < 0.8) {
      Object.assign(particle, randomParticle(index, 'meter'), { z: 20 + Math.random() * 4 })
    }

    const world = meterWorld(particle)
    const current = projection(world.x, world.y, particle.z)
    const previousZ = particle.z + 0.72
    const previousWorld = meterWorld({ z: previousZ, phase: particle.phase })
    const previous = projection(previousWorld.x, previousWorld.y, previousZ)
    context.strokeStyle = rgba(0.42 + particle.alpha * 0.46)
    context.lineWidth = Math.max(0.7, Math.min(3.2, current.scale * 0.022))
    context.beginPath()
    context.moveTo(previous.x, previous.y)
    context.lineTo(current.x, current.y)
    context.stroke()
    context.fillStyle = rgba(0.58 + particle.alpha * 0.35)
    context.beginPath()
    context.arc(current.x, current.y, Math.max(0.7, Math.min(3.4, current.scale * 0.018)), 0, Math.PI * 2)
    context.fill()
  })
}

function render(now: number) {
  animationFrame = 0
  if (!mounted) return

  const delta = Math.min((now - lastTime) / 1000, 0.04)
  const elapsed = now / 1000
  lastTime = now
  pointer.x += (pointer.tx - pointer.x) * (reducedMotion ? 1 : 0.055)
  pointer.y += (pointer.ty - pointer.y) * (reducedMotion ? 1 : 0.055)
  burst += (0 - burst) * 0.045

  context.clearRect(0, 0, width, height)
  context.save()
  context.globalCompositeOperation = 'lighter'
  if (currentPhase.value === 'core') drawCore(delta, elapsed)
  else if (currentPhase.value === 'route') drawRoute(delta, elapsed)
  else if (currentPhase.value === 'code') drawCode(delta)
  else drawMeter(delta)
  context.restore()

  if (!reducedMotion) animationFrame = window.requestAnimationFrame(render)
}

function scheduleFrame() {
  if (animationFrame) window.cancelAnimationFrame(animationFrame)
  lastTime = performance.now()
  animationFrame = window.requestAnimationFrame(render)
}

function restartTransition(duration = 820) {
  const root = hero.value
  if (!root) return

  root.classList.remove('is-changing')
  void root.offsetWidth
  root.classList.add('is-changing')
  if (phaseTimer !== undefined) window.clearTimeout(phaseTimer)
  phaseTimer = window.setTimeout(() => {
    root.classList.remove('is-changing')
    phaseTimer = undefined
  }, duration)
}

function setPhase(phase: Phase) {
  if (currentPhase.value === phase) return
  currentPhase.value = phase
  seedField(phase)
  burst = 1
  restartTransition()
  if (reducedMotion) scheduleFrame()
}

function handleWheel(event: WheelEvent) {
  event.preventDefault()
  if (wheelLocked || Math.abs(event.deltaY) < 18) return

  const currentIndex = phaseOrder.indexOf(currentPhase.value)
  const nextIndex = Math.max(
    0,
    Math.min(phaseOrder.length - 1, currentIndex + (event.deltaY > 0 ? 1 : -1))
  )
  if (nextIndex === currentIndex) return

  wheelLocked = true
  setPhase(phaseOrder[nextIndex])
  if (wheelTimer !== undefined) window.clearTimeout(wheelTimer)
  wheelTimer = window.setTimeout(() => {
    wheelLocked = false
    wheelTimer = undefined
  }, 820)
}

function handlePointerMove(event: PointerEvent) {
  const root = hero.value
  if (!root || reducedMotion) return

  const bounds = root.getBoundingClientRect()
  pointer.tx = ((event.clientX - bounds.left) / bounds.width - 0.5) * 2
  pointer.ty = ((event.clientY - bounds.top) / bounds.height - 0.5) * 2
  root.style.setProperty('--telemetry-px', `${(pointer.tx * 13).toFixed(2)}px`)
  root.style.setProperty('--telemetry-py', `${(pointer.ty * 9).toFixed(2)}px`)
}

function handlePointerDown() {
  burst = 1
  restartTransition(720)
  if (reducedMotion) scheduleFrame()
}

function closeCustomMenu() {
  customMenuOpen.value = false
}

function toggleCustomMenu() {
  customMenuOpen.value = !customMenuOpen.value
}

function handleCustomMenuClickOutside(event: PointerEvent) {
  if (customMenuRef.value && !customMenuRef.value.contains(event.target as Node)) {
    closeCustomMenu()
  }
}

function handleCustomMenuKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') closeCustomMenu()
}

function handleMotionChange(event: MediaQueryListEvent) {
  reducedMotion = event.matches
  pointer.x = 0
  pointer.y = 0
  pointer.tx = 0
  pointer.ty = 0
  hero.value?.style.setProperty('--telemetry-px', '0px')
  hero.value?.style.setProperty('--telemetry-py', '0px')
  seedField(currentPhase.value)
  scheduleFrame()
}

onMounted(() => {
  document.addEventListener('pointerdown', handleCustomMenuClickOutside)
  document.addEventListener('keydown', handleCustomMenuKeydown)

  const canvas = telemetryCanvas.value
  const root = hero.value
  if (!canvas || !root) return

  const canvasContext = canvas.getContext('2d', { alpha: true, desynchronized: true })
  if (!canvasContext) return

  context = canvasContext
  mounted = true
  motionQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
  reducedMotion = motionQuery.matches
  motionQuery.addEventListener('change', handleMotionChange)
  root.addEventListener('wheel', handleWheel, { passive: false })
  root.addEventListener('pointerdown', handlePointerDown, { passive: true })
  window.addEventListener('pointermove', handlePointerMove, { passive: true })
  window.addEventListener('resize', resizeCanvas, { passive: true })
  resizeCanvas()
  scheduleFrame()
})

onBeforeUnmount(() => {
  mounted = false
  document.removeEventListener('pointerdown', handleCustomMenuClickOutside)
  document.removeEventListener('keydown', handleCustomMenuKeydown)
  if (animationFrame) window.cancelAnimationFrame(animationFrame)
  if (phaseTimer !== undefined) window.clearTimeout(phaseTimer)
  if (wheelTimer !== undefined) window.clearTimeout(wheelTimer)
  motionQuery?.removeEventListener('change', handleMotionChange)
  hero.value?.removeEventListener('wheel', handleWheel)
  hero.value?.removeEventListener('pointerdown', handlePointerDown)
  window.removeEventListener('pointermove', handlePointerMove)
  window.removeEventListener('resize', resizeCanvas)
  animationFrame = 0
  phaseTimer = undefined
  wheelTimer = undefined
})
</script>

<style scoped>
.telemetry-home {
  --accent: #52b9ff;
  --accent-rgb: 82, 185, 255;
  --ink: #f7fbff;
  --muted: #a4b0bf;
  --telemetry-px: 0px;
  --telemetry-py: 0px;
  position: fixed;
  inset: 0;
  isolation: isolate;
  width: 100%;
  height: 100vh;
  height: 100dvh;
  overflow: hidden;
  overscroll-behavior: none;
  background: #02060b;
  color: var(--ink);
  color-scheme: dark;
  font-family: Inter, "PingFang SC", "Microsoft YaHei", sans-serif;
  touch-action: manipulation;
}

.telemetry-home[data-phase='route'] {
  --accent: #32e3b5;
  --accent-rgb: 50, 227, 181;
}

.telemetry-home[data-phase='code'] {
  --accent: #ab8cff;
  --accent-rgb: 171, 140, 255;
}

.telemetry-home[data-phase='meter'] {
  --accent: #ffc566;
  --accent-rgb: 255, 197, 102;
}

.telemetry-photos,
.telemetry-photo,
.telemetry-cinema,
.telemetry-haze,
.telemetry-grain,
.telemetry-flash {
  position: absolute;
  pointer-events: none;
}

.telemetry-photos {
  z-index: 0;
  inset: 0;
}

.telemetry-photo {
  inset: -4%;
  width: 108%;
  height: 108%;
  object-fit: cover;
  object-position: 64% 52%;
  opacity: 0;
  filter: saturate(1.35) contrast(1.25) brightness(0.42);
  transform: translate3d(calc(var(--telemetry-px) * 0.45), calc(var(--telemetry-py) * 0.45), 0)
    scale(1.08);
  transition:
    opacity 850ms ease,
    filter 1000ms ease,
    transform 1300ms cubic-bezier(0.16, 1, 0.3, 1);
  will-change: opacity, transform, filter;
}

.telemetry-photo-route {
  object-position: 62% 48%;
}

.telemetry-photo-code {
  object-position: 61% 52%;
}

.telemetry-photo-meter {
  object-position: 58% 52%;
}

.telemetry-home[data-phase='core'] .telemetry-photo-core,
.telemetry-home[data-phase='route'] .telemetry-photo-route,
.telemetry-home[data-phase='code'] .telemetry-photo-code,
.telemetry-home[data-phase='meter'] .telemetry-photo-meter {
  opacity: 1;
  transform: translate3d(var(--telemetry-px), var(--telemetry-py), 0) scale(1.02);
}

.telemetry-home[data-phase='core'] .telemetry-photo-core {
  filter: saturate(1.45) contrast(1.26) brightness(0.43);
}

.telemetry-home[data-phase='route'] .telemetry-photo-route {
  filter: saturate(1.6) contrast(1.3) brightness(0.38);
}

.telemetry-home[data-phase='code'] .telemetry-photo-code {
  filter: saturate(1.55) contrast(1.28) brightness(0.4);
}

.telemetry-home[data-phase='meter'] .telemetry-photo-meter {
  filter: saturate(1.55) contrast(1.32) brightness(0.38);
}

.telemetry-home.is-changing .telemetry-photo {
  animation: telemetry-lens-cut 900ms cubic-bezier(0.16, 1, 0.3, 1);
}

.telemetry-canvas {
  position: absolute;
  z-index: 2;
  inset: 0;
  display: block;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.telemetry-cinema,
.telemetry-haze,
.telemetry-grain,
.telemetry-flash {
  inset: 0;
}

.telemetry-cinema {
  z-index: 3;
  background:
    linear-gradient(90deg, rgba(1, 4, 10, 0.97) 0%, rgba(2, 6, 13, 0.91) 27%, rgba(2, 7, 15, 0.47) 49%, rgba(1, 4, 9, 0.12) 74%, rgba(1, 3, 8, 0.42) 100%),
    linear-gradient(180deg, rgba(0, 2, 7, 0.55), transparent 25%, transparent 68%, rgba(0, 2, 7, 0.9));
}

.telemetry-haze {
  z-index: 4;
  background:
    radial-gradient(circle at 73% 47%, rgba(var(--accent-rgb), 0.18), transparent 22%),
    radial-gradient(circle at calc(73% + var(--telemetry-px)) calc(47% + var(--telemetry-py)), transparent 0 13%, rgba(0, 0, 0, 0.11) 46%, rgba(0, 0, 0, 0.32) 100%);
  mix-blend-mode: screen;
}

.telemetry-grain {
  z-index: 28;
  opacity: 0.08;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 180 180' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='.78' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='.72'/%3E%3C/svg%3E");
  animation: telemetry-grain 260ms steps(2) infinite;
}

.telemetry-flash {
  z-index: 29;
  opacity: 0;
  background: linear-gradient(110deg, transparent 25%, rgba(var(--accent-rgb), 0.72) 50%, transparent 74%);
  mix-blend-mode: screen;
}

.telemetry-home.is-changing .telemetry-flash {
  animation: telemetry-phase-flash 760ms cubic-bezier(0.2, 0.8, 0.2, 1);
}

.telemetry-nav {
  position: absolute;
  z-index: 50;
  top: 0;
  right: 0;
  left: 0;
  display: flex;
  height: 82px;
  padding: env(safe-area-inset-top) clamp(24px, 4.5vw, 72px) 0;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.telemetry-brand,
.telemetry-nav-actions,
.telemetry-actions {
  display: flex;
  align-items: center;
}

.telemetry-brand {
  min-width: 0;
  gap: 11px;
  color: #ffffff;
  font-size: 14px;
  font-weight: 780;
  text-decoration: none;
}

.telemetry-brand-copy {
  display: grid;
  min-width: 0;
  gap: 3px;
}

.telemetry-brand-mark {
  display: grid;
  width: 26px;
  height: 26px;
  flex: 0 0 26px;
  place-items: center;
  overflow: hidden;
  background: var(--accent);
  box-shadow: 0 0 25px rgba(var(--accent-rgb), 0.28);
  color: #02060b;
  font-size: 11px;
  font-weight: 900;
  transition:
    background 500ms ease,
    box-shadow 500ms ease;
}

.telemetry-brand-mark img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.telemetry-brand-name {
  display: block;
  overflow: hidden;
  max-width: 220px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.telemetry-brand-subtitle {
  display: block;
  overflow: hidden;
  max-width: 220px;
  color: #718094;
  font: 650 8px/1 ui-monospace, SFMono-Regular, Consolas, monospace;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.telemetry-nav-actions {
  gap: clamp(14px, 1.8vw, 25px);
}

.telemetry-nav-link,
.telemetry-nav-key,
.telemetry-secondary {
  color: #aab6c4;
  font-size: 12px;
  text-decoration: none;
  transition: color 180ms ease;
}

.telemetry-nav-link:hover,
.telemetry-nav-key:hover,
.telemetry-secondary:hover {
  color: #ffffff;
}

.telemetry-nav-key {
  color: #ffffff;
  font-weight: 760;
}

.telemetry-custom-menu-direct {
  overflow: hidden;
  max-width: 104px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.telemetry-custom-menu {
  position: relative;
  display: none;
}

.telemetry-custom-menu.has-desktop-overflow {
  display: block;
}

.telemetry-custom-menu-toggle {
  display: grid;
  width: 30px;
  height: 30px;
  padding: 0;
  border: 0;
  place-items: center;
  background: transparent;
  color: #c5cfda;
  cursor: pointer;
  transition:
    background 180ms ease,
    color 180ms ease;
}

.telemetry-custom-menu-toggle:hover,
.telemetry-custom-menu-toggle:focus-visible {
  background: rgba(255, 255, 255, 0.08);
  color: #ffffff;
}

.telemetry-custom-menu-toggle:focus-visible {
  outline: 1px solid var(--accent);
  outline-offset: 3px;
}

.telemetry-custom-menu-panel {
  position: absolute;
  z-index: 70;
  top: calc(100% + 10px);
  right: 0;
  min-width: 180px;
  overflow: hidden;
  border: 1px solid rgba(164, 184, 205, 0.2);
  background: rgba(5, 11, 18, 0.94);
  box-shadow: 0 18px 55px rgba(0, 0, 0, 0.48);
  backdrop-filter: blur(18px);
}

.telemetry-custom-menu-desktop,
.telemetry-custom-menu-compact {
  display: flex;
  padding: 6px;
  flex-direction: column;
}

.telemetry-custom-menu-compact {
  display: none;
}

.telemetry-custom-menu-item {
  overflow: hidden;
  padding: 10px 11px;
  color: #c4cfda;
  font-size: 11px;
  text-decoration: none;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition:
    background 160ms ease,
    color 160ms ease;
}

.telemetry-custom-menu-item:hover,
.telemetry-custom-menu-item:focus-visible {
  background: rgba(var(--accent-rgb), 0.12);
  color: #ffffff;
  outline: none;
}

.telemetry-locale :deep(button) {
  color: #c5cfda;
}

.telemetry-phases {
  position: absolute;
  z-index: 42;
  top: 102px;
  right: clamp(24px, 4.5vw, 72px);
  display: flex;
  align-items: center;
  gap: 25px;
}

.telemetry-phase {
  position: relative;
  padding: 7px 0 11px;
  border: 0;
  background: none;
  color: #758194;
  cursor: pointer;
  font: 700 10px/1 ui-monospace, SFMono-Regular, Consolas, monospace;
  letter-spacing: 0;
  transition: color 260ms ease;
}

.telemetry-phase::after {
  position: absolute;
  bottom: 1px;
  left: 50%;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 13px rgba(var(--accent-rgb), 0.9);
  content: '';
  opacity: 0;
  transform: translateX(-50%) scale(0.5);
  transition:
    opacity 260ms ease,
    transform 260ms ease;
}

.telemetry-phase:hover,
.telemetry-phase:focus-visible,
.telemetry-phase.active {
  color: #ffffff;
}

.telemetry-phase:focus-visible {
  outline: 1px solid var(--accent);
  outline-offset: 5px;
}

.telemetry-phase.active::after {
  opacity: 1;
  transform: translateX(-50%) scale(1);
}

.telemetry-phase-number {
  margin-right: 5px;
  color: #929daf;
}

.telemetry-copy {
  position: absolute;
  z-index: 12;
  top: 23%;
  left: clamp(24px, 5.5vw, 88px);
  width: min(43vw, 650px);
  transform: translate3d(calc(var(--telemetry-px) * 0.18), calc(var(--telemetry-py) * 0.18), 0);
  transition: transform 180ms ease-out;
}

.telemetry-eyebrow {
  display: flex;
  margin: 0 0 18px;
  align-items: center;
  gap: 11px;
  color: #e0e9f3;
  font: 700 10px/1.2 ui-monospace, SFMono-Regular, Consolas, monospace;
  letter-spacing: 0;
  text-transform: uppercase;
}

.telemetry-live-dot {
  width: 6px;
  height: 6px;
  flex: 0 0 6px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 16px rgba(var(--accent-rgb), 0.95);
  animation: telemetry-blink 1.35s steps(2) infinite;
}

.telemetry-title {
  max-width: 650px;
  margin: 0;
  color: #f7fbff;
  font-size: 72px;
  font-weight: 680;
  line-height: 1;
  letter-spacing: 0;
  text-shadow: 0 9px 45px rgba(0, 0, 0, 0.58);
}

.telemetry-title span {
  display: block;
}

.telemetry-lead {
  max-width: 520px;
  min-height: 52px;
  margin: 25px 0 0;
  color: var(--muted);
  font-size: 14px;
  line-height: 1.85;
}

.telemetry-actions {
  margin-top: 30px;
  flex-wrap: wrap;
  gap: 25px;
}

.telemetry-primary {
  display: inline-flex;
  min-height: 46px;
  flex: 0 0 auto;
  padding: 0 22px;
  align-items: center;
  gap: 12px;
  background: var(--accent);
  box-shadow: 0 12px 38px rgba(var(--accent-rgb), 0.16);
  color: #02060b;
  font-size: 12px;
  font-weight: 820;
  text-decoration: none;
  white-space: nowrap;
  transition:
    transform 180ms ease,
    background 500ms ease,
    box-shadow 500ms ease;
}

.telemetry-primary:hover {
  transform: translateY(-2px);
}

.telemetry-secondary {
  flex: 0 0 auto;
  font-size: 11px;
  font-weight: 700;
  white-space: nowrap;
}

.telemetry-site-meta {
  display: flex;
  min-width: 0;
  flex: 1 1 220px;
  align-items: center;
  gap: 18px;
}

.telemetry-site-meta-item {
  display: grid;
  min-width: 0;
  max-width: 150px;
  gap: 5px;
}

.telemetry-site-meta-label {
  color: #718095;
  font: 700 7px/1 ui-monospace, SFMono-Regular, Consolas, monospace;
  text-transform: uppercase;
}

.telemetry-site-meta-value {
  overflow: hidden;
  color: #d3dce6;
  font: 650 9px/1.2 ui-monospace, SFMono-Regular, Consolas, monospace;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.telemetry-signal {
  position: absolute;
  z-index: 11;
  top: 51%;
  right: clamp(24px, 5vw, 80px);
  width: min(28vw, 390px);
  pointer-events: none;
  text-align: right;
  transform: translate3d(calc(var(--telemetry-px) * -0.34), calc(-50% + var(--telemetry-py) * -0.24), 0);
}

.telemetry-signal strong {
  display: block;
  color: rgba(255, 255, 255, 0.9);
  font: 600 58px/0.95 ui-monospace, SFMono-Regular, Consolas, monospace;
  letter-spacing: 0;
  text-shadow: 0 0 34px rgba(var(--accent-rgb), 0.3);
}

.telemetry-signal span {
  display: block;
  margin-top: 11px;
  color: var(--accent);
  font: 700 9px/1.4 ui-monospace, SFMono-Regular, Consolas, monospace;
  letter-spacing: 0;
}

.telemetry-metrics {
  position: absolute;
  z-index: 20;
  right: clamp(24px, 4.5vw, 72px);
  bottom: max(34px, env(safe-area-inset-bottom));
  left: clamp(24px, 5.5vw, 88px);
  display: grid;
  grid-template-columns: minmax(210px, 1.7fr) repeat(4, minmax(86px, 0.55fr));
  align-items: end;
  gap: clamp(18px, 3.5vw, 58px);
}

.telemetry-metric {
  min-width: 0;
}

.telemetry-metric-label {
  display: block;
  margin-bottom: 8px;
  color: #738195;
  font: 700 8px/1 ui-monospace, SFMono-Regular, Consolas, monospace;
  letter-spacing: 0;
}

.telemetry-metric-value {
  display: block;
  overflow: hidden;
  color: #eff7ff;
  font: 470 18px/1 ui-monospace, SFMono-Regular, Consolas, monospace;
  text-overflow: ellipsis;
  text-shadow: 0 0 20px rgba(0, 0, 0, 0.7);
  white-space: nowrap;
}

.telemetry-wheel-hint {
  position: absolute;
  z-index: 20;
  top: 50%;
  right: 24px;
  color: rgba(255, 255, 255, 0.44);
  font: 700 8px/1 ui-monospace, SFMono-Regular, Consolas, monospace;
  letter-spacing: 0;
  pointer-events: none;
  writing-mode: vertical-rl;
}

.telemetry-wheel-hint::after {
  color: var(--accent);
  content: ' / \2193';
  animation: telemetry-wheel-cue 1.5s ease-in-out infinite;
}

@keyframes telemetry-lens-cut {
  0% {
    clip-path: inset(0 0 0 18%);
    filter: blur(11px) saturate(1.75) brightness(0.52);
    transform: translate3d(calc(var(--telemetry-px) + 52px), var(--telemetry-py), 0) scale(1.1);
  }
  100% {
    clip-path: inset(0);
    transform: translate3d(var(--telemetry-px), var(--telemetry-py), 0) scale(1.02);
  }
}

@keyframes telemetry-phase-flash {
  0% {
    opacity: 0;
    transform: translateX(-38%) skewX(-8deg);
  }
  18% {
    opacity: 0.44;
  }
  52% {
    opacity: 0.07;
  }
  100% {
    opacity: 0;
    transform: translateX(52%) skewX(-8deg);
  }
}

@keyframes telemetry-grain {
  0% {
    transform: translate(0, 0);
  }
  50% {
    transform: translate(-1%, 1%);
  }
  100% {
    transform: translate(1%, -1%);
  }
}

@keyframes telemetry-blink {
  0%,
  42% {
    opacity: 1;
  }
  43%,
  64% {
    opacity: 0.18;
  }
  65%,
  100% {
    opacity: 1;
  }
}

@keyframes telemetry-wheel-cue {
  0%,
  100% {
    opacity: 0.35;
  }
  50% {
    opacity: 1;
  }
}

@media (max-width: 1200px) {
  .telemetry-custom-menu-direct {
    display: none;
  }

  .telemetry-custom-menu {
    display: block;
  }

  .telemetry-custom-menu-desktop {
    display: none;
  }

  .telemetry-custom-menu-compact {
    display: flex;
  }

  .telemetry-title {
    font-size: 60px;
  }

  .telemetry-signal strong {
    font-size: 48px;
  }
}

@media (max-width: 820px) {
  .telemetry-photo {
    inset: -3% -24% -3% -12%;
    width: 136%;
    height: 106%;
    object-position: 64% 52%;
  }

  .telemetry-photo-core {
    object-position: 62% 56%;
  }

  .telemetry-photo-route {
    object-position: 66% 50%;
  }

  .telemetry-photo-code {
    object-position: 62% 52%;
  }

  .telemetry-photo-meter {
    object-position: 61% 54%;
  }

  .telemetry-cinema {
    background:
      linear-gradient(180deg, rgba(1, 4, 10, 0.95) 0%, rgba(2, 6, 13, 0.78) 35%, rgba(2, 6, 13, 0.22) 59%, rgba(0, 2, 7, 0.9) 100%),
      linear-gradient(90deg, rgba(1, 4, 10, 0.72), transparent 78%);
  }

  .telemetry-haze {
    background: radial-gradient(circle at 57% 63%, rgba(var(--accent-rgb), 0.17), transparent 29%);
  }

  .telemetry-nav {
    height: 62px;
    padding-right: 18px;
    padding-left: 18px;
  }

  .telemetry-brand {
    gap: 8px;
    font-size: 13px;
  }

  .telemetry-brand-mark {
    width: 22px;
    height: 22px;
    flex-basis: 22px;
    font-size: 10px;
  }

  .telemetry-brand-name {
    max-width: 104px;
  }

  .telemetry-brand-subtitle {
    max-width: 104px;
    font-size: 6px;
  }

  .telemetry-nav-actions {
    gap: 10px;
  }

  .telemetry-nav-optional {
    display: none;
  }

  .telemetry-nav-key {
    font-size: 10px;
  }

  .telemetry-locale :deep(button) {
    padding: 6px;
  }

  .telemetry-phases {
    top: 70px;
    right: 18px;
    left: 18px;
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 3px;
  }

  .telemetry-phase {
    min-width: 0;
    font-size: 9px;
    text-align: center;
  }

  .telemetry-phase-number {
    display: none;
  }

  .telemetry-copy {
    top: 18%;
    left: 18px;
    width: calc(100% - 36px);
    transform: none;
  }

  .telemetry-eyebrow {
    margin-bottom: 12px;
    font-size: 8px;
  }

  .telemetry-title {
    max-width: 350px;
    font-size: 42px;
    line-height: 1.02;
  }

  .telemetry-lead {
    display: -webkit-box;
    max-width: 350px;
    min-height: 40px;
    margin-top: 13px;
    overflow: hidden;
    font-size: 11px;
    line-height: 1.65;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 3;
  }

  .telemetry-actions {
    margin-top: 16px;
    flex-wrap: wrap;
    gap: 18px;
  }

  .telemetry-primary {
    min-height: 39px;
    padding: 0 16px;
    font-size: 10px;
  }

  .telemetry-secondary {
    font-size: 9px;
  }

  .telemetry-site-meta {
    width: 100%;
    flex: 0 0 100%;
    gap: 14px;
  }

  .telemetry-site-meta-item {
    max-width: calc(50% - 7px);
  }

  .telemetry-signal {
    top: auto;
    right: 18px;
    bottom: 142px;
    width: 76vw;
    transform: none;
  }

  .telemetry-signal strong {
    font-size: 38px;
  }

  .telemetry-signal span {
    margin-top: 6px;
    font-size: 7px;
  }

  .telemetry-metrics {
    right: 18px;
    bottom: max(58px, calc(env(safe-area-inset-bottom) + 52px));
    left: 18px;
    grid-template-columns: minmax(0, 1.4fr) repeat(2, minmax(0, 0.75fr));
    gap: 10px;
  }

  .telemetry-metric:nth-child(4),
  .telemetry-metric:nth-child(5) {
    display: none;
  }

  .telemetry-metric-label {
    margin-bottom: 6px;
    font-size: 7px;
  }

  .telemetry-metric-value {
    font-size: 12px;
  }

  .telemetry-wheel-hint {
    display: none;
  }
}

@media (max-width: 430px) {
  .telemetry-title {
    font-size: 36px;
  }

  .telemetry-brand-name {
    max-width: 80px;
  }

  .telemetry-brand-subtitle {
    max-width: 80px;
  }

  .telemetry-secondary {
    max-width: 90px;
  }
}

@media (max-width: 820px) and (max-height: 680px) {
  .telemetry-copy {
    top: 17%;
  }

  .telemetry-eyebrow {
    margin-bottom: 8px;
  }

  .telemetry-title {
    font-size: 32px;
  }

  .telemetry-lead {
    margin-top: 9px;
    line-height: 1.45;
    -webkit-line-clamp: 2;
  }

  .telemetry-actions {
    margin-top: 10px;
  }

  .telemetry-site-meta {
    gap: 10px;
  }

  .telemetry-site-meta-item {
    gap: 3px;
  }

  .telemetry-site-meta-label {
    font-size: 6px;
  }

  .telemetry-site-meta-value {
    font-size: 8px;
  }

  .telemetry-primary {
    min-height: 35px;
  }

  .telemetry-signal {
    bottom: 108px;
  }

  .telemetry-signal strong {
    font-size: 30px;
  }

  .telemetry-metrics {
    bottom: max(48px, calc(env(safe-area-inset-bottom) + 42px));
  }
}

@media (min-width: 821px) and (max-height: 680px) {
  .telemetry-nav {
    height: 54px;
    padding-right: 28px;
    padding-left: 28px;
  }

  .telemetry-nav-actions {
    gap: 14px;
  }

  .telemetry-nav-optional {
    display: none;
  }

  .telemetry-phases {
    top: 55px;
    right: 28px;
    gap: 18px;
  }

  .telemetry-copy {
    top: 86px;
    left: 28px;
    width: 48%;
  }

  .telemetry-eyebrow {
    margin-bottom: 8px;
    font-size: 8px;
  }

  .telemetry-title {
    font-size: 38px;
  }

  .telemetry-lead {
    display: -webkit-box;
    min-height: 0;
    margin-top: 8px;
    overflow: hidden;
    font-size: 11px;
    line-height: 1.45;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
  }

  .telemetry-actions {
    margin-top: 10px;
    flex-wrap: wrap;
    gap: 18px;
  }

  .telemetry-site-meta {
    width: 100%;
    flex: 0 0 100%;
    gap: 14px;
  }

  .telemetry-site-meta-item {
    max-width: calc(50% - 7px);
    gap: 3px;
  }

  .telemetry-site-meta-label {
    font-size: 6px;
  }

  .telemetry-site-meta-value {
    font-size: 8px;
  }

  .telemetry-primary {
    min-height: 34px;
    padding: 0 16px;
    font-size: 10px;
  }

  .telemetry-secondary {
    font-size: 9px;
  }

  .telemetry-signal {
    top: 55%;
    right: 28px;
    width: 31vw;
  }

  .telemetry-signal strong {
    font-size: 30px;
  }

  .telemetry-signal span {
    margin-top: 6px;
    font-size: 7px;
  }

  .telemetry-metrics {
    right: 28px;
    bottom: max(14px, env(safe-area-inset-bottom));
    left: 28px;
    grid-template-columns: minmax(180px, 1.5fr) repeat(4, minmax(62px, 0.55fr));
    gap: 14px;
  }

  .telemetry-metric-label {
    margin-bottom: 5px;
    font-size: 7px;
  }

  .telemetry-metric-value {
    font-size: 12px;
  }

  .telemetry-wheel-hint {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .telemetry-home *,
  .telemetry-home *::before,
  .telemetry-home *::after {
    animation-duration: 0.001ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.001ms !important;
  }

  .telemetry-photo,
  .telemetry-copy,
  .telemetry-signal {
    transform: none !important;
  }
}
</style>

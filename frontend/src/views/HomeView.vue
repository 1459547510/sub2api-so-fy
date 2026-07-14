<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- SECURITY: homeContent is an administrator-only setting. -->
    <div v-else v-html="homeContent"></div>
  </div>

  <TelemetryHome
    v-else
    :site-name="siteName"
    :site-logo="siteLogo"
    :site-subtitle="siteSubtitle"
    :api-base-url="apiBaseUrl"
    :contact-info="contactInfo"
    :custom-menu-items="customMenuItems"
    :doc-url="docUrl"
    :is-authenticated="isAuthenticated"
    :dashboard-path="dashboardPath"
  />
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import TelemetryHome from '@/components/home/TelemetryHome.vue'
import { useAppStore, useAuthStore } from '@/stores'
import type { CustomMenuItem } from '@/types'
import { sanitizeUrl } from '@/utils/url'

const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle?.trim() || '')
const apiBaseUrl = computed(() => appStore.cachedPublicSettings?.api_base_url?.trim() || '')
const contactInfo = computed(() => appStore.cachedPublicSettings?.contact_info?.trim() || '')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const customMenuItems = computed<CustomMenuItem[]>(() =>
  [...(appStore.cachedPublicSettings?.custom_menu_items ?? [])]
    .filter((item) => item.visibility === 'user')
    .sort((a, b) => a.sort_order - b.sort_order)
)

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isAuthenticated = computed(() => authStore.isAuthenticated)
const dashboardPath = computed(() => authStore.isAdmin ? '/admin/dashboard' : '/dashboard')

onMounted(() => {
  authStore.checkAuth()

  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

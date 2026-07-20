<template>
  <div class="flex min-w-0 items-center gap-3">
    <span class="flex-shrink-0 font-mono text-xs font-semibold text-primary-700 dark:text-primary-300">{{ method }}</span>
    <code class="min-w-0 break-all text-sm font-semibold text-gray-900 dark:text-white">{{ path }}</code>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  method: string
  path: string
}>()
</script>

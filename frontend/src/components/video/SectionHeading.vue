<template>
  <div>
    <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ title }}</h2>
    <p class="mt-2 max-w-3xl text-sm leading-6 text-gray-500 dark:text-dark-300">{{ description }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>

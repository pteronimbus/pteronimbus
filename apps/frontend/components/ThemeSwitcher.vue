<script setup lang="ts">
const colorMode = useColorMode()

const isDark = computed({
  get () {
    return colorMode.value === 'dark'
  },
  set (value) {
    colorMode.preference = value ? 'dark' : 'light'
  }
})

const toggleTheme = () => {
  isDark.value = !isDark.value
}
</script>

<template>
  <ClientOnly>
    <button
      @click="toggleTheme"
      class="relative inline-flex items-center w-12 h-6 bg-gray-200 dark:bg-gray-700 rounded-full transition-colors duration-300 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
      aria-label="Toggle theme"
    >
      <!-- Background icons -->
      <div aria-hidden="true" class="absolute inset-0 flex items-center justify-between px-2">
        <UIcon
          name="i-heroicons-sun-20-solid"
          class="w-3 h-3 text-yellow-400 transition-opacity duration-300"
          :class="isDark ? 'opacity-0' : 'opacity-100'"
        />
        <UIcon
          name="i-heroicons-moon-20-solid"
          class="w-3 h-3 text-blue-400 transition-opacity duration-300"
          :class="isDark ? 'opacity-100' : 'opacity-0'"
        />
      </div>

      <!-- Slider -->
      <div
        aria-hidden="true"
        class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow-sm transition-transform duration-300 ease-in-out flex items-center justify-center"
        :class="{ 'translate-x-6': isDark }"
      >
        <!-- Both icons are present and we fade/scale them -->
        <UIcon
          name="i-heroicons-sun-20-solid"
          class="absolute w-3 h-3 text-yellow-500 transition-all duration-300 ease-in-out"
          :class="isDark ? 'opacity-0 scale-50' : 'opacity-100 scale-100'"
        />
        <UIcon
          name="i-heroicons-moon-20-solid"
          class="absolute w-3 h-3 text-blue-500 transition-all duration-300 ease-in-out"
          :class="isDark ? 'opacity-100 scale-100' : 'opacity-0 scale-50'"
        />
      </div>
    </button>
    <template #fallback>
      <div class="w-12 h-6 bg-gray-200 dark:bg-gray-700 rounded-full" />
    </template>
  </ClientOnly>
</template>

<style scoped>
button:hover + div {
  opacity: 1;
}
</style> 
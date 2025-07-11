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
    <div class="relative">
      <!-- Toggle Container -->
      <button
        @click="toggleTheme"
        class="relative inline-flex items-center w-12 h-6 bg-gray-200 dark:bg-gray-700 rounded-full transition-all duration-300 ease-in-out hover:bg-gray-300 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-white dark:focus:ring-offset-gray-800 shadow-md hover:shadow-lg"
        :class="{ 'bg-blue-100 dark:bg-blue-900': isDark }"
        aria-label="Toggle theme"
      >
        <!-- Toggle Slider -->
        <div
          class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow-sm transition-all duration-300 ease-in-out flex items-center justify-center"
          :class="{ 'translate-x-6': isDark, 'bg-yellow-50': !isDark, 'bg-blue-50': isDark }"
        >
          <!-- Sun Icon -->
          <UIcon
            name="i-heroicons-sun-20-solid"
            class="w-3 h-3 text-yellow-500 transition-all duration-300 ease-in-out"
            :class="{ 'opacity-0 scale-0 rotate-180': isDark, 'opacity-100 scale-100 rotate-0': !isDark }"
          />
          <!-- Moon Icon -->
          <UIcon
            name="i-heroicons-moon-20-solid"
            class="w-3 h-3 text-blue-600 absolute transition-all duration-300 ease-in-out"
            :class="{ 'opacity-100 scale-100 rotate-0': isDark, 'opacity-0 scale-0 -rotate-180': !isDark }"
          />
        </div>
        
        <!-- Background Icons -->
        <div class="absolute inset-0 flex items-center justify-between px-2">
          <!-- Sun Background -->
          <UIcon
            name="i-heroicons-sun-20-solid"
            class="w-2.5 h-2.5 text-yellow-400/40 transition-all duration-300 ease-in-out"
            :class="{ 'opacity-0': isDark, 'opacity-100': !isDark }"
          />
          <!-- Moon Background -->
          <UIcon
            name="i-heroicons-moon-20-solid"
            class="w-2.5 h-2.5 text-blue-400/40 transition-all duration-300 ease-in-out"
            :class="{ 'opacity-100': isDark, 'opacity-0': !isDark }"
          />
        </div>
      </button>
      
      <!-- Tooltip -->
      <div
        class="absolute -bottom-8 left-1/2 transform -translate-x-1/2 px-2 py-1 bg-gray-800 dark:bg-gray-200 text-white dark:text-gray-800 text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap"
      >
        {{ isDark ? 'Switch to light' : 'Switch to dark' }}
      </div>
    </div>
    
    <template #fallback>
      <div class="w-12 h-6 bg-gray-200 dark:bg-gray-700 rounded-full animate-pulse" />
    </template>
  </ClientOnly>
</template>

<style scoped>
button:hover + div {
  opacity: 1;
}
</style> 
<template>
  <div class="mb-6 flex flex-col sm:flex-row gap-4">
    <div class="flex-1">
      <UInput
        :model-value="searchQuery"
        :placeholder="searchPlaceholder"
        icon="i-heroicons-magnifying-glass-20-solid"
        size="md"
        @update:model-value="updateSearch"
      />
    </div>
    <div v-if="filters.length > 0" class="flex gap-2">
      <USelect
        v-for="filter in filters"
        :key="filter.key"
        :model-value="filter.value"
        :options="filter.options"
        size="md"
        :class="filter.class || 'w-40'"
        @update:model-value="(value) => updateFilter(filter.key, String(value || ''))"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
interface FilterOption {
  label: string
  value: string
}

interface Filter {
  key: string
  value: string
  options: FilterOption[]
  class?: string
}

interface Props {
  searchQuery: string
  searchPlaceholder?: string
  filters?: Filter[]
}

const props = withDefaults(defineProps<Props>(), {
  searchPlaceholder: 'Search...',
  filters: () => []
})

const emit = defineEmits<{
  'update:searchQuery': [value: string]
  'update:filter': [key: string, value: string]
}>()

const updateSearch = (value: string) => {
  emit('update:searchQuery', value)
}

const updateFilter = (key: string, value: string) => {
  emit('update:filter', key, value)
}
</script> 
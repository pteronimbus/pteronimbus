<script setup lang="ts">
definePageMeta({
  middleware: 'guest'
})
onMounted(() => {
  navigateTo('/login')
})
</script>
<template>
  <div>Redirecting...</div>
</template> 
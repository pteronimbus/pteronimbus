<script setup lang="ts">
definePageMeta({
  auth: false
})
onMounted(() => {
  navigateTo('/login')
})
</script>
<template>
  <div>Redirecting...</div>
</template> 
<template>
  <div class="min-h-screen flex">
    <!-- Sidebar -->
    <aside class="w-64 bg-gray-100 dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex-shrink-0">
      <nav class="mt-8 px-2">
        <ULink to="/dashboard" class="flex items-center px-4 py-2 text-gray-700 dark:text-gray-200 rounded-md hover:bg-gray-200 dark:hover:bg-gray-700" active-class="bg-gray-200 dark:bg-gray-700 font-semibold">
          <UIcon name="i-heroicons-home" class="mr-3 h-6 w-6" />
          <span>Dashboard</span>
        </ULink>
        <ULink to="/servers" class="flex items-center px-4 py-2 mt-2 text-gray-700 dark:text-gray-200 rounded-md hover:bg-gray-200 dark:hover:bg-gray-700" active-class="bg-gray-200 dark:bg-gray-700 font-semibold">
          <UIcon name="i-heroicons-server" class="mr-3 h-6 w-6" />
          <span>Servers</span>
        </ULink>
      </nav>
    </aside>

    <div class="flex-1 flex flex-col">
      <!-- User Menu Area -->
      <div class="flex justify-end p-4">
        <UserMenu />
      </div>
      <!-- Main Content -->
      <main class="p-6 flex-1">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import UserMenu from '~/components/UserMenu.vue';
</script> 
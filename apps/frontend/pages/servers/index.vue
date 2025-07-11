<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const servers = [
  { id: 1, name: 'Minecraft Survival', game: 'Minecraft', status: 'Online', players: '10/20', ip: '192.168.1.100' },
  { id: 2, name: 'Valheim Dedicated', game: 'Valheim', status: 'Online', players: '5/10', ip: '192.168.1.101' },
  { id: 3, name: 'CS:GO Competitive', game: 'CS:GO', status: 'Offline', players: '0/10', ip: '192.168.1.102' },
  { id: 4, name: 'Terraria Adventure', game: 'Terraria', status: 'Online', players: '8/16', ip: '192.168.1.103' }
]

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'game', label: 'Game' },
  { key: 'status', label: 'Status' },
  { key: 'players', label: 'Players' },
  { key: 'ip', label: 'IP Address' },
  { key: 'actions' }
]

const items = (row) => [
  [{
    label: 'Edit',
    icon: 'i-heroicons-pencil-square-20-solid'
  }, {
    label: 'Delete',
    icon: 'i-heroicons-trash-20-solid'
  }]
]
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">Servers</h1>
      <UButton icon="i-heroicons-plus-circle" size="lg">Create Server</UButton>
    </div>

    <UTable :rows="servers" :columns="columns">
      <template #status-data="{ row }">
        <UBadge :color="row.status === 'Online' ? 'green' : 'red'" variant="subtle">{{ row.status }}</UBadge>
      </template>
      <template #actions-data="{ row }">
        <UDropdown :items="items(row)">
          <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
        </UDropdown>
      </template>
    </UTable>
  </div>
</template> 
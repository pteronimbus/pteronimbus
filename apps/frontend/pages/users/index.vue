<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const router = useRouter()

// Mock user data - in real app this would come from API
const users = ref([
  { 
    id: 1, 
    name: 'John Doe', 
    email: 'john@example.com', 
    role: 'admin', 
    status: 'online', 
    lastSeen: '2 minutes ago', 
    serversAccess: 5,
    avatar: null
  },
  { 
    id: 2, 
    name: 'Jane Smith', 
    email: 'jane@example.com', 
    role: 'moderator', 
    status: 'offline', 
    lastSeen: '1 hour ago', 
    serversAccess: 3,
    avatar: null
  },
  { 
    id: 3, 
    name: 'Bob Wilson', 
    email: 'bob@example.com', 
    role: 'user', 
    status: 'online', 
    lastSeen: 'Just now', 
    serversAccess: 2,
    avatar: null
  },
  { 
    id: 4, 
    name: 'Alice Brown', 
    email: 'alice@example.com', 
    role: 'user', 
    status: 'banned', 
    lastSeen: '2 days ago', 
    serversAccess: 0,
    avatar: null
  },
  { 
    id: 5, 
    name: 'Charlie Davis', 
    email: 'charlie@example.com', 
    role: 'user', 
    status: 'suspended', 
    lastSeen: '5 hours ago', 
    serversAccess: 1,
    avatar: null
  }
])

const searchQuery = ref('')
const selectedStatus = ref('all')
const selectedRole = ref('all')

const statusOptions = [
  { value: 'all', label: 'All Status' },
  { value: 'online', label: t('users.status.online') },
  { value: 'offline', label: t('users.status.offline') },
  { value: 'banned', label: t('users.status.banned') },
  { value: 'suspended', label: t('users.status.suspended') }
]

const roleOptions = [
  { value: 'all', label: 'All Roles' },
  { value: 'admin', label: t('users.roles.admin') },
  { value: 'moderator', label: t('users.roles.moderator') },
  { value: 'user', label: t('users.roles.user') }
]

const columns = [
  { key: 'user', label: t('users.columns.name'), id: 'user' },
  { key: 'email', label: t('users.columns.email'), id: 'email' },
  { key: 'role', label: t('users.columns.role'), id: 'role' },
  { key: 'status', label: t('users.columns.status'), id: 'status' },
  { key: 'lastSeen', label: t('users.columns.lastSeen'), id: 'lastSeen' },
  { key: 'serversAccess', label: t('users.columns.serversAccess'), id: 'serversAccess' },
  { key: 'actions', label: t('users.columns.actions'), id: 'actions' }
]

// Filtered users based on search and filters
const filteredUsers = computed(() => {
  return users.value.filter(user => {
    const matchesSearch = searchQuery.value === '' || 
      user.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      user.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesStatus = selectedStatus.value === 'all' || user.status === selectedStatus.value
    const matchesRole = selectedRole.value === 'all' || user.role === selectedRole.value
    
    return matchesSearch && matchesStatus && matchesRole
  })
})

// Action items for dropdown
const getActionItems = (user) => [
  [{
    label: t('users.actions.viewDetails'),
    icon: 'i-heroicons-eye-20-solid',
    click: () => router.push(`/users/${user.id}`)
  }],
  [{
    label: t('users.actions.edit'),
    icon: 'i-heroicons-pencil-square-20-solid',
    click: () => router.push(`/users/${user.id}/edit`)
  }, {
    label: t('users.actions.resetPassword'),
    icon: 'i-heroicons-key-20-solid',
    click: () => resetPassword(user)
  }, {
    label: t('users.actions.changeRole'),
    icon: 'i-heroicons-shield-check-20-solid',
    click: () => changeRole(user)
  }],
  [{
    label: user.status === 'banned' ? t('users.actions.unban') : t('users.actions.ban'),
    icon: user.status === 'banned' ? 'i-heroicons-check-circle-20-solid' : 'i-heroicons-no-symbol-20-solid',
    click: () => toggleBan(user)
  }, {
    label: user.status === 'suspended' ? 'Unsuspend' : t('users.actions.suspend'),
    icon: user.status === 'suspended' ? 'i-heroicons-play-20-solid' : 'i-heroicons-pause-20-solid',
    click: () => toggleSuspend(user)
  }],
  [{
    label: t('users.actions.delete'),
    icon: 'i-heroicons-trash-20-solid',
    click: () => deleteUser(user)
  }]
]

// Helper functions for user actions
const resetPassword = (user) => {
  // Implementation for password reset
  console.log('Reset password for:', user.name)
}

const changeRole = (user) => {
  // Implementation for role change
  console.log('Change role for:', user.name)
}

const toggleBan = (user) => {
  const index = users.value.findIndex(u => u.id === user.id)
  if (index !== -1) {
    users.value[index].status = user.status === 'banned' ? 'offline' : 'banned'
  }
}

const toggleSuspend = (user) => {
  const index = users.value.findIndex(u => u.id === user.id)
  if (index !== -1) {
    users.value[index].status = user.status === 'suspended' ? 'offline' : 'suspended'
  }
}

const deleteUser = (user) => {
  // Implementation for user deletion
  console.log('Delete user:', user.name)
}

// Get status color
const getStatusColor = (status) => {
  switch (status) {
    case 'online': return 'green'
    case 'offline': return 'gray'
    case 'banned': return 'red'
    case 'suspended': return 'yellow'
    default: return 'gray'
  }
}

// Get role color
const getRoleColor = (role) => {
  switch (role) {
    case 'admin': return 'red'
    case 'moderator': return 'blue'
    case 'user': return 'gray'
    default: return 'gray'
  }
}

// Navigation functions
const viewUser = (user) => {
  router.push(`/users/${user.id}`)
}

const createUser = () => {
  router.push('/users/create')
}
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">{{ t('users.title') }}</h1>
        <p class="mt-1 text-gray-500 dark:text-gray-400">
          Manage users and their permissions
        </p>
      </div>
      <div class="mt-4 sm:mt-0">
        <UButton 
          icon="i-heroicons-plus-circle" 
          size="lg"
          @click="createUser"
        >
          {{ t('users.createUser') }}
        </UButton>
      </div>
    </div>

    <!-- Filters -->
    <div class="mb-6 flex flex-col sm:flex-row gap-4">
      <div class="flex-1">
        <UInput
          v-model="searchQuery"
          :placeholder="t('common.search') + ' users...'"
          icon="i-heroicons-magnifying-glass-20-solid"
          size="md"
        />
      </div>
      <div class="flex gap-2">
        <USelect
          v-model="selectedStatus"
          :options="statusOptions"
          size="md"
          class="w-40"
        />
        <USelect
          v-model="selectedRole"
          :options="roleOptions"
          size="md"
          class="w-40"
        />
      </div>
    </div>

    <!-- Users Table -->
    <UCard>
      <UTable :rows="filteredUsers" :columns="columns">
        <!-- User column with avatar and name -->
        <template #user-data="{ row }">
          <div class="flex items-center gap-3">
            <UAvatar
              :src="row.avatar"
              :alt="row.name"
              size="sm"
              :ui="{ background: 'bg-primary-100 dark:bg-primary-900' }"
            >
              <span class="text-xs font-medium text-primary-600 dark:text-primary-400">
                {{ row.name.split(' ').map(n => n[0]).join('') }}
              </span>
            </UAvatar>
            <div>
              <button
                @click="viewUser(row)"
                class="font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400 transition-colors"
              >
                {{ row.name }}
              </button>
            </div>
          </div>
        </template>

        <!-- Email column -->
        <template #email-data="{ row }">
          <span class="text-gray-600 dark:text-gray-400">{{ row.email }}</span>
        </template>

        <!-- Role column -->
        <template #role-data="{ row }">
          <UBadge 
            :color="getRoleColor(row.role)" 
            variant="subtle"
            class="capitalize"
          >
            {{ t(`users.roles.${row.role}`) }}
          </UBadge>
        </template>

        <!-- Status column -->
        <template #status-data="{ row }">
          <UBadge 
            :color="getStatusColor(row.status)" 
            variant="subtle"
            class="capitalize"
          >
            {{ t(`users.status.${row.status}`) }}
          </UBadge>
        </template>

        <!-- Last Seen column -->
        <template #lastSeen-data="{ row }">
          <span class="text-sm text-gray-500 dark:text-gray-400">{{ row.lastSeen }}</span>
        </template>

        <!-- Servers Access column -->
        <template #serversAccess-data="{ row }">
          <div class="flex items-center gap-2">
            <UIcon name="i-heroicons-server-20-solid" class="w-4 h-4 text-gray-400" />
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.serversAccess }}</span>
          </div>
        </template>

        <!-- Actions column -->
        <template #actions-data="{ row }">
          <div class="flex items-center gap-2">
            <UButton 
              color="blue" 
              variant="ghost" 
              icon="i-heroicons-eye-20-solid"
              size="sm"
              @click="viewUser(row)"
            />
            <UDropdown :items="getActionItems(row)">
              <UButton 
                color="gray" 
                variant="ghost" 
                icon="i-heroicons-ellipsis-horizontal-20-solid"
                size="sm"
              />
            </UDropdown>
          </div>
        </template>
      </UTable>

      <!-- Empty state -->
      <div v-if="filteredUsers.length === 0" class="text-center py-12">
        <UIcon name="i-heroicons-users-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
          {{ searchQuery || selectedStatus !== 'all' || selectedRole !== 'all' ? 'No users found' : t('users.noUsers') }}
        </h3>
        <p class="text-gray-500 dark:text-gray-400 mb-6">
          {{ searchQuery || selectedStatus !== 'all' || selectedRole !== 'all' 
            ? 'Try adjusting your search or filters' 
            : 'Get started by creating your first user' }}
        </p>
        <UButton 
          v-if="!searchQuery && selectedStatus === 'all' && selectedRole === 'all'"
          @click="createUser"
          icon="i-heroicons-plus-circle"
        >
          {{ t('users.createUser') }}
        </UButton>
      </div>
    </UCard>
  </div>
</template> 
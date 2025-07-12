<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const userId = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

// Mock user data - in real app this would come from API
const user = ref({
  id: parseInt(userId),
  name: 'John Doe',
  email: 'john@example.com',
  role: 'admin',
  status: 'online',
  lastSeen: '2 minutes ago',
  joinedAt: '2024-01-15',
  avatar: null as string | null,
  serversAccess: 5,
  totalPlaytime: '240 hours',
  totalSessions: 156,
  preferredServers: ['Minecraft Survival', 'Valheim Dedicated'],
  recentActivity: [
    { id: 1, action: 'Joined server', server: 'Minecraft Survival', timestamp: '2 minutes ago' },
    { id: 2, action: 'Left server', server: 'Valheim Dedicated', timestamp: '1 hour ago' },
    { id: 3, action: 'Banned user', target: 'BadPlayer', timestamp: '2 hours ago' },
    { id: 4, action: 'Started server', server: 'Terraria Adventure', timestamp: '3 hours ago' },
    { id: 5, action: 'Updated settings', server: 'Minecraft Survival', timestamp: '5 hours ago' }
  ],
  permissions: {
    canCreateServers: true,
    canManageUsers: true,
    canViewLogs: true,
    canModifySettings: true,
    canAccessConsole: true
  },
  sessions: [
    { id: 1, server: 'Minecraft Survival', duration: '2h 30m', startTime: '2024-01-20 14:30', endTime: '2024-01-20 17:00' },
    { id: 2, server: 'Valheim Dedicated', duration: '1h 45m', startTime: '2024-01-20 10:15', endTime: '2024-01-20 12:00' },
    { id: 3, server: 'CS:GO Competitive', duration: '45m', startTime: '2024-01-19 20:00', endTime: '2024-01-19 20:45' }
  ]
})

const activeTab = ref('profile')

const tabs = [
  { key: 'profile', label: t('users.details.profile'), icon: 'i-heroicons-user-20-solid' },
  { key: 'permissions', label: t('users.details.permissions'), icon: 'i-heroicons-shield-check-20-solid' },
  { key: 'activity', label: t('users.details.activity'), icon: 'i-heroicons-clock-20-solid' },
  { key: 'sessions', label: t('users.details.sessions'), icon: 'i-heroicons-computer-desktop-20-solid' }
]

// Helper functions
const getStatusColor = (status: string) => {
  switch (status) {
    case 'online': return 'success'
    case 'offline': return 'neutral'
    case 'banned': return 'error'
    case 'suspended': return 'warning'
    default: return 'neutral'
  }
}

const getRoleColor = (role: string) => {
  switch (role) {
    case 'admin': return 'error'
    case 'moderator': return 'primary'
    case 'user': return 'neutral'
    default: return 'neutral'
  }
}

const toggleBan = () => {
  user.value.status = user.value.status === 'banned' ? 'offline' : 'banned'
}

const toggleSuspend = () => {
  user.value.status = user.value.status === 'suspended' ? 'offline' : 'suspended'
}

const editUser = () => {
  router.push(`/users/${userId}/edit`)
}

const deleteUser = () => {
  // Implementation for user deletion
  console.log('Delete user:', user.value.name)
}

const resetPassword = () => {
  // Implementation for password reset
  console.log('Reset password for:', user.value.name)
}

const changeRole = () => {
  // Implementation for role change
  console.log('Change role for:', user.value.name)
}

const goBack = () => {
  router.push('/users')
}
</script>

<template>
  <div>
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center mb-4">
        <UButton 
          color="neutral" 
          variant="ghost" 
          icon="i-heroicons-arrow-left-20-solid"
          class="text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200"
          @click="goBack"
        />
        <h1 class="ml-4 text-3xl font-bold text-gray-800 dark:text-gray-100">{{ t('users.details.title') }}</h1>
      </div>
      
      <!-- User Info Card -->
      <UCard>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-6">
            <UAvatar
              :src="user.avatar || undefined"
              :alt="user.name"
              size="lg"
            >
              <span class="text-xl font-medium text-primary-600 dark:text-primary-400">
                {{ user.name.split(' ').map(n => n[0]).join('') }}
              </span>
            </UAvatar>
            <div>
              <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ user.name }}</h2>
              <p class="text-gray-600 dark:text-gray-400">{{ user.email }}</p>
              <div class="flex items-center gap-3 mt-2">
                <UBadge 
                  :color="getRoleColor(user.role)" 
                  variant="subtle"
                  class="capitalize"
                  :class="[
                    user.role === 'admin' ? 'text-red-700 dark:text-red-300' : '',
                    user.role === 'moderator' ? 'text-blue-700 dark:text-blue-300' : '',
                    user.role === 'user' ? 'text-gray-700 dark:text-gray-300' : ''
                  ]"
                >
                  {{ t(`users.roles.${user.role}`) }}
                </UBadge>
                <UBadge 
                  :color="getStatusColor(user.status)" 
                  variant="subtle"
                  class="capitalize"
                  :class="[
                    user.status === 'online' ? 'text-green-700 dark:text-green-300' : '',
                    user.status === 'offline' ? 'text-gray-700 dark:text-gray-300' : '',
                    user.status === 'banned' ? 'text-red-700 dark:text-red-300' : '',
                    user.status === 'suspended' ? 'text-yellow-700 dark:text-yellow-300' : ''
                  ]"
                >
                  {{ t(`users.status.${user.status}`) }}
                </UBadge>
              </div>
            </div>
          </div>
          
          <!-- Action buttons -->
          <div class="flex items-center gap-2">
            <UButton 
              color="primary" 
              variant="soft" 
              icon="i-heroicons-pencil-square-20-solid"
              class="text-blue-700 dark:text-blue-300"
              @click="editUser"
            >
              {{ t('users.actions.edit') }}
            </UButton>
            <UDropdownMenu :items="[
              [{
                label: t('users.actions.resetPassword'),
                icon: 'i-heroicons-key-20-solid',
                click: resetPassword
              }, {
                label: t('users.actions.changeRole'),
                icon: 'i-heroicons-shield-check-20-solid',
                click: changeRole
              }],
              [{
                label: user.status === 'banned' ? t('users.actions.unban') : t('users.actions.ban'),
                icon: user.status === 'banned' ? 'i-heroicons-check-circle-20-solid' : 'i-heroicons-no-symbol-20-solid',
                click: toggleBan
              }, {
                label: user.status === 'suspended' ? 'Unsuspend' : t('users.actions.suspend'),
                icon: user.status === 'suspended' ? 'i-heroicons-play-20-solid' : 'i-heroicons-pause-20-solid',
                click: toggleSuspend
              }],
              [{
                label: t('users.actions.delete'),
                icon: 'i-heroicons-trash-20-solid',
                click: deleteUser
              }]
            ]">
              <UButton 
                color="neutral" 
                variant="ghost" 
                icon="i-heroicons-ellipsis-horizontal-20-solid"
              />
            </UDropdownMenu>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ user.serversAccess }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ t('users.columns.serversAccess') }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ user.totalPlaytime }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">Total Playtime</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ user.totalSessions }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">Total Sessions</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ user.lastSeen }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ t('users.columns.lastSeen') }}</div>
        </div>
      </UCard>
    </div>

    <!-- Tabs -->
    <div class="mb-6">
      <div class="border-b border-gray-200 dark:border-gray-700">
        <nav class="flex space-x-8">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            @click="activeTab = tab.key"
            :class="[
              'flex items-center gap-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors',
              activeTab === tab.key
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
            ]"
          >
            <UIcon :name="tab.icon" class="w-4 h-4" />
            {{ tab.label }}
          </button>
        </nav>
      </div>
    </div>

    <!-- Tab Content -->
    <div>
      <!-- Profile Tab -->
      <div v-if="activeTab === 'profile'">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- User Information -->
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">User Information</h3>
            </template>
            <div class="space-y-4">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Name:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ user.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Email:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ user.email }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Role:</span>
                <UBadge :color="getRoleColor(user.role)" variant="subtle" class="capitalize">
                  {{ t(`users.roles.${user.role}`) }}
                </UBadge>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Status:</span>
                <UBadge :color="getStatusColor(user.status)" variant="subtle" class="capitalize">
                  {{ t(`users.status.${user.status}`) }}
                </UBadge>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Joined:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ user.joinedAt }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Last Seen:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ user.lastSeen }}</span>
              </div>
            </div>
          </UCard>

          <!-- Preferred Servers -->
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Preferred Servers</h3>
            </template>
            <div class="space-y-3">
              <div 
                v-for="server in user.preferredServers" 
                :key="server"
                class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
              >
                <div class="flex items-center gap-3">
                  <UIcon name="i-heroicons-server-20-solid" class="w-5 h-5 text-gray-400" />
                  <span class="font-medium text-gray-900 dark:text-gray-100">{{ server }}</span>
                </div>
                <UButton 
                  color="neutral" 
                  variant="ghost" 
                  icon="i-heroicons-arrow-top-right-on-square-20-solid"
                  size="sm"
                  class="text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200"
                  @click="router.push(`/servers/${server.toLowerCase().replace(/\s+/g, '-')}`)"
                />
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- Permissions Tab -->
      <div v-if="activeTab === 'permissions'">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('users.details.permissions') }}</h3>
          </template>
          <div class="space-y-4">
            <div 
              v-for="(value, key) in user.permissions" 
              :key="key"
              class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
              <span class="font-medium text-gray-900 dark:text-gray-100 capitalize">
                {{ key.replace(/([A-Z])/g, ' $1').trim() }}
              </span>
              <UToggle v-model="user.permissions[key]" />
            </div>
          </div>
        </UCard>
      </div>

      <!-- Activity Tab -->
      <div v-if="activeTab === 'activity'">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('users.details.activity') }}</h3>
          </template>
          <div class="space-y-4">
            <div 
              v-for="activity in user.recentActivity" 
              :key="activity.id"
              class="flex items-start gap-4 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
              <div class="flex-shrink-0 w-10 h-10 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
                <UIcon name="i-heroicons-bolt-20-solid" class="w-5 h-5 text-blue-600 dark:text-blue-400" />
              </div>
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-gray-100">
                  {{ activity.action }}
                  <span v-if="activity.server" class="text-blue-600 dark:text-blue-400">{{ activity.server }}</span>
                  <span v-if="activity.target" class="text-red-600 dark:text-red-400">{{ activity.target }}</span>
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ activity.timestamp }}</p>
              </div>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Sessions Tab -->
      <div v-if="activeTab === 'sessions'">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('users.details.sessions') }}</h3>
          </template>
          <div class="space-y-4">
            <div 
              v-for="session in user.sessions" 
              :key="session.id"
              class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
              <div class="flex items-center gap-4">
                <div class="flex-shrink-0 w-10 h-10 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center">
                  <UIcon name="i-heroicons-play-20-solid" class="w-5 h-5 text-green-600 dark:text-green-400" />
                </div>
                <div>
                  <p class="font-medium text-gray-900 dark:text-gray-100">{{ session.server }}</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">{{ session.startTime }} - {{ session.endTime }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ session.duration }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400">Duration</p>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template> 
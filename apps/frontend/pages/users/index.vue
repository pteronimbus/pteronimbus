<script setup lang="ts">
import { h, resolveComponent } from 'vue'

definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const router = useRouter()

// Define user interface for better type safety
interface User {
  id: number
  name: string
  email: string
  role: string
  status: string
  lastSeen: string
  serversAccess: number
  avatar: string | null
}

// Mock user data - in real app this would come from API
const users = ref<User[]>([
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

// Resolve components for use in cell renderers
const UAvatar = resolveComponent('UAvatar')
const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const UDropdownMenu = resolveComponent('UDropdownMenu')
const UIcon = resolveComponent('UIcon')

// Helper functions for user actions
const resetPassword = (user: User) => {
  // Implementation for password reset
  console.log('Reset password for:', user.name)
}

const changeRole = (user: User) => {
  // Implementation for role change
  console.log('Change role for:', user.name)
}

const toggleBan = (user: User) => {
  const index = users.value.findIndex(u => u.id === user.id)
  if (index !== -1) {
    users.value[index].status = user.status === 'banned' ? 'offline' : 'banned'
  }
}

const toggleSuspend = (user: User) => {
  const index = users.value.findIndex(u => u.id === user.id)
  if (index !== -1) {
    users.value[index].status = user.status === 'suspended' ? 'offline' : 'suspended'
  }
}

const deleteUser = (user: User) => {
  // Implementation for user deletion
  console.log('Delete user:', user.name)
}

// Get status color
const getStatusColor = (status: string) => {
  switch (status) {
    case 'online': return 'success'
    case 'offline': return 'neutral'
    case 'banned': return 'error'
    case 'suspended': return 'warning'
    default: return 'neutral'
  }
}

// Get role color
const getRoleColor = (role: string) => {
  switch (role) {
    case 'admin': return 'error'
    case 'moderator': return 'primary'
    case 'user': return 'neutral'
    default: return 'neutral'
  }
}

// Action items for dropdown
const getActionItems = (user: User) => [
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

const columns: any[] = [
  {
    accessorKey: 'name',
    header: t('users.columns.name'),
    cell: ({ row }: any) => {
      const user = row.original
      return h('div', { class: 'flex items-center gap-3' }, [
        h(UAvatar, {
          src: user.avatar || undefined,
          alt: user.name,
          size: 'sm'
        }, () => h('span', { class: 'text-xs font-medium text-primary-600 dark:text-primary-400' }, 
          user.name.split(' ').map((n: string) => n[0]).join('')
        )),
        h('div', [
          h('button', {
            onClick: () => viewUser(user),
            class: 'font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400 transition-colors'
          }, user.name)
        ])
      ])
    }
  },
  {
    accessorKey: 'email',
    header: t('users.columns.email'),
    cell: ({ row }: any) => {
      const user = row.original
      return h('span', { class: 'text-gray-600 dark:text-gray-400' }, user.email)
    }
  },
  {
    accessorKey: 'role',
    header: t('users.columns.role'),
    cell: ({ row }: any) => {
      const user = row.original
      return h(UBadge, {
        color: getRoleColor(user.role),
        variant: 'subtle',
        class: [
          'capitalize',
          user.role === 'admin' ? 'text-red-700 dark:text-red-300' : '',
          user.role === 'moderator' ? 'text-blue-700 dark:text-blue-300' : '',
          user.role === 'user' ? 'text-gray-700 dark:text-gray-300' : ''
        ]
      }, () => t(`users.roles.${user.role}`))
    }
  },
  {
    accessorKey: 'status',
    header: t('users.columns.status'),
    cell: ({ row }: any) => {
      const user = row.original
      return h(UBadge, {
        color: getStatusColor(user.status),
        variant: 'subtle',
        class: [
          'capitalize',
          user.status === 'online' ? 'text-green-700 dark:text-green-300' : '',
          user.status === 'offline' ? 'text-gray-700 dark:text-gray-300' : '',
          user.status === 'banned' ? 'text-red-700 dark:text-red-300' : '',
          user.status === 'suspended' ? 'text-yellow-700 dark:text-yellow-300' : ''
        ]
      }, () => t(`users.status.${user.status}`))
    }
  },
  {
    accessorKey: 'lastSeen',
    header: t('users.columns.lastSeen'),
    cell: ({ row }: any) => {
      const user = row.original
      return h('span', { class: 'text-sm text-gray-500 dark:text-gray-400' }, user.lastSeen)
    }
  },
  {
    accessorKey: 'serversAccess',
    header: t('users.columns.serversAccess'),
    cell: ({ row }: any) => {
      const user = row.original
      return h('div', { class: 'flex items-center gap-2' }, [
        h(UIcon, { name: 'i-heroicons-server-20-solid', class: 'w-4 h-4 text-gray-400' }),
        h('span', { class: 'text-sm text-gray-600 dark:text-gray-400' }, user.serversAccess.toString())
      ])
    }
  },
  {
    id: 'actions',
    header: t('users.columns.actions'),
    cell: ({ row }: any) => {
      const user = row.original
      return h('div', { class: 'flex items-center gap-2' }, [
        h(UButton, {
          color: 'primary',
          variant: 'ghost',
          icon: 'i-heroicons-eye-20-solid',
          size: 'sm',
          class: 'text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-200',
          onClick: () => viewUser(user)
        }),
        h(UDropdownMenu, {
          items: getActionItems(user)
        }, () => h(UButton, {
          color: 'neutral',
          variant: 'ghost',
          icon: 'i-heroicons-ellipsis-horizontal-20-solid',
          size: 'sm',
          class: 'text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200'
        }))
      ])
    }
  }
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

// Navigation functions
const viewUser = (user: User) => {
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
      <UTable :data="filteredUsers" :columns="columns" />

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
          class="text-blue-700 dark:text-blue-300"
        >
          {{ t('users.createUser') }}
        </UButton>
      </div>
    </UCard>
  </div>
</template> 
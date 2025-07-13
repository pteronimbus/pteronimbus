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

// Filter configurations for SearchAndFilters component
const filters = computed(() => [
  {
    key: 'status',
    value: selectedStatus.value,
    options: [
      { value: 'all', label: 'All Status' },
      { value: 'online', label: t('users.status.online') },
      { value: 'offline', label: t('users.status.offline') },
      { value: 'banned', label: t('users.status.banned') },
      { value: 'suspended', label: t('users.status.suspended') }
    ],
    class: 'w-40'
  },
  {
    key: 'role',
    value: selectedRole.value,
    options: [
      { value: 'all', label: 'All Roles' },
      { value: 'admin', label: t('users.roles.admin') },
      { value: 'moderator', label: t('users.roles.moderator') },
      { value: 'user', label: t('users.roles.user') }
    ],
    class: 'w-40'
  }
])

// Page header actions
const headerActions = computed(() => [
  {
    label: t('users.createUser'),
    icon: 'i-heroicons-plus-circle',
    color: 'primary' as const,
    onClick: () => router.push('/users/create')
  }
])

// Resolve components for use in cell renderers
const UAvatar = resolveComponent('UAvatar')
const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const UDropdownMenu = resolveComponent('UDropdownMenu')
const UIcon = resolveComponent('UIcon')
const StatusBadge = resolveComponent('StatusBadge')

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
      return h(StatusBadge, {
        status: user.status,
        type: 'user'
      })
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

// Handle filter updates
const handleFilterUpdate = (key: string, value: string) => {
  if (key === 'status') {
    selectedStatus.value = value
  } else if (key === 'role') {
    selectedRole.value = value
  }
}

// Check if filters are active
const hasActiveFilters = computed(() => {
  return searchQuery.value !== '' || selectedStatus.value !== 'all' || selectedRole.value !== 'all'
})
</script>

<template>
  <div>
    <!-- Page Header -->
    <PageHeader 
      :title="t('users.title')"
      description="Manage users and their permissions"
      :actions="headerActions"
    />

    <!-- Search and Filters -->
    <SearchAndFilters
      v-model:search-query="searchQuery"
      :filters="filters"
      search-placeholder="Search users..."
      @update:filter="handleFilterUpdate"
      class="mb-6"
    />

    <!-- Users Table -->
    <UCard>
      <UTable :data="filteredUsers" :columns="columns" />

      <!-- Empty state -->
      <EmptyState
        v-if="filteredUsers.length === 0"
        icon="i-heroicons-users-20-solid"
        :title="hasActiveFilters ? 'No users found' : t('users.noUsers')"
        :description="hasActiveFilters 
          ? 'Try adjusting your search or filters' 
          : 'Get started by creating your first user'"
        :action-label="!hasActiveFilters ? t('users.createUser') : undefined"
        action-icon="i-heroicons-plus-circle"
        @action="router.push('/users/create')"
      />
    </UCard>
  </div>
</template> 
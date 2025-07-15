# Frontend Authentication Integration

This document describes the authentication changes made to integrate with the Go backend instead of using direct Discord OAuth2.

## Changes Made

### Removed Dependencies

- `@sidebase/nuxt-auth` - Removed direct Discord auth
- `next-auth` - No longer needed

### New Authentication System

#### Composables

- `composables/useAuth.ts` - Main authentication composable with:
  - User state management
  - Login/logout functionality
  - Token refresh mechanism
  - API request wrapper with automatic token refresh

#### Middleware

- `middleware/auth.ts` - Protects authenticated routes
- `middleware/guest.ts` - Redirects authenticated users from login/register pages

#### Pages

- `pages/auth/callback.vue` - Handles Discord OAuth2 callback
- Updated `pages/login.vue` - Uses backend authentication flow
- All protected pages now use `middleware: 'auth'`

#### Plugins

- `plugins/auth.client.ts` - Initializes auth state from localStorage

## Authentication Flow

1. User clicks "Login with Discord" on `/login`
2. Frontend calls backend `/auth/login` endpoint
3. Backend returns Discord OAuth2 authorization URL
4. User is redirected to Discord
5. Discord redirects to `/auth/callback` with authorization code
6. Frontend calls backend `/auth/callback` with the code
7. Backend handles Discord OAuth2 exchange and returns JWT tokens
8. Frontend stores tokens and user data in localStorage
9. User is redirected to dashboard

## Configuration

Add to `.env`:

```bash
BACKEND_URL=http://localhost:8080
```

## Usage

### In Components

```vue
<script setup>
const { user, isAuthenticated, signOut } = useAuth()

// Access current user
console.log(user.value)

// Check authentication status
if (isAuthenticated.value) {
  // User is logged in
}

// Logout
await signOut()
</script>
```

### API Requests

```vue
<script setup>
const { apiRequest } = useAuth()

// Make authenticated API requests
const data = await apiRequest('/api/some-endpoint', {
  method: 'POST',
  body: { data: 'example' }
})
</script>
```

### Route Protection

```vue
<script setup>
definePageMeta({
  middleware: 'auth' // Requires authentication
})
</script>
```

## Token Management

- Access tokens expire in 1 hour
- Refresh tokens expire in 7 days
- Automatic token refresh on API requests
- Tokens stored in localStorage
- Session cleanup on logout or token expiry

## Error Handling

- Invalid/expired tokens trigger automatic refresh
- Failed refresh redirects to login
- OAuth errors shown on callback page
- Network errors handled gracefully
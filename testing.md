# Testing Documentation for Nuxt UI v3 Data Tables

## Context and Implementation Overview

This document covers testing strategies for a Nuxt UI v3 application that was converted from v2 to v3, specifically focusing on data table components that use TanStack Table.

### Key Technical Changes Made

1. **Migration from Nuxt UI v2 to v3:**
   - Changed from `:rows` prop to `:data` prop for tables
   - Replaced template slot-based column definitions with programmatic column definitions
   - Implemented TanStack Table API with `accessorKey`, `header`, and `cell` properties

2. **Vue 3 Composition API Usage:**
   - Extensive use of `h()` function for programmatic component rendering
   - Component resolution using `resolveComponent()` for UI components
   - All components resolved at setup time for use in cell renderers

3. **Fixed App Structure:**
   - Added `<UApp>` wrapper to resolve locale context errors
   - Proper component hierarchy for Nuxt UI v3

## Pages Implemented

### 1. Servers Page (`/servers/index.vue`)
- **Data Structure:** Server objects with id, name, game, status, players, IP, performance metrics
- **Key Features:** 
  - Status badges with color coding
  - Performance indicators (CPU/memory)
  - Action dropdowns with start/stop/restart/delete
  - Live status updates with mock timers
  - Search and filtering by status/game
  - Statistics cards with click-to-filter

### 2. Players Page (`/players/index.vue`)
- **Data Structure:** Player objects with avatar, name, server, status, playtime, last seen
- **Key Features:**
  - Avatar display with fallback initials
  - Server association and status
  - Time-based data (playtime, last seen)
  - Search and filtering capabilities

### 3. Users Page (`/users/index.vue`)
- **Data Structure:** User objects with profile info, roles, permissions, server access
- **Key Features:**
  - Role-based access control display
  - Server access management
  - User status indicators
  - Admin action capabilities

### 4. Alerts Page (`/alerts/index.vue`)
- **Data Structure:** Alert objects with severity, server association, timestamps, status
- **Key Features:**
  - Severity-based color coding
  - Acknowledge/dismiss actions
  - Real-time status updates
  - Server context linking

## Testing Strategy Recommendations

### 1. Component Rendering Tests
Focus on testing that components render correctly with the new TanStack Table structure:

```javascript
// Test that tables render with data
test('server table renders with mock data', () => {
  // Mount component with mock data
  // Assert table structure exists
  // Assert data is displayed correctly
})

// Test empty states
test('displays empty state when no data', () => {
  // Mount with empty data array
  // Assert empty state message and icon display
})
```

### 2. Interactive Element Tests
Test the programmatically rendered components work correctly:

```javascript
// Test dropdown menus
test('action dropdown renders and functions', () => {
  // Find dropdown trigger button
  // Click to open dropdown
  // Assert menu items are present
  // Test click handlers
})

// Test status badges
test('status badges display correct colors', () => {
  // Test each status value
  // Assert correct color classes applied
})
```

### 3. Search and Filtering Tests
Test computed properties and reactivity:

```javascript
// Test search functionality
test('search filters data correctly', () => {
  // Set search query
  // Assert filtered results match expectations
})

// Test multi-field filtering
test('status and game filters work together', () => {
  // Set multiple filter values
  // Assert compound filtering works
})
```

### 4. Data Transformation Tests
Test helper functions and computed properties:

```javascript
// Test status color mapping
test('getStatusColor returns correct colors', () => {
  expect(getStatusColor('online')).toBe('success')
  expect(getStatusColor('error')).toBe('error')
})

// Test stats calculations
test('serverStats computed property calculates correctly', () => {
  // Test with known data set
  // Assert correct counts for each status
})
```

## Mock Data Structure

### Server Mock Data
```typescript
interface Server {
  id: number
  name: string
  game: string
  status: 'online' | 'offline' | 'starting' | 'stopping' | 'error'
  players: string // format: "current/max"
  ip: string
  port: number
  version: string
  uptime: string
  cpu: number // percentage
  memory: number // percentage
  createdAt: string
}
```

### Key Mock Data Patterns
- **Status values:** 'online', 'offline', 'starting', 'stopping', 'error'
- **Performance metrics:** CPU and memory as percentages (0-100)
- **Players format:** "current/max" (e.g., "10/20")
- **Time formats:** Human-readable strings (e.g., "2d 14h 30m")

## Critical Testing Areas

### 1. Component Resolution
- Ensure `resolveComponent()` calls work in test environment
- Test that UI components (UIcon, UBadge, UButton, etc.) render correctly in cells

### 2. Cell Renderer Functions
- Test that `h()` function calls create expected DOM structure
- Verify event handlers attached to programmatically created elements work
- Test complex cell content (nested components, conditional rendering)

### 3. Reactivity and State Updates
- Test that search query updates filter results immediately
- Test that status changes trigger UI updates
- Test that mock action functions (start/stop server) update state correctly

### 4. Accessibility and UX
- Test keyboard navigation through tables
- Test screen reader compatibility with programmatically rendered content
- Test responsive behavior of tables

## Testing Gotchas and Considerations

### 1. Component Resolution in Tests
- UI components need to be properly mocked or available in test environment
- `resolveComponent()` may need special handling in testing framework

### 2. Mock Timer Functions
- Several components use `setTimeout` for status transitions
- Tests may need to handle or mock these async operations

### 3. Router Integration
- Many actions involve router navigation
- Tests should mock router or test navigation behavior

### 4. i18n Integration
- All text uses translation keys
- Tests should either mock i18n or provide test translations

## Recommended Testing Libraries

- **Vue Test Utils:** For component mounting and interaction
- **Vitest:** For test runner (compatible with Vite/Nuxt)
- **Testing Library:** For user-centric testing approaches
- **MSW (Mock Service Worker):** For API mocking when real endpoints are implemented

## Future Testing Considerations

### When Real APIs Are Implemented
- Replace mock data with API integration tests
- Test error states and loading states
- Test real-time updates and WebSocket connections
- Add E2E tests for complete user workflows

### Performance Testing
- Test table performance with large datasets
- Test filtering performance with many records
- Test memory usage with dynamic content updates

### Security Testing
- Test that action permissions are properly enforced
- Test that sensitive data is properly protected
- Test that user roles correctly limit available actions

## Test Data Management

### Static Test Data
Use the existing mock data structure but with predictable values for testing:

```typescript
const testServers = [
  {
    id: 1,
    name: 'Test Server 1',
    game: 'Minecraft',
    status: 'online',
    players: '5/10',
    // ... other properties
  },
  // ... more test data
]
```

### Dynamic Test Scenarios
Create test data builders for different scenarios:
- All servers online
- All servers offline
- Mixed status scenarios
- Empty data sets
- Large data sets for performance testing

This testing approach ensures the TanStack Table implementation is robust and maintainable while providing comprehensive coverage of the interactive features and data transformations. 
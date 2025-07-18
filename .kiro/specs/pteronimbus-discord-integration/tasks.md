# Implementation Tasks

## Phase 1: Authentication & Core Infrastructure

- [x] 1. Complete Discord OAuth2 Backend Integration
  - Implement Discord OAuth2 flow in Go backend (replace frontend direct Discord auth)
  - Create backend authentication endpoints (/auth/login, /auth/callback, /auth/refresh)
  - Implement JWT token generation and validation middleware
  - Add refresh token mechanism with secure storage in Redis
  - Create authentication middleware for API routes
  - Update frontend to authenticate via Go backend instead of direct Discord
  - Write unit and integration tests
  - _Requirements: 1.1, 1.2, 1.4, 1.5_

- [x] 2. Refactor Frontend Authentication Architecture
  - Remove direct Discord OAuth2 integration from frontend (@sidebase/nuxt-auth)
  - Implement custom authentication composables for Go backend integration
  - Create login flow that redirects to Go backend auth endpoints
  - Update authentication middleware to work with backend JWT tokens
  - Modify session management to use backend-issued tokens
  - Write unit and integration tests
  - _Requirements: 1.1, 1.2, 1.5_

- [x] 3. Implement Multi-Tenant Infrastructure







  - Create database schema for tenants, users, and Discord roles
  - Implement Discord server discovery via Discord API
  - Create tenant creation and configuration endpoints
  - Add tenant context middleware for API requests
  - Build tenant selection UI component in frontend
  - Add tenant context management to frontend (store, routing, API calls)
  - Update existing frontend pages to be tenant-aware
  - Modify dashboard and all components to work within tenant context
  - Write unit and integration tests
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.6_

- [x] 4. Fix Multi-Tenant Test Suite



  - Fix failing dashboard tests that were testing wrong component behavior
  - Test the moved dashboard functionaltiy under tenant/xxx/dashboard
  - _Requirements: Testing Requirements 1, 3_

- [ ] 5. Set up Kubernetes Controller Foundation

  - Implement GameTemplate and GameServer CRD definitions
  - Create controller reconciliation loop structure using controller-runtime
  - Implement secure communication between backend and controller
  - Add controller health checks and status reporting
  - Create backend API endpoints for controller communication
  - Write unit and integration tests
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.7, 5.8_

## Phase 2: RBAC & Permission System

- [ ] 6. Implement Discord Role-Based Access Control
  - Create Discord role synchronization service
  - Implement permission mapping system between Discord roles and Pteronimbus permissions
  - Add RBAC enforcement middleware for all API endpoints
  - Create role management UI components for authorized users
  - Implement super admin functionality override system
  - Write unit and integration tests
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_

- [ ] 7. Build Fine-Grained Permission System
  - Implement CRUD permission checks for game servers
  - Add permission validation for server control operations
  - Create read-only access controls for logs and metrics
  - Implement user management permissions within tenants
  - Add permission enforcement at controller level
  - Write unit and integration tests
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7_

## Phase 3: Game Server Management

- [ ] 8. Implement Game Server Templates
  - Create GameTemplate CRD handling in controller
  - Implement Pterodactyl egg import and conversion functionality
  - Build template management API endpoints
  - Create template management UI components
  - Add template versioning and rollback system
  - Write unit and integration tests
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6_

- [ ] 9. Build Game Server Lifecycle Management
  - Implement GameServer CRD creation and management
  - Create server deployment logic in Kubernetes controller
  - Build server management API endpoints (start, stop, restart, delete)
  - Implement server status monitoring and reporting
  - Create server management UI with real-time status updates
  - Write unit and integration tests
  - _Requirements: 4.1, 4.2, 4.3, 4.5, 4.6, 4.7_

## Phase 4: Discord Bot Integration

- [ ] 10. Develop Discord Bot Service
  - Set up Discord bot application and register slash commands
  - Implement bot command handlers with RBAC integration
  - Create server management commands (start, stop, restart, status)
  - Add real-time notifications to Discord channels
  - Implement audit logging for all bot actions
  - Write unit and integration tests
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_

## Phase 5: Real-Time Features & Polish

- [ ] 11. Implement Real-Time Updates
  - Add WebSocket connections for live server status updates
  - Create event system for server state changes
  - Implement real-time log streaming functionality
  - Update frontend components to handle live data
  - Add notification system for important events
  - Write unit and integration tests
  - _Requirements: 4.7, 7.3, 7.4_

- [ ] 12. Add Advanced Game Server Features
  - Implement server resource usage monitoring
  - Create backup and restore functionality for game data
  - Add scheduled task system (automatic restarts, backups)
  - Implement server performance analytics and metrics
  - Create custom configuration options for different game types
  - Write unit and integration tests
  - _Requirements: 4.1, 4.2, 8.4_

## Testing Requirements

Each implementation task must include:

- [ ] 13. Backend API Testing
  - Unit tests for all service functions (80% coverage minimum)
  - Integration tests for Discord API interactions
  - API endpoint testing with authentication scenarios
  - Database operation testing with test containers
  - Mock Discord API responses for testing
  - Write unit and integration tests
  - _Requirements: Testing Requirements 1, 2, 6_

- [ ] 14. Frontend Component Testing
  - Component tests for all UI components (70% coverage minimum)
  - Integration tests for authentication flows
  - End-to-end tests for critical user workflows
  - Mock API responses for frontend testing
  - Accessibility testing for all components
  - Write unit and integration tests
  - _Requirements: Testing Requirements 1, 3, 7_

- [ ] 15. Controller Testing
  - Unit tests for reconciliation logic (80% coverage minimum)
  - Integration tests with test Kubernetes cluster
  - CRD validation and status update testing
  - Error handling and retry logic testing
  - Resource cleanup testing
  - Write unit and integration tests
  - _Requirements: Testing Requirements 1, 4, 9_

## Deployment & Infrastructure
3.1
- [ ] 16. Production Deployment Setup
  - Create Helm charts for all components
  - Implement CI/CD pipeline with automated testing
  - Set up monitoring and observability stack
  - Configure production database and Redis
  - Implement backup and disaster recovery procedures
  - Write unit and integration tests
  - _Requirements: 5.6_

## Notes

- Each task should be implemented incrementally with comprehensive testing, not thorugh scripts, but in the language and projects that they test.
- All Discord API interactions must handle rate limiting and errors gracefully
- Security considerations must be addressed at every level
- Performance testing should be conducted for multi-tenant scenarios
- Documentation should be updated as features are implemented
- Any changes to the frontend must respect the dark and light themes
- Any changes to the frontend should be done to the quality of a staff ui developer, shortcuts should be avoided
- Any changes to the frontend should respect i18n
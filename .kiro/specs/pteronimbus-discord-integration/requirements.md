# Requirements Document

## Introduction

Pteronimbus is a multi-tenant, Kubernetes-native game server management platform with deep Discord integration. The system enables Discord server owners to deploy, manage, and monitor game servers through both a web interface and Discord bot commands. Unlike traditional game server panels, Pteronimbus leverages Discord's existing role system for authentication and authorization, creating a seamless experience where Discord servers become tenants with their own isolated game server environments.

The platform follows a pull-based architecture where a backend API stores user intent and a Kubernetes controller reconciles the desired state, ensuring secure isolation between tenants while providing fine-grained permission controls.

## Requirements

### Requirement 1: Discord OAuth2 Authentication

**User Story:** As a Discord user, I want to log in to Pteronimbus using my Discord account, so that I can access game server management without creating separate credentials.

#### Acceptance Criteria

1. WHEN a user visits the login page THEN the system SHALL present a Discord OAuth2 login option
2. WHEN a user completes Discord OAuth2 flow THEN the system SHALL issue a JWT token with 1-hour expiration
3. WHEN a JWT token expires THEN the system SHALL provide refresh token functionality following OAuth2 standards
4. IF a user's Discord account is deactivated THEN the system SHALL revoke access to all tenants
5. WHEN a user logs in THEN the system SHALL store and manage their Discord tokens securely in the backend

### Requirement 2: Multi-Tenant Discord Server Management

**User Story:** As a Discord server owner, I want to install Pteronimbus to my Discord server, so that my community can manage game servers through our existing Discord infrastructure.

#### Acceptance Criteria

1. WHEN a user logs in THEN the system SHALL display a list of Discord servers (tenants) they have access to
2. IF a user has 'Manage Server' permissions on a Discord server THEN the system SHALL offer the option to install Pteronimbus to that server
3. WHEN Pteronimbus is installed to a Discord server THEN the system SHALL create a new tenant with that Discord server as the identifier
4. WHEN a user accesses an existing tenant THEN the system SHALL automatically map their Discord roles to Pteronimbus permissions
5. WHEN a Discord bot is installed THEN the system SHALL enable Discord command functionality for game server management
6. IF a user joins a Discord server with Pteronimbus already installed THEN the system SHALL automatically grant appropriate access based on their Discord roles

### Requirement 3: Discord Role-Based Access Control (RBAC)

**User Story:** As a Discord server administrator, I want to manage Pteronimbus access control through Discord roles while maintaining granular permission control within Pteronimbus, so that I can leverage my existing Discord structure while having fine-grained game server management.

#### Acceptance Criteria

1. WHEN Pteronimbus is installed to a Discord server THEN the system SHALL pull all users and roles from that Discord server
2. WHEN Discord data is pulled THEN the system SHALL create a representation of users and roles in the Pteronimbus database
3. IF a user has 'Manage Server' or 'Manage Roles' Discord permissions THEN the system SHALL allow them to manage RBAC within Pteronimbus
4. WHEN managing RBAC THEN authorized users SHALL be able to correlate Discord roles with Pteronimbus roles
5. WHEN Discord roles change THEN the system SHALL update user permissions based on Discord role membership rather than internal Pteronimbus assignments
6. WHEN users are added or removed from Discord roles THEN the system SHALL automatically update their Pteronimbus permissions within 5 minutes
7. IF a server host exists THEN the system SHALL provide super admin permissions that override all tenant-level restrictions

### Requirement 4: Game Server Lifecycle Management

**User Story:** As a game server administrator, I want to create, configure, and manage game servers through both web interface and Discord commands, so that I can efficiently operate game servers for my community.

#### Acceptance Criteria

1. WHEN creating a game server THEN the system SHALL support multiple game types including Minecraft, CS2, Valheim, and Terraria
2. WHEN selecting a game type THEN the system SHALL provide templates based on GameTemplate CRDs
3. WHEN importing Pterodactyl eggs THEN the system SHALL convert them to GameTemplate CRD format
4. WHEN managing servers through Discord THEN the system SHALL support start, stop, and restart commands via bot
5. WHEN a game server is created THEN the system SHALL deploy it as Kubernetes resources managed by the controller
6. IF a user has appropriate permissions THEN the system SHALL allow server configuration updates
7. WHEN server state changes THEN the system SHALL provide real-time status updates

### Requirement 5: Kubernetes-Native Architecture

**User Story:** As a platform operator, I want Pteronimbus to leverage Kubernetes-native patterns, so that game servers are managed consistently with cloud-native best practices.

#### Acceptance Criteria

1. WHEN game servers are created THEN the system SHALL use Custom Resource Definitions (CRDs) to model them
2. WHEN the controller runs THEN it SHALL reach out to the backend API to determine desired state
3. WHEN desired state differs from actual state THEN the controller SHALL reconcile Kubernetes resources accordingly
4. WHEN users interact with the system THEN the web UI and API SHALL never directly interact with Kubernetes
5. WHEN the system is deployed as SaaS THEN users SHALL NOT have direct cluster access
7. WHEN CRDs are deployed THEN the system SHALL support GameTemplate CRDs for defining game server templates
8. WHEN CRDs are deployed THEN the system SHALL support GameServer CRDs that reference GameTemplate CRDs to drive their configuration

### Requirement 6: Fine-Grained Permission System

**User Story:** As a tenant administrator, I want to control exactly what each Discord role can do with game servers, so that I can maintain appropriate access controls for my community.

#### Acceptance Criteria

1. WHEN defining permissions THEN the system SHALL support CRUD operations for game servers (create, read, update, delete)
2. WHEN defining permissions THEN the system SHALL support CRUD operations for server configurations
3. WHEN defining permissions THEN the system SHALL support read-only access to server logs
4. WHEN defining permissions THEN the system SHALL support server control operations (start, stop, restart)
5. WHEN defining permissions THEN the system SHALL support user management within the tenant
6. IF a resource is read-only by nature THEN the system SHALL only provide read permissions (e.g., logs, metrics)
7. WHEN permissions are checked THEN the system SHALL enforce them at both API and controller levels

### Requirement 7: Discord Bot Integration

**User Story:** As a Discord server member, I want to manage game servers directly through Discord commands, so that I can perform common operations without leaving Discord.

#### Acceptance Criteria

1. WHEN the bot is installed THEN it SHALL provide slash commands for server management
2. WHEN using Discord commands THEN the system SHALL respect the same RBAC as the web interface
3. WHEN server status changes THEN the bot SHALL provide notifications to configured Discord channels
4. WHEN executing commands THEN the bot SHALL provide immediate feedback and status updates
5. IF a command requires elevated permissions THEN the bot SHALL provide clear error messages
6. WHEN commands are executed THEN the system SHALL log all actions for audit purposes

### Requirement 8: Multi-Game Support and Templating

**User Story:** As a game server administrator, I want to deploy different types of game servers using standardized templates, so that I can quickly set up servers for various games my community plays.

#### Acceptance Criteria

1. WHEN selecting a game type THEN the system SHALL provide pre-configured GameTemplate CRDs for Minecraft, CS2, Valheim, and Terraria
2. WHEN importing Pterodactyl eggs THEN the system SHALL automatically convert them to GameTemplate CRD format
3. WHEN creating custom templates THEN the system SHALL allow saving them as reusable GameTemplate CRDs
4. WHEN deploying from templates THEN the system SHALL allow parameter customization before creating GameServer CRDs
5. IF a GameTemplate is updated THEN the system SHALL provide options to update existing GameServer instances or keep current configuration
6. WHEN GameTemplates are managed THEN the system SHALL support versioning and rollback capabilities
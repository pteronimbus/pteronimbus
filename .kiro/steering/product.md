# Product Overview

Pteronimbus is a Kubernetes-native game server hosting platform inspired by Pterodactyl. It enables self-service game server creation and management through a modern web interface while maintaining security and infrastructure privacy.

## Core Purpose
- Self-hosted game server management for homelabbers and hobbyists
- Kubernetes-native architecture using Custom Resource Definitions (CRDs)
- Fine-grained RBAC authorization with OIDC authentication
- Power-user support via kubectl for GameServer resources

## Key Components
- **Frontend**: Modern Nuxt-based web interface for server management
- **Backend**: Go API server handling authentication, user management, and manifest rendering
- **Controller**: Kubernetes controller that watches CRDs and manages game server pods
- **CRDs**: Custom Kubernetes resources defining desired game server state

## Target Users
- Homelabbers running private Kubernetes clusters
- Hobbyists wanting self-hosted game server management
- Cloud-native enthusiasts preferring Kubernetes workflows
- Power users comfortable with kubectl for advanced operations

## Authentication & Security
- OIDC-compatible authentication (designed for Authentik)
- Discord OAuth2 integration for development/testing
- RBAC ensures users only perform authorized operations
- Separation between hosted (restricted kubectl) and self-hosted (trusted admin) environments
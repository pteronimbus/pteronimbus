# 🦖 Pteronimbus

**Pteronimbus** is a Kubernetes-native game server hosting platform inspired by Pterodactyl, designed for homelabbers, hobbyists, and cloud-native workflows. It lets users create and manage game servers securely while keeping infrastructure private.

---

## 🚀 Features

- Self-service game server creation through a modern web UI
- Fine-grained RBAC authorization for all operations
- Authentication via Authentik (OIDC-compatible)
- Kubernetes-native with Custom Resource Definitions (CRDs)
- Controller watches CRDs and manages game server pods
- Power-user support via `kubectl` for `GameServer` resources

---

## 🧩 Components

- **Frontend:** Nuxt-based web interface
- **Backend:** Go API server handling user, auth, and manifest rendering
- **Controller:** Kubernetes controller managing game servers via CRDs
- **CRDs:** Kubernetes resources defining the desired state of game servers

---

## ⚙️ Quick Start

### Prerequisites

- Kubernetes cluster (k3s, kind, or cloud provider)
- Helm 3.x
- `kubectl`
- Authentik (or compatible OIDC provider)

### Deploy

```bash
helm install pteronimbus charts/pteronimbus/
```

Access the frontend, authenticate via Authentik, and create your first game server.

---

## 🎮 Managing Game Servers

- Create new servers via the web UI
- Each server corresponds to a `GameServer` CRD managed by the controller
- Monitor status via web or:

```bash
kubectl get gameservers
kubectl describe gameserver <name>
```

---

## 🛠️ Advanced Usage

- Manually edit `GameServer` CRDs with `kubectl edit`
- Extend or customize game server types via CRDs
- Use the API directly for integration or automation

---

## 🔒 Security & Permissions

- RBAC ensures users only perform allowed operations
- Hosted environments restrict direct `kubectl` access
- Self-hosted environments expect trusted admin users

---

## 📄 License

Core platform is MIT licensed. Additional pro features may be commercial.

---

## 📞 Support & Community

Join our community channels or open GitHub issues for help and feature requests.

export default defineI18nConfig(() => ({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      welcome: 'Welcome',
      nav: {
        dashboard: 'Dashboard',
        servers: 'Servers',
        users: 'Users',
        settings: 'Settings'
      },
      dashboard: {
        welcome: 'Welcome, {name}',
        overview: 'Here\'s a quick overview of your game server environment.',
        stats: {
          activeServers: 'Active Servers',
          totalPlayers: 'Total Players',
          cpuUsage: 'CPU Usage',
          memoryUsage: 'Memory Usage',
          diskUsage: 'Disk Usage',
          networkIO: 'Network I/O',
          uptime: 'Uptime',
          totalUsers: 'Total Users',
          onlineUsers: 'Online Users',
          alertsActive: 'Active Alerts',
          recentActivity: 'Recent Activity',
          resourceMonitoring: 'Resource Monitoring'
        },
        activity: {
          title: 'Recent Activity',
          noActivity: 'No recent activity',
          serverStarted: 'Server {name} started',
          serverStopped: 'Server {name} stopped',
          userJoined: 'User {name} joined {server}',
          userLeft: 'User {name} left {server}',
          serverCreated: 'Server {name} created',
          userBanned: 'User {name} banned from {server}'
        },
        alerts: {
          title: 'System Alerts',
          noAlerts: 'No active alerts',
          highCpu: 'High CPU usage detected',
          highMemory: 'High memory usage detected',
          serverDown: 'Server {name} is down',
          diskSpace: 'Low disk space warning'
        }
      },
      servers: {
        title: 'Servers',
        createServer: 'Create Server',
        noServers: 'No servers available',
        status: {
          online: 'Online',
          offline: 'Offline',
          starting: 'Starting',
          stopping: 'Stopping',
          error: 'Error'
        },
        columns: {
          name: 'Name',
          game: 'Game',
          status: 'Status',
          players: 'Players',
          ip: 'IP Address',
          uptime: 'Uptime',
          actions: 'Actions'
        },
        actions: {
          start: 'Start',
          stop: 'Stop',
          restart: 'Restart',
          edit: 'Edit',
          delete: 'Delete',
          console: 'Console',
          viewDetails: 'View Details',
          managePlayers: 'Manage Players'
        },
        details: {
          title: 'Server Details',
          console: 'Console',
          logs: 'Logs',
          files: 'Files',
          settings: 'Settings',
          players: 'Players',
          performance: 'Performance',
          backups: 'Backups'
        }
      },
      users: {
        title: 'Users',
        createUser: 'Create User',
        noUsers: 'No users available',
        status: {
          online: 'Online',
          offline: 'Offline',
          banned: 'Banned',
          suspended: 'Suspended'
        },
        columns: {
          name: 'Name',
          email: 'Email',
          role: 'Role',
          status: 'Status',
          lastSeen: 'Last Seen',
          serversAccess: 'Servers Access',
          actions: 'Actions'
        },
        actions: {
          edit: 'Edit',
          ban: 'Ban',
          unban: 'Unban',
          suspend: 'Suspend',
          delete: 'Delete',
          viewDetails: 'View Details',
          resetPassword: 'Reset Password',
          changeRole: 'Change Role'
        },
        roles: {
          admin: 'Admin',
          moderator: 'Moderator',
          user: 'User'
        },
        details: {
          title: 'User Details',
          profile: 'Profile',
          permissions: 'Permissions',
          activity: 'Activity',
          sessions: 'Sessions'
        }
      },
      common: {
        actions: 'Actions',
        save: 'Save',
        cancel: 'Cancel',
        delete: 'Delete',
        edit: 'Edit',
        create: 'Create',
        update: 'Update',
        confirm: 'Confirm',
        loading: 'Loading...',
        error: 'Error',
        success: 'Success',
        warning: 'Warning',
        info: 'Info',
        viewAll: 'View All',
        refresh: 'Refresh',
        search: 'Search',
        filter: 'Filter',
        export: 'Export',
        import: 'Import'
      }
    },
    fr: {
      welcome: 'Bienvenue',
      nav: {
        dashboard: 'Tableau de bord',
        servers: 'Serveurs',
        users: 'Utilisateurs',
        settings: 'Paramètres'
      },
      dashboard: {
        welcome: 'Bienvenue, {name}!',
        overview: 'Voici un aperçu rapide de votre environnement de serveur de jeu.',
        stats: {
          activeServers: 'Serveurs actifs',
          totalPlayers: 'Joueurs totaux',
          cpuUsage: 'Utilisation CPU',
          memoryUsage: 'Utilisation mémoire',
          diskUsage: 'Utilisation disque',
          networkIO: 'E/S réseau',
          uptime: 'Temps de fonctionnement',
          totalUsers: 'Utilisateurs totaux',
          onlineUsers: 'Utilisateurs en ligne',
          alertsActive: 'Alertes actives',
          recentActivity: 'Activité récente',
          resourceMonitoring: 'Surveillance des ressources'
        }
      },
      servers: {
        title: 'Serveurs',
        createServer: 'Créer un serveur',
        status: {
          online: 'En ligne',
          offline: 'Hors ligne',
          starting: 'Démarrage',
          stopping: 'Arrêt',
          error: 'Erreur'
        }
      },
      users: {
        title: 'Utilisateurs',
        createUser: 'Créer un utilisateur'
      },
      common: {
        actions: 'Actions',
        save: 'Sauvegarder',
        cancel: 'Annuler',
        delete: 'Supprimer',
        edit: 'Modifier',
        create: 'Créer',
        update: 'Mettre à jour',
        confirm: 'Confirmer',
        loading: 'Chargement...',
        viewAll: 'Voir tout',
        refresh: 'Actualiser'
      }
    }
  }
})) 
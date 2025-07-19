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
        quickActions: 'Quick Actions',
        gameServers: 'Game Servers',
        discordIntegration: 'Discord Integration',
        manageRoles: 'Manage Roles',
        viewLogs: 'View Logs',
        settings: 'Settings',
        sync: 'Sync',
        noGameServers: 'No Game Servers',
        noGameServersDesc: 'Create your first game server to get started',
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
          resourceMonitoring: 'Resource Monitoring',
          gameServers: 'Game Servers',
          discordMembers: 'Discord Members'
        },
        activity: {
          title: 'Recent Activity',
          noActivity: 'No recent activity',
          noActivityDesc: 'Activity will appear here once you start using the platform',
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
        create: {
          title: 'Create Server',
          description: 'Create a new game server.',
          comingSoon: 'Server Creation Coming Soon',
          comingSoonDesc: 'The ability to create new game servers will be available here shortly.'
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
        },
        modals: {
          add: {
            title: 'Create New Server',
            create: 'Create Server',
            fields: {
              name: 'Server Name',
              game: 'Game Type',
              maxPlayers: 'Max Players',
              port: 'Port',
              description: 'Description',
              autoStart: 'Auto-start server'
            },
            placeholders: {
              name: 'Enter server name',
              game: 'Select a game',
              maxPlayers: '20',
              port: '25565',
              description: 'Optional server description'
            },
            validation: {
              nameRequired: 'Server name is required',
              gameRequired: 'Please select a game type'
            },
            errors: {
              createFailed: 'Failed to create server'
            }
          },
          delete: {
            title: 'Delete Server',
            confirm: 'Delete Server',
            confirmMessage: 'Are you sure you want to delete "{name}"?',
            warningMessage: 'This will permanently delete the server, all its data, configurations, and backups. This action cannot be undone.',
            confirmationLabel: 'Type "{name}" to confirm deletion',
            validation: {
              confirmationRequired: 'Please enter the server name to confirm',
              confirmationMismatch: 'Server name does not match'
            },
            errors: {
              deleteFailed: 'Failed to delete server'
            }
          }
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
      roles: {
        title: 'Role Management',
        description: 'Manage roles and permissions for your tenant.',
        comingSoon: 'Role Management Coming Soon',
        comingSoonDesc: 'The ability to manage roles and permissions will be available here shortly.'
      },
      logs: {
        title: 'Audit Logs',
        description: 'Review all activity that has occurred within your tenant.',
        comingSoon: 'Audit Logs Coming Soon',
        comingSoonDesc: 'The ability to view audit logs will be available here shortly.'
      },
      settings: {
        title: 'Tenant Settings',
        description: 'Manage your tenant settings and preferences.',
        comingSoon: 'Tenant Settings Coming Soon',
        comingSoonDesc: 'The ability to manage tenant settings will be available here shortly.'
      },
      admin: {
        title: 'Admin Panel',
        controllers: {
          title: 'Controller Management',
          description: 'Monitor and manage Kubernetes controllers across all clusters',
          totalControllers: 'Total Controllers',
          online: 'Online',
          offline: 'Offline',
          errors: 'Errors',
          refreshControllers: 'Refresh Controllers',
          cleanupInactive: 'Cleanup Inactive',
          viewDetails: 'View Details',
          viewLogs: 'View Logs',
          restartController: 'Restart Controller',
          removeController: 'Remove Controller',
          noControllers: 'No controllers registered',
          noControllersDesc: 'Controllers will appear here once they register with the backend',
          clusterName: 'Cluster Name',
          status: 'Status',
          version: 'Version',
          uptime: 'Uptime',
          lastHeartbeat: 'Last Heartbeat',
          actions: 'Actions'
        }
      },
      tenants: {
        title: 'Discord Servers',
        addServer: 'Add Server',
        noServers: 'No servers available',
        loading: 'Loading your Discord servers...',
        modals: {
          add: {
            title: 'Add Discord Server',
            description: 'Select a Discord server where you have "Manage Server" permissions. This will allow Pteronimbus to integrate with your Discord server for game server management.',
            infoTitle: 'What happens when you add a server?',
            infoItems: [
              'Your Discord roles will be synced for permission management',
              'You can manage game servers through Discord commands',
              'Server notifications will be sent to Discord channels'
            ],
            availableServers: 'Available Servers',
            loadingServers: 'Loading your Discord servers...',
            noAvailableServers: 'No Available Servers',
            noServersDescription: 'You need "Manage Server" permissions to add a Discord server to Pteronimbus. Make sure you\'re an administrator or have the required permissions.',
            refreshList: 'Refresh List',
            addButton: 'Add Server',
            owner: 'Owner',
            manager: 'Manager',
            successTitle: 'Server Added',
            successDescription: 'The Discord server {serverName} was added. Next, invite the bot to your server to enable full functionality.',
            errors: {
              loadFailed: 'Failed to load available guilds',
              addFailed: 'Failed to add Discord server'
            }
          },
          delete: {
            title: 'Remove Server',
            confirmMessage: 'Are you sure you want to remove "{name}"?',
            warningMessage: 'This will permanently delete all game servers, configurations, and data associated with this Discord server. This action cannot be undone.',
            confirmButton: 'Remove Server',
            errors: {
              deleteFailed: 'Failed to remove server'
            }
          }
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
        import: 'Import',
        back: 'Back'
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
        quickActions: 'Actions rapides',
        gameServers: 'Serveurs de jeu',
        discordIntegration: 'Intégration Discord',
        manageRoles: 'Gérer les rôles',
        viewLogs: 'Voir les journaux',
        settings: 'Paramètres',
        sync: 'Synchroniser',
        noGameServers: 'Aucun serveur de jeu',
        noGameServersDesc: 'Créez votre premier serveur de jeu pour commencer',
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
          resourceMonitoring: 'Surveillance des ressources',
          gameServers: 'Serveurs de jeu',
          discordMembers: 'Membres Discord'
        },
        activity: {
          title: 'Activité récente',
          noActivity: 'Aucune activité récente',
          noActivityDesc: 'L\'activité apparaîtra ici une fois que vous commencerez à utiliser la plateforme'
        }
      },
      servers: {
        title: 'Serveurs',
        createServer: 'Créer un serveur',
        noServers: 'Aucun serveur disponible',
        status: {
          online: 'En ligne',
          offline: 'Hors ligne',
          starting: 'Démarrage',
          stopping: 'Arrêt',
          error: 'Erreur'
        },
        create: {
          title: 'Créer un serveur',
          description: 'Créez un nouveau serveur de jeu.',
          comingSoon: 'Création de serveur bientôt disponible',
          comingSoonDesc: 'La possibilité de créer de nouveaux serveurs de jeu sera bientôt disponible ici.'
        },
        columns: {
          name: 'Nom',
          game: 'Jeu',
          status: 'Statut',
          players: 'Joueurs',
          ip: 'Adresse IP',
          uptime: 'Temps de fonctionnement',
          actions: 'Actions'
        },
        actions: {
          start: 'Démarrer',
          stop: 'Arrêter',
          restart: 'Redémarrer',
          edit: 'Modifier',
          delete: 'Supprimer',
          console: 'Console',
          viewDetails: 'Voir les détails',
          managePlayers: 'Gérer les joueurs'
        },
        modals: {
          add: {
            title: 'Créer un nouveau serveur',
            create: 'Créer le serveur',
            fields: {
              name: 'Nom du serveur',
              game: 'Type de jeu',
              maxPlayers: 'Joueurs max',
              port: 'Port',
              description: 'Description',
              autoStart: 'Démarrage automatique'
            },
            placeholders: {
              name: 'Entrez le nom du serveur',
              game: 'Sélectionnez un jeu',
              maxPlayers: '20',
              port: '25565',
              description: 'Description optionnelle du serveur'
            },
            validation: {
              nameRequired: 'Le nom du serveur est requis',
              gameRequired: 'Veuillez sélectionner un type de jeu'
            },
            errors: {
              createFailed: 'Échec de la création du serveur'
            }
          },
          delete: {
            title: 'Supprimer le serveur',
            confirm: 'Supprimer le serveur',
            confirmMessage: 'Êtes-vous sûr de vouloir supprimer "{name}" ?',
            warningMessage: 'Cela supprimera définitivement le serveur, toutes ses données, configurations et sauvegardes. Cette action ne peut pas être annulée.',
            confirmationLabel: 'Tapez "{name}" pour confirmer la suppression',
            validation: {
              confirmationRequired: 'Veuillez entrer le nom du serveur pour confirmer',
              confirmationMismatch: 'Le nom du serveur ne correspond pas'
            },
            errors: {
              deleteFailed: 'Échec de la suppression du serveur'
            }
          }
        }
      },
      users: {
        title: 'Utilisateurs',
        createUser: 'Créer un utilisateur'
      },
      tenants: {
        title: 'Serveurs Discord',
        addServer: 'Ajouter un serveur',
        noServers: 'Aucun serveur disponible',
        loading: 'Chargement de vos serveurs Discord...',
        modals: {
          add: {
            title: 'Ajouter un serveur Discord',
            description: 'Sélectionnez un serveur Discord où vous avez les permissions "Gérer le serveur". Cela permettra à Pteronimbus de s\'intégrer avec votre serveur Discord pour la gestion des serveurs de jeu.',
            infoTitle: 'Que se passe-t-il quand vous ajoutez un serveur ?',
            infoItems: [
              'Vos rôles Discord seront synchronisés pour la gestion des permissions',
              'Vous pourrez gérer les serveurs de jeu via les commandes Discord',
              'Les notifications du serveur seront envoyées aux canaux Discord'
            ],
            availableServers: 'Serveurs disponibles',
            loadingServers: 'Chargement de vos serveurs Discord...',
            noAvailableServers: 'Aucun serveur disponible',
            noServersDescription: 'Vous avez besoin des permissions "Gérer le serveur" pour ajouter un serveur Discord à Pteronimbus. Assurez-vous d\'être administrateur ou d\'avoir les permissions requises.',
            refreshList: 'Actualiser la liste',
            addButton: 'Ajouter le serveur',
            owner: 'Propriétaire',
            manager: 'Gestionnaire',
            successTitle: 'Serveur Ajouté',
            successDescription: 'Le serveur Discord {serverName} a été ajouté. Ensuite, invitez le bot sur votre serveur pour activer la fonctionnalité complète.',
            errors: {
              loadFailed: 'Échec du chargement des guildes disponibles',
              addFailed: 'Échec de l\'ajout du serveur Discord'
            }
          },
          delete: {
            title: 'Supprimer le serveur',
            confirmMessage: 'Êtes-vous sûr de vouloir supprimer "{name}" ?',
            warningMessage: 'Cela supprimera définitivement tous les serveurs de jeu, configurations et données associées à ce serveur Discord. Cette action ne peut pas être annulée.',
            confirmButton: 'Supprimer le serveur',
            errors: {
              deleteFailed: 'Échec de la suppression du serveur'
            }
          }
        }
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
        refresh: 'Actualiser',
        back: 'Retour'
      }
    }
  }
})) 
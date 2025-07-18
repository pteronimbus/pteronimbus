package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "A simple ping command",
		},
		{
			Name:        "sync",
			Description: "Sync roles and users from Discord",
		},
		{
			Name:        "say",
			Description: "Make the bot say something in a channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel to send the message to",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The message to send",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

type Bot struct {
	Session      *discordgo.Session
	token        string
	syncService  Syncer
	auditService Auditer
}

type Auditer interface {
	Log(event string, details map[string]interface{})
}

type Syncer interface {
	SyncRoles(tenantID string, guildID string) error
	SyncUsers(tenantID string, guildID string) error
}

func NewBot(token string, syncService Syncer, auditService Auditer) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}
	bot := &Bot{
		Session:      s,
		token:        token,
		syncService:  syncService,
		auditService: auditService,
	}
	bot.setupHandlers()
	return bot, nil
}

func (b *Bot) setupHandlers() {
	commandHandlers["ping"] = b.handlePing
	commandHandlers["sync"] = b.handleSync
	commandHandlers["say"] = b.handleSay
}

func (b *Bot) handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	b.auditService.Log("ping_command", map[string]interface{}{
		"user_id":  i.Member.User.ID,
		"guild_id": i.GuildID,
	})
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
}

func (b *Bot) handleSay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	channel := optionMap["channel"].ChannelValue(s)
	message := optionMap["message"].StringValue()

	b.auditService.Log("say_command", map[string]interface{}{
		"user_id":    i.Member.User.ID,
		"guild_id":   i.GuildID,
		"channel_id": channel.ID,
		"message":    message,
	})

	_, err := b.SendMessage(channel.ID, message)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to send message.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Message sent!",
		},
	})
}

func (b *Bot) handleSync(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Acknowledge the interaction immediately.
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Printf("Error acknowledging interaction: %v\n", err)
		return
	}

	b.auditService.Log("sync_command", map[string]interface{}{
		"user_id":  i.Member.User.ID,
		"guild_id": i.GuildID,
	})

	// For now, hardcoding tenant and guild ID.
	// This will be replaced with dynamic values later.
	tenantID := "f9c1b2a9-7c9c-4f7d-8e4a-5f0e1c2b3a4d" // Example UUID
	guildID := i.GuildID

	err = b.syncService.SyncRoles(tenantID, guildID)
	if err != nil {
		content := fmt.Sprintf("Error syncing roles: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}

	err = b.syncService.SyncUsers(tenantID, guildID)
	if err != nil {
		content := fmt.Sprintf("Error syncing users: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}

	content := "Roles and users synced successfully!"
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

func (b *Bot) Start() error {
	b.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	fmt.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", v)
		if err != nil {
			return fmt.Errorf("cannot create '%v' command: %w", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	fmt.Println("Bot is now running. The application will handle shutdown.")
	// Keep the connection alive
	select {}
}

func (b *Bot) Stop() {
	fmt.Println("Removing commands...")
	for _, v := range commands {
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, "", v.ID)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v\n", v.Name, err)
		}
	}
}

func (b *Bot) SendMessage(channelID string, message string) (*discordgo.Message, error) {
	return b.Session.ChannelMessageSend(channelID, message)
}

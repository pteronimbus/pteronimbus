package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"golang.org/x/oauth2"
)

// DiscordService handles Discord OAuth2 and API operations
type DiscordService struct {
	config      *config.DiscordConfig
	oauthConfig *oauth2.Config
	httpClient  *http.Client
}

// NewDiscordService creates a new Discord service
func NewDiscordService(cfg *config.Config) *DiscordService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Discord.ClientID,
		ClientSecret: cfg.Discord.ClientSecret,
		RedirectURL:  cfg.Discord.RedirectURL,
		Scopes:       []string{"identify", "email", "guilds"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}

	return &DiscordService{
		config:      &cfg.Discord,
		oauthConfig: oauthConfig,
		httpClient:  &http.Client{},
	}
}

// GetAuthURL generates Discord OAuth2 authorization URL
func (d *DiscordService) GetAuthURL(state string) string {
	return d.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCodeForToken exchanges authorization code for access token
func (d *DiscordService) ExchangeCodeForToken(ctx context.Context, code string) (*models.DiscordTokenResponse, error) {
	token, err := d.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	discordToken := &models.DiscordTokenResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Scope:        strings.Join(d.oauthConfig.Scopes, " "),
	}

	if token.Expiry.IsZero() {
		discordToken.ExpiresIn = 3600 // Default 1 hour if not provided
	} else {
		discordToken.ExpiresIn = int(token.Expiry.Sub(token.Expiry).Seconds())
	}

	return discordToken, nil
}

// GetUserInfo retrieves user information from Discord API
func (d *DiscordService) GetUserInfo(ctx context.Context, accessToken string) (*models.DiscordUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.config.APIBaseURL+"/users/@me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Use correct Authorization header for bot tokens
	if strings.HasPrefix(accessToken, "Bot ") {
		req.Header.Set("Authorization", accessToken)
	} else if accessToken == d.config.BotToken {
		req.Header.Set("Authorization", "Bot "+accessToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var discordUser models.DiscordUser
	err = json.Unmarshal(body, &discordUser)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &discordUser, nil
}

// RefreshToken refreshes a Discord access token
func (d *DiscordService) RefreshToken(ctx context.Context, refreshToken string) (*models.DiscordTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", d.config.ClientID)
	data.Set("client_secret", d.config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make refresh request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord refresh token error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response body: %w", err)
	}

	var tokenResponse models.DiscordTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal refresh token response: %w", err)
	}

	return &tokenResponse, nil
}

// GetUserGuilds retrieves the guilds/servers a user has access to
func (d *DiscordService) GetUserGuilds(ctx context.Context, accessToken string) ([]models.DiscordGuild, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.config.APIBaseURL+"/users/@me/guilds", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create guilds request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make guilds request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord guilds API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read guilds response body: %w", err)
	}

	var guilds []models.DiscordGuild
	err = json.Unmarshal(body, &guilds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal guilds data: %w", err)
	}

	return guilds, nil
}

// GetGuildRoles retrieves roles for a specific guild
func (d *DiscordService) GetGuildRoles(ctx context.Context, botToken, guildID string) ([]models.DiscordRole, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", d.config.APIBaseURL+"/guilds/"+guildID+"/roles", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create guild roles request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make guild roles request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord guild roles API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read guild roles response body: %w", err)
	}

	var roles []models.DiscordRole
	err = json.Unmarshal(body, &roles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal guild roles data: %w", err)
	}

	return roles, nil
}

// GetGuildMembers retrieves members for a specific guild
func (d *DiscordService) GetGuildMembers(ctx context.Context, botToken, guildID string, limit int) ([]models.DiscordMember, error) {
	url := fmt.Sprintf("%s/guilds/%s/members?limit=%d", d.config.APIBaseURL, guildID, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create guild members request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make guild members request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord guild members API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read guild members response body: %w", err)
	}

	var members []models.DiscordMember
	err = json.Unmarshal(body, &members)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal guild members data: %w", err)
	}

	return members, nil
}

// GetGuildMember retrieves a specific member from a guild
func (d *DiscordService) GetGuildMember(ctx context.Context, botToken, guildID, userID string) (*models.DiscordMember, error) {
	url := fmt.Sprintf("%s/guilds/%s/members/%s", d.config.APIBaseURL, guildID, userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create guild member request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make guild member request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("member not found in guild")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord guild member API error: %d - %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read guild member response body: %w", err)
	}

	var member models.DiscordMember
	err = json.Unmarshal(body, &member)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal guild member data: %w", err)
	}

	return &member, nil
}

// Add this method to allow access to the bot token
func (d *DiscordService) BotToken() string {
	return d.config.BotToken
}

package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// WebhookConfig holds configuration for Discord/Slack webhooks
type WebhookConfig struct {
	URL       string   `mapstructure:"url"`
	MinLevel  LogLevel `mapstructure:"min_level"`
	Enable    bool     `mapstructure:"enable"`
	ChannelID string   `mapstructure:"channel_id,omitempty"` // For Discord
}

// LoggingConfig for log settings
type LoggingConfig struct {
	Level   string        `mapstructure:"level"`
	Format  string        `mapstructure:"format"`
	Output  string        `mapstructure:"output"`
	Discord WebhookConfig `mapstructure:"discord"`
	Slack   WebhookConfig `mapstructure:"slack"`
}

// multiLogger sends logs to multiple destinations
type multiLogger struct {
	consoleLogger Logger
	discordLogger *webhookLogger
	slackLogger   *webhookLogger
}

// webhookLogger implements webhook-based logging
type webhookLogger struct {
	webhookURL string
	channelID  string
	minLevel   LogLevel
	client     *http.Client
	isDiscord  bool
}

// NewMultiLogger creates a logger that sends logs to console and webhooks
func NewMultiLogger(cfg LoggingConfig) Logger {
	consoleLogger := New(LogLevel(cfg.Level), nil)

	var discordLogger *webhookLogger
	if cfg.Discord.Enable && cfg.Discord.URL != "" {
		discordLogger = &webhookLogger{
			webhookURL: cfg.Discord.URL,
			channelID:  cfg.Discord.ChannelID,
			minLevel:   cfg.Discord.MinLevel,
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
			isDiscord: true,
		}
	}

	var slackLogger *webhookLogger
	if cfg.Slack.Enable && cfg.Slack.URL != "" {
		slackLogger = &webhookLogger{
			webhookURL: cfg.Slack.URL,
			minLevel:   cfg.Slack.MinLevel,
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
			isDiscord: false,
		}
	}

	return &multiLogger{
		consoleLogger: consoleLogger,
		discordLogger: discordLogger,
		slackLogger:   slackLogger,
	}
}

// shouldLog determines if log should be sent to webhook based on level
func (w *webhookLogger) shouldLog(level LogLevel) bool {
	if w == nil {
		return false
	}

	levelOrder := map[LogLevel]int{
		DebugLevel: 0,
		InfoLevel:  1,
		WarnLevel:  2,
		ErrorLevel: 3,
	}

	return levelOrder[level] >= levelOrder[w.minLevel]
}

// sendToDiscord sends a log message to Discord webhook
func (w *webhookLogger) sendToDiscord(level LogLevel, msg string, args ...any) {
	if !w.shouldLog(level) {
		return
	}

	// Format args as key-value pairs
	fields := make([]map[string]string, 0)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields = append(fields, map[string]string{
				"name":   fmt.Sprintf("%v", args[i]),
				"value":  fmt.Sprintf("%v", args[i+1]),
				"inline": "true",
			})
		}
	}

	color := map[LogLevel]int{
		DebugLevel: 7506394,  // Gray
		InfoLevel:  3447003,  // Blue
		WarnLevel:  16776960, // Yellow
		ErrorLevel: 15158332, // Red
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":     fmt.Sprintf("%s: %s", level, msg),
				"color":     color[level],
				"timestamp": time.Now().Format(time.RFC3339),
				"fields":    fields,
				"footer": map[string]string{
					"text": "Application Log",
				},
			},
		},
	}

	go w.sendWebhook(payload)
}

// sendToSlack sends a log message to Slack webhook
func (w *webhookLogger) sendToSlack(level LogLevel, msg string, args ...any) {
	if !w.shouldLog(level) {
		return
	}

	// Format args as key-value pairs
	fields := make([]map[string]string, 0)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields = append(fields, map[string]string{
				"title": fmt.Sprintf("%v", args[i]),
				"value": fmt.Sprintf("%v", args[i+1]),
				"short": "true",
			})
		}
	}

	color := map[LogLevel]string{
		DebugLevel: "#808080", // Gray
		InfoLevel:  "#3498db", // Blue
		WarnLevel:  "#f1c40f", // Yellow
		ErrorLevel: "#e74c3c", // Red
	}

	payload := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"title":  fmt.Sprintf("%s: %s", level, msg),
				"color":  color[level],
				"fields": fields,
				"footer": "Application Log",
				"ts":     time.Now().Unix(),
			},
		},
	}

	go w.sendWebhook(payload)
}

// sendWebhook makes the HTTP request to webhook URL
func (w *webhookLogger) sendWebhook(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling webhook payload: %v\n", err)
		return
	}

	resp, err := w.client.Post(w.webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending webhook: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Webhook error: status=%d, body=%s\n", resp.StatusCode, string(body))
	}
}

// Implement Logger interface for multiLogger
func (m *multiLogger) Debug(msg string, args ...any) {
	m.consoleLogger.Debug(msg, args...)
	if m.discordLogger != nil {
		m.discordLogger.sendToDiscord(DebugLevel, msg, args...)
	}
	if m.slackLogger != nil {
		m.slackLogger.sendToSlack(DebugLevel, msg, args...)
	}
}

func (m *multiLogger) Info(msg string, args ...any) {
	m.consoleLogger.Info(msg, args...)
	if m.discordLogger != nil {
		m.discordLogger.sendToDiscord(InfoLevel, msg, args...)
	}
	if m.slackLogger != nil {
		m.slackLogger.sendToSlack(InfoLevel, msg, args...)
	}
}

func (m *multiLogger) Warn(msg string, args ...any) {
	m.consoleLogger.Warn(msg, args...)
	if m.discordLogger != nil {
		m.discordLogger.sendToDiscord(WarnLevel, msg, args...)
	}
	if m.slackLogger != nil {
		m.slackLogger.sendToSlack(WarnLevel, msg, args...)
	}
}

func (m *multiLogger) Error(msg string, args ...any) {
	m.consoleLogger.Error(msg, args...)
	if m.discordLogger != nil {
		m.discordLogger.sendToDiscord(ErrorLevel, msg, args...)
	}
	if m.slackLogger != nil {
		m.slackLogger.sendToSlack(ErrorLevel, msg, args...)
	}
}

func (m *multiLogger) With(args ...any) Logger {
	return &multiLogger{
		consoleLogger: m.consoleLogger.With(args...),
		discordLogger: m.discordLogger,
		slackLogger:   m.slackLogger,
	}
}

func (m *multiLogger) WithContext(ctx context.Context) Logger {
	return &multiLogger{
		consoleLogger: m.consoleLogger.WithContext(ctx),
		discordLogger: m.discordLogger,
		slackLogger:   m.slackLogger,
	}
}

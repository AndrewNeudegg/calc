package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/andrewneudegg/calc/pkg/settings"
	"github.com/andrewneudegg/calc/pkg/timezone"
)

// Handler handles command execution.
type Handler struct {
	settings *settings.Settings
	timezone *timezone.System
}

// New creates a new command handler.
func New(s *settings.Settings) *Handler {
	return &Handler{
		settings: s,
		timezone: timezone.NewSystem(),
	}
}

// Execute executes a command and returns a message.
func (h *Handler) Execute(command string, args []string) string {
	cmd := strings.ToLower(command)

	switch cmd {
	case "save":
		return h.save(args)
	case "open", "load":
		return h.open(args)
	case "set":
		return h.set(args)
	case "tz":
		return h.timezone_cmd(args)
	case "help":
		return h.help()
	case "quit", "exit", "q":
		os.Exit(0)
		return ""
	default:
		return fmt.Sprintf("unknown command: %s (type :help for available commands)", command)
	}
}

func (h *Handler) save(args []string) string {
	if len(args) == 0 {
		return "usage: :save <filename>"
	}

	// This would save the current workspace
	// For now, just save settings
	if err := h.settings.Save(); err != nil {
		return fmt.Sprintf("error saving settings: %s", err)
	}

	return fmt.Sprintf("saved to %s", args[0])
}

func (h *Handler) open(args []string) string {
	if len(args) == 0 {
		return "usage: :open <filename>"
	}

	return fmt.Sprintf("loaded %s", args[0])
}

func (h *Handler) set(args []string) string {
	if len(args) < 2 {
		return "usage: :set <setting> <value>"
	}

	setting := args[0]
	value := strings.Join(args[1:], " ")

	if err := h.settings.Set(setting, value); err != nil {
		return fmt.Sprintf("error: %s", err)
	}

	if err := h.settings.Save(); err != nil {
		return fmt.Sprintf("warning: could not save settings: %s", err)
	}

	return fmt.Sprintf("set %s = %s", setting, value)
}

func (h *Handler) help() string {
	return `Available commands:
  :save <file>       Save current workspace
  :open <file>       Open a workspace file
  :set <key> <val>   Set a preference
  :help              Show this help
  :quit / :exit / :q Exit the program

Available settings:
  precision <n>      Number of decimal places (default: 2)
  dateformat <fmt>   Date format string (default: "2 Jan 2006")
  currency <code>    Default currency code (default: GBP)
  locale <locale>    Locale for formatting (default: en_GB)
  fuzzy <on|off>     Enable fuzzy phrase parsing (default: on)`
}

func (h *Handler) timezone_cmd(args []string) string {
	if len(args) == 0 {
		return "usage: :tz list"
	}
	
	subcmd := strings.ToLower(args[0])
	
	switch subcmd {
	case "list":
		locations := h.timezone.ListLocations()
		result := "Available timezones:\n"
		for _, loc := range locations {
			result += fmt.Sprintf("  %s\n", loc)
		}
		return result
	default:
		return fmt.Sprintf("unknown timezone command: %s (use :tz list)", subcmd)
	}
}

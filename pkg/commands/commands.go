package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/andrewneudegg/calc/pkg/constants"
	"github.com/andrewneudegg/calc/pkg/settings"
	"github.com/andrewneudegg/calc/pkg/timezone"
	"github.com/andrewneudegg/calc/pkg/units"
)

// Handler handles command execution.
type Handler struct {
	settings  *settings.Settings
	timezone  *timezone.System
	constants *constants.System
	units     *units.System
	// Optional workspace operations provided by the REPL
	SaveWorkspace  func(filename string) error
	LoadWorkspace  func(filename string) error
	ClearWorkspace func() error
	// Quiet mode controls provided by the REPL
	SetQuiet    func(enabled bool)
	ToggleQuiet func() bool
	GetQuiet    func() bool
	// Unit management
	DefineUnit func(name string, value float64, baseUnit string) error
	// shouldQuit is set to true when the quit command is executed
	shouldQuit bool
}

// New creates a new command handler.
func New(s *settings.Settings) *Handler {
	return &Handler{
		settings:  s,
		timezone:  timezone.NewSystem(),
		constants: constants.NewSystem(),
	}
}

// SetUnits sets the units system for the command handler.
func (h *Handler) SetUnits(u *units.System) {
	h.units = u
}

// ShouldQuit returns true if the quit command has been executed.
func (h *Handler) ShouldQuit() bool {
	return h.shouldQuit
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
	case "const":
		return h.const_cmd(args)
	case "unit":
		return h.unit_cmd(args)
	case "help":
		return h.help()
	case "clear", "cls":
		return h.clear()
	case "quiet":
		return h.quiet(args)
	case "quit", "exit", "q":
		h.shouldQuit = true
		return ""
	default:
		return fmt.Sprintf("unknown command: %s (type :help for available commands)", command)
	}
}

func (h *Handler) save(args []string) string {
	if len(args) == 0 {
		return "usage: :save <filename>"
	}

	// Save settings first (preferences)
	if err := h.settings.Save(); err != nil {
		return fmt.Sprintf("error saving settings: %s", err)
	}

	// Save the current workspace if a handler is available
	if h.SaveWorkspace != nil {
		if err := h.SaveWorkspace(args[0]); err != nil {
			return fmt.Sprintf("error saving workspace: %s", err)
		}
	}

	return fmt.Sprintf("saved to %s", args[0])
}

func (h *Handler) open(args []string) string {
	if len(args) == 0 {
		return "usage: :open <filename>"
	}

	if h.LoadWorkspace != nil {
		if err := h.LoadWorkspace(args[0]); err != nil {
			return fmt.Sprintf("error loading %s: %v", args[0], err)
		}
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
	:clear             Clear screen and reset current session
	:quiet [on|off]    Toggle or set quiet mode (suppress assignment output)
  :const list        List all physical constants
  :const show <name> Show details of a specific constant
  :unit define <name> = <value> <unit>  Define a custom unit
  :unit list         List all units (custom and standard)
  :unit show <name>  Show details of a specific unit
  :unit delete <name> Delete a custom unit
  :help              Show this help
  :quit / :exit / :q Exit the program

Available settings:
  precision <n>         Number of decimal places (default: 2)
  dateformat <fmt>      Date format string (default: "2 Jan 2006")
  currency <code>       Default currency code (default: GBP)
  locale <locale>       Locale for formatting (default: en_GB)
  fuzzy <on|off>        Enable fuzzy phrase parsing (default: on)
  autocomplete <on|off> Enable autocomplete suggestions (default: on)`
}

func (h *Handler) clear() string {
	// Reset workspace/session state if provided by REPL
	if h.ClearWorkspace != nil {
		if err := h.ClearWorkspace(); err != nil {
			return fmt.Sprintf("error clearing session: %s", err)
		}
	}
	// Return ANSI clear-screen sequence and home cursor
	// This will be printed directly by the REPL handler
	return "\x1b[2J\x1b[H"
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

func (h *Handler) quiet(args []string) string {
	// Require REPL to wire quiet controls
	if h.SetQuiet == nil && h.ToggleQuiet == nil && h.GetQuiet == nil {
		return "quiet mode not supported in this context"
	}

	// No args: toggle
	if len(args) == 0 {
		var on bool
		if h.ToggleQuiet != nil {
			on = h.ToggleQuiet()
		} else if h.GetQuiet != nil && h.SetQuiet != nil {
			// Fallback toggle via get/set
			cur := h.GetQuiet()
			on = !cur
			h.SetQuiet(on)
		}
		if on {
			return "quiet: on"
		}
		return "quiet: off"
	}

	// With arg: on/off
	v := strings.ToLower(args[0])
	switch v {
	case "on", "true", "1", "yes", "y":
		if h.SetQuiet != nil {
			h.SetQuiet(true)
		}
		return "quiet: on"
	case "off", "false", "0", "no", "n":
		if h.SetQuiet != nil {
			h.SetQuiet(false)
		}
		return "quiet: off"
	default:
		return "usage: :quiet [on|off]"
	}
}

func (h *Handler) const_cmd(args []string) string {
	if len(args) == 0 {
		return "usage: :const list | :const show <name>"
	}

	subcmd := strings.ToLower(args[0])

	switch subcmd {
	case "list":
		return h.constList(args[1:])
	case "show":
		if len(args) < 2 {
			return "usage: :const show <name>"
		}
		return h.constShow(args[1])
	default:
		return fmt.Sprintf("unknown const command: %s (use :const list or :const show <name>)", subcmd)
	}
}

func (h *Handler) constList(args []string) string {
	var consts []*constants.Constant

	// If category specified, filter by category
	if len(args) > 0 {
		category := args[0]
		consts = h.constants.ListByCategory(category)
		if len(consts) == 0 {
			cats := h.constants.GetCategories()
			sort.Strings(cats)
			return fmt.Sprintf("Unknown category: %s\nAvailable categories: %s", category, strings.Join(cats, ", "))
		}
	} else {
		consts = h.constants.ListConstants()
	}

	// Sort by name
	sort.Slice(consts, func(i, j int) bool {
		return consts[i].Name < consts[j].Name
	})

	result := "Physical Constants:\n"
	for _, c := range consts {
		symbol := ""
		if c.Symbol != "" && c.Symbol != c.Name {
			symbol = fmt.Sprintf(" (%s)", c.Symbol)
		}
		result += fmt.Sprintf("  %-25s = %e %s\n", c.Name+symbol, c.Value, c.Unit)
	}

	// Show categories if no category specified
	if len(args) == 0 {
		cats := h.constants.GetCategories()
		sort.Strings(cats)
		result += fmt.Sprintf("\nCategories: %s\n", strings.Join(cats, ", "))
		result += "Use :const list <category> to filter by category\n"
	}

	return result
}

func (h *Handler) constShow(name string) string {
	c, err := h.constants.GetConstant(name)
	if err != nil {
		return fmt.Sprintf("Unknown constant: %s\nUse :const list to see all constants", name)
	}

	result := fmt.Sprintf("Constant: %s\n", c.Name)
	if c.Symbol != "" && c.Symbol != c.Name {
		result += fmt.Sprintf("Symbol: %s\n", c.Symbol)
	}
	result += fmt.Sprintf("Value: %e\n", c.Value)
	if c.Unit != "" {
		result += fmt.Sprintf("Unit: %s\n", c.Unit)
	}
	result += fmt.Sprintf("Category: %s\n", c.Category)
	if c.Description != "" {
		result += fmt.Sprintf("Description: %s\n", c.Description)
	}

	return result
}

func (h *Handler) unit_cmd(args []string) string {
	if len(args) == 0 {
		return "usage: :unit define|list|show|delete"
	}

	subcmd := strings.ToLower(args[0])

	switch subcmd {
	case "define":
		return h.unitDefine(args[1:])
	case "list":
		return h.unitList()
	case "show":
		if len(args) < 2 {
			return "usage: :unit show <name>"
		}
		return h.unitShow(args[1])
	case "delete":
		if len(args) < 2 {
			return "usage: :unit delete <name>"
		}
		return h.unitDelete(args[1])
	default:
		return fmt.Sprintf("unknown unit command: %s (use :unit define|list|show|delete)", subcmd)
	}
}

func (h *Handler) unitDefine(args []string) string {
	if h.units == nil {
		return "unit system not available"
	}

	if h.DefineUnit == nil {
		return "unit definition not available in this context"
	}

	// Parse: <name> = <value> <unit>
	// Join all args and parse
	input := strings.Join(args, " ")
	parts := strings.Split(input, "=")
	if len(parts) != 2 {
		return "usage: :unit define <name> = <value> <unit>\nexample: :unit define spoon = 15 ml"
	}

	name := strings.TrimSpace(parts[0])
	if name == "" {
		return "unit name cannot be empty"
	}

	// Check if unit already exists
	if h.units.UnitExists(name) {
		return fmt.Sprintf("unit '%s' already exists (use :unit delete first to redefine)", name)
	}

	// Parse the value and unit from the right side
	valueUnitStr := strings.TrimSpace(parts[1])
	
	// Call the DefineUnit callback which will parse and evaluate the expression
	// For now, we'll return an instruction for the REPL to handle this
	return fmt.Sprintf("DEFINE_UNIT:%s:%s", name, valueUnitStr)
}

func (h *Handler) unitList() string {
	if h.units == nil {
		return "unit system not available"
	}

	customUnits := h.units.ListCustomUnits()
	
	if len(customUnits) == 0 {
		return "No custom units defined\nUse :unit define <name> = <value> <unit> to create one"
	}

	result := "Custom Units:\n"
	for _, u := range customUnits {
		result += fmt.Sprintf("  %-15s (dimension: %d, base: %s)\n", u.Name, u.Dimension, u.BaseUnit)
	}

	return result
}

func (h *Handler) unitShow(name string) string {
	if h.units == nil {
		return "unit system not available"
	}

	unit, ok := h.units.GetUnit(name)
	if !ok {
		return fmt.Sprintf("unknown unit: %s\nUse :unit list to see custom units", name)
	}

	result := fmt.Sprintf("Unit: %s\n", unit.Name)
	result += fmt.Sprintf("Dimension: %d\n", unit.Dimension)
	result += fmt.Sprintf("Base Unit: %s\n", unit.BaseUnit)
	result += fmt.Sprintf("Conversion Factor: %e\n", unit.ToBase)
	
	if unit.IsCustom {
		result += "Type: Custom\n"
	} else {
		result += "Type: Standard\n"
	}

	return result
}

func (h *Handler) unitDelete(name string) string {
	if h.units == nil {
		return "unit system not available"
	}

	err := h.units.DeleteCustomUnit(name)
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}

	return fmt.Sprintf("deleted custom unit: %s", name)
}

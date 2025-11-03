# Calc - Terminal Notepad Calculator

A local-only, dependency-free terminal calculator inspired by Soulver. Mix free-form text with arithmetic, units, dates, currencies, and variables—all evaluated in real-time with keyboard-driven simplicity.

## Features

### Core Functionality
- **Arithmetic & Variables**: Standard operations with operator precedence, parentheses, and named variables
- **Unit Conversion**: SI units (length, mass, time, volume, temperature) with automatic conversion
- **Custom Units**: Define your own units for domain-specific calculations
- **Currency Operations**: Built-in currency conversion with static exchange rates
- **Percentages**: Intuitive percentage calculations including "X% of Y", "increase/decrease by Y%"
- **Date & Time Arithmetic**: Natural language date operations (today, tomorrow, next week, etc.)
- **Fuzzy Phrases**: English phrases like "half of X", "double X", "three quarters of Y"
- **Functions**: sum(), average(), mean(), and natural language equivalents
- **Dependency Tracking**: Automatic recalculation when variables change

### Design Principles
- **Local & Deterministic**: No network calls, no AI, completely reproducible
- **Zero Dependencies**: Built entirely with Go standard library
- **Keyboard-Driven**: Fast, focused workflow with command mode
- **British English**: All identifiers, comments, and documentation use British English

## Installation

```bash
go build ./cmd/calc
```

## Usage

### Basic Arithmetic
```
1> 3 + 4 * 5
   = 23.00

2> (10 + 5) / 2
   = 7.50
```

### Variables
```
3> hourly = £23.50
   = £23.50

4> weekly = hourly * 37.5
   = £881.25
```

### Unit Conversions
```
5> 10 m in cm
   = 1,000.00 cm

6> 100 km in miles
   = 62.14 miles

7> 70 kg in lb
   = 154.32 lb
```

### Currency
```
8> £120 + $30
   = £143.50

9> $50 in GBP
   = £39.50
```

### Percentages
```
10> 30 + 20%
   = 36.00

11> 20% of 50
   = 10.00

12> increase 100 by 10%
   = 110.00

13> 20 is what % of 50
   = 40.00%
```

### Fuzzy Phrases
```
14> half of 80
   = 40.00

15> double 15
   = 30.00

16> three quarters of 200
   = 150.00
```

### Fuzzy Phrases with Variables
Fuzzy phrases work seamlessly with variables and assignments:

```
17> fo = half of 99
   = 49.50

18> half of fo
   = 24.75

19> result = double fo
   = 99.00

20> x = 100
   = 100.00

21> y = half of x
   = 50.00

22> z = double y
   = 100.00
```

### Time Arithmetic
Times in `HH:MM` format are recognised and maintain their format through calculations:

```
23> 11:00 - 09:00
   = 02:00

24> 14:00 + 2
   = 16:00

25> start = 09:30
   = 09:30

26> end = 17:45
   = 17:45

27> end - start
   = 08:15

28> meeting = start + 2
   = 11:30
```

Times are stored as time units and displayed in `HH:MM` format. You can add or subtract hours (as numbers) or other times.

### Functions
```
29> sum(10, 20, 30)
   = 60.00

30> average(3, 4, 5)
   = 4.00
```

### Date Arithmetic
```
31> today + 3 weeks
   = 22 Nov 2025

32> tomorrow - 2 days
   = 31 Oct 2025
```

### Commands
```
:help              Show available commands
:set precision 3   Set decimal precision
:set currency GBP  Set default currency
:save file.txt     Save workspace
:quit              Exit
```

## Project Structure

```
calc/
├── cmd/
│   └── calc/          # Main application entry point
├── pkg/
│   ├── lexer/         # Tokenisation
│   ├── parser/        # AST generation
│   ├── evaluator/     # Expression evaluation
│   ├── units/         # Unit conversion system
│   ├── currency/      # Currency handling
│   ├── formatter/     # Output formatting
│   ├── settings/      # User preferences
│   ├── commands/      # Command mode
│   ├── graph/         # Dependency tracking
│   └── display/       # REPL interface
└── docs/
    ├── briefs/        # Project brief
    └── specifications/# Detailed specifications
```

## Testing

Run the full test suite:

```bash
go test ./...
```

Run with race detection:

```bash
go test ./... -race
```

View test coverage:

```bash
go test ./... -cover
```

## Supported Units

### Length
metre (m), centimetre (cm), millimetre (mm), kilometre (km), foot (ft), inch (in), yard (yd), mile (mi)

### Mass
kilogram (kg), gram (g), milligram (mg), pound (lb), ounce (oz)

### Time
second (s), minute (min), hour (h), day, week, month, year

### Volume
litre (l), millilitre (ml), gallon (gal)

### Temperature
celsius (c), fahrenheit (f)

## Currency Support

- USD ($)
- GBP (£)
- EUR (€)
- JPY (¥)

Exchange rates can be customised via the settings system.

## Settings

Settings are stored in `~/.config/calc/settings.json`:

```json
{
  "precision": 2,
  "date_format": "2 Jan 2006",
  "currency": "GBP",
  "locale": "en_GB",
  "fuzzy_mode": true
}
```

## Implementation Notes

- **Parsing**: Hand-written recursive descent parser for deterministic, fast parsing
- **Evaluation**: Tree-walking interpreter with environment for variable storage
- **Units**: Dimension-based system prevents incompatible unit operations
- **No Dependencies**: Maintains project constraint of zero external dependencies

## Licence

This is a demonstration project built to specification.

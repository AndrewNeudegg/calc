# calc

Terminal calculator with units, currency conversion, and natural language expressions. Zero dependencies, runs locally.

## Features

- Arithmetic with operator precedence and parentheses
- Named variables with dependency tracking
- Unit conversions (length, mass, time, volume, temperature, etc.)
- Currency conversion with postfix notation support
- Date and time arithmetic
- Percentage calculations
- Natural language phrases ("half of", "double", etc.)
- Built-in functions (sum, average, mean)
- REPL with command mode, syntax highlighting, and themes
- Line comments with // (ignored by the lexer)
- Save/load workspace files from the REPL

## Installation

### Quick Install (Latest Release)

```bash
curl -fsSL https://raw.githubusercontent.com/AndrewNeudegg/calc/main/install.sh | sh
```

Or download manually from [releases](https://github.com/AndrewNeudegg/calc/releases).

### Build from Source

```bash
go build -o calc ./cmd/calc
./calc
```

Or run directly:
```bash
go run ./cmd/calc
```

### CLI Usage

Single calculation mode:
```bash
./calc -c "12 gbp in dollars"
```

Execute a script file:
```bash
./calc -f examples/k8s-cluster.calc
```

Read from stdin (use '-' as the file):
```bash
cat examples/k8s-cluster.calc | ./calc -f -
```

Show help:
```bash
./calc -h
```

## Quick Reference

### Operators

| Operator | Description | Example | Result |
|----------|-------------|---------|--------|
| `+` | Addition | `5 + 3` | `8.00` |
| `-` | Subtraction | `10 - 4` | `6.00` |
| `*` | Multiplication | `6 * 7` | `42.00` |
| `/` | Division | `20 / 4` | `5.00` |
| `%` | Percentage | `20%` | `0.20` |
| `=` | Assignment | `x = 10` | `10.00` |
| `in` | Unit conversion | `10 m in cm` | `1,000.00 cm` |

### Currency Formats

| Format | Example | Display |
|--------|---------|---------|
| Symbol prefix | `£12`, `$50`, `€100`, `¥1000` | Currency symbol shown |
| Code postfix | `12 gbp`, `50 usd`, `100 eur` | Converted to symbol |
| Name postfix | `50 dollars`, `25 euros`, `1000 yen` | Converted to symbol |

Supported: USD ($), GBP (£), EUR (€), JPY (¥), and many more codes including: AUD, CAD, NZD, CHF, CNY, HKD, SGD, INR, KRW, TWD, SEK, NOK, DKK, TRY, RUB, PLN, CZK, HUF, RON, ILS, AED, SAR, THB, MYR, IDR, PHP, ZAR, MXN, BRL

**Note:** "pound" and "pounds" refer to weight (lb). Use "gbp" or "£" for currency.

### Time Format

Times in `HH:MM` format are recognized automatically:

| Expression | Result |
|------------|--------|
| `14:00 + 2` | `16:00` |
| `11:00 - 09:00` | `02:00` |
| `17:45 - 09:30` | `08:15` |

### Natural Language

| Phrase | Example | Result |
|--------|---------|--------|
| `half of` | `half of 100` | `50.00` |
| `double` | `double 25` | `50.00` |
| `three quarters of` | `three quarters of 80` | `60.00` |
| `X% of Y` | `20% of 50` | `10.00` |
| `increase X by Y%` | `increase 100 by 10%` | `110.00` |
| `decrease X by Y%` | `decrease 100 by 10%` | `90.00` |
| `X is what % of Y` | `20 is what % of 50` | `40.00%` |

### Functions

| Function | Description | Example |
|----------|-------------|---------|
| `sum(...)` | Sum of all arguments | `sum(10, 20, 30)` → `60.00` |
| `average(...)` | Average of arguments | `average(10, 20, 30)` → `20.00` |
| `mean(...)` | Alias for average | `mean(5, 10, 15)` → `10.00` |
| `min(...)` | Minimum of arguments | `min(3, 7, 2, 9)` → `2.00` |
| `max(...)` | Maximum of arguments | `max(3, 7, 2, 9)` → `9.00` |
| `print("...")` | Interpolate `{var}` placeholders and return the string | `tt = 55` then `print("foo: {tt}")` → `foo: 55` |

Notes:
- Function arguments can be expressions.
- `sum()` with no arguments returns `0`.
- `average()` requires at least one argument; calling it with none is an error.
- `min()`/`max()` require at least one argument; calling either with none is an error.
- Functions return plain numbers. If you need to preserve units or currency, convert to a common unit first or use explicit operators (e.g., `a + b` instead of `sum(a, b)`).

### Strings and Print

- String literals are written with double quotes: `"hello world"`.
- Use `print("...")` to interpolate variables inside `{}` and produce a string result that is printed by the REPL/CLI:

```
1> tt = 55
   = 55.00
2> print("foo bar: {tt}, foo bar")
   = foo bar: 55, foo bar
```

Notes:
- Placeholders must be simple variable names, e.g., `{rate}`, `{total_cost}`.
- If a placeholder variable is undefined, `print` returns an error.
- Interpolated values use sensible defaults for formatting; date/time values render in a readable form.

### Date Keywords

| Keyword | Description |
|---------|-------------|
| `now` | Current date/time |
| `today` | Current date |
| `tomorrow` | Today + 1 day |
| `yesterday` | Today - 1 day |
| `next week` | Today + 7 days |
| `last week` | Today - 7 days |
| `next month` | Today + 30 days |

Also supported in date arithmetic: smaller units including hours, minutes, and seconds (e.g., `today + 3 days + 2 hours`).

### REPL Commands

| Command | Description |
|---------|-------------|
| `:help` | Show available commands |
| `:save <file>` | Save current workspace to the current directory |
| `:open <file>` | Open a workspace file and restore variables |
| `:set <key> <value>` | Update a preference (see below) |
| `:clear` | Clear screen and reset current session |
| `:quit` / `:exit` / `:q` | Exit |
| `:tz list` | List available timezones |

Settings keys for `:set`:
- `precision <n>` – Number of decimal places (default: 2)
- `dateformat <fmt>` – Date format string (default: `2 Jan 2006`)
- `currency <CODE>` – Default currency code (GBP, USD, EUR, JPY)
- `locale <locale>` – Locale for formatting (default: `en_GB`)
- `fuzzy <on|off>` – Enable/disable natural-language parsing

Tips:
- Press Ctrl-C to cancel the current input line.
- Press Ctrl-D to exit (same as `:quit`).
- Type `:help` any time to see the command summary.

Notes on saving:
- `:save <file>` writes a plain-text workspace file in your current working directory. Only expressions are saved (commands are skipped).
- Preferences are stored separately at `~/.config/calc/settings.json` and are also saved when you run `:save`.

### Comments

Use `//` for line comments. Everything after `//` on a line is ignored:

```
x = 10 // hourly rate
rate = 2.5 // per hour
```

## Supported Units

### Length

| Unit | Aliases | Symbol |
|------|---------|--------|
| metre | meter, metres, meters | m |
| centimetre | centimeter, centimetres, centimeters | cm |
| millimetre | millimeter, millimetres, millimeters | mm |
| kilometre | kilometer, kilometres, kilometers | km |
| foot | feet | ft |
| inch | inches | in |
| yard | yards | yd |
| mile | miles | mi |

### Mass

| Unit | Aliases | Symbol |
|------|---------|--------|
| kilogram | kilograms | kg |
| gram | grams | g |
| milligram | milligrams | mg |
| microgram | micrograms | µg, ug |
| pound | pounds | lb, lbs |
| ounce | ounces | oz |
| stone | stones | st |
| carat | carats | ct |
| tonne | tonnes, ton, tons | - |

### Time

| Unit | Aliases | Symbol |
|------|---------|--------|
| nanosecond | nanoseconds | ns |
| microsecond | microseconds | µs, us |
| millisecond | milliseconds | ms |
| second | seconds, sec | s |
| minute | minutes | min |
| hour | hours | h, hr |
| day | days | - |
| week | weeks | - |
| fortnight | fortnights | - |
| month | months | - |
| quarter | quarters | - |
| semester | semesters | - |
| year | years | - |

### Volume

| Unit | Aliases | Symbol |
|------|---------|--------|
| litre | liter, litres, liters | l |
| millilitre | milliliter, millilitres, milliliters | ml |
| centilitre | centiliter, centilitres, centiliters | cl |
| decilitre | deciliter, decilitres, deciliters | dl |
| gallon | gallons | gal |
| US gallon | usgallon, usgallons | usgal |
| UK gallon | ukgallon, ukgallons, imperialgallon | ukgal, impgal |
| quart | quarts | qt |
| pint | pints | pt |
| cup | cups | - |
| fluid ounce | fluidounce, fluidounces | floz |
| tablespoon | tablespoons | tbsp |
| teaspoon | teaspoons | tsp |

### Temperature

| Unit | Aliases | Symbol |
|------|---------|--------|
| celsius | - | c, °c |
| fahrenheit | - | f, °f |
| kelvin | - | k |
| rankine | - | r, °r |

### Speed

| Unit | Description | Symbol |
|------|-------------|--------|
| metres per second | - | mps |
| kilometres per hour | - | kph, kmh |
| miles per hour | - | mph |
| feet per second | - | fps |
| knot | knots | kn |

### Area

| Unit | Aliases | Symbol |
|------|---------|--------|
| square metre | squaremeter, squaremetres, squaremeters | sqm, m², m2 |
| square kilometre | squarekilometer, squarekilometres, squarekilometers | sqkm, km², km2 |
| square centimetre | squarecentimeter, squarecentimetres, squarecentimeters | sqcm, cm², cm2 |
| square millimetre | squaremillimeter, squaremillimetres, squaremillimeters | sqmm, mm², mm2 |
| square foot | squarefoot, squarefeet | sqft, ft², ft2 |
| square inch | squareinch, squareinches | sqin, in², in2 |
| square yard | squareyard, squareyards | sqyd, yd², yd2 |
| square mile | squaremile, squaremiles | sqmi, mi², mi2 |
| acre | acres | - |
| hectare | hectares | ha |

### Pressure

| Unit | Aliases | Symbol |
|------|---------|--------|
| pascal | pascals | pa |
| kilopascal | kilopascals | kpa |
| megapascal | megapascals | mpa |
| bar | bars | - |
| millibar | millibars | mbar |
| atmosphere | atmospheres | atm |
| pound per square inch | - | psi |
| torr | - | - |
| millimetres of mercury | - | mmhg |
| inches of mercury | - | inhg |

### Force

| Unit | Aliases | Symbol |
|------|---------|--------|
| newton | newtons | n |
| kilonewton | kilonewtons | kn |
| meganewton | meganewtons | mn |
| pound force | poundsforce | lbf |
| kilogram force | - | kgf |
| dyne | dynes | - |

### Angle

| Unit | Aliases | Symbol |
|------|---------|--------|
| degree | degrees | deg, ° |
| radian | radians | rad |
| gradian | gradians | grad, gon |
| turn | turns, revolution, revolutions | - |

### Frequency

| Unit | Aliases | Symbol |
|------|---------|--------|
| hertz | - | hz |
| kilohertz | - | khz |
| megahertz | - | mhz |
| gigahertz | - | ghz |
| terahertz | - | thz |
| revolutions per minute | - | rpm |

### Data Storage (Bytes)

| Unit | Aliases | Symbol |
|------|---------|--------|
| byte | bytes | b |
| kilobyte | kilobytes | kb |
| megabyte | megabytes | mb |
| gigabyte | gigabytes | gb |
| terabyte | terabytes | tb |
| petabyte | petabytes | pb |

### Data Storage (Bits)

| Unit | Aliases | Symbol |
|------|---------|--------|
| bit | bits | - |
| kilobit | kilobits | kbit |
| megabit | megabits | mbit |
| gigabit | gigabits | gbit |
| terabit | terabits | tbit |
| petabit | petabits | pbit |

### Data Rate

| Unit | Description |
|------|-------------|
| bps, kbps, mbps, gbps, tbps | Bits per second |
| Bps, KBps, MBps, GBps, TBps | Bytes per second |


## Examples

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

10> 12 gbp in dollars
   = $15.24

11> 100 usd in euros
   = €90.91

12> 50 dollars + 25 euros
   = $77.50
```

Note: Currency can be written with symbols (£, $, €, ¥) before the number, or with codes/names (gbp, usd, dollars, euros, yen) after the number.

### Percentages
```
13> 30 + 20%
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

## Testing

```bash
# Run all tests
go test ./...

# With race detection
go test ./... -race

# With coverage
go test ./... -cover
```

## Releases

To create a new release:
```bash
git tag v0.1.0
git push origin v0.1.0
```

Binaries are built automatically for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)

The install script downloads the appropriate binary for your platform.

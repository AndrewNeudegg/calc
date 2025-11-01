# 🧾 Calc Specification

### Goal

A local-only, dependency-free, terminal-based “notepad calculator” inspired by Soulver — where the user types plain text mixed with maths, units, dates, and time zones, and results are shown in real-time on the right side.

---

## 1. Overview

Calc is a text-first workspace for everyday reasoning — doing quick calculations, unit conversions, time comparisons, and date manipulations — without needing to remember programming syntax or spreadsheet formulas.

Each line can contain free-form natural text that includes quantities, dates, units, currencies, or named variables.
The system interprets what it can, evaluates it, and displays a clear result inline.

The interface is minimalist and typable — think Vim meets REPL — prioritising focus, keyboard flow, and clarity.

---

## 2. Key Principles

1. **Local and deterministic** — No internet access, APIs, or external dependencies. All evaluations are local and reproducible.
2. **Typable, immediate, reversible** — Everything happens via keystrokes; edits re-evaluate instantly.
3. **Predictable natural language** — Phrases map to defined patterns (no AI). “Half of 40” → `20`. “Tomorrow + 2 weeks” → `date`.
4. **Composable** — Results can be reused as variables. `hourly = £23.50`, then `weekly = hourly * 37.5`.
5. **Transparent** — The system always shows how it interpreted your input.

---

## 3. Interface Summary

* Each line of input is one **expression**.
* Left side: user input.
* Right side: evaluated result (auto-updated).
* Lines can be edited, deleted, or inserted.
* A command mode (`:` prefix) provides meta-commands like saving, loading, or setting preferences.

### Example screen

```
1  rent per week = £450                    £450
2  yearly rent = rent per week * 52        £23,400
3  utilities = £120 per month              £120/month
4  total yearly cost = yearly rent + utilities * 12   £24,840
5  half of yearly cost                    £12,420
```

---

## 4. User Stories

### 4.1 Arithmetic & Variables

**As a user**, I can perform simple arithmetic and store results for reuse.

**Examples:**

```
3 + 4 * 5                     23
x = 10
y = x + 3                     13
monthly = £120
yearly = monthly * 12         £1,440
(10 + 5) / 2                  7.5
```

**Acceptance tests:**

* `3 + 4 * 5` → `23`
* Parentheses override precedence.
* Re-editing a variable line re-evaluates dependent lines.

---

### 4.2 Unit Conversion

**As a user**, I can convert between compatible physical units.

**Examples:**

```
10 m in cm                    1,000 cm
100 km in miles               62.14 miles
70 kg in lb                   154.32 lb
2 hours in minutes            120 minutes
```

**Acceptance tests:**

* Handles SI conversions.
* Handles mixed case units (`Km`, `M`).
* Raises error on incompatible units (`10 kg in metres`).

---

### 4.3 Custom Units

**As a user**, I can define my own units or groupings.

**Examples:**

```
1 box = 20 apples
15 boxes in apples            300 apples
300 apples in boxes           15 boxes
```

**Acceptance tests:**

* Custom units persist for session.
* Circular unit definitions are disallowed.

---

### 4.4 Rates & Compound Units

**As a user**, I can divide or multiply units to create rates.

**Examples:**

```
100 km / 2 h                  50 km/h
200 km / 4 hours in m/s       13.89 m/s
speed = 100 m / 9.58 s        10.44 m/s
```

**Acceptance tests:**

* Simplifies derived units correctly.
* Converts compound units when requested (`in m/s`).

---

### 4.5 Currency & Conversion

**As a user**, I can perform arithmetic with currencies and convert between them using a local table.

**Examples:**

```
$50 in GBP                    £39.50
£120 + $30                    £143.50
1 USD = 0.79 GBP
$10 in GBP                    £7.90
```

**Acceptance tests:**

* Reads static conversion table.
* Rounds to two decimals.
* Disallows arithmetic between currencies without conversion (if disabled).

---

### 4.6 Percentages

**As a user**, I can use percentages in intuitive ways.

**Examples:**

```
30 + 20%                      36
20% of 50                     10
20 is what % of 50            40%
increase 100 by 10%           110
decrease 100 by 10%           90
```

**Acceptance tests:**

* Recognises “of”, “increase/decrease by”.
* Shows `%` symbol on outputs.

---

### 4.7 Date and Time Arithmetic

**As a user**, I can perform arithmetic on dates and times.

**Examples:**

```
today + 3 weeks               (date 3 weeks ahead)
1 June 2025 - 3 days          29 May 2025
next monday                   3 Nov 2025
2 hours after noon            14:00
3 hours before midnight       21:00
```

**Acceptance tests:**

* Supports “today”, “tomorrow”, “yesterday”.
* Parses natural weekday names.
* Handles month/year rollovers.

---

### 4.8 Time Zone Operations

**As a user**, I can compare or convert between time zones.

**Examples:**

```
London - Singapore            8 hours
10:00 London in Singapore     18:00
now in New York               07:12
```

**Acceptance tests:**

* Time zone offsets correct for DST.
* Friendly names like “London” map to correct IANA zones.

---

### 4.9 Fuzzy English Arithmetic

**As a user**, I can type simple natural phrases.

**Examples:**

```
half of 80                    40
double 15                     30
three quarters of 200         150
10% off £120                  £108
twice 4 hours                 8 hours
```

**Acceptance tests:**

* Each phrase maps to deterministic arithmetic.
* Order of words can vary (`half of 80` = `80 halved`).

---

### 4.10 Fuzzy Time & Date

**As a user**, I can type informal time phrases.

**Examples:**

```
3 days ago                    (date 3 days before today)
in 5 weeks                    (date 5 weeks ahead)
2 hours from now              (time + 2 hours)
10am London in Tokyo          6pm
```

**Acceptance tests:**

* Phrases like “ago”, “from now”, “in N days” normalise correctly.
* Produces consistent date/time formatting.

---

### 4.11 Totals and Functions

**As a user**, I can compute totals and summaries.

**Examples:**

```
sum(10, 20, 30)               60
average(3, 4, 5)              4
mean of 4, 6, 10              6.67
total of 5 + 10 + 15          30
```

**Acceptance tests:**

* Function syntax `sum()`, `average()` works.
* “sum of …” phrase also works.

---

### 4.12 Contextual Recalculation

**As a user**, when I change a value, all dependent lines update automatically.

**Example:**

```
hourly = 20
weekly = hourly * 37.5        750
```

→ edit line 1: `hourly = 25`
→ line 2 auto-updates: `weekly = 937.5`

**Acceptance tests:**

* Recalculates dependents in topological order.
* Circular references show error.

---

### 4.13 Command Mode

**As a user**, I can issue commands to manage the session.

**Examples:**

```
:save budget.txt
:open budget.txt
:set format = uk
:set fuzzy = on
:rates load ./rates.json
:tz list
```

**Acceptance tests:**

* Each command runs synchronously and prints confirmation.
* Invalid commands show syntax help.

---

### 4.14 Formatting & Display

**As a user**, I can control output precision and formats.

**Examples:**

```
:set precision = 3
:set dateformat = "02 Jan 2006"
:set currency = GBP
```

**Acceptance tests:**

* Applies settings globally.
* Persist settings in config file.

---

## 5. Fuzzy Phrase Catalogue (baseline)

| Phrase                  | Meaning                     | Example → Result                |
| ----------------------- | --------------------------- | ------------------------------- |
| half of X               | 0.5 * X                     | half of 80 → 40                 |
| twice X / double X      | 2 * X                       | double 15 → 30                  |
| three quarters of X     | 0.75 * X                    | three quarters of 200 → 150     |
| increase X by Y%        | X * (1 + Y/100)             | increase 100 by 10% → 110       |
| decrease X by Y%        | X * (1 - Y/100)             | decrease 100 by 10% → 90        |
| X per Y                 | X / Y                       | £120 per month → £120/month     |
| in N days/weeks/etc     | today + N                   | in 3 days → 3 days later        |
| N days/weeks/etc ago    | today - N                   | 3 days ago → 3 days earlier     |
| X after Y               | Y + X                       | 2 hours after noon → 14:00      |
| X before Y              | Y - X                       | 3 hours before midnight → 21:00 |
| convert X to Y / X in Y | unit or currency conversion | 10 m in cm → 1000 cm            |
| what % of               | (x/y)*100                   | 20 is what % of 50 → 40%        |

---

## 6. Error Handling

Errors are inline, in red, and never crash the session.

| Type               | Message                            |
| ------------------ | ---------------------------------- |
| Unknown unit       | `unknown unit 'lightyearz'`        |
| Incompatible units | `cannot convert kg to metres`      |
| Circular reference | `circular dependency: a -> b -> a` |
| Syntax error       | `unexpected token 'of'`            |
| Unknown command    | `unknown command :fly`             |
| Missing timezone   | `unknown timezone 'moonbase'`      |

---

## 7. Non-Goals

* No live currency or stock data.
* No internet calls, LLMs, or AI interpretation.
* No graphical interface (TUI only).
* No spreadsheet-like cell addressing or formulas referencing line numbers.
* No persistent background daemon.

---

## 8. Future Extensions (later phases)

* `plot(expr)` inline ASCII graphing.
* Named sheets or tabs.
* `:export csv` and `:export md`.
* Persistent variable store between sessions.
* Weekday/holiday awareness (`next working day`).

---

## 9. Test Matrix Summary

Each example in this document is a **canonical test**.
A future test suite will be built from these expressions, asserting:

* Parsed tokens
* Evaluated value (numerical or datetime)
* Formatted string
* No error


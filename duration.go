package iso8601

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Period represents an amount of time in years, months, weeks and days.
// A period is not a accurate measurement of time, but rather a way to express a duration in human terms.
type Period struct {
	Years  float64
	Months float64
	Weeks  float64
	Days   float64
}

func (p Period) String() string {
	if p.Years == 0 && p.Months == 0 && p.Weeks == 0 && p.Days == 0 {
		return "P0D"
	}

	out := "P"
	if p.Years != 0 {
		out += strconv.FormatFloat(math.Abs(p.Years), 'f', -1, 64) + "Y"
	}

	if p.Months != 0 {
		out += strconv.FormatFloat(math.Abs(p.Months), 'f', -1, 64) + "M"
	}

	if p.Weeks != 0 {
		out += strconv.FormatFloat(math.Abs(p.Weeks), 'f', -1, 64) + "W"
	}

	if p.Days != 0 {
		out += strconv.FormatFloat(math.Abs(p.Days), 'f', -1, 64) + "D"
	}
	return out
}

// Parse parses from a string in iso-8601 period format. The string must start with a "P" and contain
// one or more of the following: "Y" for years, "M" for months, "W" for weeks, and "D" for days.
// The values can be floating point numbers.
// The duration part (starts with T) will be ignored
func (p *Period) Parse(s string) error {
	if len(s) < 2 || s[0] != 'P' {
		return fmt.Errorf("invalid period format: must start with P")
	}

	// Reset the period values
	p.Years = 0
	p.Months = 0
	p.Weeks = 0
	p.Days = 0

	// Remove the duration part if present (everything after T)
	if idx := strings.Index(s, "T"); idx != -1 {
		s = s[:idx]
	}

	// Handle special case P0D
	if s == "P0D" {
		return nil
	}

	// Current number being parsed
	var num string
	var negative bool

	for i := 1; i < len(s); i++ {
		c := s[i]

		switch {
		case c == '-':
			negative = true
		case c == '+':
			negative = false
		case c >= '0' && c <= '9' || c == '.':
			num += string(c)
		case c == 'Y' || c == 'M' || c == 'W' || c == 'D':
			if num == "" {
				continue
			}

			val, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return fmt.Errorf("invalid number format: %s", num)
			}

			if negative {
				val = -val
			}

			switch c {
			case 'Y':
				p.Years = val
			case 'M':
				p.Months = val
			case 'W':
				p.Weeks = val
			case 'D':
				p.Days = val
			}

			// Reset for next number
			num = ""
			negative = false
		default:
			return fmt.Errorf("invalid character in period: %c", c)
		}
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (p Period) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", p.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *Period) UnmarshalJSON(data []byte) error {
	// Remove quotes from the string
	s := string(data)
	if len(s) < 2 {
		return fmt.Errorf("invalid JSON string length")
	}
	s = s[1 : len(s)-1] // Remove surrounding quotes

	return p.Parse(s)
}

// ParsePeriod parses a string in iso-8601 period format.
func ParsePeriod(s string) (*Period, error) {
	p := &Period{}
	if err := p.Parse(s); err != nil {
		return nil, err
	}
	return p, nil
}

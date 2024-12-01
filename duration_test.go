package iso8601

import (
	"encoding/json"
	"testing"
)

func TestPeriodString(t *testing.T) {
	tests := []struct {
		name     string
		period   Period
		expected string
	}{
		{
			name:     "zero period",
			period:   Period{},
			expected: "P0D",
		},
		{
			name:     "full period",
			period:   Period{Years: 1, Months: 2, Weeks: 3, Days: 4},
			expected: "P1Y2M3W4D",
		},
		{
			name:     "negative values",
			period:   Period{Years: -1, Months: -2, Weeks: -3, Days: -4},
			expected: "P1Y2M3W4D",
		},
		{
			name:     "decimal values",
			period:   Period{Years: 1.5, Months: 2.5},
			expected: "P1.5Y2.5M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.period.String(); got != tt.expected {
				t.Errorf("Period.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPeriodParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Period
		wantErr bool
	}{
		{
			name:    "basic period",
			input:   "P1Y2M3W4D",
			want:    Period{Years: 1, Months: 2, Weeks: 3, Days: 4},
			wantErr: false,
		},
		{
			name:    "zero period",
			input:   "P0D",
			want:    Period{},
			wantErr: false,
		},
		{
			name:    "decimal values",
			input:   "P1.5Y2.5M",
			want:    Period{Years: 1.5, Months: 2.5},
			wantErr: false,
		},
		{
			name:    "negative values",
			input:   "P-1Y-2M-3W-4D",
			want:    Period{Years: -1, Months: -2, Weeks: -3, Days: -4},
			wantErr: false,
		},
		{
			name:    "with time component",
			input:   "P1YT2H",
			want:    Period{Years: 1},
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "1Y2M",
			wantErr: true,
		},
		{
			name:    "invalid character",
			input:   "P1Y2X",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Period
			err := p.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Period.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (p != tt.want) {
				t.Errorf("Period.Parse() = %v, want %v", p, tt.want)
			}
		})
	}
}

func TestPeriodJSON(t *testing.T) {
	tests := []struct {
		name     string
		period   Period
		expected string
	}{
		{
			name:     "zero period",
			period:   Period{},
			expected: `"P0D"`,
		},
		{
			name:     "full period",
			period:   Period{Years: 1, Months: 2, Weeks: 3, Days: 4},
			expected: `"P1Y2M3W4D"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			got, err := json.Marshal(tt.period)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			if string(got) != tt.expected {
				t.Errorf("json.Marshal() = %v, want %v", string(got), tt.expected)
			}

			// Test unmarshaling
			var p Period
			err = json.Unmarshal(got, &p)
			if err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}
			if p != tt.period {
				t.Errorf("json.Unmarshal() = %v, want %v", p, tt.period)
			}
		})
	}
}

func TestParsePeriod(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Period
		wantErr bool
	}{
		{
			name:    "valid period",
			input:   "P1Y2M3W4D",
			want:    &Period{Years: 1, Months: 2, Weeks: 3, Days: 4},
			wantErr: false,
		},
		{
			name:    "invalid period",
			input:   "invalid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePeriod(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && *got != *tt.want {
				t.Errorf("ParsePeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

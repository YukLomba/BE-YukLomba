package common

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type Datetime time.Time

func (d Datetime) MarshalJSON() ([]byte, error) {
	// Format the time as "2006-01-02"
	t := time.Time(d)
	formatted := t.Format("2006-01-02")

	// Marshal the formatted string into JSON
	return json.Marshal(formatted)
}

func (d *Datetime) UnmarshalJSON(b []byte) error {
	// Unmarshal the JSON string to a Go string
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal to a string: %w", err)
	}

	// Parse the date string using the "2006-01-02" layout
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	date := t.Format("2006-01-02")
	slog.Info("Parsed time:", "date", date)

	// Assign the parsed time to *d
	*d = Datetime(t)
	return nil
}
func (d Datetime) Validate() error {
	t := time.Time(d)
	if t.IsZero() {
		return fmt.Errorf("date must be set")
	}
	// Optional: Check for future/past constraints
	// if t.Before(time.Now()) {
	//     return fmt.Errorf("date must not be in the past")
	// }
	return nil
}

func (d Datetime) ToTime() time.Time {
	return time.Time(d)
}

package date

import (
	"encoding/json"
	"fmt"
	"time"
)

// Date is a JSON serialisable time
type Date struct{ time.Time }

// UnmarshalJSON converts a json value to a Date
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("date should be a string, got %s", data)
	}
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", s)
	if err != nil {
		return fmt.Errorf("invalid date: %v", err)
	}
	d.Time = t
	return nil
}

// MarshalJSON Converts a Date json value
func (d *Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d.Time).Format("2006-01-02 15:04:05.999999999 -0700 MST")
	return json.Marshal(t)
}

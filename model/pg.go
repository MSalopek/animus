package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// PgTime represents time values as they are stored in postgres.
type PgTime struct {
	Infinity    bool
	NegInfinity bool
	Time        time.Time
}

func (t PgTime) Value() (driver.Value, error) {
	if t.Infinity {
		return "infinity", nil
	}

	if t.NegInfinity {
		return "-infinity", nil
	}

	return t.Time, nil
}

func (t *PgTime) Scan(src interface{}) error {
	switch d := src.(type) {
	case time.Time:
		t.Time = d
		return nil
	case string:
		switch strings.ToLower(d) {
		case "infinity":
			t.Infinity = true
			return nil
		case "-infinity":
			t.NegInfinity = true
			return nil
		default:
			return fmt.Errorf("unable to parse %v as time", d)
		}
	default:
		return fmt.Errorf("unable to parse %v as time", src)
	}
}

package dirty

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type Update struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

type Updates map[string]Update

func NewUpdates() Updates {
	return make(Updates)
}

func (u Updates) DiffBool(name string, field bool, value Bool) bool {
	if !value.Dirty || !value.Valid || field == value.Bool.Bool {
		return field
	}
	u[name] = Update{From: field, To: value.Bool.Bool}
	return value.Bool.Bool
}

func (u Updates) DiffNullBool(name string, field null.Bool, value Bool) null.Bool {
	if !value.Dirty || field == value.Bool {
		return field
	}
	u[name] = Update{From: field, To: value.Bool}
	return value.Bool
}

func (u Updates) DiffFloat(name string, field float64, value Float) float64 {
	if !value.Dirty || !value.Valid || field == value.Float.Float64 {
		return field
	}
	u[name] = Update{From: field, To: value.Float.Float64}
	return value.Float.Float64
}

func (u Updates) DiffNullFloat(name string, field null.Float, value Float) null.Float {
	if !value.Dirty || field == value.Float {
		return field
	}
	u[name] = Update{From: field, To: value.Float}
	return value.Float
}

func (u Updates) DiffInt(name string, field int64, value Int) int64 {
	if !value.Dirty || !value.Valid || field == value.Int.Int64 {
		return field
	}
	u[name] = Update{From: field, To: value.Int.Int64}
	return value.Int.Int64
}

func (u Updates) DiffNullInt(name string, field null.Int, value Int) null.Int {
	if !value.Dirty || field == value.Int {
		return field
	}
	u[name] = Update{From: field, To: value.Int}
	return value.Int
}

func (u Updates) DiffString(name string, field string, value String) string {
	if !value.Dirty || !value.Valid || field == value.String.String {
		return field
	}
	u[name] = Update{From: field, To: value.String.String}
	return value.String.String
}

func (u Updates) DiffNullString(name string, field null.String, value String) null.String {
	if !value.Dirty || field == value.String {
		return field
	}
	u[name] = Update{From: field, To: value.String}
	return value.String
}

func (u Updates) DiffTime(name string, field time.Time, value Time) time.Time {
	if !value.Dirty || !value.Valid || field.Equal(value.Time.Time) {
		return field
	}
	u[name] = Update{From: field, To: value.Time.Time}
	return value.Time.Time
}

func (u Updates) DiffNullTime(name string, field null.Time, value Time) null.Time {
	if !value.Dirty || field.Time.Equal(value.Time.Time) {
		return field
	}
	u[name] = Update{From: field, To: value.Time}
	return value.Time
}

func (u Updates) DiffNullTimeBool(name string, field null.Time, value Bool) null.Time {
	if !value.Dirty || !value.Valid {
		return field
	}
	updated := field
	if value.ValueOrZero() {
		if !field.Valid {
			updated = null.TimeFrom(time.Now())
		}
	} else {
		updated = null.Time{}
	}
	if updated.Valid == field.Valid {
		return field
	}
	u[name] = Update{From: field, To: updated}
	return updated
}

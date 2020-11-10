// Package dirty tracks when values have been unmarshalled in order to differentiate between values that are absent and values that are null.
package dirty

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type Bool struct {
	null.Bool
	Dirty bool
}

func BoolFrom(b bool) Bool {
	return Bool{
		Bool:  null.BoolFrom(b),
		Dirty: true,
	}
}

func (b *Bool) Scan(value interface{}) error {
	b.Dirty = true
	return b.Bool.Scan(value)
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	b.Dirty = true
	return b.Bool.UnmarshalJSON(data)
}

func (b *Bool) UnmarshalText(text []byte) error {
	b.Dirty = true
	return b.Bool.UnmarshalText(text)
}

type Float struct {
	null.Float
	Dirty bool
}

func FloatFrom(f float64) Float {
	return Float{
		Float: null.FloatFrom(f),
		Dirty: true,
	}
}

func (f *Float) Scan(value interface{}) error {
	f.Dirty = true
	return f.Float.Scan(value)
}

func (f *Float) UnmarshalJSON(data []byte) error {
	f.Dirty = true
	return f.Float.UnmarshalJSON(data)
}

func (f *Float) UnmarshalText(text []byte) error {
	f.Dirty = true
	return f.Float.UnmarshalText(text)
}

type Int struct {
	null.Int
	Dirty bool
}

func IntFrom(i int64) Int {
	return Int{
		Int:   null.IntFrom(i),
		Dirty: true,
	}
}

func (i *Int) Scan(value interface{}) error {
	i.Dirty = true
	return i.Int.Scan(value)
}

func (i *Int) UnmarshalJSON(data []byte) error {
	i.Dirty = true
	return i.Int.UnmarshalJSON(data)
}

func (i *Int) UnmarshalText(text []byte) error {
	i.Dirty = true
	return i.Int.UnmarshalText(text)
}

type String struct {
	null.String
	Dirty bool
}

func StringFrom(s string) String {
	return String{
		String: null.StringFrom(s),
		Dirty:  true,
	}
}

func (s *String) Scan(value interface{}) error {
	s.Dirty = true
	return s.String.Scan(value)
}

func (s *String) UnmarshalJSON(data []byte) error {
	s.Dirty = true
	return s.String.UnmarshalJSON(data)
}

func (s *String) UnmarshalText(text []byte) error {
	s.Dirty = true
	return s.String.UnmarshalText(text)
}

type Time struct {
	null.Time
	Dirty bool
}

func TimeFrom(t time.Time) Time {
	return Time{
		Time:  null.TimeFrom(t),
		Dirty: true,
	}
}

func (t *Time) Scan(value interface{}) error {
	t.Dirty = true
	return t.Time.Scan(value)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	t.Dirty = true
	return t.Time.UnmarshalJSON(data)
}

func (t *Time) UnmarshalText(text []byte) error {
	t.Dirty = true
	return t.Time.UnmarshalText(text)
}

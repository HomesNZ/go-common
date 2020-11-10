package dirty

import (
	"encoding/json"
	"testing"
)

func TestBoolUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		Name  string
		JSON  string
		Err   error
		Dirty bool
		Valid bool
	}{
		{Name: "not present", JSON: `{}`, Err: nil, Dirty: false, Valid: false},
		{Name: "present", JSON: `{"Foo":true}`, Err: nil, Dirty: true, Valid: true},
		{Name: "present null", JSON: `{"Foo":null}`, Err: nil, Dirty: true, Valid: false},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var v struct {
				Foo Bool
			}
			err := json.Unmarshal([]byte(tc.JSON), &v)
			if err != tc.Err {
				t.Fatalf("err == %v, want %v", err, tc.Err)
			}
			if v.Foo.Dirty != tc.Dirty {
				t.Errorf("v.Foo.Dirty == %v, want %v", v.Foo.Dirty, tc.Dirty)
			}
			if v.Foo.Valid != tc.Valid {
				t.Errorf("v.Foo.Valid == %v, want %v", v.Foo.Valid, tc.Valid)
			}
		})
	}
}

func TestFloatUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		Name  string
		JSON  string
		Err   error
		Dirty bool
		Valid bool
	}{
		{Name: "not present", JSON: `{}`, Err: nil, Dirty: false, Valid: false},
		{Name: "present", JSON: `{"Foo":1.2345}`, Err: nil, Dirty: true, Valid: true},
		{Name: "present null", JSON: `{"Foo":null}`, Err: nil, Dirty: true, Valid: false},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var v struct {
				Foo Float
			}
			err := json.Unmarshal([]byte(tc.JSON), &v)
			if err != tc.Err {
				t.Fatalf("err == %v, want %v", err, tc.Err)
			}
			if v.Foo.Dirty != tc.Dirty {
				t.Errorf("v.Foo.Dirty == %v, want %v", v.Foo.Dirty, tc.Dirty)
			}
			if v.Foo.Valid != tc.Valid {
				t.Errorf("v.Foo.Valid == %v, want %v", v.Foo.Valid, tc.Valid)
			}
		})
	}
}

func TestIntUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		Name  string
		JSON  string
		Err   error
		Dirty bool
		Valid bool
	}{
		{Name: "not present", JSON: `{}`, Err: nil, Dirty: false, Valid: false},
		{Name: "present", JSON: `{"Foo":123}`, Err: nil, Dirty: true, Valid: true},
		{Name: "present null", JSON: `{"Foo":null}`, Err: nil, Dirty: true, Valid: false},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var v struct {
				Foo Int
			}
			err := json.Unmarshal([]byte(tc.JSON), &v)
			if err != tc.Err {
				t.Fatalf("err == %v, want %v", err, tc.Err)
			}
			if v.Foo.Dirty != tc.Dirty {
				t.Errorf("v.Foo.Dirty == %v, want %v", v.Foo.Dirty, tc.Dirty)
			}
			if v.Foo.Valid != tc.Valid {
				t.Errorf("v.Foo.Valid == %v, want %v", v.Foo.Valid, tc.Valid)
			}
		})
	}
}

func TestStringUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		Name  string
		JSON  string
		Err   error
		Dirty bool
		Valid bool
	}{
		{Name: "not present", JSON: `{}`, Err: nil, Dirty: false, Valid: false},
		{Name: "present", JSON: `{"Foo":"Bar"}`, Err: nil, Dirty: true, Valid: true},
		{Name: "present null", JSON: `{"Foo":null}`, Err: nil, Dirty: true, Valid: false},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var v struct {
				Foo String
			}
			err := json.Unmarshal([]byte(tc.JSON), &v)
			if err != tc.Err {
				t.Fatalf("err == %v, want %v", err, tc.Err)
			}
			if v.Foo.Dirty != tc.Dirty {
				t.Errorf("v.Foo.Dirty == %v, want %v", v.Foo.Dirty, tc.Dirty)
			}
			if v.Foo.Valid != tc.Valid {
				t.Errorf("v.Foo.Valid == %v, want %v", v.Foo.Valid, tc.Valid)
			}
		})
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		Name  string
		JSON  string
		Err   error
		Dirty bool
		Valid bool
	}{
		{Name: "not present", JSON: `{}`, Err: nil, Dirty: false, Valid: false},
		{Name: "present", JSON: `{"Foo":"2012-12-21T21:21:21Z"}`, Err: nil, Dirty: true, Valid: true},
		{Name: "present null", JSON: `{"Foo":null}`, Err: nil, Dirty: true, Valid: false},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			var v struct {
				Foo Time
			}
			err := json.Unmarshal([]byte(tc.JSON), &v)
			if err != tc.Err {
				t.Fatalf("err == %v, want %v", err, tc.Err)
			}
			if v.Foo.Dirty != tc.Dirty {
				t.Errorf("v.Foo.Dirty == %v, want %v", v.Foo.Dirty, tc.Dirty)
			}
			if v.Foo.Valid != tc.Valid {
				t.Errorf("v.Foo.Valid == %v, want %v", v.Foo.Valid, tc.Valid)
			}
		})
	}
}

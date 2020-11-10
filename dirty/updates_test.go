package dirty

import (
	"reflect"
	"testing"
	"time"

	"gopkg.in/guregu/null.v3"
)

func TestUpdates_DiffString(t *testing.T) {
	model := struct {
		DirtySame string
		DirtyDiff string
		Clean     string
	}{
		DirtySame: "apple",
		DirtyDiff: "apple",
		Clean:     "apple",
	}
	request := struct {
		DirtySame String
		DirtyDiff String
		Clean     String
	}{
		DirtySame: String{String: null.StringFrom("apple"), Dirty: true},
		DirtyDiff: String{String: null.StringFrom("banana"), Dirty: true},
		Clean:     String{String: null.StringFrom("banana"), Dirty: false},
	}
	expected := Updates{
		"dirty_diff": Update{
			From: "apple",
			To:   "banana",
		},
	}

	u := NewUpdates()
	model.DirtySame = u.DiffString("dirty_same", model.DirtySame, request.DirtySame)
	model.DirtyDiff = u.DiffString("dirty_diff", model.DirtyDiff, request.DirtyDiff)
	model.Clean = u.DiffString("clean", model.Clean, request.Clean)
	if !reflect.DeepEqual(u, expected) {
		t.Errorf("u == %v, want %v", u, expected)
	}
	if model.DirtyDiff != request.DirtyDiff.String.String {
		t.Errorf("model.DirtyDiff = %v, want %v", model.DirtyDiff, request.DirtyDiff.String.String)
	}
}

func TestUpdates_DiffNullString(t *testing.T) {
	model := struct {
		DirtySame null.String
		DirtyDiff null.String
		Clean     null.String
	}{
		DirtySame: null.StringFrom("apple"),
		DirtyDiff: null.StringFrom("apple"),
		Clean:     null.StringFrom("apple"),
	}
	request := struct {
		DirtySame String
		DirtyDiff String
		Clean     String
	}{
		DirtySame: String{String: null.StringFrom("apple"), Dirty: true},
		DirtyDiff: String{String: null.StringFrom("banana"), Dirty: true},
		Clean:     String{String: null.StringFrom("banana"), Dirty: false},
	}
	expected := Updates{
		"dirty_diff": Update{
			From: null.StringFrom("apple"),
			To:   null.StringFrom("banana"),
		},
	}

	u := NewUpdates()
	model.DirtySame = u.DiffNullString("dirty_same", model.DirtySame, request.DirtySame)
	model.DirtyDiff = u.DiffNullString("dirty_diff", model.DirtyDiff, request.DirtyDiff)
	model.Clean = u.DiffNullString("clean", model.Clean, request.Clean)
	if !reflect.DeepEqual(u, expected) {
		t.Errorf("u == %v, want %v", u, expected)
	}
	if model.DirtyDiff != request.DirtyDiff.String {
		t.Errorf("model.DirtyDiff = %v, want %v", model.DirtyDiff, request.DirtyDiff.String)
	}
}

func TestUpdates_DiffNullTimeBool(t *testing.T) {
	model := struct {
		DirtyUnchanged null.Time
		DirtyUnset     null.Time
		DirtySet       null.Time
		Clean          null.Time
		CleanNull      null.Time
	}{
		DirtyUnchanged: null.TimeFrom(time.Now()),
		DirtyUnset:     null.TimeFrom(time.Now()),
		DirtySet:       null.Time{},
		Clean:          null.TimeFrom(time.Now()),
		CleanNull:      null.TimeFrom(time.Now()),
	}
	request := struct {
		DirtyUnchanged Bool
		DirtyUnset     Bool
		DirtySet       Bool
		Clean          Bool
		CleanNull      Bool
	}{
		DirtyUnchanged: Bool{Bool: null.BoolFrom(true), Dirty: true},
		DirtyUnset:     Bool{Bool: null.BoolFrom(false), Dirty: true},
		DirtySet:       Bool{Bool: null.BoolFrom(true), Dirty: true},
		Clean:          Bool{Bool: null.Bool{}, Dirty: false},
		CleanNull:      Bool{Bool: null.Bool{}, Dirty: true},
	}

	u := NewUpdates()
	model.DirtyUnchanged = u.DiffNullTimeBool("DirtyUnchanged", model.DirtyUnchanged, request.DirtyUnchanged)
	model.DirtyUnset = u.DiffNullTimeBool("DirtyUnset", model.DirtyUnset, request.DirtyUnset)
	model.DirtySet = u.DiffNullTimeBool("DirtySet", model.DirtySet, request.DirtySet)
	model.Clean = u.DiffNullTimeBool("Clean", model.Clean, request.Clean)
	model.CleanNull = u.DiffNullTimeBool("CleanNull", model.CleanNull, request.CleanNull)

	if !model.DirtyUnchanged.Valid {
		t.Errorf("model.DirtyUnchanged.Valid = %v, want %v", model.DirtyUnchanged.Valid, true)
	}
	if model.DirtyUnset.Valid {
		t.Errorf("model.DirtyUnset.Valid = %v, want %v", model.DirtyUnset.Valid, false)
	}
	if !model.DirtySet.Valid {
		t.Errorf("model.DirtySet.Valid = %v, want %v", model.DirtySet.Valid, true)
	}
	if !model.Clean.Valid {
		t.Errorf("model.Clean.Valid = %v, want %v", model.Clean.Valid, true)
	}
	if !model.CleanNull.Valid {
		t.Errorf("model.CleanNull.Valid = %v, want %v", model.CleanNull.Valid, true)
	}
	if _, ok := u["DirtyUnchanged"]; ok {
		t.Errorf("DirtyUnchanged should not be updated")
	}
	if _, ok := u["DirtyUnset"]; !ok {
		t.Errorf("DirtyUnset should be updated")
	}
	if _, ok := u["DirtySet"]; !ok {
		t.Errorf("DirtySet should be updated")
	}
	if _, ok := u["Clean"]; ok {
		t.Errorf("Clean should not be updated")
	}
	if _, ok := u["CleanNull"]; ok {
		t.Errorf("CleanNull should not be updated")
	}
}

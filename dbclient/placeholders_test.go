package dbclient

import (
	"reflect"
	"testing"
)

func Test_generatePlaceholders(t *testing.T) {
	type args struct {
		args [][]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generates valid placeholders for 0 length slice",
			args: args{
				args: [][]interface{}{},
			},
			want: "",
		},
		{
			name: "generates valid placeholders for 1 length slice without internal commas",
			args: args{
				args: [][]interface{}{
					{struct{ Foo string }{"a"}},
				},
			},
			want: "($1)",
		},
		{
			name: "generates valid placeholders for 3 length slice without internal commas",
			args: args{
				args: [][]interface{}{
					{struct{ Foo string }{"a"}},
					{struct{ Foo string }{"b"}},
					{struct{ Foo string }{"c"}},
				},
			},
			want: "($1),($2),($3)",
		},
		{
			name: "generates valid placeholders for 3 length slice",
			args: args{
				args: [][]interface{}{
					{struct{ Foo string }{"a"}, struct{ Foo string }{"A"}},
					{struct{ Foo string }{"b"}, struct{ Foo string }{"B"}},
					{struct{ Foo string }{"c"}, struct{ Foo string }{"C"}},
				},
			},
			want: "($1,$2),($3,$4),($5,$6)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generatePlaceholders(tt.args.args); got != tt.want {
				t.Errorf("generatePlaceholders() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}

func Test_extractArgs(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		m    map[string]bool
		want [][]interface{}
	}{
		{
			name: "generates valid args for 0 length slice",
			args: args{
				args: []interface{}{},
			},
			want: [][]interface{}(nil),
		},
		{
			name: "generates valid args for 1 length slice",
			args: args{
				args: []interface{}{
					struct{ Foo string }{"a"},
				},
			},
			want: [][]interface{}{
				{"a"},
			},
		},
		{
			name: "generates valid placeholders for 3 slices of 1 length",
			args: args{
				args: []interface{}{
					struct{ Foo string }{"a"},
					struct{ Foo string }{"b"},
					struct{ Foo string }{"c"},
				},
			},
			want: [][]interface{}{
				{"a"},
				{"b"},
				{"c"},
			},
		},
		{
			name: "generates valid placeholders for 3 slices of 2 length",
			args: args{
				args: []interface{}{
					struct{ Foo, Bar string }{"a", "A"},
					struct{ Foo, Bar string }{"b", "B"},
					struct{ Foo, Bar string }{"c", "C"},
				},
			},
			want: [][]interface{}{
				{"a", "A"},
				{"b", "B"},
				{"c", "C"},
			},
		},
		{
			name: "generates valid placeholders for 3 slices of 2 length, ignoring private fields",
			args: args{
				args: []interface{}{
					struct{ Foo, Bar, baz string }{"a", "A", "e"},
					struct{ Foo, Bar, baz string }{"b", "B", "e"},
					struct{ Foo, Bar, baz string }{"c", "C", "e"},
				},
			},
			want: [][]interface{}{
				{"a", "A"},
				{"b", "B"},
				{"c", "C"},
			},
		},
		{
			name: "generates valid placeholders for 3 slices of 2 length, ignoring internal private fields",
			args: args{
				args: []interface{}{
					struct{ Foo, baz, Bar string }{"a", "e", "A"},
					struct{ Foo, baz, Bar string }{"b", "e", "B"},
					struct{ Foo, baz, Bar string }{"c", "e", "C"},
				},
			},
			want: [][]interface{}{
				{"a", "A"},
				{"b", "B"},
				{"c", "C"},
			},
		},
		{
			name: "excludes one field and generates valid placeholders for 3 slices of 1 length",
			args: args{
				args: []interface{}{
					struct{ Foo, Bar string }{"a", "A"},
					struct{ Foo, Bar string }{"b", "B"},
					struct{ Foo, Bar string }{"c", "C"},
				},
			},
			m: map[string]bool{
				"foo": true,
			},
			want: [][]interface{}{
				{"A"},
				{"B"},
				{"C"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractArgs(tt.args.args, tt.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractArgs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

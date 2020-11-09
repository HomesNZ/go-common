package dbclient

import (
	"reflect"
	"strconv"
	"strings"
)

const SQL_MAX_PLACEHOLDERS = 65535

type Placeholder struct {
	Placeholders string
	Args         []interface{}
}

// Placeholders converts a given slice of interfaces into a set of postgresql insertion instructions
func Placeholders(rawArgs []interface{}) []Placeholder {
	if len(rawArgs) == 0 {
		return nil
	}

	var placeholders []Placeholder
	fields := reflect.ValueOf(rawArgs[0]).NumField()
	batchSize := SQL_MAX_PLACEHOLDERS / fields
	for i := 0; i < len(rawArgs); i += batchSize {
		to := i + batchSize
		if len(rawArgs) < to {
			to = len(rawArgs) - 1
		}

		structArgs := extractArgs(rawArgs[i:to])
		placeholders = append(placeholders, Placeholder{
			Placeholders: generatePlaceholders(structArgs),
			Args:         flattenArgs(structArgs),
		})
	}

	return placeholders
}

func flattenArgs(args [][]interface{}) []interface{} {
	var flattened []interface{}
	for _, arr := range args {
		for _, v := range arr {
			flattened = append(flattened, v)
		}
	}
	return flattened
}

// extractArgs converts a slice of structs to a slice of slices containing the public fields of the given structs
func extractArgs(args []interface{}) [][]interface{} {
	var sqlArgs [][]interface{}

	for _, arg := range args {
		var fields []interface{}
		v := reflect.ValueOf(arg)
		for i := 0; i < v.NumField(); i++ {
			a := v.Field(i)
			if !a.CanInterface() {
				continue
			}
			fields = append(fields, a.Interface())
		}

		sqlArgs = append(sqlArgs, fields)
	}

	return sqlArgs
}

func generatePlaceholders(args [][]interface{}) string {
	b := strings.Builder{}
	total := 1
	for i, arg := range args {
		b.WriteString("(")
		var placeholders []string

		for j := range arg {
			placeholders = append(placeholders, "$"+strconv.Itoa(total+j))
		}

		b.WriteString(strings.Join(placeholders, ","))
		b.WriteString(")")
		total += len(arg)
		// if not last element, add ','
		if i+1 < len(args) {
			b.WriteString(",")
		}
	}

	return b.String()
}

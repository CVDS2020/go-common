package log

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
)

func yamlTypeError(value *yaml.Node, typ reflect.Type, err error) error {
	v := value.Value
	if value.Tag != "!!seq" && value.Tag != "!!map" {
		if len(v) > 10 {
			v = " `" + v[:7] + "...`"
		} else {
			v = " `" + v + "`"
		}
	}
	return &yaml.TypeError{Errors: []string{
		fmt.Sprintf("line %d: cannot unmarshal %s%s into %s, cause: %s", value.Line, value.Tag, v, typ, err.Error()),
	}}
}

package config

import (
	"encoding/json"
	"gitee.com/sy_183/common/id"
	"gopkg.in/yaml.v3"
	"strings"
)

type Type struct {
	Id          uint64
	Name        string
	Suffixes    []string
	Unmarshaler func([]byte, interface{}) error
}

var typeIdCxt uint64

func NewType(Name string, Suffixes []string, Unmarshaler func([]byte, interface{}) error) Type {
	return Type{
		Id:          id.Uint64Id(&typeIdCxt),
		Name:        Name,
		Suffixes:    Suffixes,
		Unmarshaler: Unmarshaler,
	}
}

var (
	TypeUnknown = Type{Id: id.Uint64Id(&typeIdCxt), Name: "unknown"}
	TypeYaml    = Type{Id: id.Uint64Id(&typeIdCxt), Name: "yaml", Suffixes: []string{"yaml", "yml"}, Unmarshaler: yaml.Unmarshal}
	TypeJson    = Type{Id: id.Uint64Id(&typeIdCxt), Name: "json", Suffixes: []string{"json"}, Unmarshaler: json.Unmarshal}
)

var types = []*Type{&TypeYaml, &TypeJson}

func ProbeType(path string) Type {
	for _, tpy := range types {
		for _, suffix := range tpy.Suffixes {
			if strings.HasSuffix(path, "."+suffix) {
				return *tpy
			}
		}
	}
	return TypeUnknown
}

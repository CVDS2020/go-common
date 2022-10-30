package config

import (
	"testing"
)

type A struct {
	S  string  `default:"hello"`
	I  int     `default:"-1"`
	U  uint    `default:"1"`
	F  float64 `default:"3.14"`
	B  bool    `default:"true"`
	SP *string
	M  map[string]any  `default:"{username: admin, password: admin@123}"`
	SS []string        `default:"[192.168.1.11, 192.168.1.12]"`
	SA [3]string       `default:"[test, user1, user2]"`
	MP *map[string]any `default:"{}"`
	ST struct {
		S  string         `yaml:"s" default:"ttt"`
		SP *string        `yaml:"sp"`
		M  map[string]any `yaml:"m"`
		SS []string       `yaml:"ss"`
	} `default:"{s: tt, sp: ll, m: , ss: [1, 2, 3]}"`
	STP *struct {
		S  string
		SP *string `default:"ming"`
		M  map[string]any
		SS []string
	} `default:""`
}

func TestHandleA(t *testing.T) {
	a := new(A)
	err := Handle(a)
	if err != nil {
		t.Fatal(err)
	}
}

type C struct {
	S  string
	SP *string
	M  map[string]any
	Mp *map[string]any
}

func (c C) PreModify() (nc any, modified bool) {
	return C{S: "123"}, true
}

type D struct {
	S  string
	SP *string
	M  map[string]any
	Mp *map[string]any
}

func (d *D) PreModify() (nc any, modified bool) {
	d.S = "123"
	return d, true
}

func (d *D) PreHandle() {
	s := "12345"
	d.SP = &s
}

type B struct {
	S  string
	SP *string
	M  map[string]C  `default:"{a: , b: }"`
	Mp *map[string]D `default:"{a: , b: }"`
	C  C
	CP *C `default:""`
	D  D
	DP *D `default:""`
}

func TestHandleB(t *testing.T) {
	b := new(B)
	err := Handle(b)
	if err != nil {
		t.Fatal(err)
	}
}

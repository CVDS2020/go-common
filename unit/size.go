package unit

import (
	"bytes"
	"fmt"
	"gitee.com/sy_183/common/uns"
	"gopkg.in/yaml.v3"
	"reflect"
	"strconv"
	"strings"
)

type UnknownSizeUnitError struct {
	Unit string
}

func (e *UnknownSizeUnitError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("unknown size unit '%s'", e.Unit)
}

const (
	Bit      = 0.125
	BitDiv   = 8
	Byte     = 1
	KiloByte = 1000
	KiBiByte = 1 << 10
	KiloBit  = KiloByte / 8
	KiBiBit  = KiBiByte / 8
	MegaByte = 1000 * 1000
	MeBiByte = 1 << 20
	MegaBit  = MegaByte / 8
	MeBiBit  = MeBiByte / 8
	GigaByte = 1000 * 1000 * 1000
	GiBiByte = 1 << 30
	GigaBit  = GigaByte / 8
	GiBiBit  = GiBiByte / 8
	TeraByte = 1000 * 1000 * 1000 * 1000
	TeBiByte = 1 << 40
	TeraBit  = TeraByte / 8
	TeBiBit  = TeBiByte / 8
	PetaByte = 1000 * 1000 * 1000 * 1000 * 1000
	PeBiByte = 1 << 50
	PetaBit  = PetaByte / 8
	PeBiBit  = PeBiByte / 8
	ExaByte  = 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	EbiByte  = 1 << 60
	ExaBit   = ExaByte / 8
	EbiBit   = EbiByte / 8
)

type Size uint64

func (s Size) Uint64() uint64 {
	return uint64(s)
}

func (s Size) Uint() uint {
	return uint(s)
}

func (s Size) Int64() int64 {
	return int64(s)
}

func (s Size) Int() int {
	return int(s)
}

func (s Size) Float64() float64 {
	return float64(s)
}

func (s *Size) UnmarshalText(text []byte) error {
	text = bytes.TrimSpace(text)
	if len(text) == 0 {
		return nil
	}
	ns := uns.BytesToString(text)
	nsu := strings.ToUpper(ns)
	us, usu := "", ""
	i := strings.IndexAny(nsu, "BKMGTPE")
	if i >= 0 {
		us = strings.TrimSpace(ns[i:])
		usu = strings.TrimSpace(nsu[i:])
		ns = strings.TrimSpace(nsu[:i])
	}
	mul := uint64(1)
	div := uint64(1)
	if usu != "" {
		switch usu[0] {
		case 'B':
			switch usu[1:] {
			case "", "YTE":
			case "IT":
				div = 8
			default:
				return &UnknownSizeUnitError{Unit: us}
			}
		case 'K':
			switch usu[1:] {
			case "", "B", "BYTE", "ILOBYTE":
				mul = 1000
			case "IB", "IBYTE", "IBIBYTE":
				mul = 1 << 10
			case "BIT", "ILOBIT":
				mul, div = 1000, 8
			case "IBIT", "IBIBIT":
				mul, div = 1<<10, 8
			default:
				return &UnknownSizeUnitError{Unit: us}
			}
		case 'M':
			switch usu[1:] {
			case "", "B", "BYTE", "EGABYTE":
				mul = 1000 * 1000
			case "IB", "IBYTE", "EBIBYTE":
				mul = 1 << 20
			case "BIT", "EGABIT":
				mul, div = 1000*1000, 8
			case "IBIT", "EBIBIT":
				mul, div = 1<<20, 8
			default:
				return &UnknownSizeUnitError{Unit: us}
			}
		case 'G':
			switch usu[1:] {
			case "", "B", "BYTE", "IGABYTE":
				mul = 1000 * 1000 * 1000
			case "IB", "IBYTE", "IBIBYTE":
				mul = 1 << 30
			case "BIT", "IGABIT":
				mul, div = 1000*1000*1000, 8
			case "IBIT", "IBIBIT":
				mul, div = 1<<20, 8
			default:
				return &UnknownSizeUnitError{Unit: us}
			}
		case 'T':
			switch usu[1:] {
			case "", "B", "BYTE", "ERABYTE":
				mul = 1000 * 1000 * 1000 * 1000
			case "IB", "IBYTE", "EBIBYTE":
				mul = 1 << 40
			case "BIT", "ERABIT":
				mul, div = 1000*1000*1000*1000, 8
			case "IBIT", "EBIBIT":
				mul, div = 1<<40, 8
			default:
				return &UnknownSizeUnitError{Unit: us}
			}
		case 'P':
			switch usu[1:] {
			case "", "B", "BYTE", "ETABYTE":
				mul = 1000 * 1000 * 1000 * 1000 * 1000
			case "IB", "IBYTE", "EBIBYTE":
				mul = 1 << 50
			case "BIT", "ETABIT":
				mul, div = 1000*1000*1000*1000*1000, 8
			case "IBIT", "EBIBIT":
				mul, div = 1<<50, 8
			default:
				return &UnknownSizeUnitError{Unit: us}
			}
		case 'E':
			switch usu[1:] {
			case "", "B", "BYTE", "XABYTE":
				mul = 1000 * 1000 * 1000 * 1000 * 1000 * 1000
			case "IB", "IBYTE", "BIBYTE":
				mul = 1 << 60
			case "BIT", "XABIT":
				mul, div = 1000*1000*1000*1000*1000*1000, 8
			case "IBIT", "BIBIT":
				mul, div = 1<<60, 8
			default:
				return &UnknownSizeUnitError{Unit: us}
			}
		default:
			panic("impossible")
		}
	}
	n, err := strconv.ParseUint(ns, 10, 64)
	if err != nil {
		f, err := strconv.ParseFloat(ns, 64)
		if err != nil {
			return err
		}
		*s = (Size)(f * float64(mul) / float64(div))
		return nil
	}
	*s = (Size)(n * mul / div)
	return nil
}

func (s *Size) yamlTypeError(value *yaml.Node, err error) error {
	v := value.Value
	if value.Tag != "!!seq" && value.Tag != "!!map" {
		if len(v) > 10 {
			v = " `" + v[:7] + "...`"
		} else {
			v = " `" + v + "`"
		}
	}
	return &yaml.TypeError{Errors: []string{
		fmt.Sprintf("line %d: cannot unmarshal %s%s into %s, cause: %s", value.Line, value.Tag, v, reflect.TypeOf(s).Elem(), err.Error()),
	}}
}

func (s *Size) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		if err := s.UnmarshalText(uns.StringToBytes(value.Value)); err != nil {
			return s.yamlTypeError(value, err)
		}
		return nil
	}
	return s.yamlTypeError(value, nil)
}

func (s *Size) UnmarshalJSON(bytes []byte) error {
	return s.UnmarshalText(bytes)
}

func (s Size) String() string {
	switch {
	case s < KiloByte:
		return fmt.Sprintf("%dB", s)
	case s < MegaByte:
		return fmt.Sprintf("%fKB", float64(s)/KiloByte)
	case s < GigaByte:
		return fmt.Sprintf("%fMB", float64(s)/MegaByte)
	case s < TeraByte:
		return fmt.Sprintf("%fGB", float64(s)/GigaByte)
	case s < PetaByte:
		return fmt.Sprintf("%fTB", float64(s)/TeraByte)
	case s < ExaByte:
		return fmt.Sprintf("%fPB", float64(s)/PetaByte)
	default:
		return fmt.Sprintf("%fEB", float64(s)/ExaByte)
	}
}

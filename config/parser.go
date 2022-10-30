package config

import (
	"os"
)

type parser interface {
	Unmarshal(c interface{}) error
}

type Parser struct {
	parsers []parser
}

func (p *Parser) AddBytes(bs []byte, typ Type) {
	p.parsers = append(p.parsers, &bytesParser{
		bytes: bs,
		typ:   typ,
	})
}

func (p *Parser) SetBytes(bs []byte, typ Type) {
	p.parsers = p.parsers[:0]
	p.AddBytes(bs, typ)
}

func (p *Parser) AddFile(path string, typ *Type) {
	var t Type
	if typ == nil {
		t = ProbeType(path)
	} else {
		t = *typ
	}
	p.parsers = append(p.parsers, &fileParser{
		path: path,
		typ:  t,
	})
}

func (p *Parser) SetFile(path string, typ *Type) {
	p.parsers = p.parsers[:0]
	p.AddFile(path, typ)
}

func (p *Parser) AddFilePrefix(prefix string, types ...Type) {
	group := &parserGroup{errorIgnore: os.IsNotExist}
	for _, typ := range types {
		for _, suffix := range typ.Suffixes {
			group.parsers = append(group.parsers, &fileParser{
				path: prefix + "." + suffix,
				typ:  typ,
			})
		}
	}
	p.parsers = append(p.parsers, group)
}

func (p *Parser) SetFilePrefix(prefix string, types ...Type) {
	p.parsers = p.parsers[:0]
	p.AddFilePrefix(prefix, types...)
}

func (p *Parser) Unmarshal(c interface{}) error {
	if err := HandleDefault(c); err != nil {
		return err
	}
	if err := PreHandle(c); err != nil {
		return err
	}
	for _, parser := range p.parsers {
		err := parser.Unmarshal(c)
		if err != nil {
			return err
		}
	}
	if err := PostHandle(c); err != nil {
		return err
	}
	return nil
}

//
//type configFile struct {
//	path           string
//	data           []byte
//	typ            Type
//	ignoreNotExist bool
//}
//
//
//
//type Parser struct {
//	configs []configFile
//}
//
//func Exist(path string) bool {
//	stat, err := os.Stat(path)
//	if err != nil {
//		return false
//	}
//	return stat.Mode().IsRegular()
//}
//
//func (p *Parser) SetConfigFilePrefix(prefix string, types ...*Type) {
//	p.SetConfigFilePrefixIgnoreNotExist(prefix, types...)
//	p.configs[len(p.configs)-1].ignoreNotExist = false
//}
//
//func (p *Parser) SetConfigFilePrefixIgnoreNotExist(prefix string, types ...*Type) {
//	p.configs = p.configs[:0]
//	for _, typ := range types {
//		for _, suffix := range typ.Suffixes {
//			p.configs = append(p.configs, configFile{
//				path:           prefix + "." + suffix,
//				typ:            *typ,
//				ignoreNotExist: true,
//			})
//		}
//	}
//	p.configs = append(p.configs, configFile{
//		path: prefix,
//	})
//}
//
//func (p *Parser) SetConfigFile(path string, typ *Type) {
//	p.setConfigFile(path, typ, false)
//}
//
//func (p *Parser) SetConfigFileIgnoreNotExist(path string, typ *Type) {
//	p.setConfigFile(path, typ, true)
//}
//
//func (p *Parser) setConfigFile(path string, typ *Type, ignoreNotExist bool) {
//	var t Type
//	if typ == nil {
//		t = ProbeType(path)
//	} else {
//		t = *typ
//	}
//	p.configs = []configFile{{path: path, typ: t, ignoreNotExist: ignoreNotExist}}
//}
//
//func (p *Parser) AddConfigFile(path string, typ *Type) {
//	p.addConfigFile(path, typ, false)
//}
//
//func (p *Parser) AddConfigFileIgnoreNotExist(path string, typ *Type) {
//	p.addConfigFile(path, typ, true)
//}
//
//func (p *Parser) addConfigFile(path string, typ *Type, ignoreNotExist bool) {
//	var t Type
//	if typ == nil {
//		t = ProbeType(path)
//	} else {
//		t = *typ
//	}
//	p.configs = append(p.configs, configFile{path: path, typ: t, ignoreNotExist: ignoreNotExist})
//}
//
//func (p *Parser) Unmarshal(c interface{}) error {
//	// invoke config pre handle
//	handle(preHandler{}, c, make(map[interface{}]struct{}))
//	var notExistErrors error
//	for i, config := range p.configs {
//		info, err := os.Stat(config.path)
//		if err != nil {
//			if os.IsNotExist(err) {
//				notExistErrors = errors.Append(notExistErrors, err)
//				if config.ignoreNotExist {
//					continue
//				}
//				return notExistErrors
//			}
//			return err
//		}
//		if info.Size() > MaxConfigFileSize {
//			return errors.New("config file size too large")
//		}
//		data, err := os.ReadFile(config.path)
//		if err != nil {
//			return err
//		}
//		p.configs[i].data = data
//	}
//
//out:
//	// parse config
//	for _, config := range p.configs {
//		if config.typ.Id == TypeUnknown.Id {
//			var es error
//			for _, tpy := range types {
//				if tpy.Unmarshaler != nil {
//					if err := config.typ.Unmarshaler(config.data, c); err != nil {
//						// retry next unmarshaler parse config
//						errors.Append(es, err)
//						continue
//					}
//					// parse config success
//					continue out
//				}
//			}
//			// all unmarshaler parse config failed
//			return es
//		} else {
//			if err := config.typ.Unmarshaler(config.data, c); err != nil {
//				return err
//			}
//		}
//	}
//
//	// invoke config post handle
//	_, _, err := handle(postHandler{}, c, make(map[interface{}]struct{}))
//	return err
//}

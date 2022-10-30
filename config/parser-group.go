package config

import "gitee.com/sy_183/common/errors"

type parserGroup struct {
	parsers     []parser
	errorIgnore func(err error) bool
}

func (p *parserGroup) Unmarshal(c interface{}) error {
	var es error
	for _, parser := range p.parsers {
		err := parser.Unmarshal(c)
		if err == nil {
			return nil
		}
		errors.Append(es, err)
		if p.errorIgnore(err) {
			continue
		}
		return err
	}
	return es
}

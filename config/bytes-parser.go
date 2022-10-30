package config

import (
	"gitee.com/sy_183/common/errors"
)

type bytesParser struct {
	bytes []byte
	typ   Type
}

func (p *bytesParser) Unmarshal(c interface{}) error {
	if p.typ.Id == TypeUnknown.Id {
		var es error
		for _, typ := range types {
			if typ.Unmarshaler != nil {
				if err := typ.Unmarshaler(p.bytes, c); err != nil {
					// retry next unmarshaler parse config
					errors.Append(es, err)
					continue
				}
				// parse config success
				break
			}
		}
	} else {
		if err := p.typ.Unmarshaler(p.bytes, c); err != nil {
			return err
		}
	}

	return nil
}

package config

import (
	"gitee.com/sy_183/common/errors"
	"os"
)

const MaxConfigFileSize = 32 * 1024 * 1024

var ConfigSizeTooLargeError = errors.New("config file size too large")

type fileParser struct {
	path string
	typ  Type
}

func (p *fileParser) Unmarshal(c interface{}) error {
	info, err := os.Stat(p.path)
	if err != nil {
		return err
	}
	if info.Size() > MaxConfigFileSize {
		return errors.New("config file size too large")
	}
	data, err := os.ReadFile(p.path)
	if err != nil {
		return err
	}
	bp := &bytesParser{bytes: data, typ: p.typ}
	return bp.Unmarshal(c)
}

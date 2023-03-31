package config

import (
	"fmt"
	"gitee.com/sy_183/common/unit"
	"os"
)

const MaxConfigFileSize = 32 * 1024 * 1024

var ConfigSizeTooLargeError = fmt.Errorf("配置文件大小过大，超过限定大小(%s)", unit.Size(MaxConfigFileSize))

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
		return ConfigSizeTooLargeError
	}
	data, err := os.ReadFile(p.path)
	if err != nil {
		return err
	}
	bp := &bytesParser{bytes: data, typ: p.typ}
	return bp.Unmarshal(c)
}

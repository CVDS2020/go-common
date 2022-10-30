package log

import "fmt"

type levelFilterCore struct {
	core  Core
	level LevelEnabler
}

// NewIncreaseLevelCore creates a core that can be used to increase the level of
// an existing Core. It cannot be used to decrease the logging level, as it acts
// as a filter before calling the underlying core. If level decreases the log level,
// an error is returned.
func NewIncreaseLevelCore(core Core, level LevelEnabler) (Core, error) {
	for l := _maxLevel; l >= _minLevel; l-- {
		if !core.Enabled(l) && level.Enabled(l) {
			return nil, fmt.Errorf("invalid increase level, as level %q is allowed by increased level, but not by existing core", l)
		}
	}

	return &levelFilterCore{core, level}, nil
}

func (c *levelFilterCore) Enabled(lvl Level) bool {
	return c.level.Enabled(lvl)
}

func (c *levelFilterCore) With(fields []Field) Core {
	return &levelFilterCore{c.core.With(fields), c.level}
}

func (c *levelFilterCore) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	if !c.Enabled(ent.Level) {
		return ce
	}

	return c.core.Check(ent, ce)
}

func (c *levelFilterCore) Write(ent Entry, fields []Field) error {
	return c.core.Write(ent, fields)
}

func (c *levelFilterCore) Sync() error {
	return c.core.Sync()
}

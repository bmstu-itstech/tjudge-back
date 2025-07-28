package tjudge

import (
	"context"
	"errors"
	"time"
)

var ErrConfigNotExist = errors.New("config doesn't exist")
var ErrInvalidConfig = errors.New("invalid config passed")

type Config struct {
	title        string
	allowLateReg bool
	maintenance  bool
	judgeTimeout time.Duration
	fileLimitKb  int
}

func (c Config) Title() string {
	return c.title
}

func (c Config) AllowLateReg() bool {
	return c.allowLateReg
}

func (c Config) Maintenance() bool {
	return c.maintenance
}

func (c Config) JudgeTimeout() time.Duration {
	return c.judgeTimeout
}

func (c Config) FileLimitKb() int {
	return c.fileLimitKb
}

func NewConfig(title string, allowLateReg bool, m bool, jt time.Duration, lim int) (Config, error) {
	if title == "" || jt.Seconds() <= 0 || lim <= 0 {
		return Config{}, ErrInvalidConfig
	}
	return Config{
		title,
		allowLateReg,
		m,
		jt,
		lim,
	}, nil
}

func MustNewConfig(title string, allowLateReg bool, m bool, jt time.Duration, lim int) Config {
	c, err := NewConfig(title, allowLateReg, m, jt, lim)
	if err != nil {
		panic(err)
	}
	return c
}

type ConfigRepository interface {
	Get(context.Context) (Config, error)
	Save(context.Context, Config) error
}

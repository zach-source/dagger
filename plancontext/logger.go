package plancontext

import (
	"fmt"
	"sync"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog"
	"go.dagger.io/dagger/compiler"
	"go.dagger.io/dagger/pkg"
)

var loggerIDPath = cue.MakePath(
	cue.Str("$dagger"),
	cue.Str("logger"),
	cue.Hid("_id", pkg.DaggerPackage),
)

func IsLoggerValue(v *compiler.Value) bool {
	return v.LookupPath(loggerIDPath).Exists()
}

// TODO: support using context to pull data out
type Logger struct {
	id     string
	logger *zerolog.Logger
	level  zerolog.Level
}

func (c *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	if level > c.level {
		return c.logger.WithLevel(level)
	}

	return c.logger.WithLevel(zerolog.Disabled)
}

func (c *Logger) Event() *zerolog.Event {
	return c.logger.WithLevel(c.level)
}

func (c *Logger) MarshalCUE() *compiler.Value {
	v := compiler.NewValue()
	if err := v.FillPath(loggerIDPath, c.id); err != nil {
		panic(err)
	}
	return v
}

type loggerContext struct {
	l     sync.RWMutex
	store map[string]*Logger
}

func (c *loggerContext) New(name string, lg zerolog.Logger, lvl zerolog.Level) *Logger {
	c.l.Lock()
	defer c.l.Unlock()

	l := &Logger{
		id:     hashID(name),
		logger: &lg,
		level:  lvl,
	}

	c.store[l.id] = l
	return l
}

func (c *loggerContext) Get(id string) *Logger {
	c.l.RLock()
	defer c.l.RUnlock()

	return c.store[id]
}

func (c *loggerContext) FromValue(v *compiler.Value) (*Logger, error) {
	c.l.RLock()
	defer c.l.RUnlock()

	if !v.LookupPath(loggerIDPath).IsConcrete() {
		return nil, fmt.Errorf("invalid Logger at path %q: Logger is not set", v.Path())
	}

	id, err := v.LookupPath(loggerIDPath).String()
	if err != nil {
		return nil, fmt.Errorf("invalid Logger at path %q: %w", v.Path(), err)
	}

	logger, ok := c.store[id]
	if !ok {
		return nil, fmt.Errorf("Logger %q not found", id)
	}

	return logger, nil
}

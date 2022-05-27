package task

import (
	"context"
	"errors"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.dagger.io/dagger/compiler"
	"go.dagger.io/dagger/plancontext"
	"go.dagger.io/dagger/solver"
)

func init() {
	Register("Log", func() Task { return &logTask{} })
}

type logTask struct {
}

func (l *logTask) Run(ctx context.Context, pctx *plancontext.Context, s *solver.Solver, v *compiler.Value) (*compiler.Value, error) {
	lg := log.Ctx(ctx)

	message, err := v.Lookup("message").String()
	if err != nil {
		return nil, err
	}

	level, err := v.Lookup("level").String()
	if err != nil {
		return nil, err
	}

	var fields map[string]interface{}
	if i := v.Lookup("fields"); i.Exists() {
		defFields, err := i.Fields(cue.All())
		if err != nil {
			return nil, err
		}

		fields = make(map[string]interface{})
		for i := 0; i < len(defFields); i = i + 2 {
			keyStr, err := l.getString(pctx, defFields[i].Value)
			if err != nil {
				return nil, err
			}
			fieldStr, err := l.getString(pctx, defFields[i+1].Value)
			if err != nil {
				return nil, err
			}
			fields[keyStr] = fieldStr
		}
	}

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	lg.WithLevel(lvl).Fields(fields).Msg(message)

	// Discuss, should panic/fatal/error logs fail our exec?
	// Should we even allow them?
	// Should this be a separate task or func? Or an option on the log?
	if lvl == zerolog.ErrorLevel {
		return nil, errors.New("error encountered")
	}

	if lvl == zerolog.PanicLevel {
		return nil, errors.New("panic error encountered")
	}

	if lvl == zerolog.FatalLevel {
		return nil, errors.New("fatal error encountered")
	}

	return compiler.NewValue(), nil
}

func (l logTask) getString(pctx *plancontext.Context, v *compiler.Value) (string, error) {
	if plancontext.IsSecretValue(v) {
		return "****", nil
	}

	s, err := v.String()
	if err != nil {
		return "", err
	}

	return s, nil
}

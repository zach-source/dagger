package task

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog"
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
	lg, err := pctx.Loggers.FromValue(v.Lookup("input"))
	if err != nil {
		return nil, err
	}

	message, err := v.Lookup("message").String()
	if err != nil {
		return nil, err
	}

	var level string

	if l := v.Lookup("level"); l.Exists() {
		level, err = l.String()
		if err != nil {
			return nil, err
		}
	}

	fields := make(map[string]interface{})
	if i := v.Lookup("fields"); i.Exists() {
		allFields, err := i.Fields(cue.All())
		if err != nil { return nil, err }

		for _, f := range allFields {
			selectors := f.Value.Path().Selectors()
			name := selectors[len(selectors)-1].String()

			var value string
			switch f.Value.Kind() {
			case cue.BoolKind:
				v, err := f.Value.Bool()
				if err != nil { return nil, err }
				value = fmt.Sprintf("%v", v)
			case cue.StringKind:
				v, err := f.Value.String()
				if err != nil { return nil, err }
				value = fmt.Sprintf("%v", v)
			case cue.IntKind:
				v, err := f.Value.Int64()
				if err != nil { return nil, err }
				value = fmt.Sprintf("%v", v)
			case cue.FloatKind:
				b, err := f.Value.Cue().Float64()
				if err != nil { return nil, err }
				value = fmt.Sprintf("%v", b)
			default:
				return nil, fmt.Errorf("value provided for '%s' has an unsupported '%v' kind for logging", f.Value.Path().String(), f.Value.Kind())
			}

			fields[name] = value
		}
	}

	event := lg.Event()

	if level != "" {
		lvl, err := zerolog.ParseLevel(level)
		if err != nil {
			return nil, err
		}

		event = lg.WithLevel(lvl)
	}

	event.Str("task", v.Path().String()).Fields(fields).Msg(message)
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

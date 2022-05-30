package task

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.dagger.io/dagger/compiler"
	"go.dagger.io/dagger/plancontext"
	"go.dagger.io/dagger/solver"
)

func init() {
	Register("ClientLogger", func() Task { return &clientLogger{} })
}

type clientLogger struct{}

func (t clientLogger) Run(ctx context.Context, pctx *plancontext.Context, s *solver.Solver, v *compiler.Value) (*compiler.Value, error) {
	lg := log.Ctx(ctx)

	name, err := v.Lookup("name").String()
	if err != nil {
		return nil, err
	}

	var level zerolog.Level

	if v2 := v.Lookup("level"); v2.Exists() {
		levelStr, err := v2.String()
		if err != nil {
			return nil, err
		}
		level, err = zerolog.ParseLevel(levelStr)
		if err != nil {
			return nil, err
		}
	} else {
		// use cli default if not set
		levelStr := viper.GetString("log-level")
		level, err = zerolog.ParseLevel(levelStr)
		if err != nil {
			return nil, err
		}
	}

	logger := pctx.Loggers.New(name, lg.With().Logger(), level)
	if err != nil {
		return nil, err
	}

	return compiler.NewValue().FillFields(map[string]interface{}{
		"output": logger.MarshalCUE(),
	})
}

package job

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/accuknox/rinc/internal/report/rabbitmq"
	"github.com/accuknox/rinc/internal/util"
)

// GenerateRMQReport generates a RabbitMQ status and metrics report.
func (j Job) GenerateRMQReport(ctx context.Context, now time.Time) error {
	stamp := now.Format(util.IsosecLayout)
	err := j.initStamp(ctx, stamp)
	if err != nil {
		return fmt.Errorf("initializing %q stamp: %w", stamp, err)
	}

	path := filepath.Join(j.conf.Output, stamp, "rabbitmq.html")
	f, err := os.Create(path)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			fmt.Sprintf("creating %q file", path),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("creating %q file: %w", path, err)
	}
	defer f.Close()

	rmqR := rabbitmq.NewReporter(j.conf.RabbitMQ, j.kubeClient)
	err = rmqR.Report(ctx, f, now)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"generating rabbitmq report",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("generating rabbitmq report: %w", err)
	}

	return nil
}

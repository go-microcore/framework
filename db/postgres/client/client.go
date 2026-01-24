package client // import "go.microcore.dev/framework/db/postgres/client"

import (
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type config struct {
	postgres postgres.Config
	gorm     *gorm.Config
}

var logger = log.New(pkg)

func New(opts ...Option) *gorm.DB {
	config := &config{
		postgres: postgres.Config{},
		gorm:     &gorm.Config{},
	}

	for _, opt := range opts {
		opt(config)
	}

	client, err := gorm.Open(
		postgres.New(config.postgres),
		config.gorm,
	)
	if err != nil {
		logger.Error(
			"failed to create client",
			slog.Any("error", err),
		)
		shutdown.Exit(shutdown.ExitUnavailable)
	}

	logger.Debug(
		"client created",
		slog.Group("postgres",
			slog.String("driver_name", config.postgres.DriverName),
			slog.String("dsn", MaskDSN(config.postgres.DSN)),
			slog.Bool("without_quoting_check", config.postgres.WithoutQuotingCheck),
			slog.Bool("prefer_simple_protocol", config.postgres.PreferSimpleProtocol),
			slog.Bool("without_returning", config.postgres.WithoutReturning),
			slog.Any("conn", config.postgres.Conn),
		),
		slog.Group("gorm",
			slog.Bool("skip_default_transaction", config.gorm.SkipDefaultTransaction),
			slog.Duration("default_transaction_timeout", config.gorm.DefaultTransactionTimeout),
			slog.Duration("default_context_timeout", config.gorm.DefaultContextTimeout),
			slog.Any("naming_strategy", config.gorm.NamingStrategy),
			slog.Bool("full_save_associations", config.gorm.FullSaveAssociations),
			slog.Any("logger", config.gorm.Logger),
			slog.Any("now_func", config.gorm.NowFunc),
			slog.Bool("dry_run", config.gorm.DryRun),
			slog.Bool("prepare_stmt", config.gorm.PrepareStmt),
			slog.Int("prepare_stmt_max_size", config.gorm.PrepareStmtMaxSize),
			slog.Duration("prepare_stmt_ttl", config.gorm.PrepareStmtTTL),
			slog.Bool("disable_automatic_ping", config.gorm.DisableAutomaticPing),
			slog.Bool("disable_foreign_key_constraint_when_migrating", config.gorm.DisableForeignKeyConstraintWhenMigrating),
			slog.Bool("ignore_relationships_when_migrating", config.gorm.IgnoreRelationshipsWhenMigrating),
			slog.Bool("disable_nested_transaction", config.gorm.DisableNestedTransaction),
			slog.Bool("allow_global_update", config.gorm.AllowGlobalUpdate),
			slog.Bool("query_fields", config.gorm.QueryFields),
			slog.Int("create_batch_size", config.gorm.CreateBatchSize),
			slog.Bool("translate_error", config.gorm.TranslateError),
			slog.Bool("propagate_unscoped", config.gorm.PropagateUnscoped),
			slog.Any("clause_builders", config.gorm.ClauseBuilders),
			slog.Any("conn_pool", config.gorm.ConnPool),
			slog.Any("plugins", config.gorm.Plugins),
		),
	)

	return client
}

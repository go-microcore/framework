package client // import "go.microcore.dev/framework/db/postgres/client"

import (
	"fmt"
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type client struct {
	dsn    *dsn
	config *gorm.Config
}

type dsn struct {
	host            string
	port            int
	user            string
	password        string
	db              string
	ssl             string
	searchPath      string
	applicationName string
}

var logger = log.New(pkg)

func New(opts ...Option) *gorm.DB {
	client := &client{
		dsn: &dsn{
			host:            defaultDsnHost,
			port:            defaultDsnPort,
			user:            defaultDsnUser,
			password:        defaultDsnPassword,
			db:              defaultDsnDb,
			ssl:             defaultDsnSsl,
			searchPath:      defaultDsnSearchPath,
			applicationName: defaultDsnApplicationName,
		},
		config: &gorm.Config{},
	}

	for _, opt := range opts {
		opt(client)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s application_name=%s",
		client.dsn.host,
		client.dsn.port,
		client.dsn.user,
		client.dsn.password,
		client.dsn.db,
		client.dsn.ssl,
		client.dsn.searchPath,
		client.dsn.applicationName,
	)

	c, err := gorm.Open(
		postgres.Open(dsn),
		client.config,
	)
	if err != nil {
		logger.Error(
			"failed to create client",
			slog.Any("error", err),
		)
		panic(err)
	}

	logger.Info(
		"client has been successfully created",
		slog.Group("dsn",
			slog.String("host", client.dsn.host),
			slog.Int("port", client.dsn.port),
			slog.String("db", client.dsn.db),
			slog.String("ssl", client.dsn.ssl),
			slog.String("searchPath", client.dsn.searchPath),
			slog.String("applicationName", client.dsn.applicationName),
		),
		slog.Group("config",
			slog.Bool("skip_default_transaction", client.config.SkipDefaultTransaction),
			slog.Duration("default_transaction_timeout", client.config.DefaultTransactionTimeout),
			slog.Duration("default_context_timeout", client.config.DefaultContextTimeout),
			slog.Bool("full_save_associations", client.config.FullSaveAssociations),
			slog.Bool("dry_run", client.config.DryRun),
			slog.Bool("prepare_stmt", client.config.PrepareStmt),
			slog.Int("prepare_stmt_max_size", client.config.PrepareStmtMaxSize),
			slog.Duration("prepare_stmt_ttl", client.config.PrepareStmtTTL),
			slog.Bool("disable_automatic_ping", client.config.DisableAutomaticPing),
			slog.Bool("disable_foreign_key_constraint_when_migrating", client.config.DisableForeignKeyConstraintWhenMigrating),
			slog.Bool("ignore_relationships_when_migrating", client.config.IgnoreRelationshipsWhenMigrating),
			slog.Bool("disable_nested_transaction", client.config.DisableNestedTransaction),
			slog.Bool("allow_global_update", client.config.AllowGlobalUpdate),
			slog.Bool("query_fields", client.config.QueryFields),
			slog.Int("create_batch_size", client.config.CreateBatchSize),
			slog.Bool("translate_error", client.config.TranslateError),
			slog.Bool("propagate_unscoped", client.config.PropagateUnscoped),
		),
	)

	return c
}

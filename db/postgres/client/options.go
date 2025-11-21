package client // import "go.microcore.dev/framework/db/postgres/client"

import (
	"time"

	_ "go.microcore.dev/framework"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	loggerGorm "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Option func(*config)

// -------------------- Postgres Options --------------------

// WithPostgresDriverName sets the driver name for the Postgres Dialector
func WithPostgresDriverName(driverName string) Option {
	return func(c *config) {
		c.postgres.DriverName = driverName
	}
}

// WithPostgresDSN sets the Data Source Name (connection string) for Postgres
func WithPostgresDSN(dsn string) Option {
	return func(c *config) {
		c.postgres.DSN = dsn
	}
}

// WithPostgresWithoutQuotingCheck disables column/table quoting checks
func WithPostgresWithoutQuotingCheck(withoutQuotingCheck bool) Option {
	return func(c *config) {
		c.postgres.WithoutQuotingCheck = withoutQuotingCheck
	}
}

// WithPostgresPreferSimpleProtocol enables Postgres simple protocol instead of extended protocol
func WithPostgresPreferSimpleProtocol(preferSimpleProtocol bool) Option {
	return func(c *config) {
		c.postgres.PreferSimpleProtocol = preferSimpleProtocol
	}
}

// WithPostgresWithoutReturning disables generating RETURNING clause in INSERT/UPDATE/DELETE
func WithPostgresWithoutReturning(withoutReturning bool) Option {
	return func(c *config) {
		c.postgres.WithoutReturning = withoutReturning
	}
}

// WithPostgresConn sets a custom connection pool for Postgres (e.g., *sql.DB or sqlmock)
func WithPostgresConn(conn gorm.ConnPool) Option {
	return func(c *config) {
		c.postgres.Conn = conn
	}
}

// -------------------- GORM Options --------------------

// WithGormSkipDefaultTransaction disables automatic transactions for single create/update/delete operations
func WithGormSkipDefaultTransaction(skipDefaultTransaction bool) Option {
	return func(c *config) {
		c.gorm.SkipDefaultTransaction = skipDefaultTransaction
	}
}

// WithGormDefaultTransactionTimeout sets the default timeout for transactions
func WithGormDefaultTransactionTimeout(defaultTransactionTimeout time.Duration) Option {
	return func(c *config) {
		c.gorm.DefaultTransactionTimeout = defaultTransactionTimeout
	}
}

// WithGormDefaultContextTimeout sets the default timeout for context in DB operations
func WithGormDefaultContextTimeout(defaultContextTimeout time.Duration) Option {
	return func(c *config) {
		c.gorm.DefaultContextTimeout = defaultContextTimeout
	}
}

// WithGormNamingStrategy sets a custom naming strategy for tables and columns
func WithGormNamingStrategy(namingStrategy schema.Namer) Option {
	return func(c *config) {
		c.gorm.NamingStrategy = namingStrategy
	}
}

// WithGormFullSaveAssociations enables saving all associations when saving a record
func WithGormFullSaveAssociations(fullSaveAssociations bool) Option {
	return func(c *config) {
		c.gorm.FullSaveAssociations = fullSaveAssociations
	}
}

// WithGormLogger sets a custom logger for GORM
func WithGormLogger(logger loggerGorm.Interface) Option {
	return func(c *config) {
		c.gorm.Logger = logger
	}
}

// WithGormNowFunc sets a custom function to generate timestamps for new records
func WithGormNowFunc(nowFunc func() time.Time) Option {
	return func(c *config) {
		c.gorm.NowFunc = nowFunc
	}
}

// WithGormDryRun enables generating SQL statements without executing them
func WithGormDryRun(dryRun bool) Option {
	return func(c *config) {
		c.gorm.DryRun = dryRun
	}
}

// WithGormPrepareStmt enables caching prepared statements for faster execution
func WithGormPrepareStmt(prepareStmt bool) Option {
	return func(c *config) {
		c.gorm.PrepareStmt = prepareStmt
	}
}

// WithGormPrepareStmtMaxSize sets the maximum size of the prepared statement cache
func WithGormPrepareStmtMaxSize(prepareStmtMaxSize int) Option {
	return func(c *config) {
		c.gorm.PrepareStmtMaxSize = prepareStmtMaxSize
	}
}

// WithGormPrepareStmtTTL sets the TTL (time-to-live) for prepared statements in the cache
func WithGormPrepareStmtTTL(prepareStmtTTL time.Duration) Option {
	return func(c *config) {
		c.gorm.PrepareStmtTTL = prepareStmtTTL
	}
}

// WithGormDisableAutomaticPing disables automatic database ping on initialization
func WithGormDisableAutomaticPing(disableAutomaticPing bool) Option {
	return func(c *config) {
		c.gorm.DisableAutomaticPing = disableAutomaticPing
	}
}

// WithGormDisableForeignKeyConstraintWhenMigrating disables creating foreign key constraints during migrations
func WithGormDisableForeignKeyConstraintWhenMigrating(disableForeignKeyConstraintWhenMigrating bool) Option {
	return func(c *config) {
		c.gorm.DisableForeignKeyConstraintWhenMigrating = disableForeignKeyConstraintWhenMigrating
	}
}

// WithGormIgnoreRelationshipsWhenMigrating disables saving relationships when performing migrations
func WithGormIgnoreRelationshipsWhenMigrating(ignoreRelationshipsWhenMigrating bool) Option {
	return func(c *config) {
		c.gorm.IgnoreRelationshipsWhenMigrating = ignoreRelationshipsWhenMigrating
	}
}

// WithGormDisableNestedTransaction disables nested transactions
func WithGormDisableNestedTransaction(disableNestedTransaction bool) Option {
	return func(c *config) {
		c.gorm.DisableNestedTransaction = disableNestedTransaction
	}
}

// WithGormAllowGlobalUpdate allows updates/deletes without WHERE conditions
func WithGormAllowGlobalUpdate(allowGlobalUpdate bool) Option {
	return func(c *config) {
		c.gorm.AllowGlobalUpdate = allowGlobalUpdate
	}
}

// WithGormQueryFields enables querying all fields of a table
func WithGormQueryFields(queryFields bool) Option {
	return func(c *config) {
		c.gorm.QueryFields = queryFields
	}
}

// WithGormCreateBatchSize sets default batch size for bulk inserts
func WithGormCreateBatchSize(createBatchSize int) Option {
	return func(c *config) {
		c.gorm.CreateBatchSize = createBatchSize
	}
}

// WithGormTranslateError enables translating database errors using Dialector
func WithGormTranslateError(translateError bool) Option {
	return func(c *config) {
		c.gorm.TranslateError = translateError
	}
}

// WithGormPropagateUnscoped propagates Unscoped mode to nested statements
func WithGormPropagateUnscoped(propagateUnscoped bool) Option {
	return func(c *config) {
		c.gorm.PropagateUnscoped = propagateUnscoped
	}
}

// WithGormClauseBuilders sets custom clause builders for SQL generation
func WithGormClauseBuilders(clauseBuilders map[string]clause.ClauseBuilder) Option {
	return func(c *config) {
		c.gorm.ClauseBuilders = clauseBuilders
	}
}

// WithGormConnPool sets a custom connection pool for GORM operations
func WithGormConnPool(connPool gorm.ConnPool) Option {
	return func(c *config) {
		c.gorm.ConnPool = connPool
	}
}

// WithGormPlugins sets custom GORM plugins
func WithGormPlugins(plugins map[string]gorm.Plugin) Option {
	return func(c *config) {
		c.gorm.Plugins = plugins
	}
}

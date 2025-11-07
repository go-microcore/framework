package client // import "go.microcore.dev/framework/db/postgres/client"

import (
	"time"

	_ "go.microcore.dev/framework"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	loggerGorm "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Option func(*client)

// DSN

func WithHost(host string) Option {
	return func(c *client) {
		c.dsn.host = host
	}
}

func WithPort(port int) Option {
	return func(c *client) {
		c.dsn.port = port
	}
}

func WithUser(user string) Option {
	return func(c *client) {
		c.dsn.user = user
	}
}

func WithPassword(password string) Option {
	return func(c *client) {
		c.dsn.password = password
	}
}

func WithDb(db string) Option {
	return func(c *client) {
		c.dsn.db = db
	}
}

func WithSsl(ssl string) Option {
	return func(c *client) {
		c.dsn.ssl = ssl
	}
}

func WithSearchPath(searchPath string) Option {
	return func(c *client) {
		c.dsn.searchPath = searchPath
	}
}

func WithApplicationName(applicationName string) Option {
	return func(c *client) {
		c.dsn.applicationName = applicationName
	}
}

// Config

// GORM perform single create, update, delete operations in
// transactions by default to ensure database data integrity
func WithSkipDefaultTransaction(skipDefaultTransaction bool) Option {
	return func(c *client) {
		c.config.SkipDefaultTransaction = skipDefaultTransaction
	}
}

func WithDefaultTransactionTimeout(defaultTransactionTimeout time.Duration) Option {
	return func(c *client) {
		c.config.DefaultTransactionTimeout = defaultTransactionTimeout
	}
}

func WithDefaultContextTimeout(defaultContextTimeout time.Duration) Option {
	return func(c *client) {
		c.config.DefaultContextTimeout = defaultContextTimeout
	}
}

// NamingStrategy tables, columns naming strategy
func WithNamingStrategy(namingStrategy schema.Namer) Option {
	return func(c *client) {
		c.config.NamingStrategy = namingStrategy
	}
}

// FullSaveAssociations full save associations
func WithFullSaveAssociations(fullSaveAssociations bool) Option {
	return func(c *client) {
		c.config.FullSaveAssociations = fullSaveAssociations
	}
}

func WithLogger(logger loggerGorm.Interface) Option {
	return func(c *client) {
		c.config.Logger = logger
	}
}

func WithNowFunc(nowFunc func() time.Time) Option {
	return func(c *client) {
		c.config.NowFunc = nowFunc
	}
}

// DryRun generate sql without execute
func WithDryRun(dryRun bool) Option {
	return func(c *client) {
		c.config.DryRun = dryRun
	}
}

// PrepareStmt executes the given query in cached statement
func WithPrepareStmt(prepareStmt bool) Option {
	return func(c *client) {
		c.config.PrepareStmt = prepareStmt
	}
}

// PrepareStmt cache support LRU expired,
// default maxsize=int64 Max value and ttl=1h
func WithPrepareStmtMaxSize(prepareStmtMaxSize int) Option {
	return func(c *client) {
		c.config.PrepareStmtMaxSize = prepareStmtMaxSize
	}
}
func WithPrepareStmtTTL(prepareStmtTTL time.Duration) Option {
	return func(c *client) {
		c.config.PrepareStmtTTL = prepareStmtTTL
	}
}

func WithDisableAutomaticPing(disableAutomaticPing bool) Option {
	return func(c *client) {
		c.config.DisableAutomaticPing = disableAutomaticPing
	}
}

func WithDisableForeignKeyConstraintWhenMigrating(disableForeignKeyConstraintWhenMigrating bool) Option {
	return func(c *client) {
		c.config.DisableForeignKeyConstraintWhenMigrating = disableForeignKeyConstraintWhenMigrating
	}
}

func WithIgnoreRelationshipsWhenMigrating(ignoreRelationshipsWhenMigrating bool) Option {
	return func(c *client) {
		c.config.IgnoreRelationshipsWhenMigrating = ignoreRelationshipsWhenMigrating
	}
}

func WithDisableNestedTransaction(disableNestedTransaction bool) Option {
	return func(c *client) {
		c.config.DisableNestedTransaction = disableNestedTransaction
	}
}

func WithAllowGlobalUpdate(allowGlobalUpdate bool) Option {
	return func(c *client) {
		c.config.AllowGlobalUpdate = allowGlobalUpdate
	}
}

// QueryFields executes the SQL query with all fields of the table
func WithQueryFields(queryFields bool) Option {
	return func(c *client) {
		c.config.QueryFields = queryFields
	}
}

// CreateBatchSize default create batch size
func WithCreateBatchSize(createBatchSize int) Option {
	return func(c *client) {
		c.config.CreateBatchSize = createBatchSize
	}
}

// TranslateError enabling error translation
func WithTranslateError(translateError bool) Option {
	return func(c *client) {
		c.config.TranslateError = translateError
	}
}

// PropagateUnscoped propagate Unscoped to every other nested statement
func WithPropagateUnscoped(propagateUnscoped bool) Option {
	return func(c *client) {
		c.config.PropagateUnscoped = propagateUnscoped
	}
}

func WithClauseBuilders(clauseBuilders map[string]clause.ClauseBuilder) Option {
	return func(c *client) {
		c.config.ClauseBuilders = clauseBuilders
	}
}

func WithConnPool(connPool gorm.ConnPool) Option {
	return func(c *client) {
		c.config.ConnPool = connPool
	}
}

func WithPlugins(plugins map[string]gorm.Plugin) Option {
	return func(c *client) {
		c.config.Plugins = plugins
	}
}

package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/SuperB1aze/catalog-service/internal/app/config/section"
	"github.com/SuperB1aze/catalog-service/migration"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()

	log.Printf("Connecting to postgres: read_timeout=%d, write_timeout=%d", cfg.ReadTimeout, cfg.WriteTimeout)

	connector := pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithReadTimeout(cfg.ReadTimeout),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout),
	)

	sqlDB := sql.OpenDB(connector)
	sqlDB.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := bunDB.PingContext(pingCtx); err != nil {
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Printf("postgres close error: %v", closeErr)
		}
		return nil, fmt.Errorf("postgres error: %w", err)
	}

	return &Client{
		_bunDB:   bunDB,
		rawBunDB: bunDB,
		cfg:      cfg,
	}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	sub, err := fs.Sub(migration.Postgres, "postgres")
	if err != nil {
		return 0, 0, fmt.Errorf("postgres error: %w", err)
	}

	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sub); err != nil {
		return 0, 0, fmt.Errorf("migrations error: %w", err)
	}

	migrator := migrate.NewMigrator(c.rawBunDB, migrations, migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable+"_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	)
	if err := migrator.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("init migrations error: %w", err)
	}

	if err := migrator.Lock(ctx); err != nil {
		return 0, 0, fmt.Errorf("lock migrations error: %w", err)
	}
	defer func() {
		if unlockErr := migrator.Unlock(ctx); unlockErr != nil {
			log.Printf("migrator unlock error: %v", unlockErr)
		}
	}()

	appliedMigrations, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("applied migrations error: %w", err)
	}
	for _, mg := range appliedMigrations {
		v, _ := strconv.ParseInt(mg.Name, 10, 64)
		if v > oldVer {
			oldVer = v
		}
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("migration error: %w", err)
	}

	newVer = oldVer
	for _, mg := range group.Migrations {
		v, _ := strconv.ParseInt(mg.Name, 10, 64)
		if v > newVer {
			newVer = v
		}
	}

	return oldVer, newVer, nil
}

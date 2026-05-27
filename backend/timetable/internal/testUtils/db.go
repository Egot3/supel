package testutils

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func NewTestInjector(t *testing.T) do.Injector {
	t.Helper()

	dsn := fmt.Sprintf("file::memory:?cache=private", t.Name())

	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	require.NoError(t, err)

	sqldb.SetMaxOpenConns(1)

	db := bun.NewDB(sqldb, sqlitedialect.New())

	err = RunMigrations(t.Context(), db)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Logf("failed to close test db(memory leaks go brrr): %v", err)
		}
	})

	i := do.New()

	do.Provide(i, func(i do.Injector) (*bun.DB, error) {
		return db, nil
	})

	return i
}

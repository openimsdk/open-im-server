package relation

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMaybeCreateTable(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		err := maybeCreateTable(&option{
			Username:      "root",
			Password:      "openIM123",
			Address:       []string{"172.28.0.1:13306"},
			Database:      "openIM_v3",
			LogLevel:      4,
			SlowThreshold: 500,
			MaxOpenConn:   1000,
			MaxIdleConn:   100,
			MaxLifeTime:   60,
			Connect: connect(expectExec{
				query: "CREATE DATABASE IF NOT EXISTS `openIM_v3` default charset utf8mb4 COLLATE utf8mb4_unicode_ci",
				args:  nil,
			}),
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("im-db", func(t *testing.T) {
		err := maybeCreateTable(&option{
			Username:      "root",
			Password:      "openIM123",
			Address:       []string{"172.28.0.1:13306"},
			Database:      "im-db",
			LogLevel:      4,
			SlowThreshold: 500,
			MaxOpenConn:   1000,
			MaxIdleConn:   100,
			MaxLifeTime:   60,
			Connect: connect(expectExec{
				query: "CREATE DATABASE IF NOT EXISTS `im-db` default charset utf8mb4 COLLATE utf8mb4_unicode_ci",
				args:  nil,
			}),
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("err", func(t *testing.T) {
		e := errors.New("e")
		err := maybeCreateTable(&option{
			Username:      "root",
			Password:      "openIM123",
			Address:       []string{"172.28.0.1:13306"},
			Database:      "openIM_v3",
			LogLevel:      4,
			SlowThreshold: 500,
			MaxOpenConn:   1000,
			MaxIdleConn:   100,
			MaxLifeTime:   60,
			Connect: connect(expectExec{
				err: e,
			}),
		})
		if !errors.Is(err, e) {
			t.Fatalf("err not is e: %v", err)
		}
	})
}

func connect(e expectExec) func(string, int) (*gorm.DB, error) {
	return func(string, int) (*gorm.DB, error) {
		return gorm.Open(mysql.New(mysql.Config{
			SkipInitializeWithVersion: true,
			Conn:                      sql.OpenDB(e),
		}), &gorm.Config{
			Logger: logger.Discard,
		})
	}
}

type expectExec struct {
	err   error
	query string
	args  []driver.NamedValue
}

func (c expectExec) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.err != nil {
		return nil, c.err
	}
	if query != c.query {
		return nil, fmt.Errorf("query mismatch. expect: %s, got: %s", c.query, query)
	}
	if reflect.DeepEqual(args, c.args) {
		return nil, fmt.Errorf("args mismatch. expect: %v, got: %v", c.args, args)
	}
	return noEffectResult{}, nil
}

func (e expectExec) Connect(context.Context) (driver.Conn, error) { return e, nil }
func (expectExec) Driver() driver.Driver                          { panic("not implemented") }
func (expectExec) Prepare(query string) (driver.Stmt, error)      { panic("not implemented") }
func (expectExec) Close() (e error)                               { return }
func (expectExec) Begin() (driver.Tx, error)                      { panic("not implemented") }

type noEffectResult struct{}

func (noEffectResult) LastInsertId() (i int64, e error) { return }
func (noEffectResult) RowsAffected() (i int64, e error) { return }

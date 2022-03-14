package database

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/smithyworks/FaDO/util"
)

var dbPool *pgxpool.Pool

var ctx = context.Background()

func Connect(connectionString string) (err error) {
	dbPool, err = pgxpool.Connect(ctx, connectionString)
	return util.ProcessErr(err)
}

func Close() {
	dbPool.Close()
}

type DBConn interface {
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}

func Acquire() (conn *pgxpool.Conn, err error) {
	conn, err = dbPool.Acquire(ctx)
	return conn, util.ProcessErr(err)
}

func Begin() (tx pgx.Tx, err error) {
	tx, err = dbPool.Begin(ctx)
	return tx, util.ProcessErr(err)
}

func Query(conn DBConn, sql string, args ...interface{}) (rows pgx.Rows, err error) {
	rows, err = conn.Query(ctx, sql, args...)
	return rows, util.ProcessErr(err)
}

func Exec(conn DBConn, sql string, args ...interface{}) (ct pgconn.CommandTag, err error) {
	ct, err = conn.Exec(ctx, sql, args...)
	return ct, util.ProcessErr(err)
}

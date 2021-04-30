// Copyright © 2021 Kaleido, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlcommon

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/kaleido-io/firefly/internal/fftypes"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/internal/persistence"
)

type SQLCommon struct {
	db      *sql.DB
	options *SQLCommonOptions
}

type SQLCommonOptions struct {
	PlaceholderFormat sq.PlaceholderFormat
}

func InitSQLCommon(ctx context.Context, s *SQLCommon, db *sql.DB, options *SQLCommonOptions) (*persistence.Capabilities, error) {
	s.db = db
	if options == nil {
		options = &SQLCommonOptions{
			PlaceholderFormat: sq.Dollar,
		}
	}
	s.options = options
	return &persistence.Capabilities{}, nil
}

func (s *SQLCommon) beginTx(ctx context.Context) (context.Context, *sql.Tx, error) {
	l := log.L(ctx).WithField("dbtx", fftypes.ShortID())
	ctx = log.WithLogger(ctx, l)
	l.Debugf("SQL-> begin")
	tx, err := s.db.Begin()
	if err != nil {
		return ctx, nil, i18n.WrapError(ctx, err, i18n.MsgDBBeginFailed)
	}
	l.Debugf("SQL<- begin")
	return ctx, tx, err
}

func (s *SQLCommon) queryTx(ctx context.Context, tx *sql.Tx, q sq.SelectBuilder) (*sql.Rows, error) {
	l := log.L(ctx)
	sqlQuery, args, err := q.PlaceholderFormat(s.options.PlaceholderFormat).ToSql()
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBQueryBuildFailed)
	}
	l.Debugf(`SQL-> query sql=[ %s ] args=%+v`, sqlQuery, args)
	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.QueryContext(ctx, sqlQuery, args...)
	} else {
		rows, err = s.db.QueryContext(ctx, sqlQuery, args...)
	}
	if err != nil {
		l.Errorf(`SQL query failed: %s sql=[ %s ]`, err, sqlQuery)
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBQueryFailed)
	}
	l.Debugf(`SQL<- query`)
	return rows, nil
}

func (s *SQLCommon) query(ctx context.Context, q sq.SelectBuilder) (*sql.Rows, error) {
	return s.queryTx(ctx, nil, q)
}

func (s *SQLCommon) insertTx(ctx context.Context, tx *sql.Tx, q sq.InsertBuilder) (sql.Result, error) {
	l := log.L(ctx)
	sqlQuery, args, err := q.PlaceholderFormat(s.options.PlaceholderFormat).ToSql()
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBQueryBuildFailed)
	}
	l.Debugf(`SQL-> insert sql=[ %s ] args=%+v`, sqlQuery, args)
	res, err := tx.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		l.Errorf(`SQL insert failed: %s sql=[ %s ]: %s`, err, sqlQuery, err)
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBUpdateFailed)
	}
	ra, _ := res.RowsAffected() // currently only used for debugging
	l.Debugf(`SQL<- insert affected=%d`, ra)
	return res, nil
}

func (s *SQLCommon) deleteTx(ctx context.Context, tx *sql.Tx, q sq.DeleteBuilder) (sql.Result, error) {
	l := log.L(ctx)
	sqlQuery, args, err := q.PlaceholderFormat(s.options.PlaceholderFormat).ToSql()
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBQueryBuildFailed)
	}
	l.Debugf(`SQL-> delete sql=[ %s ] args=%+v`, sqlQuery, args)
	res, err := tx.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		l.Errorf(`SQL delete failed: %s sql=[ %s ]: %s`, err, sqlQuery, err)
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBUpdateFailed)
	}
	ra, _ := res.RowsAffected() // currently only used for debugging
	l.Debugf(`SQL<- delete affected=%d`, ra)
	return res, nil
}

func (s *SQLCommon) updateTx(ctx context.Context, tx *sql.Tx, q sq.UpdateBuilder) (sql.Result, error) {
	l := log.L(ctx)
	sqlQuery, args, err := q.PlaceholderFormat(s.options.PlaceholderFormat).ToSql()
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBQueryBuildFailed)
	}
	l.Debugf(`SQL-> update sql=[ %s ] args=%+v`, sqlQuery, args)
	res, err := tx.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		l.Errorf(`SQL update failed: %s sql=[ %s ]`, err, sqlQuery)
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBUpdateFailed)
	}
	ra, _ := res.RowsAffected() // currently only used for debugging
	l.Debugf(`SQL<- update affected=%d`, ra)
	return res, nil
}

// rollbackTx be safely called as a defer, as it is a cheap no-op if the transaction is complete
func (s *SQLCommon) rollbackTx(ctx context.Context, tx *sql.Tx) {
	err := tx.Rollback()
	if err == nil {
		log.L(ctx).Warnf("SQL! transaction rollback")
	}
	if err != nil && err != sql.ErrTxDone {
		log.L(ctx).Errorf(`SQL rollback failed: %s`, err)
	}
}

func (s *SQLCommon) commitTx(ctx context.Context, tx *sql.Tx) error {
	l := log.L(ctx)
	l.Debugf(`SQL-> commit`)
	err := tx.Commit()
	if err != nil {
		l.Errorf(`SQL commit failed: %s`, err)
		return i18n.WrapError(ctx, err, i18n.MsgDBCommitFailed)
	}
	l.Debugf(`SQL<- commit`)
	return nil
}
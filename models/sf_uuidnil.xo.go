package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"

	"github.com/google/uuid"
)

// UUIDNil calls the stored function 'public.uuid_nil() uuid' on db.
func UUIDNil(ctx context.Context, db DB) (uuid.UUID, error) {
	// call public.uuid_nil
	const sqlstr = `SELECT * FROM public.uuid_nil()`
	// run
	var r0 uuid.UUID
	logf(sqlstr)
	if err := db.QueryRowContext(ctx, sqlstr).Scan(&r0); err != nil {
		return uuid.UUID{}, logerror(err)
	}
	return r0, nil
}

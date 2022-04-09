package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"

	"github.com/google/uuid"
)

// UUIDNsX500 calls the stored function 'public.uuid_ns_x500() uuid' on db.
func UUIDNsX500(ctx context.Context, db DB) (uuid.UUID, error) {
	// call public.uuid_ns_x500
	const sqlstr = `SELECT * FROM public.uuid_ns_x500()`
	// run
	var r0 uuid.UUID
	logf(sqlstr)
	if err := db.QueryRowContext(ctx, sqlstr).Scan(&r0); err != nil {
		return uuid.UUID{}, logerror(err)
	}
	return r0, nil
}

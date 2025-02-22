package db

import (
	"database/sql"
	"github.com/jackc/pgx/v5/pgtype"
)

// sql.NullInt32 to *int
func NullInt32ToPtr(n sql.NullInt32) *int {
	if n.Valid {
		val := int(n.Int32)
		return &val
	}
	return nil
}

// *int into sql.NullInt32
func Int32ToNullInt32(i *int) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: int32(*i), Valid: true}
	}
	return sql.NullInt32{Valid: false}
}

// pgtype.Int4 to *int
func NullInt4ToPtr(n pgtype.Int4) *int {
	if n.Valid {
		val := int(n.Int32)
		return &val
	}
	return nil
}

// *int to pgtype.Int4
func Int32ToNullInt4(i *int) pgtype.Int4 {
	if i != nil {
		return pgtype.Int4{Int32: int32(*i), Valid: true}
	}
	return pgtype.Int4{Valid: false}
}

// pgtype.Text to *string
func NullStringToPtr(n pgtype.Text) *string {
	if n.Valid {
		return &n.String
	}
	return nil
}

// *string to pgtype.Text
func StringToNullString(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{Valid: false}
}

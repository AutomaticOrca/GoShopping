package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
	"time"
)

// pgtype.Int4 -> *int
func NullInt4ToPtr(n pgtype.Int4) *int {
	if n.Valid {
		val := int(n.Int32)
		return &val
	}
	return nil
}

// *int -> pgtype.Int4
func Int32ToNullInt4(i *int) pgtype.Int4 {
	if i != nil {
		return pgtype.Int4{Int32: int32(*i), Valid: true}
	}
	return pgtype.Int4{Valid: false}
}

// pgtype.Text -> *string
func NullStringToPtr(n pgtype.Text) *string {
	if n.Valid {
		return &n.String
	}
	return nil
}

// *string -> pgtype.Text
func StringToNullString(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{Valid: false}
}

// pgtype.Text -> string
func TextToString(text pgtype.Text) string {
	if text.Valid {
		return text.String
	}
	return ""
}

// pgtype.Numeric -> float64 (with error)
func NumericToFloat64(numeric pgtype.Numeric) float64 {
	val, err := numeric.Float64Value()
	if err != nil {
		return 0.0
	}
	return val.Float64
}

// pgtype.Numeric -> float64 (no error, safer for API)
func MustNumericToFloat64(numeric pgtype.Numeric) float64 {
	val, err := numeric.Float64Value()
	if err != nil {
		return 0.0 // Default value if conversion fails
	}
	return val.Float64
}

// pgtype.Timestamp -> string
func TimestampToString(ts pgtype.Timestamp) string {
	if ts.Valid {
		return ts.Time.Format(time.RFC3339)
	}
	return ""
}

// pgtype.Int4 -> int
func Int4ToInt(n pgtype.Int4) int {
	if n.Valid {
		return int(n.Int32)
	}
	return 0 // Default to 0 if not valid
}

// int -> pgtype.Int4
func IntToInt4(i int) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(i), Valid: true}
}

// pgtype.Int4 -> int32
func Int4ToInt32(n pgtype.Int4) int32 {
	if n.Valid {
		return n.Int32
	}
	return 0
}

// int32 -> pgtype.Int4
func Int32ToInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func StringToInt32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

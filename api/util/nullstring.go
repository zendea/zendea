package util

import "database/sql"

func SqlNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  len(value) > 0,
	}
}

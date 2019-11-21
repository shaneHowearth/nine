package datastore

import "database/sql"

// SetSQLOpenForTest - Allow test functions to set the SQL Open function
func SetSQLOpenForTest(s func(driver, dataSource string) (*sql.DB, error)) {
	sqlOpen = s
}

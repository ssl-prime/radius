package util

import "database/sql"

//ConnectDB ...
func ConnectDB() (*sql.DB, error) {
	var (
		dbUser     = "root"
		dbPassword = "root"
		dbHost     = "localhost"
		dbPort     = "3306"
		dbName     = "radius"
	)
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	return db, err
}

package golangdatabase

import (
	"database/sql"
	"time"
)

/*
(DB) SetMaxIdleConns(number) 		-> Pengaturan berapa jumlah koneksi minimal yang dibuat
(DB) SetMaxOpenConns(number) 		-> Pengaturan berapa jumlah koneksi maksimal yang dibuat
(DB) SetConnMaxIdleTime(duration) 	-> Pengaturan berapa lama koneksi yang sudah tidak digunakan akan dihapus
(DB) SetConnMaxLifetime(duration) 	-> Pengaturan berapa lama koneksi boleh digunakan
*/

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/belajar_golang_database?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}

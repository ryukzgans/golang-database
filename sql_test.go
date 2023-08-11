package golangdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// Di Golang menyediakan function yang bisa kita gunakan untuk mengirim perintah SQL ke database menggunakan function (DB) ExecContext(context, sql, params)
func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUES('tarmiji', 'Tarmiji')"
	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Success Insert Data")
}

// Function untuk melakukan query ke database, bisa menggunakan function (DB) QueryContext(context, sql, params)
func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	/*
		hasil dari rows adalah data struct sql.Rows
		menggunakan function (Rows) Next() (boolean) untuk melakukan iterasi terhadap data query, jika return data false, artinya sudah tidak ada data lagi didalam result
		untuk membaca tiap data menggunakan (Rows) Scan(columns...)
		dan jgn lupa untuk menutup (Rows) Close()
	*/

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)

		if err != nil {
			panic(err)
		}

		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
	}
}

/*
Tipe Data Database						Tipe Data Golang				Tipe Data Nullable
VARCHAR, CHAR							string							database/sql.NullString
INT, BIGINT								int32, int64					database/sql.NullInt32 || database/sql.NullInt64
FLOAT, DOUBLE							float32, float64				database/sql.NullFloat64
BOOLEAN									bool							database/sql.NullBool
DATE, DATETIME, TIME, TIMESTAMP			time.Time						database/sql.NullTime
*/

// menambahkan ?parseTime=true pada databaseConnection adalah untuk mengkonversi secara otamatis data DATE, DATETIME, TIME TIMESTAMP yang ada pada database menjadi time.Time pada tipe data di golang
func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString // contoh mengatasi null pada email tipe data string
		var balance int32
		var rating float64
		var created_at time.Time
		var birth_date sql.NullTime // contoh mengatasi null pada birth_date tipe data time
		var married bool

		err := rows.Scan(&id, &name, &email, &balance, &rating, &birth_date, &married, &created_at)

		if err != nil {
			panic(err)
		}

		fmt.Println("=====================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)

		if email.Valid {
			fmt.Println("Email:", email.String)
		}

		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)

		if birth_date.Valid {
			fmt.Println("Birth Date:", birth_date.Time)
		}

		fmt.Println("Married:", married)
		fmt.Println("Created At:", created_at)
		fmt.Println("=====================")
	}
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	// default login admin:admin
	username := "admin'; #" // contoh sql injection
	password := "admin123"  // passwordnya jelas disini salah

	sqlQuery := "select username from user where username = '" + username + "' and password = '" + password + "' limit 1"
	rows, err := db.QueryContext(ctx, sqlQuery)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Sukses Login", username)
	} else {
		fmt.Println("Gagal Login")
	}

}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	// default login admin:admin
	username := "admin"
	password := "admin"

	// untuk mengatasi hal diatas kita bisa menggunakan sql dgn parameter yaitu dgn hanya menggunakan tanda "?" pada query sql

	sqlQuery := "select username from user where username = ? and password = ? limit 1"
	rows, err := db.QueryContext(ctx, sqlQuery, username, password) // tambahan parameter sesuai dengan parameter yg ada pada query
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Sukses Login", username)
	} else {
		fmt.Println("Gagal Login")
	}

}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	execScript := "insert into user(username, password) values (?, ?)" // menambah user baru pada table user

	username := "ilham"
	password := "ilham123"
	_, err := db.QueryContext(ctx, execScript, username, password)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Success insert new user")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	execScript := "insert into comments(email, comments) values (?, ?)"

	email := "ilhamgod@gmail.com"
	comment := "Test Comment"

	result, err := db.ExecContext(ctx, execScript, email, comment)
	if err != nil {
		panic(err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		panic(nil)
	}

	fmt.Println("Success insert new comment with id", lastId)
}

// digunakan untuk menghindari pengulangan pada Function Query atau Exec yang menggunakan parameter, jadi solusinya kita menerapkan statement dlu sebelum execquery
func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	execScript := "insert into comments(email, comments) values (?, ?)"

	stmt, err := db.PrepareContext(ctx, execScript)
	if err != nil {
		panic(err)
	}

	for i := 0; i <= 10; i++ {
		email := "ilhamgod" + strconv.Itoa(i) + "@gmail.com"
		comment := "comment " + strconv.Itoa(i)

		result, err := stmt.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		resultId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment ID :", resultId)

	}
}

/*
Secara default, semua perintah SQL yang kita kirim menggunakan Golang akan otomatis di commit, atau istilahnya auto commit
Namun kita bisa menggunakan fitur transaksi sehingga SQL yang kita kirim tidak secara otomatis di commit ke database
Untuk memulai transaksi, kita bisa menggunakan function (DB) Begin(), dimana akan menghasilkan struct Tx yang merupakan representasi Transaction
Struct Tx ini yang kita gunakan sebagai pengganti DB untuk melakukan transaksi, dimana hampir semua function di DB ada di Tx, seperti Exec, Query atau Prepare
Setelah selesai proses transaksi, kita bisa gunakan function (Tx) Commit() untuk melakukan commit atau Rollback()
*/
func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	querySql := "insert into comments(email, comments) values (?, ?)"
	// do transaction
	for i := 0; i <= 10; i++ {
		email := "ilhamgod" + strconv.Itoa(i) + "@gmail.com"
		comment := "comment " + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, querySql, email, comment)
		if err != nil {
			panic(err)
		}

		resultId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment ID :", resultId)

	}

	// err = tx.Commit()   // tx.commit untuk commit data
	err = tx.Rollback() // tx.rollbackk untuk membatalkan mengirim data, jdi data yg diatas insert into tidak akan terkirim ke database
	if err != nil {
		panic(err)
	}
}

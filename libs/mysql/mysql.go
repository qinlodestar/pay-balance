package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var MyDb *sql.DB

func Query(str string) (result []map[string]string) {
	rows, err := MyDb.Query(str)
	//rows, err := MyDb.Query("SELECT * FROM user where id=20")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var tmp = make(map[string]string, len(columns))
	//ret := []map[string]string{}

	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Fetch rows
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//fmt.Println(columns[i], ": ", value)
			tmp[columns[i]] = value
			//tmp = map[string]string{columns[i]: value}
		}
		result = append(result, tmp)
		//fmt.Println("-----------------------------------")
	}
	return
}

func Exec(str string, sUserId string, sMoney string) (insertId int64, affectedRows int64, err error) {
	stmtIns, _ := MyDb.Prepare(str)
	res, _ := stmtIns.Exec(sMoney, sUserId)
	insertId, err = res.LastInsertId()
	affectedRows, err = res.RowsAffected()
	return
}

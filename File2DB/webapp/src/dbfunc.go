package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

var userid, password, server, database *string

func initFlags() {
	user := envData.msSQLUser
	pass := envData.msSQLPass
	svr := envData.msSQLServer
	dbname := envData.msSQLDBName

	userid = flag.String("U", user, "login_id")
	password = flag.String("P", pass, "password")
	server = flag.String("S", svr, "server_name[\\instance_name]")
	database = flag.String("d", dbname, "db_name")

}

func getDBConnection() (*sql.DB, error) {
	// if flags are already initialized, no need to initialize them again
	if flag.Lookup("U") == nil {
		initFlags()
	}
	flag.Parse()

	dsn := "server=" + *server + ";user id=" + *userid + ";password=" + *password + ";database=" + *database
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Cannot connect DB: Error Message: " + err.Error())
		return db, err
	} else {
		fmt.Println("DB Connected.")
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		fmt.Println("Unable to ping DB: Error Message: " + err.Error())
		return db, err
	} else {
		fmt.Println("Ping successful.")
	}

	return db, nil
}

func writeDataIntoDB(data string) error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}

	defer db.Close()

	// insert data into DB table
	query := "INSERT INTO " + envData.msSQLDBName + ".dbo.FILE2DBDATA(DATA) VALUES(?)"
	err = execute(db, query, data)

	return err
}

func readDatafromDB() (string, error) {
	retString := ""
	db, err := getDBConnection()
	if err != nil {
		return retString, err
	}
	//fmt.Println(str)

	defer db.Close()

	query := "SELECT ID, DATA FROM " + envData.msSQLDBName + ".dbo.FILE2DBDATA"
	retString, err = queryDB(db, query)
	if err != nil {
		return retString, err
	}

	return retString, nil
}

func execute(db *sql.DB, query string, data string) error {
	result, err := db.Exec(query, data)
	if err != nil {
		return err
	}

	rowcount, _ := result.RowsAffected()
	fmt.Println("No of rows afftected: ", rowcount)

	return nil
}

func queryDB(db *sql.DB, query string) (string, error) {
	retString := ""
	rows, err := db.Query(query)
	if err != nil {
		return retString, err
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return retString, err
	}

	if cols == nil {
		return retString, nil
	}

	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	// structure for db data
	type File2DBData struct {
		ID   string `json:"ID"`
		Data string `json:"DATA"`
	}
	// slice of file2DBData structure
	dbData := make([]File2DBData, 0)

	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var data File2DBData
		for i := 0; i < len(vals); i++ {
			strData := getRowValue(vals[i].(*interface{}))

			// var data File2DBData
			switch i {
			case 0: // ID
				data.ID = strData
			case 1: // DATA
				data.Data = strData
			}
		}
		dbData = append(dbData, data)
	}

	jsonBytes, err := json.MarshalIndent(dbData, "", " ")
	if err != nil {
		return retString, err
	}
	retString = string(jsonBytes)

	// if rows.Err() != nil {
	// 	return retString, rows.Err()
	// }

	return retString, nil
}

func getRowValue(pval *interface{}) string {
	retString := ""

	switch v := (*pval).(type) {
	// case nil:
	// 	fmt.Print("NULL")
	// case bool:
	// 	if v {
	// 		fmt.Print("1")
	// 	} else {
	// 		fmt.Print("0")
	// 	}
	// case []byte:
	// 	fmt.Print(string(v))
	// 	retString = string(v)
	// case int:
	// 	fmt.Print(v)
	// case time.Time:
	// 	fmt.Print(v.Format("2006-01-02 15:04:05.999"))
	default:
		//fmt.Print(v)
		retString = fmt.Sprint(v)
		//fmt.Print(retString)
	}
	return retString
}

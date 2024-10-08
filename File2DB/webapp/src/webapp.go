package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

type envVarData struct {
	envName, datafilepath, port                    string
	msSQLServer, msSQLDBName, msSQLUser, msSQLPass string
}

var envData envVarData

func main() {
	// Load Environment variable data
	initData(&envData)

	// fmt.Println("ENV_NAME = ", envData.envName)
	//fmt.Println("DATA_FILE_PATH = ", envData.datafilepath)
	// fmt.Println("PORT = ", envData.port)
	// fmt.Println("MSSQL_SERVER = ", envData.msSQLServer)
	// fmt.Println("MSSQL_DBNAME = ", envData.msSQLDBName)
	// fmt.Println("MSSQL_USER = ", envData.msSQLUser)
	// fmt.Println("MSSQL_PASS = ", envData.msSQLPass)

	http.HandleFunc("/", getRoot)                //This is just to check, if application is up and running
	http.HandleFunc("/processdata", processData) // processdata API
	http.HandleFunc("/showdbdata", showDBData)   // showdbdata API

	fmt.Println("Listening at port :", envData.port)
	http.ListenAndServe(":"+envData.port, nil)
}

// check the OS
func isWindowsOS() bool {
	return runtime.GOOS == "windows"
}

// Data will be read from .env file on windows and from env variables on ubuntu
func initData(envData *envVarData) {
	if isWindowsOS() {
		envFile, _ := godotenv.Read("..\\config\\.env")
		envData.envName = envFile["ENV_NAME"]
		envData.datafilepath = envFile["DATA_FILE_PATH"]
		envData.port = envFile["PORT"]
		envData.msSQLServer = envFile["MSSQL_SERVER"]
		envData.msSQLDBName = envFile["MSSQL_DBNAME"]
		envData.msSQLUser = envFile["MSSQL_USER"]
		envData.msSQLPass = envFile["MSSQL_PASS"]
	} else {
		envData.envName = os.Getenv("ENV_NAME")
		envData.datafilepath = os.Getenv("DATA_FILE_PATH")
		envData.port = os.Getenv("PORT")
		envData.msSQLServer = os.Getenv("MSSQL_SERVER")
		envData.msSQLDBName = os.Getenv("MSSQL_DBNAME")
		envData.msSQLUser = os.Getenv("MSSQL_USER")
		envData.msSQLPass = os.Getenv("MSSQL_PASS")
	}

}

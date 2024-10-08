package main

import (
	"fmt"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Congratulations!! DataInputApp is running.\n\n")
	fmt.Fprintf(w, "ENV_NAME=%s\n", envData.envName)
}

func writeData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ENV_NAME=%s\n\n", envData.envName)

	val := r.URL.Query()
	if val.Has("filename") && val.Has("data") {
		//fmt.Fprintf(w, "%s\n\n", val.Encode())

		filename := val.Get("filename")
		//fmt.Fprintf(w, "filename = %s\n", filename)

		data := val.Get("data")
		//fmt.Fprintf(w, "data = %s\n", data)

		// create dir if doesn't exist already
		if _, err := os.Stat(envData.datafilepath); os.IsNotExist(err) {
			err := os.Mkdir(envData.datafilepath, 0777)
			if err != nil {
				fmt.Println(err)
			}
		}

		// Create file
		completeFilePath := envData.datafilepath + filename
		file, err := os.Create(completeFilePath)
		if err != nil {
			fmt.Printf("Unable to create file at %s. Error Message: %s\n", completeFilePath, err)
			return
		}
		defer file.Close()

		// write data into the file
		nums, err := file.WriteString(data)
		if err != nil {
			fmt.Println("Unable to write data into file. Error Message: ", err)
			return
		}
		fmt.Println("Size of data written - ", nums)

		fmt.Fprintln(w, "Successfully written data into file ", completeFilePath)
	} else {
		fmt.Fprintf(w, "Please provide data like '/writedata?filename=xyz.txt&data=Hello!!User'.\n")

		fmt.Println("Error: API is not called properly. Please call it like /writedata?filename=xyz.txt&data=Hello!!User'.")
	}

}

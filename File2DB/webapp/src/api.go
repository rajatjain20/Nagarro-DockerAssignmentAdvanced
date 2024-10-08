package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Congratulations!! WebApp is running.\n\n")
	fmt.Fprintf(w, "ENV_NAME=%s\n", envData.envName)
}

func processData(w http.ResponseWriter, r *http.Request) {
	// fetch all the files exist in the directory to get processed
	dirEntry, err := os.ReadDir(envData.datafilepath)
	if err != nil {
		fmt.Println("Error in processData - ReadDir func: Error Message: ", err.Error())
		fmt.Fprintf(w, "processdata not completed. Please check logs.")
		return
	}

	countFileProcessed := 0
	// loop throgh each file and write data into DB
	for _, dir := range dirEntry {
		// if it is directory, continue. We only want to process files
		if dir.IsDir() {
			continue
		}

		filename := dir.Name()
		byteData, err := os.ReadFile(envData.datafilepath + filename)
		if err != nil {
			fmt.Println("Error in processData - ReadFile func: Error Message: ", err.Error())
			fmt.Fprintf(w, "processdata not completed. Please check logs.")
			return
		}

		data := string(byteData)

		err = writeDataIntoDB(data)
		if err != nil {
			fmt.Println("Error in processData - writeDataIntoDB func: Error Message: ", err.Error())
			fmt.Fprintf(w, "processdata not completed. Please check logs.")
			return
		}

		// Move the processed file to processed dir
		processedDirName := "processed/"
		if isWindowsOS() {
			processedDirName = "processed\\"
		}
		dstPath := envData.datafilepath + processedDirName
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			err := os.Mkdir(dstPath, 0777)
			if err != nil {
				fmt.Println(err)
			}
		}
		moveFile(envData.datafilepath+filename, dstPath+filename)

		countFileProcessed++
	}

	fmt.Fprintf(w, "%d files are processed.\n", countFileProcessed)
}

func showDBData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("showDBData Started....")

	jsondata, err := readDatafromDB()
	if err != nil {
		fmt.Println("Error in showDBData - readDatafromDB func: Error Message: ", err.Error())
	} else {
		fmt.Fprintln(w, jsondata)
		fmt.Println(jsondata)
	}

	fmt.Println("showDBData Completed....")
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		fmt.Printf("Couldn't open source file: %v", err)
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("Couldn't open dest file: %v", err)
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		fmt.Printf("Couldn't copy to dest from source: %v", err)
		return err
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		fmt.Printf("Couldn't remove source file: %v", err)
		return err
	}
	return nil
}

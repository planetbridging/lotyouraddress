package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readLocal(path string) []string {
	var lstContents []string
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		lstContents = append(lstContents, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	//runtime.GC()
	return lstContents
}

func readFindAllPostcodes(path string) []string {
	//var lstContents []string
	var lstData []string
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		//lstContents = append(lstContents, scanner.Text())
		line := scanner.Text()
		var r_split = strings.Split(line, "|")
		if !strings.Contains(r_split[26], "POSTCODE") {
			found := contains(lstData, r_split[26]+","+r_split[24])
			if !found {
				//lstContents = append(lstContents, r_split[26])
				lstData = append(lstData, r_split[26]+","+r_split[24])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	return lstData
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func saveFile(d []string, location string) {
	f, err := os.Create(location)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	//d := []string{"Welcome to the world of Go1.", "Go is a compiled language.", "It is easy to learn Go."}

	for _, v := range d {
		fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file written successfully")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func folderExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func createFolder(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
}

func getFolders(path string) []string {
	var lst []string
	files, err := ioutil.ReadDir(path)

	if err != nil {

		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			lst = append(lst, f.Name())
		}

		//fmt.Println(f.Name())
	}
	return lst
}

func getFiles(path string) []string {
	var lst []string
	files, err := ioutil.ReadDir(path)

	if err != nil {

		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			lst = append(lst, f.Name())
		}

		//fmt.Println(f.Name())
	}
	return lst
}

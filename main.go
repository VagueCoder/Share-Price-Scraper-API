package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"bufio"
	"strings"
	"strconv"
	"syscall"
	"net/http"
	"os/signal"

	"github.com/gorilla/mux"
)

type File struct {
	directory	string
	name		string
}

func filename() string {
	dtlayout := "02-01-2006 3.04.05 PM"
	return "Share-Price-Scraper-API Export " + time.Now().Format(dtlayout) + ".csv"
}

func getStatus(statusfile string) string {
	file, err := os.Open(statusfile)
    if err != nil {
        fmt.Printf("OS Error: Failed to read file %s with error: %v\n", statusfile, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, strings.TrimSpace(scanner.Text()))
    }
	
	pointer := 1
	for true {
		if lines[len(lines)-pointer] == "" {
			pointer++
		} else {
			break
		}
	}
   	return lines[len(lines)-pointer]
}

func main() {
	file := &File{directory:"/datastore/"}
	var filenames []string

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		for _, JSON_file := range filenames {
			file.name = JSON_file
			deleteCSV(file)
		}
		fmt.Println("\rClosing Gracefully!")
		os.Exit(0)
	}()
	
	router := mux.NewRouter()
	router.HandleFunc("/", 
		func(writer http.ResponseWriter, request *http.Request) {
			file.name = filename()
			filenames = append(filenames, file.name)
			createCSV(file)

			html := "<div style=\"font-family:Apple Chancery,cursive;text-align:center;border-radius:25px;padding:20px;margin:150px;width:1000px;border: 5px solid #004080;\"><p style=\"color:004080;margin:0;padding:0;font-size:30px;\">Click: <a href=\"/download\">%s</a></p><br><button id=\"status\" style=\"font-size:15px;border-radius:5px;margin:10px;padding:10px;background-color:#b3d9ff;border:3px solid #004080;\">Click Here for Current Scraper Status</button><script>document.getElementById(\"status\").onclick = function(){alert(\"%s\");}</script><br><p style=\"color:red;font-family:Apple Chancery,cursive;text-align:center;margin:0;padding:0;font-size:20px;\">Note: You can download the file multiple times till you reload. Reloading the page gives latest file.</p></div>"
			fmt.Fprintf(writer, html, file.name, getStatus(file.directory + "stats.txt"))
		
			return
		},
	)

	router.HandleFunc("/download", 
		func(writer http.ResponseWriter, request *http.Request) {
			fmt.Printf("\r%s%s\t\tDownloaded\t\t%v\n", file.directory, file.name, time.Now())
			Openfile, err := os.Open(file.directory+file.name)
			defer Openfile.Close()
			if err != nil {
				http.Error(writer, "File not found.", 404)
				return
			}
		
			FileHeader := make([]byte, 512)
			Openfile.Read(FileHeader)
			FileContentType := http.DetectContentType(FileHeader)
		
			FileStat, _ := Openfile.Stat()
			FileSize := strconv.FormatInt(FileStat.Size(), 10)
		
			writer.Header().Set("Content-Disposition", "attachment; filename="+file.name)
			writer.Header().Set("Content-Type", FileContentType)
			writer.Header().Set("Content-Length", FileSize)
		
			Openfile.Seek(0, 0)
			io.Copy(writer, Openfile)
			return
		},
	)
	fmt.Println("The application is up and listening on port 8000.")
	fmt.Println("\rFilename\t\tRequest Type\t\tTimestamp")
	fmt.Printf("HTTP Error: Error at Listen and Server: %v", http.ListenAndServe(":8000", router))
}
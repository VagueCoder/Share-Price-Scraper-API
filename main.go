package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"strconv"
	"syscall"
	"net/http"
	"os/signal"

	"github.com/gorilla/mux"
)

type File struct {
	name string
}

func filename() string {
	dtlayout := "02-01-2006 3.04.05 PM"
	return "Share-Price-Scraper-API Export " + time.Now().Format(dtlayout) + ".csv"
}

func main() {
	file := &File{}
	var filenames []string

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rCtrl+C pressed in Terminal")
		for _, file := range filenames {
			deleteCSV("DataStore/"+file)
		}
		os.Exit(0)
	}()
	
	fmt.Println("Filename\t\tRequest Type\t\tTimestamp")
	router := mux.NewRouter()
	router.HandleFunc("/", 
		func(writer http.ResponseWriter, request *http.Request) {
			file.name = filename()
			filenames = append(filenames, file.name)
			createCSV("DataStore/"+file.name)
			fmt.Fprintf(writer, "<div style=\"border-radius: 25px;padding:20px;margin:150px;width:1000px;background-color:#b3d9ff;border: 5px solid #004080;\"><p style=\"color:004080;margin:0;padding:0;font-family:Apple Chancery,cursive;text-align:center;font-size:30px;\">Click: <a href=\"/download\">%s</a></p><br><p style=\"color:red;font-family:Apple Chancery,cursive;text-align:center;margin:0;padding:0;font-size:20px;\">Note: You can download the file multiple times. Reload the Page to download the latest file.</p></div>", file.name)
		
			return
		},
	)

	router.HandleFunc("/download", 
		func(writer http.ResponseWriter, request *http.Request) {
			fmt.Printf("DataStore/%+v\t\tDownloaded\t\t%v\n", file.name, time.Now())
			Openfile, err := os.Open("DataStore/"+file.name)
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
	fmt.Printf("HTTP Error: Error at Listen and Server: %v", http.ListenAndServe(":8000", router))

}
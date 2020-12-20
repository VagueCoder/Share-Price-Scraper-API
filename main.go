package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"strconv"
	"net/http"

	"github.com/gorilla/mux"
)

type File struct {
	name string
}

func filename() string {
	dtlayout := "02-01-2006 3:04:05 PM"
	return "Share-Price-Scraper-API Export " + time.Now().Format(dtlayout) + ".csv"
}

func (file *File) HomePage(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("%+v", file)
	_, err := os.Create(file.name)
	if err != nil {
		fmt.Printf("OS Error: File Creation Failed, %s\n", file.name)
	}

	fmt.Fprintf(writer, "<div style=\"border-radius: 25px;padding:20px;margin:150px;width:1000px;background-color:#b3d9ff;border: 5px solid #004080;\"><h1 style=\"color:004080;font-family:Apple Chancery,cursive;text-align:center;\">The file '%s' should start downloading shortly.</h1></div>", file.name)
	time.Sleep(8 * time.Second)
	http.Redirect(writer, request, "/download", 200)
}

func (file *File) FileDownloadClient(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("%+v", file)
	Openfile, err := os.Open(file.name)
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
}

func main() {
	file := &File{name: filename()}
	router := mux.NewRouter()
	router.HandleFunc("/", file.HomePage)
	router.HandleFunc("/download", file.FileDownloadClient)
	err := http.ListenAndServe(":8000", router)

	if err != nil {
		fmt.Println(err)
	}
}
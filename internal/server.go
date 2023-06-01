package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func parse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	queryParams := r.URL.Query()
	templateName := queryParams.Get("loader_id")

	err := r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("http parse: error retrieving the file: %v", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	templateFile, handler, err := r.FormFile("template")
	if err != nil {
		log.Printf("http parse: error retrieving the file: %v", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer templateFile.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("http parse: cannot read file: %v", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	templateBytes, err := io.ReadAll(templateFile)
	if err != nil {
		log.Printf("http parse: cannot template read file: %v", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = os.Stat(os.Getenv("UPLOAD_PATH") + templateName)
	if err != nil {
		err = os.Mkdir(os.Getenv("UPLOAD_PATH")+templateName, 775)
		if err != nil {
			log.Printf("http parse: cannot create dir: %v", err.Error())
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}

	tempFile, err := ioutil.TempFile(os.Getenv("UPLOAD_PATH")+templateName, "*.xlsx")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	tempFile.Write(fileBytes)

	p := NewFileParser(tempFile.Name(), templateBytes)
	jsonString, err := p.Parse()
	if err != nil {
		log.Printf("http parse: cannot parse file: %v", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonString)
}

func NewServer(port string) {
	http.HandleFunc("/", parse)
	log.Println("server started on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

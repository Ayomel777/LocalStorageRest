package main

import (
	"log"
	"net/http"

	"Lab3_KSIS/internal/handler"
	"Lab3_KSIS/internal/storage"
)

func main() {
	store, err := storage.NewStorage("./storage")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/files/", handler.NewServer(store))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//PUT — создать файл
//curl -X PUT http://localhost:8080/files/hello.txt -d "Hello, World!"

//GET — прочитать
//curl http://localhost:8080/files/hello.txt

//POST — дописать в конец
//curl -X POST http://localhost:8080/files/hello.txt -d " Appended"

//COPY
//curl -X COPY http://localhost:8080/files/hello.txt -H "Destination: /files/backup.txt"

//MOVE
//curl.exe -X MOVE http://localhost:8080/files/backup.txt -H "Destination: /files/moved.txt"

//DELETE
//curl -X DELETE http://localhost:8080/files/hello.txt

//GET когда удалили (404)
//curl -v http://localhost:8080/files/hello.txt

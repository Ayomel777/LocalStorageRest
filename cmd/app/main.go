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
//curl.exe -X PUT http://localhost:8080/files/hello.txt -d "Hello, World!"

//PUT - проверка на разбитие чанков
//curl.exe -X PUT --data-binary "@D:\BSUIR\interface.mp4" http://localhost:8080/files/interface.mp4
//GET — прочитать
//curl.exe http://localhost:8080/files/hello.txt

//POST — дописать в конец
//curl.exe -X POST http://localhost:8080/files/hello.txt -d " Appended"

//COPY
//curl.exe -X COPY http://localhost:8080/files/hello.txt -H "Destination: /files/backup.txt"

//MOVE
//curl.exe -X MOVE http://localhost:8080/files/backup.txt -H "Destination: /files/moved.txt"

//DELETE
//curl.exe -X DELETE http://localhost:8080/files/hello.txt

//GET когда удалили (404)
//curl.exe -v http://localhost:8080/files/hello.txt

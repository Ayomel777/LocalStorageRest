package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const base = "http://localhost:8080"
//работа с большими файлами, разбиение на чанки
func main() {
	fmt.Println("PUT - создаём файлы")
	put("/files/hello.txt", "Hello, World!")
	put("/files/notes/note.txt", "This is a note.")
	put("/files/docs/report.txt", "Report content line 1.")

	fmt.Println("GET - читаем файлы")
	get("/files/hello.txt")
	get("/files/notes/note.txt")

	fmt.Println("GET - файл не существует (ожидаем 404)")
	get("/files/nonexistent.txt")

	fmt.Println("POST - дописываем в конец")
	post("/files/hello.txt", " Appended text.")
	get("/files/hello.txt")

	fmt.Println("PUT - перезаписываем")
	put("/files/hello.txt", "Completely new content.")
	get("/files/hello.txt")

	fmt.Println("COPY - копируем hello.txt -> backup/hello.txt")
	copyFile("/files/hello.txt", "/files/backup/hello.txt")
	get("/files/backup/hello.txt")

	fmt.Println("MOVE - перемещаем notes/note.txt -> archive/note.txt")
	moveFile("/files/notes/note.txt", "/files/archive/note.txt")
	get("/files/archive/note.txt")

	fmt.Println("GET - notes/note.txt после перемещения (ожидаем 404)")
	get("/files/notes/note.txt")

	fmt.Println("DELETE - удаляем hello.txt")
	deleteFile("/files/hello.txt")
	get("/files/hello.txt")
}

func get(path string) {
	req, _ := http.NewRequest(http.MethodGet, base+path, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("GET %-40s → ERROR: %v\n", path, err)
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Printf("GET %-40s → %s | %s\n", path, resp.Status, strings.TrimSpace(string(b)))
}

func put(path, body string) {
	req, _ := http.NewRequest(http.MethodPut, base+path, strings.NewReader(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("PUT %-40s → ERROR: %v\n", path, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("PUT %-40s → %s\n", path, resp.Status)
}

func post(path, body string) {
	req, _ := http.NewRequest(http.MethodPost, base+path, strings.NewReader(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("POST %-40s → ERROR: %v\n", path, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("POST %-40s → %s\n", path, resp.Status)
}

func deleteFile(path string) {
	req, _ := http.NewRequest(http.MethodDelete, base+path, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("DELETE %-40s → ERROR: %v\n", path, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("DELETE %-40s → %s\n", path, resp.Status)
}

func copyFile(src, dst string) {
	req, _ := http.NewRequest("COPY", base+src, nil)
	req.Header.Set("Destination", dst)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("COPY %-40s → ERROR: %v\n", src, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("COPY %-40s → %s\n", src+" → "+dst, resp.Status)
}

func moveFile(src, dst string) {
	req, _ := http.NewRequest("MOVE", base+src, nil)
	req.Header.Set("Destination", dst)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("MOVE %-40s → ERROR: %v\n", src, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("MOVE %-40s → %s\n", src+" → "+dst, resp.Status)
}

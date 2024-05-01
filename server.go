package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	//http.HandleFunc("/", handler)
	fs := http.FileServer(http.Dir(""))
	http.Handle("/", Auth(fs, "charo", "12345"))
	http.HandleFunc("/files", FileList)
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
func FileList(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading %v", err)
	}
	for _, file := range files {
		w.Write([]byte(file.Name() + "\n"))
	}
}
func Auth(next http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}

		Log(username + " fetched " + r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func Log(s string) {
	date := time.Now()
	file, err := os.OpenFile("log/"+date.Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	logger := log.New(file, "-", log.Ldate|log.Ltime)
	logger.Println(s)
}

/*func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
	}
	email := r.Form.Get("email")
	fmt.Println(email)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Form submitted")
}
func handler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("sales.json")
	if err != nil {
		http.Error(w, "Error opening JSON file", http.StatusInternalServerError)
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error transmitting file", http.StatusInternalServerError)
		return
	}
}*/

package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/query", withQueryParam)
	http.HandleFunc("/form", withFormValue)
	http.HandleFunc("/create", personCreate)
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Println(http.ListenAndServe(":8090", nil))

}

// func (srv *Server) ListenAndServe() error {
//     // ...

//     ln, err := net.Listen("tcp", addr)
//     if err != nil {
//         return err
//     }
//     return srv.Serve(ln)
// }

func pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func withQueryParam(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("name")
	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(param))
}

func withFormValue(w http.ResponseWriter, r *http.Request) {
	param := r.FormValue("name")
	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(param))
}

func personCreate(w http.ResponseWriter, r *http.Request) {
	var p Person
	w.Header().Set("Content-type", "application/json")
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Eğer kullanıcı decode etmek istediğimiz struct'ta bulunmayan değerler gönderiyorsa(name ve email hariç) o zaman response olarak 400 bad request kodu ile json: unknown field "field_name" mesajını dönebiliriz
	err := dec.Decode(&p)

	if r.Method != "POST" {
		http.Error(w, "You can use POST method only", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, _ := json.Marshal(p)

	w.Write(resp)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

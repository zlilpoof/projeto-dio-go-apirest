package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Cliente struct {
	Id           string    `json:"id,omitempty"`
	Nome         string    `json:"nome,omitempty"`
	Idade        string    `json:"idade,omitempty"`
	BoloFavorito string    `json:"bolofavorito,omitempty"`
	Endereco     *Endereco `json:"endereco,omitempty"`
}

type Endereco struct {
	Cidade string `json:"cidade,omitempty"`
	Estado string `json:"estado,omitempty"`
}

var clientes []Cliente

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", AbrirPagina).Methods("GET")
	router.HandleFunc("/clientes", CreateCliente).Methods("POST")
	router.HandleFunc("/clientes/{id}", GetCliente).Methods("GET")
	router.HandleFunc("/clientes", GetClientes).Methods("GET")
	router.HandleFunc("/clientes/{id}", DeletarCliente).Methods("DELETE")

	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", router)
}

func generateUUID() string {
	return uuid.New().String()
}

func AbrirPagina(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pagina.html"))
	tmpl.Execute(w, nil)
}

func CreateCliente(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	nome := r.FormValue("nome")
	idade := r.FormValue("idade")
	boloFavorito := r.FormValue("bolofavorito")
	cidade := r.FormValue("endereco.cidade")
	estado := r.FormValue("endereco.estado")

	cliente := Cliente{
		Id:           generateUUID(),
		Nome:         nome,
		Idade:        idade,
		BoloFavorito: boloFavorito,
		Endereco: &Endereco{
			Cidade: cidade,
			Estado: estado,
		},
	}

	clientes = append(clientes, cliente)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetCliente(w http.ResponseWriter, r *http.Request) {
	parametro := mux.Vars(r)
	for _, pos := range clientes {
		if pos.Id == parametro["id"] {
			json.NewEncoder(w).Encode(pos)
			return
		}
	}

	json.NewEncoder(w).Encode(&Cliente{})
}

func GetClientes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

func DeletarCliente(w http.ResponseWriter, r *http.Request) {
	parametro := mux.Vars(r)
	for i, pos := range clientes {
		if pos.Id == parametro["id"] {
			clientes = append(clientes[:i], clientes[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

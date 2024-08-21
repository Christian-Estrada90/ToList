package main

import (
	userService "crud-golang/services"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Si es una solicitud OPTIONS, solo respondemos con el encabezado de CORS y devolvemos
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", userService.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users", userService.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", userService.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", userService.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", userService.DeleteUser).Methods(http.MethodDelete)
	//Actividades
	router.HandleFunc("/actividades", userService.CreateActividad).Methods(http.MethodPost)
	router.HandleFunc("/actividades/{id_usuario}", userService.GetActividad).Methods(http.MethodGet)
	router.HandleFunc("/actividades/{id}", userService.UpdateActividad).Methods(http.MethodPut)
	router.HandleFunc("/actividades/{id}", userService.DeleteActividad).Methods(http.MethodDelete)
	router.HandleFunc("/FinalizarActividad/{id}", userService.FinalizarActividad).Methods(http.MethodPut)

	//Categorias
	router.HandleFunc("/categorias", userService.GetCategorias).Methods(http.MethodGet)

	//Estados
	router.HandleFunc("/estados", userService.GetEstados).Methods(http.MethodGet)

	//IMPORTANCIAS
	router.HandleFunc("/importancias", userService.GetImportancias).Methods(http.MethodGet)

	corsRouter := enableCors(router)

	fmt.Print("Listenning on 5000")
	log.Fatal(http.ListenAndServe(":5000", corsRouter))
}

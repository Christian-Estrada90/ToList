package userService

import (
	"crud-golang/database"
	"encoding/json"
	"net/http"
)

type Categoria struct {
	Id     uint32 `json:"id"`
	Nombre string `json:"nombre"`
}

func GetCategorias(w http.ResponseWriter, r *http.Request) {
	db, err := database.DbConnection()
	if handleGenericErrorCat(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nombre FROM categoria WHERE activo = 'S'")
	if handleGenericErrorCat(w, "Failed to retrieve actividad!", err) {
		return
	}
	defer rows.Close()

	var categorias []Categoria
	for rows.Next() {
		var categoria Categoria
		if err := rows.Scan(&categoria.Id, &categoria.Nombre); handleGenericErrorCat(w, "Failed to scan actividad!", err) {
			return
		}
		categorias = append(categorias, categoria)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(categorias); handleGenericErrorCat(w, "Failed to convert actividades to JSON!", err) {
		return
	}
}

func handleGenericErrorCat(w http.ResponseWriter, errorMessage string, err error) bool {
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return true
	}
	return false
}

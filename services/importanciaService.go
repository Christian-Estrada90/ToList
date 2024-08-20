package userService

import (
	"crud-golang/database"
	"encoding/json"
	"net/http"
)

type Importancia struct {
	Id     uint32 `json:"id"`
	Nombre string `json:"nombre"`
}

func GetImportancias(w http.ResponseWriter, r *http.Request) {
	db, err := database.DbConnection()
	if handleGenericErrorImp(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nombre FROM importancia WHERE activo = 'S'")
	if handleGenericErrorImp(w, "Failed to retrieve actividad!", err) {
		return
	}
	defer rows.Close()

	var importancias []Importancia
	for rows.Next() {
		var importancia Importancia
		if err := rows.Scan(&importancia.Id, &importancia.Nombre); handleGenericErrorImp(w, "Failed to scan actividad!", err) {
			return
		}
		importancias = append(importancias, importancia)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(importancias); handleGenericErrorImp(w, "Failed to convert actividades to JSON!", err) {
		return
	}
}

func handleGenericErrorImp(w http.ResponseWriter, errorMessage string, err error) bool {
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return true
	}
	return false
}

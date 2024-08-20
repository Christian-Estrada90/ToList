package userService

import (
	"crud-golang/database"
	"encoding/json"
	"net/http"
)

type Estado struct {
	Id     uint32 `json:"id"`
	Nombre string `json:"nombre"`
}

func GetEstados(w http.ResponseWriter, r *http.Request) {
	db, err := database.DbConnection()
	if handleGenericErrorEst(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nombre FROM estado WHERE activo = 'S'")
	if handleGenericErrorEst(w, "Failed to retrieve actividad!", err) {
		return
	}
	defer rows.Close()

	var estados []Estado
	for rows.Next() {
		var estado Estado
		if err := rows.Scan(&estado.Id, &estado.Nombre); handleGenericErrorEst(w, "Failed to scan actividad!", err) {
			return
		}
		estados = append(estados, estado)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(estados); handleGenericErrorEst(w, "Failed to convert actividades to JSON!", err) {
		return
	}
}

func handleGenericErrorEst(w http.ResponseWriter, errorMessage string, err error) bool {
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return true
	}
	return false
}

package userService

import (
	"crud-golang/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Actividad struct {
	ID                uint32 `json:"id"`
	IDUsuario         int    `json:"idusuario"`
	IDCategoria       int    `json:"idcategoria"`
	IDImportancia     int    `json:"idimportancia"`
	IDEstado          int    `json:"idestado"`
	Nombre            string `json:"nombre"`
	Descripcion       string `json:"descripcion"`
	FechaFinaliza     string `json:"fechafinaliza"`
	FechaRealFinaliza string `json:"fecharealfinaliza"`
	Activo            string `json:"activo"`
	FechaRegistro     string `json:"fecharegistro"`
}

// CreateActividad inserts a user into the database.
func CreateActividad(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if handleGenericErrorAc(w, "Failed to read request body!", err) {
		return
	}

	var newActividad Actividad
	if err = json.Unmarshal(body, &newActividad); handleGenericErrorAc(w, "Failed to unmarshal request body!", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericErrorAc(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO actividad (id_usuario, id_categoria, id_importancia, id_estado, nombre, descripcion, fecha_finaliza, fecha_real_finaliza, activo, fecha_registro) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, SYSDATE())")
	if handleGenericErrorAc(w, "Failed to create statement!", err) {
		return
	}
	defer statement.Close()

	result, err := statement.Exec(
		newActividad.IDUsuario,
		newActividad.IDCategoria,
		newActividad.IDImportancia,
		newActividad.IDEstado,
		newActividad.Nombre,
		newActividad.Descripcion,
		newActividad.FechaFinaliza,
		newActividad.FechaRealFinaliza,
		newActividad.Activo,
	)
	if handleGenericErrorAc(w, "Failed to execute statement!", err) {
		return
	}

	createdID, err := result.LastInsertId()
	if handleGenericErrorAc(w, "Failed to retrieve last insert ID!", err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Successfully created actividad %d", createdID)
}

// GetActividad retrieves all users from the database.
func GetActividad(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	IDUsuario, err := strconv.ParseUint(params["id_usuario"], 10, 32)
	if handleGenericErrorAc(w, "Failed to convert parameter to integer", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericErrorAc(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM actividad WHERE id_usuario = ? ORDER BY id_estado asc", IDUsuario)
	if handleGenericErrorAc(w, "Failed to retrieve actividad!", err) {
		return
	}
	defer rows.Close()

	var actividades []Actividad
	for rows.Next() {
		var actividad Actividad
		if err := rows.Scan(
			&actividad.ID,
			&actividad.IDUsuario,
			&actividad.IDCategoria,
			&actividad.IDImportancia,
			&actividad.IDEstado,
			&actividad.Nombre,
			&actividad.Descripcion,
			&actividad.FechaFinaliza,
			&actividad.FechaRealFinaliza,
			&actividad.Activo,
			&actividad.FechaRegistro); handleGenericErrorAc(w, "Failed to scan actividad!", err) {
			return
		}
		actividades = append(actividades, actividad)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(actividades); handleGenericErrorAc(w, "Failed to convert actividades to JSON!", err) {
		return
	}
}

// UpdateActividad updates the data of a user in the database.
func UpdateActividad(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if handleGenericErrorAc(w, "Failed to convert parameter to integer", err) {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if handleGenericErrorAc(w, "Failed to read request body!", err) {
		return
	}

	var updatedActivity Actividad
	if err = json.Unmarshal(body, &updatedActivity); handleGenericErrorAc(w, "Failed to unmarshal request body!", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericErrorAc(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE actividad SET id_usuario= ?, id_categoria = ?, id_importancia = ?, id_estado = ?, nombre = ?, descripcion = ?, fecha_finaliza = ?, fecha_real_finaliza = ?, activo = ?, fecha_registro = sysdate() WHERE id = ?")
	if handleGenericErrorAc(w, "Failed to create statement!", err) {
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(
		updatedActivity.IDUsuario,
		updatedActivity.IDCategoria,
		updatedActivity.IDImportancia,
		updatedActivity.IDEstado,
		updatedActivity.Nombre,
		updatedActivity.Descripcion,
		updatedActivity.FechaFinaliza,
		updatedActivity.FechaRealFinaliza,
		updatedActivity.Activo,
		ID); handleGenericErrorAc(w, "Failed to execute statement!", err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteActividad deletes a user from the database.
func DeleteActividad(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if handleGenericErrorAc(w, "Failed to convert parameter to integer", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericErrorAc(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM actividad WHERE id = ?")
	if handleGenericErrorAc(w, "Failed to create statement!", err) {
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); handleGenericErrorAc(w, "Failed to execute statement!", err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleGenericErrorAc(w http.ResponseWriter, errorMessage string, err error) bool {
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return true
	}
	return false
}

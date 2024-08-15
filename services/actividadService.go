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
	if handleGenericError(w, "Failed to read request body!", err) {
		return
	}

	var newActividad Actividad
	if err = json.Unmarshal(body, &newActividad); handleGenericError(w, "Failed to unmarshal request body!", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericError(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO actividad (idusuario, idcategoria, idimportancia, idestado, nombre, descripcion, fechafinaliza, fecharealfinaliza, activo, fecharegistro) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if handleGenericError(w, "Failed to create statement!", err) {
		return
	}
	defer statement.Close()

	result, err := statement.Exec(newActividad.IDUsuario, newActividad.IDCategoria, newActividad.IDImportancia, newActividad.IDEstado, newActividad.Nombre, newActividad.Descripcion, newActividad.FechaFinaliza, newActividad.FechaRealFinaliza, newActividad.Activo, newActividad.FechaRegistro)
	if handleGenericError(w, "Failed to execute statement!", err) {
		return
	}

	createdID, err := result.LastInsertId()
	if handleGenericError(w, "Failed to retrieve last insert ID!", err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Successfully created actividad %d", createdID)
}

// GetActividad retrieves all users from the database.
func GetActividad(w http.ResponseWriter, r *http.Request) {
	db, err := database.DbConnection()
	if handleGenericError(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM actividad")
	if handleGenericError(w, "Failed to retrieve actividad!", err) {
		return
	}
	defer rows.Close()

	var actividades []Actividad
	for rows.Next() {
		var actividad Actividad
		if err := rows.Scan(&actividad.ID, &actividad.IDUsuario, &actividad.IDCategoria, &actividad.IDImportancia, &actividad.IDEstado, &actividad.Nombre, &actividad.Descripcion, &actividad.FechaFinaliza, &actividad.FechaRealFinaliza, &actividad.Activo, &actividad.FechaRegistro); handleGenericError(w, "Failed to scan actividad!", err) {
			return
		}
		actividades = append(actividades, actividad)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(actividades); handleGenericError(w, "Failed to convert actividades to JSON!", err) {
		return
	}
}

// GetActividadByID retrieves a user from the database by ID.
func GetActividadByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if handleGenericError(w, "Failed to convert parameter to integer", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericError(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM usuario WHERE id = ?", ID)
	if handleGenericError(w, "Failed to retrieve user "+strconv.FormatUint(ID, 10), err) {
		return
	}
	defer row.Close()

	var user User
	if row.Next() {
		if err := row.Scan(&user.ID, &user.Nombre, &user.Email, &user.Password, &user.Activo, &user.Fecha); handleGenericError(w, "Failed to scan users!", err) {
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); handleGenericError(w, "Failed to convert user to JSON!", err) {
		return
	}
}

// UpdateActividad updates the data of a user in the database.
func UpdateActividad(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if handleGenericError(w, "Failed to convert parameter to integer", err) {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if handleGenericError(w, "Failed to read request body!", err) {
		return
	}

	var updatedUser User
	if err = json.Unmarshal(body, &updatedUser); handleGenericError(w, "Failed to unmarshal request body!", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericError(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE usuario SET name = ?, email = ?, password = ?, activo = ? WHERE id = ?")
	if handleGenericError(w, "Failed to create statement!", err) {
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(updatedUser.Nombre, updatedUser.Email, updatedUser.Password, updatedUser.Activo, updatedUser.Fecha, ID); handleGenericError(w, "Failed to execute statement!", err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteActividad deletes a user from the database.
func DeleteActividad(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if handleGenericError(w, "Failed to convert parameter to integer", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericError(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	statement, err := db.Prepare("DELETE FROM usuario WHERE id = ?")
	if handleGenericError(w, "Failed to create statement!", err) {
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); handleGenericError(w, "Failed to execute statement!", err) {
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

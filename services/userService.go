package userService

import (
	"crud-golang/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	ID       uint32 `json:"id"`
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Activo   string `json:"activo"`
	Fecha    string `json:"fecha"`
}

// CreateUser inserts a user into the database.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if handleGenericError(w, "Failed to read request body!", err) {
		return
	}

	var newUser User
	if err = json.Unmarshal(body, &newUser); handleGenericError(w, "Failed to unmarshal request body!", err) {
		return
	}

	db, err := database.DbConnection()
	if handleGenericError(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO usuario (nombre, email, password, activo, fecha) VALUES (?, ?, ?, ?, ?)")
	if handleGenericError(w, "Failed to create statement!", err) {
		return
	}
	defer statement.Close()

	result, err := statement.Exec(newUser.Nombre, newUser.Email, newUser.Password, newUser.Activo, newUser.Fecha)
	if handleGenericError(w, "Failed to execute statement!", err) {
		return
	}

	createdID, err := result.LastInsertId()
	if handleGenericError(w, "Failed to retrieve last insert ID!", err) {
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Successfully created user %d", createdID)
}

// GetUsers retrieves all users from the database.
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.DbConnection()
	if handleGenericError(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM usuario")
	if handleGenericError(w, "Failed to retrieve users!", err) {
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Nombre, &user.Email, &user.Password, &user.Activo, &user.Fecha); handleGenericError(w, "Failed to scan users!", err) {
			return
		}
		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); handleGenericError(w, "Failed to convert users to JSON!", err) {
		return
	}
}

// GetUserByID retrieves a user from the database by ID.
func GetUserByID(w http.ResponseWriter, r *http.Request) {
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

// LoginUser authenticates a user using their email and password.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse the email and password from the request body.
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	db, err := database.DbConnection()
	if handleGenericError(w, "Failed to connect to the database!", err) {
		return
	}
	defer db.Close()

	// Query the database for the user with the provided email.
	var user User
	row := db.QueryRow("SELECT id, email, password FROM usuario WHERE email = ?", creds.Email)
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	// Compare the provided password with the stored password.
	if user.Password != creds.Password {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Respond with the user's ID if the authentication is successful.
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(struct {
		ID uint64 `json:"id"`
	}{
		ID: uint64(user.ID), // Convertir uint32 a uint64
	})
	if err != nil {
		http.Error(w, "Failed to convert user to JSON", http.StatusInternalServerError)
		return
	}
}

// UpdateUser updates the data of a user in the database.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	statement, err := db.Prepare("UPDATE usuario SET nombre = ?, email = ?, password = ?, activo = ?, fecha = ? WHERE id = ?")
	if handleGenericError(w, "Failed to create statement!", err) {
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(updatedUser.Nombre, updatedUser.Email, updatedUser.Password, updatedUser.Activo, updatedUser.Fecha, ID); handleGenericError(w, "Failed to execute statement!", err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser deletes a user from the database.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
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

func handleGenericError(w http.ResponseWriter, errorMessage string, err error) bool {
	if err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return true
	}
	return false
}

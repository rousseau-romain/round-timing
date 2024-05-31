package api

// import (
// 	"encoding/json"
// 	"examples/model"
// 	"io"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/gorilla/mux"
// )

// type UserRequest struct {
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }
// type UserRequestUpdate struct {
// 	Name     *string `json:"name"`
// 	Email    *string `json:"email"`
// 	Password *string `json:"password"`
// }

// func HandlerApi(r *mux.Router) http.Handler {
// 	r.HandleFunc("/user", getUsers).Methods("GET")
// 	r.HandleFunc("/user", createUser).Methods("POST")
// 	r.HandleFunc("/user/{id}", getUser).Methods("GET")
// 	r.HandleFunc("/user/{id}", updateUser).Methods("PATCH")
// 	r.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")

// 	return r
// }

// func getUsers(w http.ResponseWriter, r *http.Request) {
// 	users, err := model.GetUsers()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	jsonData, err := json.Marshal(users)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	w.Write(jsonData)
// }

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.ParseInt(vars["id"], 10, 64)

// 	user, err := model.GetUser(int(userId))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	jsonData, err := json.Marshal(user)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	w.Write(jsonData)
// }

// func createUser(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		log.Println(err.Error())
// 		return
// 	}

// 	emailFromForm := r.FormValue("email")
// 	userAlreadyExists, err := model.UserExists(emailFromForm)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if userAlreadyExists {
// 		http.Error(w, "Email already exists", http.StatusBadRequest)
// 		return
// 	}

// 	hashedPassword, _ := HashPassword(r.FormValue("password"))

// 	userId, err := model.CreateUser(model.UserCreate{
// 		Name:  r.FormValue("name"),
// 		Email: &emailFromForm,
// 		Hash:  &hashedPassword,
// 	}, false)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	user, err := model.GetUser(int(userId))

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	jsonData, err := json.Marshal(user)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(jsonData)
// }

// func updateUser(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.ParseInt(vars["id"], 10, 64)

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	var userR UserRequestUpdate
// 	err = json.Unmarshal(body, &userR)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	modelUpdate := model.UserUpdate{}
// 	modelUpdate.Name = userR.Name
// 	modelUpdate.Email = userR.Email

// 	if userR.Password != nil {
// 		hashedPassword, _ := HashPassword(*userR.Password)
// 		modelUpdate.Hash = &hashedPassword
// 	}
// 	err = model.UpdateUser(userId, modelUpdate)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	user, err := model.GetUser(int(userId))

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	jsonData, err := json.Marshal(user)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonData)

// }

// func deleteUser(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.ParseInt(vars["id"], 10, 64)

// 	err := model.DeleteUser(userId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

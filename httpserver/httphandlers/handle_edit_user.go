package httphandlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/nkhamm-spb/red_soft_test/storage"
)

type HandlerEditUser struct {
	Storage storage.StorageInterface
}

// @Summary Изменить пользователя
// @Description Изменить пользователя
// @Tags example
// @Accept   json
// @Produce  json
// @Param   id path int true "id пользователя"
// @Param   input body   schemas.EditUser true  "Данные для редактирования"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/users/{id}/edit_user [put]
func (h *HandlerEditUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Тут специально нет десериализации в EditUser что бы было возможно заменить данные на пустые

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Request to edit user with id: %d data to edit: %v\n", id, data)

	editedUser, err := h.Storage.EditUser(r.Context(), id, data)

	if err != nil {
		log.Printf("Error in edit user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(editedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

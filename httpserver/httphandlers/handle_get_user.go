package httphandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/nkhamm-spb/red_soft_test/storage"
)

type HandlerGetUser struct {
	Storage storage.StorageInterface
}

// @Summary Получить данные пользователя по id
// @Description Получить данные пользователя по id
// @Tags example
// @Accept  json
// @Produce  json
// @Param   id path int true "id пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/users/{id}/get_user [get]
func (h *HandlerGetUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Request to get user with id: %d\n", id)

	user, err := h.Storage.GetUserById(r.Context(), id)

	if err != nil {
		log.Printf("Error in get user by id request: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

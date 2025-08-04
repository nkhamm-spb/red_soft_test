package httphandlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nkhamm-spb/red_soft_test/storage"
)

type HandlerGetBySurname struct {
	Storage storage.StorageInterface
}

// @Summary Получить данные пользователя по фамилии
// @Description Получить данные пользователя по фамилии
// @Tags example
// @Accept  json
// @Produce  json
// @Param   surname path string true "Фамилия пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/users/get_by_surname/{surname} [get]
func (h *HandlerGetBySurname) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	surname := vars["surname"]
	log.Printf("Request to get user with surname: %s", surname)

	user, err := h.Storage.GetUserBySurname(r.Context(), surname)

	if err != nil {
		log.Printf("Error in get user by surname request: %v\n", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

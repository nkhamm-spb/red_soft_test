package httphandlers

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/nkhamm-spb/red_soft_test/storage"
)

type HandlerGetAll struct {
	Storage storage.StorageInterface
}

// @Summary Получить всех пользователей
// @Description Получить всех пользователей
// @Tags example
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/users/get_all [get]
func (h *HandlerGetAll) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get all users")

	users, err := h.Storage.GetAll(r.Context())

	if err != nil {
		log.Printf("Error in get all users: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

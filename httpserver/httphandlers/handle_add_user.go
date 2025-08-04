package httphandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/nkhamm-spb/red_soft_test/metadata"
	"github.com/nkhamm-spb/red_soft_test/schemas"
	"github.com/nkhamm-spb/red_soft_test/storage"
)

type HandlerAddUser struct {
	Storage storage.StorageInterface
}

// @Summary Добавить пользователя
// @Description Добавить пользователя
// @Tags example
// @Accept   json
// @Produce  json
// @Param   input body   schemas.NewUser true  "Данные пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/users/add_user [post]
func (h *HandlerAddUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var newUser schemas.NewUser

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Request to add new user: %v\n", newUser)

	user := schemas.User{
		Name:    newUser.Name,
		Surname: newUser.Surname,
		Emails:  newUser.Emails,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	ch := make(chan error, 3)

	go func() {
		defer wg.Done()
		user.Age, err = metadata.GetAge(r.Context(), user.Name, user.Surname)
		if err != nil {
			ch <- err
		}
	}()
	go func() {
		defer wg.Done()
		user.Gender, _ = metadata.GetGender(r.Context(), user.Name, user.Surname)
		if err != nil {
			ch <- err
		}
	}()
	go func() {
		defer wg.Done()
		user.Nationalize, _ = metadata.GetNationalize(r.Context(), user.Name, user.Surname)
		if err != nil {
			ch <- err
		}
	}()

	wg.Wait()
	close(ch)

	for err := range ch {
		log.Printf("Error in add new user: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addedUser, err := h.Storage.AddUser(r.Context(), &user)

	if err != nil {
		log.Printf("Error in add new user: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(addedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

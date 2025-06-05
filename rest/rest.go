//go:generate mockgen -source=rest.go -destination=../mocks/storage.go -package=mocks

package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"networkCommunicationMin/models"
	secondary "networkCommunicationMin/secondary_function"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type Service struct {
	Storage Storage
}

type Storage interface {
	GetUsers() (map[int]models.User, error)
	GetUserById(id int) (models.User, error)
	SaveUser(name string, age int, friends []int64) (int, error)
	UpdateUser(id int, name string, age int, friends []int64) error
	DeleteUser(id int) error
}

func NewService(storage Storage) *Service {
	return &Service{Storage: storage}
}

func (s *Service) Ping(w http.ResponseWriter, r *http.Request) {
	// Ping - метод проверки соединения с сервером
	render.Status(r, http.StatusOK)
	_ = render.Render(w, r, &ResponsePayload{Result: "ping successful"})
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	// GetAll - метод, которая возвращает данные по всем пользователям
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	usersStor, err := s.Storage.GetUsers()
	if err != nil {
		err = fmt.Errorf("GetAll: get users: %w", err)
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	var keys []int
	for u := range usersStor {
		keys = append(keys, u)
	}
	sort.Ints(keys)

	var users []string
	for _, k := range keys {
		user := usersStor[k]
		users = append(users, user.ToSting())
	}

	_ = render.Render(w, r, &ResponsePayload{Result: users})
}

func (s *Service) GetFriends(w http.ResponseWriter, r *http.Request) {
	// GetFriends - метод, которая возвращает всех друзей пользователя
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		err = fmt.Errorf("invalid url id parameter, must be a number")
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}
	friends := []models.User{}

	user, err := s.Storage.GetUserById(userID)
	if err != nil {
		err = fmt.Errorf("GetFriends: get user: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	friendsList := user.Friends
	for _, idFriends := range friendsList {
		f, err := s.Storage.GetUserById(int(idFriends))
		if err != nil {
			err = fmt.Errorf("GetFriends: get friend: %w", err)
			render.Status(r, http.StatusBadRequest)
			_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
			return
		} else {
			friends = append(friends, f)
		}
	}

	_ = render.Render(w, r, &ResponsePayload{Result: friends})
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	// Create - метод, которая создаёт нового пользователя
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("Create: read content: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}
	defer r.Body.Close()

	var u models.User
	if err = json.Unmarshal(content, &u); err != nil {
		err = fmt.Errorf("Create: unmarshal content: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	id, err := s.Storage.SaveUser(u.Name, u.Age, u.Friends)
	if err != nil {
		err = fmt.Errorf("Create: save user in db: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, &ResponsePayload{Result: "user created with id: " + strconv.Itoa(id)})
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	// MakeFriends - метод, которая делает друзей из двух пользователей
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("MakeFriends: read content: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}
	defer r.Body.Close()

	var f models.Friends
	if err = json.Unmarshal(content, &f); err != nil {
		err = fmt.Errorf("MakeFriends: unmarshal content: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	sourceUser, err := s.Storage.GetUserById(f.SourceId)
	if err != nil {
		err = fmt.Errorf("MakeFriends: get source user: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	targetUser, err := s.Storage.GetUserById(f.TargetId)
	if err != nil {
		err = fmt.Errorf("MakeFriends: get target user: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	var checkFriend = false
	for _, idFriend := range sourceUser.Friends {
		if idFriend == int64(targetUser.ID) {
			checkFriend = true
		}
	}

	if !checkFriend {
		sourceUser.Friends = append(sourceUser.Friends, int64(targetUser.ID))
		if err = s.Storage.UpdateUser(sourceUser.ID, sourceUser.Name, sourceUser.Age, sourceUser.Friends); err != nil {
			err = fmt.Errorf("MakeFriends: update source user: %w", err)
			render.Status(r, http.StatusInternalServerError)
			_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
			return
		}

		targetUser.Friends = append(targetUser.Friends, int64(sourceUser.ID))
		if err = s.Storage.UpdateUser(targetUser.ID, targetUser.Name, targetUser.Age, targetUser.Friends); err != nil {
			err = fmt.Errorf("MakeFriends: update target user: %w", err)
			render.Status(r, http.StatusInternalServerError)
			_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
			return
		}

		res := []string{sourceUser.Name, "and", targetUser.Name, "are now friends"}
		_ = render.Render(w, r, &ResponsePayload{Result: strings.Join(res, " ")})
	} else {
		res := []string{sourceUser.Name, "and", targetUser.Name, "are already friends"}
		_ = render.Render(w, r, &ResponsePayload{Result: strings.Join(res, " ")})
	}
}

func (s *Service) UpdateUserAge(w http.ResponseWriter, r *http.Request) {
	// UpdateUser - метод, которая обновляет данные пользователя
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("UpdateUserAge: read content: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}
	defer r.Body.Close()

	var c models.ChangeAge
	if err = json.Unmarshal(content, &c); err != nil {
		err = fmt.Errorf("UpdateUserAge: unmarshal content: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		err = fmt.Errorf("invalid url id parameter, must be a number")
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	user, err := s.Storage.GetUserById(userID)
	if err != nil {
		err = fmt.Errorf("UpdateUserAge: get user: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	if err = s.Storage.UpdateUser(user.ID, user.Name, c.NewAge, user.Friends); err != nil {
		err = fmt.Errorf("UpdateUserAge: update user: %w", err)
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	_ = render.Render(w, r, &ResponsePayload{Result: "user age updated successfully"})
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	// Delete - метод, которая принимает ID пользователя и удаляет его из хранилища,
	// а также стирает его из массива friends у всех его друзей.
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		err = fmt.Errorf("invalid url id parameter, must be a number")
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	userToDelete, err := s.Storage.GetUserById(userID)
	if err != nil {
		err = fmt.Errorf("Delete: get user: %w", err)
		render.Status(r, http.StatusBadRequest)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	dataUsers, err := s.Storage.GetUsers()
	if err != nil {
		err = fmt.Errorf("Delete: get users data: %w", err)
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	for _, user := range dataUsers {
		for _, friendID := range user.Friends {
			if friendID == int64(userToDelete.ID) {
				indexUserToDelete := secondary.FindUser(user.Friends, userToDelete)
				user.Friends = slices.Delete(user.Friends, indexUserToDelete, indexUserToDelete+1)
				if err = s.Storage.UpdateUser(user.ID, user.Name, user.Age, user.Friends); err != nil {
					err = fmt.Errorf("Delete: update user: %w", err)
					render.Status(r, http.StatusInternalServerError)
					_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
					return
				}
			}
		}
	}

	if err = s.Storage.DeleteUser(userID); err != nil {
		err = fmt.Errorf("Delete: delete user: %w", err)
		render.Status(r, http.StatusInternalServerError)
		_ = render.Render(w, r, &ResponsePayload{Errors: []string{err.Error()}})
		return
	}

	res := []string{"user", userToDelete.Name, "successfully deleted"}
	_ = render.Render(w, r, &ResponsePayload{Result: strings.Join(res, " ")})
}

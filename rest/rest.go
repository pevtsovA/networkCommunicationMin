package rest

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"networkCommunicationMin/contract"
	"networkCommunicationMin/models"
	"networkCommunicationMin/secondary_function"
	"strconv"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type Service struct {
	Storage contract.Storage
}

func NewService(storage contract.Storage) Service {
	return Service{Storage: storage}
}

func (s *Service) Ping(w http.ResponseWriter, r *http.Request) {
	// Ping - функция проверки соединения с сервером
	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	// GetAll - функция, которая возвращает данные по всем пользователям
	log.Info("================= GetAll ==================")
	users, err := s.Storage.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get users"))
		log.Error("method 'GetAll', ", err)
		return
	}
	response := ""
	for _, user := range users {
		response += user.ToSting()
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (s *Service) GetFriends(w http.ResponseWriter, r *http.Request) {
	// GetFriends - функция, которая возвращает всех друзей пользователя
	log.Info("================= GetFriends ==================")
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)
	friends := []models.User{}

	user, err := s.Storage.GetUserById(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))
		log.Error("method 'GetFriends', ", err)
		return
	}
	log.Info("User received: ", user)

	friendsList := user.Friends
	for i := 0; i < len(friendsList); i++ {
		f, err := s.Storage.GetUserById(int(friendsList[i]))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("User friends does not exist"))
			log.Error("method 'GetFriends', ", err)
		} else {
			friends = append(friends, f)
		}
	}

	body, _ := json.Marshal(friends)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	// Create - функция, которая создаёт нового пользователя
	log.Info("================= Create ==================")
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to retrieve data"))
		log.Error("method 'Create', ", err)
		return
	}
	defer r.Body.Close()

	var u models.User
	if err = json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to retrieve data"))
		log.Error("method 'Create', ", err)
		return
	}

	id, err := s.Storage.SaveUser(u.Name, u.Age, u.Friends)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to save user"))
		log.Error("method 'Create', ", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created with id: " + strconv.Itoa(id)))
	log.Info("user created with id: " + strconv.Itoa(id))
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	// MakeFriends - функция, которая делает друзей из двух пользователей
	log.Info("=============== MakeFriends ===============")
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to retrieve data"))
		log.Error("method 'MakeFriends', ", err)
		return
	}
	defer r.Body.Close()

	var f models.Friends
	if err = json.Unmarshal(content, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to retrieve data"))
		log.Error("method 'MakeFriends', ", err)
		return
	}

	sourceUser, err := s.Storage.GetUserById(f.SourceId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Source user does not exist"))
		log.Error("method 'MakeFriends', ", err)
		return
	}
	log.Info("source user received: ", sourceUser)

	targetUser, err := s.Storage.GetUserById(f.TargetId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Target user does not exist"))
		log.Error("method 'MakeFriends', ", err)
		return
	}
	log.Info("target user received: ", targetUser)

	var checkFriend = false
	for _, idFriend := range sourceUser.Friends {
		if idFriend == int64(targetUser.ID) {
			checkFriend = true
		}
	}

	if checkFriend == false {
		sourceUser.Friends = append(sourceUser.Friends, int64(targetUser.ID))
		err = s.Storage.UpdateUser(sourceUser.ID, sourceUser.Name, sourceUser.Age, sourceUser.Friends)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to update source user"))
			log.Error("method 'MakeFriends', ", err)
			return
		}

		targetUser.Friends = append(targetUser.Friends, int64(sourceUser.ID))
		err = s.Storage.UpdateUser(targetUser.ID, targetUser.Name, targetUser.Age, targetUser.Friends)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to update target user"))
			log.Error("method 'MakeFriends', ", err)
			return
		}

		w.Write([]byte(sourceUser.Name + " and " + targetUser.Name + " are now friends"))
		log.Info(sourceUser.Name + " and " + targetUser.Name + " are now friends")
	} else {
		w.Write([]byte(sourceUser.Name + " and " + targetUser.Name + " are already friends"))
		log.Warning(sourceUser.Name + " and " + targetUser.Name + " are already friends")
	}
}

func (s *Service) UpdateUserAge(w http.ResponseWriter, r *http.Request) {
	// UpdateUser - функция, которая обновляет данные пользователя
	log.Info("============== UpdateUserAge ==============")
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to retrieve data"))
		log.Error("method 'UpdateUserAge', ", err)
		return
	}
	defer r.Body.Close()

	var c models.ChangeAge
	if err = json.Unmarshal(content, &c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to retrieve data"))
		log.Error("method 'UpdateUserAge', ", err)
		return
	}

	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	user, err := s.Storage.GetUserById(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))
		log.Error("method 'UpdateUserAge', ", err)
		return
	}
	log.Info("user received:", user)

	err = s.Storage.UpdateUser(user.ID, user.Name, c.NewAge, user.Friends)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update user"))
		log.Error("method 'UpdateUserAge', ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User age updated successfully"))
	log.Info("user age updated successfully")
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	// Delete - функция, которая принимает ID пользователя и удаляет его из хранилища,
	// а также стирает его из массива friends у всех его друзей.
	log.Info("================= DELETE ==================")
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to retrieve data"))
		log.Error("method 'Delete', ", err)
		return
	}
	defer r.Body.Close()

	var f models.Friends
	if err = json.Unmarshal(content, &f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to retrieve data"))
		log.Error("method 'Delete', ", err)
		return
	}
	uID := f.TargetId

	userToDelete, err := s.Storage.GetUserById(uID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))
		log.Error("method 'Delete', ", err)
		return
	}
	log.Info("user received for delete:", userToDelete)
	remoteUsername := userToDelete.Name

	dataUsers, err := s.Storage.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get users"))
		log.Error("method 'Delete', ", err)
		return
	}

	for _, user := range dataUsers {
		for _, friendID := range user.Friends {
			if friendID == int64(userToDelete.ID) {
				indexUserToDelete := secondary_function.FindUser(user.Friends, userToDelete)
				user.Friends = append((user.Friends)[:indexUserToDelete], (user.Friends)[indexUserToDelete+1:]...)
				err = s.Storage.UpdateUser(user.ID, user.Name, user.Age, user.Friends)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Failed to update user"))
					log.Error("method 'Delete', ", err)
					return
				}
			}
		}
	}

	err = s.Storage.DeleteUser(uID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete user"))
		log.Error("method 'Delete', ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("remote username: " + remoteUsername))
	log.Info("remote username:", remoteUsername)
}

package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"networkCommunicationMin/models"
	"strconv"
)

type Service struct {
	Storage Storage
}

func (s *Service) Ping(w http.ResponseWriter, r *http.Request) {
	// Ping - функция проверки соединения с сервером
	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	// GetAll - функция, которая возвращает данные по всем пользователям
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := s.Storage.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Method GetAll. Failed to get users: " + err.Error()))
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
	if r.Method == "GET" {
		id := chi.URLParam(r, "id")
		userID, _ := strconv.Atoi(id)
		response := []models.User{}
		dataUsers, err := s.Storage.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Method GetFriends. Failed to get users: " + err.Error()))
			return
		}

		if len(dataUsers) > 0 && userID >= models.MinUserID(dataUsers) && userID <= models.MaxUserID(dataUsers) {
			for _, userID := range dataUsers[userID].Friends {
				response = append(response, dataUsers[int(userID)])
			}
			body, _ := json.Marshal(response)

			w.WriteHeader(http.StatusOK)
			w.Write(body)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: wrong data" + "; http code: " + strconv.Itoa(http.StatusBadRequest)))
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	// Create - функция, которая создаёт нового пользователя
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u models.User
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		id, err := s.Storage.SaveUser(u.Name, u.Age, u.Friends)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Method Create. Failed to save user: " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created with id: " + strconv.Itoa(id) + "; http code: " + strconv.Itoa(http.StatusCreated)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	// MakeFriends - функция, которая делает друзей из двух пользователей
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var f models.Friends
		if err := json.Unmarshal(content, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		dataUsers, err := s.Storage.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Method MakeFriends. Failed to get users: " + err.Error()))
			return
		}

		if len(dataUsers) > 1 && f.SourceId >= models.MinUserID(dataUsers) && f.SourceId <= models.MaxUserID(dataUsers) &&
			f.TargetId >= models.MinUserID(dataUsers) && f.TargetId <= models.MaxUserID(dataUsers) {
			var checkFriend = false
			for _, value := range dataUsers {
				if value.ID == f.SourceId {
					for _, friendID := range value.Friends {
						if friendID == int64(f.TargetId) {
							checkFriend = true
						}
					}
				}
			}

			source := dataUsers[f.SourceId]
			target := dataUsers[f.TargetId]
			if checkFriend == false {
				source.Friends = append(source.Friends, int64(f.TargetId))
				target.Friends = append(target.Friends, int64(f.SourceId))
				dataUsers[source.ID] = source
				dataUsers[target.ID] = target

				err = s.Storage.UpdateUser(source.ID, dataUsers[source.ID].Name, dataUsers[source.ID].Age, dataUsers[source.ID].Friends)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Method MakeFriends. Failed to update source user: " + err.Error()))
					return
				}
				err = s.Storage.UpdateUser(target.ID, dataUsers[target.ID].Name, dataUsers[target.ID].Age, dataUsers[target.ID].Friends)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Method MakeFriends. Failed to update target user: " + err.Error()))
					return
				}

				w.Write([]byte(source.Name + " and " + target.Name + " are now friends" + "; http code: " + strconv.Itoa(http.StatusOK)))
			} else {
				w.Write([]byte(source.Name + " and " + target.Name + " are already friends" + "; http code: " + strconv.Itoa(http.StatusOK)))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: wrong data" + "; http code: " + strconv.Itoa(http.StatusBadRequest)))
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// UpdateUser - функция, которая обновляет возраст пользователя
	if r.Method == "PUT" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var c models.ChangeAge
		if err := json.Unmarshal(content, &c); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		id := chi.URLParam(r, "id")
		userID, _ := strconv.Atoi(id)
		dataUsers, err := s.Storage.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Method UpdateUser. Failed to get user: " + err.Error()))
			return
		}

		if len(dataUsers) > 0 && userID >= models.MinUserID(dataUsers) && userID <= models.MaxUserID(dataUsers) {
			u := dataUsers[userID]
			u.Age = c.NewAge
			dataUsers[userID] = u
			err = s.Storage.UpdateUser(dataUsers[userID].ID, dataUsers[userID].Name, dataUsers[userID].Age, dataUsers[userID].Friends)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Method UpdateUser. Failed to update user: " + err.Error()))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("user age updated successfully; http code: " + strconv.Itoa(http.StatusOK)))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: wrong data" + "; http code: " + strconv.Itoa(http.StatusBadRequest)))
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	// Delete - функция, которая принимает ID пользователя и удаляет его из хранилища,
	// а также стирает его из массива friends у всех его друзей.
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var f models.Friends
		if err := json.Unmarshal(content, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		dataUsers, err := s.Storage.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Method Delete. Failed to get users: " + err.Error()))
			return
		}
		uID := f.TargetId
		remoteUsername := dataUsers[uID].Name

		if len(dataUsers) > 0 && uID >= models.MinUserID(dataUsers) && uID <= models.MaxUserID(dataUsers) {
			for _, value := range dataUsers {
				for _, friendID := range value.Friends {
					if friendID == int64(dataUsers[uID].ID) {
						indexUser := models.FindUser(value.Friends, dataUsers[uID])
						value.Friends = append((value.Friends)[:indexUser], (value.Friends)[indexUser+1:]...)
						dataUsers[value.ID] = value
						err = s.Storage.UpdateUser(dataUsers[value.ID].ID, dataUsers[value.ID].Name, dataUsers[value.ID].Age, dataUsers[value.ID].Friends)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte("Method Delete. Failed to update user: " + err.Error()))
							return
						}
					}
				}
			}

			delete(dataUsers, uID)
			err = s.Storage.DeleteUser(uID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Method Delete. Failed to delete user: " + err.Error()))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("remote username: " + remoteUsername + "; http code: " + strconv.Itoa(http.StatusOK)))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: wrong data" + "; http code: " + strconv.Itoa(http.StatusBadRequest)))
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

package dataBase

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"networkCommunicationMin/models"
	"strconv"
)

func (s *Service) Ping(w http.ResponseWriter, r *http.Request) {
	// Ping - функция проверки соединения с сервером
	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	// GetAll - функция, которая возвращает данные по всем пользователям
	if r.Method == "GET" {
		response := ""
		for _, user := range s.DataReading() {
			response += user.ToSting()
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetFriends(w http.ResponseWriter, r *http.Request) {
	// GetFriends - функция, которая возвращает всех друзей пользователя
	if r.Method == "GET" {
		id := chi.URLParam(r, "id")
		userID, _ := strconv.Atoi(id)
		response := []models.User{}
		dataUsers := s.DataReading()

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

		id := s.DataRecording(u.Name, u.Age, u.Friends)

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
		dataUsers := s.DataReading()

		if len(dataUsers) > 1 && f.Source_id >= models.MinUserID(dataUsers) && f.Source_id <= models.MaxUserID(dataUsers) &&
			f.Target_id >= models.MinUserID(dataUsers) && f.Target_id <= models.MaxUserID(dataUsers) {
			var checkFriend = false
			for _, value := range dataUsers {
				if value.ID == f.Source_id {
					for _, friendID := range value.Friends {
						if friendID == int64(f.Target_id) {
							checkFriend = true
						}
					}
				}
			}

			source := dataUsers[f.Source_id]
			target := dataUsers[f.Target_id]
			if checkFriend == false {
				source.Friends = append(source.Friends, int64(f.Target_id))
				target.Friends = append(target.Friends, int64(f.Source_id))
				dataUsers[source.ID] = source
				dataUsers[target.ID] = target

				s.DataUpdating(source.ID, dataUsers[source.ID].Name, dataUsers[source.ID].Age, dataUsers[source.ID].Friends)
				s.DataUpdating(target.ID, dataUsers[target.ID].Name, dataUsers[target.ID].Age, dataUsers[target.ID].Friends)

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
		dataUsers := s.DataReading()

		if len(dataUsers) > 0 && userID >= models.MinUserID(dataUsers) && userID <= models.MaxUserID(dataUsers) {
			u := dataUsers[userID]
			u.Age = c.New_age
			dataUsers[userID] = u
			s.DataUpdating(dataUsers[userID].ID, dataUsers[userID].Name, dataUsers[userID].Age, dataUsers[userID].Friends)

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
		dataUsers := s.DataReading()
		uID := f.Target_id
		remoteUsername := dataUsers[uID].Name

		if len(dataUsers) > 0 && uID >= models.MinUserID(dataUsers) && uID <= models.MaxUserID(dataUsers) {
			for _, value := range dataUsers {
				for _, friendID := range value.Friends {
					if friendID == int64(dataUsers[uID].ID) {
						indexUser := models.FindUser(value.Friends, dataUsers[uID])
						value.Friends = append((value.Friends)[:indexUser], (value.Friends)[indexUser+1:]...)
						dataUsers[value.ID] = value
						s.DataUpdating(dataUsers[value.ID].ID, dataUsers[value.ID].Name, dataUsers[value.ID].Age, dataUsers[value.ID].Friends)
					}
				}
			}

			delete(dataUsers, uID)
			s.DataDeleting(uID)

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

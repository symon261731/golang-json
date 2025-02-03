package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main/package/types"
	"net/http"
	"os"
	"strconv"
)

var PORT = ":8080"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			jsonData, err := os.ReadFile("mockDB/mockDB.json")

			if err != nil {
				log.Println("Произошла ошибка при чтении json")
			}

			log.Println(string(jsonData))
		}

	})
	r.HandleFunc("/create", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			jsonData, err := os.ReadFile("mockDB/mockDB.json")

			if err != nil {
				log.Println("Произошла ошибка при чтении json")
				log.Println(err)
			}

			var bodyData types.CreateUserData

			var parseBodyErr = json.NewDecoder(request.Body).Decode(&bodyData)
			if parseBodyErr != nil {
				log.Println("Произошла ошибка при парсинге body")
				log.Println(parseBodyErr)
			}

			var users types.UserListMap
			errParseJson := json.Unmarshal(jsonData, &users)

			if errParseJson != nil {
				log.Println("Произошла ошибка при парсинге json")
				log.Println(errParseJson)
			}

			newId := len(users) + 1

			newUser := types.User{Id: newId, Name: bodyData.Name, Age: bodyData.Age, Friends: make([]types.UserFriends, 0)}

			//TODO вынести в отдельный алгоритм
			var newUserMap = types.UserListMap{}
			for s, user := range users {
				newUserMap[s] = user
			}
			newUserMap[strconv.Itoa(newId)] = newUser

			decodeUserMap, errMarshalJson := json.Marshal(newUserMap)

			if errMarshalJson != nil {
				log.Println("Произошла ошибка при шифровании json")
				log.Println(err)
				return
			}

			errWriteFile := os.WriteFile("mockDB/mockDB.json", decodeUserMap, 0644)

			if errWriteFile != nil {
				log.Println("Произошла ошибка при записи новых данных")
				log.Println(errWriteFile)
			}

		}
	})

	log.Printf("Веб-сервер запущен на http://127.0.0.1%s", PORT)
	err := http.ListenAndServe(PORT, r)

	log.Fatal(err)
}

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

var PORT = ":4000"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			http.Error(writer, "invalid get request", http.StatusNotFound)
			return
		}
		jsonData, err := os.ReadFile("mockDB/mockDB.json")

		if err != nil {
			log.Println("Произошла ошибка при чтении json")
		}

		log.Println(string(jsonData))

	})

	r.HandleFunc("/create", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "POST" {
			http.Error(writer, "invalid create request", http.StatusNotFound)
			return
		}

		var bodyData types.CreateUserData

		var parseBodyErr = json.NewDecoder(request.Body).Decode(&bodyData)
		if parseBodyErr != nil {
			log.Println("Произошла ошибка при парсинге body")
			log.Println(parseBodyErr)
		}

		jsonData, err := os.ReadFile("mockDB/mockDB.json")
		if err != nil {
			log.Println("Произошла ошибка при чтении json")
			log.Println(err)
		}

		var users types.UserListMap
		errParseJson := json.Unmarshal(jsonData, &users)

		if errParseJson != nil {
			log.Println("Произошла ошибка при парсинге json")
			log.Println(errParseJson)
		}

		newId := len(users) + 1
		newUser := types.User{Id: newId, Name: bodyData.Name, Age: bodyData.Age, Friends: make([]types.UserFriends, 0)}

		users[strconv.Itoa(newId)] = newUser

		decodeUserMap, errMarshalJson := json.Marshal(users)
		if errMarshalJson != nil {
			log.Println("Произошла ошибка при шифровании json")
			log.Println(err)
		}

		errWriteFile := os.WriteFile("mockDB/mockDB.json", decodeUserMap, 0644)
		if errWriteFile != nil {
			log.Println("Произошла ошибка при записи новых данных")
			log.Println(errWriteFile)
		}

	})

	r.HandleFunc("/{user_id}", func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != "PUT" {
			http.Error(writer, "invalid request", http.StatusNotFound)
			return
		}

		vars := mux.Vars(request)
		userId := vars["user_id"]

		var bodyData types.PutNewAgeJson

		parseBodyErr := json.NewDecoder(request.Body).Decode(&bodyData)
		if parseBodyErr != nil {
			log.Println("Произошла ошибка при парсинге body")
			log.Println(parseBodyErr)
		}

		newAge, errPrepareAge := strconv.Atoi(bodyData.NewAge)
		if errPrepareAge != nil {
			log.Println("Произошла ошибка при парсинге возраста из body")
			log.Println(errPrepareAge)
		}

		mockJsonData, err := os.ReadFile("mockDB/mockDB.json")
		if err != nil {
			log.Println("Произошла ошибка при чтении json")
			log.Println(err)
		}

		var users types.UserListMap
		marshalError := json.Unmarshal(mockJsonData, &users)
		if marshalError != nil {
			log.Println("Произошла ошибка при парсинге json")
			log.Println(marshalError)
		}

		var neededUser types.User

		if entry, ok := users[userId]; ok {
			neededUser = entry
		}

		neededUser.Age = newAge
		users[userId] = neededUser

		decodeUserMap, marshalError := json.Marshal(users)

		if marshalError != nil {
			log.Println("Произошла ошибки при marshal json")
			log.Println(marshalError)
		}

		errorWriteFile := os.WriteFile("mockDB/mockDB.json", decodeUserMap, 0644)

		if errorWriteFile != nil {
			log.Println("Произошла ошибка при записи в файл")
			log.Println(errorWriteFile)
		}

	})

	r.HandleFunc("/make_friends", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "POST" {
			http.Error(writer, "invalid make friends request", http.StatusNotFound)
			return
		}

		var bodyData types.PostIdsFriends
		errJsonDecode := json.NewDecoder(request.Body).Decode(&bodyData)

		if errJsonDecode != nil {
			log.Println("Произошла ошибка при парсинге body")
			log.Println(errJsonDecode)
		}

		mockJsonData, err := os.ReadFile("mockDB/mockDB.json")
		if err != nil {
			log.Println("Произошла ошибка при чтении json")
			log.Println(err)
		}

		var users types.UserListMap
		marshalError := json.Unmarshal(mockJsonData, &users)

		if marshalError != nil {
			log.Println("Произошла ошибка при парсинге json файла")
			log.Println(err)
		}

		var (
			sourceUser types.User
			targetUser types.User
		)

		if entry, ok := users[bodyData.Source_id]; ok {
			sourceUser = entry
		} else {
			log.Println("Пользователь c таким source_id не существует")
		}

		if entry, ok := users[bodyData.Target_id]; ok {
			targetUser = entry
		} else {
			log.Println("Пользователь c таким source_id не существует")
		}

		sourceUser.Friends = append(sourceUser.Friends, types.UserFriends{Id: targetUser.Id, Name: targetUser.Name})
		targetUser.Friends = append(targetUser.Friends, types.UserFriends{Id: sourceUser.Id, Name: sourceUser.Name})

		log.Println(sourceUser)
		log.Println(targetUser)

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Пользователи стали друзьями"))

	})

	log.Printf("Веб-сервер запущен на http://127.0.0.1%s", PORT)
	err := http.ListenAndServe(PORT, r)

	log.Fatal(err)
}

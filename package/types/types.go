package types

type CreateUserData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []UserFriends
}

type UserListMap map[string]User

type UserFriends struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PutNewAgeJson struct {
	NewAge string `json:"new age"`
}

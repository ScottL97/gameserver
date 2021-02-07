package models

type ServerInfo struct {
	Users map[string]string `json:"users"`
	//Messages []ClientMsg `json:messages`
	//Time string `json:time`
}

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type ClientMsg struct {
	Id       string `json:"id"`
	Msg      string `json:"msg"`
	UserName string `json:"username"`
}
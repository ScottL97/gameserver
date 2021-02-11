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

type UserCheck struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

type ClientMsg struct {
	Id       string `json:"id"`
	Msg      string `json:"msg"`
	UserName string `json:"username"`
}

type StatusReq struct {
	Running string `json:"running"`
	Message string `json:"message"` // 游戏胜利/失败
}

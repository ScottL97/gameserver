package models

type GameInfo struct {
	Map      [7][7]Square `json:"map"`
	Players  []Player     `json:"players"`
	Round    int          `json:"round"`
	Progress int          `json:"progress"`
}

type Square struct {
	// Player    string `json:"player"`
	Virus     int `json:"virus"`     // 0代表没有，1-3代表3个级别的病毒
	Institute int `json:"institute"` // 0代表没有，1代表有
}

type Player struct {
	Username   string  `json:"username"`
	Occupation int     `json:"occupation"`
	PosX       int     `json:"posx"`
	PosY       int     `json:"posy"`
	Energy     float32 `json:"energy"`
	Finish     int     `json:"finish"` // 0代表没有完成，1代表完成
}

const (
	ADD_VIRUS    = iota // 添加病毒
	SUB_VIRUS           // 清除病毒
	MOV_PLAYER          // 移动玩家
	SUB_PLAYER          // 删除玩家，如断开连接时执行
	ADD_INSTITUE        // 建造研究所，工程师可以执行
	DO_RESEARCH         // 进行研究，科学家可以进行
)

// http客户端->服务器，清除病毒
type KillVirusReq struct {
	Username string `json:"username"`
	PosX     int    `json:"posx"`
	PosY     int    `json:"posy"`
}

// ws服务器->客户端，同步清除病毒
type KillVirusWSReq struct {
	Type int `json:"type"`
	PosX int `json:"posx"`
	PosY int `json:"posy"`
}

type BuildInstituteReq struct {
	Username string `json:"username"`
	PosX     int    `json:"posx"`
	PosY     int    `json:"posy"`
}

type BuildInstituteWSReq struct {
	Type int `json:"type"`
	PosX int `json:"posx"`
	PosY int `json:"posy"`
}

type DoResearchReq struct {
	Username string `json:"username"`
	PosX     int    `json:"posx"`
	PosY     int    `json:"posy"`
}
type DoResearchWSReq struct {
	Type     int `json:"type"`
	Progress int `json:"progress"`
}

// http客户端->服务器，客户端检查路径可行后才能请求，服务器返回路径是否改变成功
type MovPlayerReq struct {
	Username string `json:"username"`
	PosX     int    `json:"posx"`
	PosY     int    `json:"posy"`
	Drive    string `json:"drive"`
}

// ws服务器->客户端，同步玩家移动
type MovPlayerWSReq struct {
	Type     int     `json:"type"`
	Username string  `json:"username"`
	PosX     int     `json:"posx"`
	PosY     int     `json:"posy"`
	Energy   float32 `json:"energy"`
}

// ws服务器->客户端，客户端清理用户
type SubPlayerWSReq struct {
}

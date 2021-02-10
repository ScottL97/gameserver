package models

type GameInfo struct {
	Map     [7][7]Square `json:"map"`
	Players []Player     `json:"players"`
	Round   int          `json:"round"`
}

type Square struct {
	Player    string `json:"player"`
	Virus     int    `json:"virus"`     // 0代表没有，1-3代表3个级别的病毒
	Institute int    `json:"institute"` // 0代表没有，1代表有
}

type Player struct {
	Username   string `json:"username"`
	Occupation int    `json:"occupation"`
	PosX       int    `json:"posx"`
	PosY       int    `json:"posy"`
}

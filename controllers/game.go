package controllers

import (
	"fmt"
	"gameserver/models"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

type GameController struct {
	beego.Controller
}

type Game struct {
	Status       int
	Squares      [7][7]models.Square // 用于发送给玩家
	Players      []models.Player
	PlayersMutex sync.Mutex
	Round        int
}

func getGame() *Game {
	gameInstanceOnce.Do(func() {
		gameInstance = new(Game)
	})
	return gameInstance
}

func init() {
	getGame().Status = PREPARE
	// 配置每轮会产生的病毒数量
	for i := 0; i < 10; i++ {
		virus1Num[i] = 5
	}
	for i := 10; i < 20; i++ {
		virus1Num[i] = 6
		virus2Num[i] = i - 10
	}
	for i := 20; i < 30; i++ {
		virus1Num[i] = 7
		virus2Num[i] = i - 20
		virus3Num[i] = i - 20
	}
}

type ReadyUsernames struct {
	Usernames []string `json:"readyUsernames"`
}

var (
	gameInstanceOnce sync.Once
	gameInstance     *Game
	occupations      map[string]int
	// 序号为轮数-1，值为病毒数目
	virus1Num         [30]int
	virus2Num         [30]int
	virus3Num         [30]int
	gameInfoBroadcast = make(chan int, 1)
)

const (
	PREPARE = iota
	RUNNING
)

const (
	SCIENTIST = iota
	ENGINEER
	DOCTOR
	DRIVER
)

func (c *GameController) UserReady() {
	usernameParam := c.Ctx.Input.Param(":username")
	if _, ok := getUserManager().loginUsers[usernameParam]; !ok {
		fmt.Printf("[UserReady]%s is not in loginUsers\n", usernameParam)
		c.Ctx.WriteString("fail")
		return
	}
	// 游戏正在进行中，不允许准备
	if getGame().Status == RUNNING {
		fmt.Printf("[UserReady]%s is not allowed to be ready cause the game has begun\n", usernameParam)
		c.Ctx.WriteString("fail")
		return
	}
	player := models.Player{Username: usernameParam}
	insertPlayer(player)
	c.Ctx.WriteString("ok")
	// readyUsernames := []string{}
	// for _, p := range getGame().Players {
	// 	readyUsernames = append(readyUsernames, p.Username)
	// }
	// c.Data["json"] = &ReadyUsernames{Usernames: readyUsernames}
	// c.ServeJSON()
	if len(getGame().Players) >= 3 {
		GameInit()
	}
}

func (c *GameController) UserCancel() {
	usernameParam := c.Ctx.Input.Param(":username")
	if _, ok := getUserManager().loginUsers[usernameParam]; !ok {
		fmt.Printf("[UserCancel]%s is not in loginUsers\n", usernameParam)
		c.Ctx.WriteString("fail")
		return
	}
	// 游戏正在进行中，不允许取消准备
	if getGame().Status == RUNNING {
		fmt.Printf("[UserCancel]%s is not allowed to cancel cause the game has begun\n", usernameParam)
		c.Ctx.WriteString("fail")
		return
	}
	deletePlayer(usernameParam)
	c.Ctx.WriteString("ok")
}

func insertPlayer(player models.Player) {
	getGame().PlayersMutex.Lock()
	defer getGame().PlayersMutex.Unlock()
	getGame().Players = append(getGame().Players, player)
}
func deletePlayer(username string) {
	getGame().PlayersMutex.Lock()
	defer getGame().PlayersMutex.Unlock()
	for i := 0; i < len(getGame().Players); i++ {
		if getGame().Players[i].Username == username {
			getGame().Players = append(getGame().Players[:i], getGame().Players[i+1:]...)
		}
	}
}

// 玩家准备完成，初始化游戏，向玩家广播游戏初始化信息
func GameInit() {
	getGame().Status = RUNNING
	getGame().Round = 1
	rand.Seed(time.Now().UnixNano())
	// 初始化病毒位置
	for i := 0; i < virus1Num[getGame().Round-1]; i++ {
		x := rand.Intn(7)
		y := rand.Intn(7)
		if getGame().Squares[x][y].Virus == 0 {
			getGame().Squares[x][y].Virus = 1
		} else {
			// 如果已经设置，需要索引减1，多随机一次
			i--
		}
	}
	for i := 0; i < virus2Num[getGame().Round-1]; i++ {
		x := rand.Intn(7)
		y := rand.Intn(7)
		if getGame().Squares[x][y].Virus == 0 {
			getGame().Squares[x][y].Virus = 2
		} else {
			// 如果已经设置，需要索引减1，多随机一次
			i--
		}
	}
	for i := 0; i < virus3Num[getGame().Round-1]; i++ {
		x := rand.Intn(7)
		y := rand.Intn(7)
		if getGame().Squares[x][y].Virus == 0 {
			getGame().Squares[x][y].Virus = 3
		} else {
			// 如果已经设置，需要索引减1，多随机一次
			i--
		}
	}
	// 初始化玩家职业及初始位置
	occupations = map[string]int{"scientist": 0, "engineer": 1, "doctor": 2, "driver": 3}
	occupationsNotChosen := []string{"scientist", "engineer", "doctor", "driver"}
	leftNum := 4
	for i := 0; i < len(getGame().Players); i++ {
		randNum := rand.Intn(leftNum)
		leftNum--
		chosen := occupationsNotChosen[randNum]
		occupationsNotChosen = append(occupationsNotChosen[:randNum], occupationsNotChosen[randNum+1:]...)

		getGame().Players[i].Occupation = occupations[chosen]
		getGame().Players[i].PosX = rand.Intn(7)
		getGame().Players[i].PosY = rand.Intn(7)
		getGame().Squares[getGame().Players[i].PosX][getGame().Players[i].PosY].Player = getGame().Players[i].Username
	}
	go sendGameInfo()
	sendGameStatus(0)
}

func sendGameStatus(status int) {
	gameInfoBroadcast <- status
	var statusReq models.StatusReq
	if status == 0 {
		statusReq = models.StatusReq{Running: "yes"}
	} else {
		statusReq = models.StatusReq{Running: "no"}
	}
	for client := range getWebSocketManager().activeWebSockets {
		err := client.WriteJSON(statusReq)
		if err != nil {
			log.Printf("[startGame]client.WriteJSON error: %v", err)
		}
	}
}

// 不足3人时结束游戏

func sendGameInfo() {
	// 游戏信息变动时，向客户端广播地图信息，直到游戏结束
	for {
		status := <-gameInfoBroadcast
		if status == 1 {
			return
		}
		gameInfo := models.GameInfo{Map: getGame().Squares, Players: getGame().Players, Round: getGame().Round}
		fmt.Println(gameInfo)

		for client := range getWebSocketManager().activeWebSockets {
			err := client.WriteJSON(gameInfo)
			if err != nil {
				log.Printf("[sendMapInfo]client.WriteJSON error: %v", err)
			}
		}
	}
}

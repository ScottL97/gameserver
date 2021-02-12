package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gameserver/models"
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
	FinishNum    int
	VirusNum     int
	Progress     int
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
		virus2Num[i] = 6
	}
	for i := 20; i < 30; i++ {
		virus3Num[i] = 7
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

const PLAYER_NUM = 2

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
	player := models.Player{Username: usernameParam, Energy: 4}
	insertPlayer(player)
	c.Ctx.WriteString("ok")
	// todo: 如果还有未准备的玩家，进行等待，限制玩家数为3-4
	fmt.Println("len(getGame().Players): ", len(getGame().Players))
	if len(getGame().Players) >= PLAYER_NUM {
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
	fmt.Println("[deletePlayer]delete:", username)
	getGame().PlayersMutex.Lock()
	for i := 0; i < len(getGame().Players); i++ {
		if getGame().Players[i].Username == username {
			getGame().Players = append(getGame().Players[:i], getGame().Players[i+1:]...)
			break
		}
	}
	getGame().PlayersMutex.Unlock()
	if getGame().Status != RUNNING {
		return
	}
	// 如果游戏正在进行中，玩家数少于2，游戏结束
	if len(getGame().Players) < 2 {
		sendGameStatus(1)
		stopGame()
	}
}
func finishPlayer(username string) error {
	getGame().PlayersMutex.Lock()
	defer getGame().PlayersMutex.Unlock()
	for i := 0; i < len(getGame().Players); i++ {
		if getGame().Players[i].Username == username {
			if getGame().Players[i].Finish == 0 {
				getGame().Players[i].Finish = 1
				getGame().FinishNum++
				return nil
			}
		}
	}
	return errors.New("player has finished before")
}
func getPlayer(username string) (models.Player, error) {
	getGame().PlayersMutex.Lock()
	defer getGame().PlayersMutex.Unlock()
	for i := 0; i < len(getGame().Players); i++ {
		if getGame().Players[i].Username == username {
			return getGame().Players[i], nil
		}
	}
	return models.Player{}, errors.New("can't find player")
}
func movPlayer(username string, posx, posy int, energy float32) {
	getGame().PlayersMutex.Lock()
	defer getGame().PlayersMutex.Unlock()
	for i := 0; i < len(getGame().Players); i++ {
		if getGame().Players[i].Username == username {
			getGame().Players[i].PosX = posx
			getGame().Players[i].PosY = posy
			getGame().Players[i].Energy -= energy
			break
		}
	}
}
func cleanPlayers() {
	getGame().PlayersMutex.Lock()
	defer getGame().PlayersMutex.Unlock()

	getGame().Players = []models.Player{}
}

func (c *GameController) UserMove() {
	var req models.MovPlayerReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &req)
	if err != nil {
		fmt.Println("[UserMove]json.Unmarshal is err:", err.Error())
		c.Ctx.WriteString("fail")
		return
	}
	player, err := getPlayer(req.Username)
	if err != nil {
		fmt.Printf("[UserMove]%s is not in Players\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	energyNeed, err := getEnergyNeed(player.PosX, player.PosY, req.PosX, req.PosY, player.Occupation)
	if err != nil {
		fmt.Printf("[UserMove]%s getEnergyNeed fail\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	if player.Energy-energyNeed < 0.0 {
		fmt.Printf("[UserMove]%s energy is not enough\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	movPlayer(req.Username, req.PosX, req.PosY, energyNeed)
	movReq := models.MovPlayerWSReq{Type: models.MOV_PLAYER, Username: req.Username, PosX: req.PosX, PosY: req.PosY, Energy: energyNeed}
	fmt.Println("UserMove:", movReq)
	err = sendJSON(movReq)
	if err != nil {
		fmt.Println("[UserMove]sendJSON fail")
	}
	if req.Drive != "" {
		movPlayer(req.Drive, req.PosX, req.PosY, 0)
		movReq = models.MovPlayerWSReq{Type: models.MOV_PLAYER, Username: req.Drive, PosX: req.PosX, PosY: req.PosY, Energy: 0.0}
		err = sendJSON(movReq)
		if err != nil {
			fmt.Println("[UserMove]sendJSON fail")
		}
	}

	c.Ctx.WriteString("ok")
}

func (c *GameController) UserFinish() {
	usernameParam := c.Ctx.Input.Param(":username")
	if _, ok := getUserManager().loginUsers[usernameParam]; !ok {
		fmt.Printf("[UserFinish]%s is not in loginUsers\n", usernameParam)
		c.Ctx.WriteString("fail")
		return
	}
	// 游戏未在进行中，不允许执行完成操作
	if getGame().Status != RUNNING {
		fmt.Printf("[UserFinish]%s is not allowed to finish cause the game is not running\n", usernameParam)
		c.Ctx.WriteString("fail")
		return
	}
	err := finishPlayer(usernameParam)
	if err != nil {
		c.Ctx.WriteString("warn")
	} else {
		c.Ctx.WriteString("ok")
	}
	if getGame().FinishNum >= len(getGame().Players) {
		newRound()
	}
}

func myAbs(a, b int) int {
	if a-b > 0 {
		return a - b
	}
	return b - a
}
func getEnergyNeed(curx, cury, tarx, tary, occupation int) (float32, error) {
	if tarx > 6 || tarx < 0 || tary > 6 || tary < 0 {
		return 0, errors.New("position wrong")
	}
	energy := float32(myAbs(curx, tarx) + myAbs(cury, tary))
	if occupation == DRIVER {
		energy /= 2.0
	}

	return energy, nil
}

func (c *GameController) KillVirus() {
	var req models.KillVirusReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &req)
	if err != nil {
		fmt.Println("[KillVirus]json.Unmarshal is err:", err.Error())
		c.Ctx.WriteString("fail")
		return
	}
	player, err := getPlayer(req.Username)
	if err != nil {
		fmt.Printf("[KillVirus]%s is not in Players\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	if getGame().Squares[player.PosX][player.PosY].Virus == 0 {
		fmt.Printf("[KillVirus]%s is not in a square including virus\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	var energyNeed float32
	if player.Occupation == DOCTOR {
		energyNeed = 1.0
	} else {
		energyNeed = float32(getGame().Squares[player.PosX][player.PosY].Virus) + 1.0
	}
	if player.Energy-energyNeed < 0.0 {
		fmt.Printf("[KillVirus]%s energy is not enough\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	getGame().Squares[player.PosX][player.PosY].Virus = 0
	player.Energy -= energyNeed
	getGame().VirusNum--
	killReq := models.KillVirusWSReq{Type: models.SUB_VIRUS, PosX: req.PosX, PosY: req.PosY}
	err = sendJSON(killReq)
	if err != nil {
		fmt.Println("[KillVirus]sendJSON fail")
	}

	c.Ctx.WriteString("ok")
	if getGame().VirusNum == 0 {
		sendGameStatus(3)
		stopGame()
	}
}

func (c *GameController) BuildInstitute() {
	var req models.BuildInstituteReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &req)
	if err != nil {
		fmt.Println("[BuildInstitute]json.Unmarshal is err:", err.Error())
		c.Ctx.WriteString("fail")
		return
	}
	player, err := getPlayer(req.Username)
	if err != nil {
		fmt.Printf("[BuildInstitute]%s is not in Players\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	if player.Occupation != ENGINEER {
		fmt.Printf("[BuildInstitute]%s is not an engineer\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	if player.Energy-4.0 < 0.0 {
		fmt.Printf("[BuildInstitute]%s energy is not enough\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	player.Energy -= 4.0
	getGame().Squares[player.PosX][player.PosY].Institute = 1
	buildReq := models.BuildInstituteWSReq{Type: models.ADD_INSTITUE, PosX: req.PosX, PosY: req.PosY}
	err = sendJSON(buildReq)
	if err != nil {
		fmt.Println("[BuildInstitute]sendJSON fail")
	}

	c.Ctx.WriteString("ok")
}

func (c *GameController) DoResearch() {
	var req models.DoResearchReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &req)
	if err != nil {
		fmt.Println("[DoResearch]json.Unmarshal is err:", err.Error())
		c.Ctx.WriteString("fail")
		return
	}
	player, err := getPlayer(req.Username)
	if err != nil {
		fmt.Printf("[DoResearch]%s is not in Players\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	if player.Occupation != SCIENTIST {
		fmt.Printf("[DoResearch]%s is not an scientist\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	if player.Energy-2.0 < 0.0 {
		fmt.Printf("[DoResearch]%s energy is not enough\n", req.Username)
		c.Ctx.WriteString("fail")
		return
	}
	player.Energy -= 2.0
	getGame().Progress += 5
	researchReq := models.DoResearchWSReq{Type: models.DO_RESEARCH, Progress: getGame().Progress}
	err = sendJSON(researchReq)
	if err != nil {
		fmt.Println("[BuildInstitute]sendJSON fail")
	}

	c.Ctx.WriteString("ok")
	if getGame().Progress == 100 {
		sendGameStatus(3)
		stopGame()
	}
}

// 玩家准备完成，初始化游戏，向玩家广播游戏初始化信息
func GameInit() {
	fmt.Println("[GameInit]game on!")
	getGame().Status = RUNNING
	getGame().Round = 1
	getGame().VirusNum = virus1Num[getGame().Round-1] + virus2Num[getGame().Round-1] + virus3Num[getGame().Round-1]
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
	}
	go sendGameInfo()
	time.Sleep(time.Millisecond * 100)
	sendGameStatus(0)
}

func sendGameStatus(status int) {
	if status == 2 { // 游戏状态更新，只是为了sendGameInfo()
		gameInfoBroadcast <- status
		return
	}
	var statusReq models.StatusReq
	if status == 0 { // 游戏开始
		statusReq = models.StatusReq{Running: "yes", Message: ""}
	} else if status == 1 { // 游戏失败，结束sendGameInfo协程
		statusReq = models.StatusReq{Running: "no", Message: "lose"}
	} else if status == 3 { // 游戏胜利
		statusReq = models.StatusReq{Running: "no", Message: "win"}
	}
	for client := range getWebSocketManager().activeWebSockets {
		getWebSocketManager().webSocketMutex.Lock()
		err := client.WriteJSON(statusReq)
		getWebSocketManager().webSocketMutex.Unlock()
		if err != nil {
			fmt.Printf("[startGame]client.WriteJSON error: %v\n", err)
		}
	}
	gameInfoBroadcast <- status
}

func stopGame() {
	fmt.Println("[stopGame]game end!")
	getGame().Status = PREPARE
	getGame().Round = 1
	getGame().FinishNum = 0
	getGame().VirusNum = 0
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			getGame().Squares[i][j].Virus = 0
			getGame().Squares[i][j].Institute = 0
		}
	}
	cleanPlayers()
}

// todo: 不足3人时结束游戏
// todo: 有人断开连接时如果其他人均为回合结束状态，开始下一轮

func sendGameInfo() {
	// 游戏信息变动时，向客户端广播地图信息，直到游戏结束
	for {
		status := <-gameInfoBroadcast
		if status == 1 { // 游戏失败，不再更新游戏信息
			return
		}
		gameInfo := models.GameInfo{Map: getGame().Squares, Players: getGame().Players,
			Round: getGame().Round, Progress: getGame().Progress}
		err := sendJSON(gameInfo)
		if err != nil {
			fmt.Println("[sendGameInfo]sendJSON fail")
		}
	}
}

func newRound() {
	// time.Sleep(time.Second) // todo: 等待finishUser逻辑处理完成，应该用通道
	getGame().Round++
	if getGame().Round > 30 {
		sendGameStatus(1)
		stopGame()
		return
	}
	getGame().VirusNum = getGame().VirusNum + virus1Num[getGame().Round-1] + virus2Num[getGame().Round-1] + virus3Num[getGame().Round-1]
	if getGame().VirusNum > 49 {
		sendGameStatus(1)
		stopGame()
		return
	}
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
	// 恢复玩家行动力
	for i := 0; i < len(getGame().Players); i++ {
		getGame().Players[i].Energy = 4
		getGame().Players[i].Finish = 0
		getGame().FinishNum = 0
	}
	sendGameStatus(2)
}

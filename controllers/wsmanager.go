package controllers

import (
	"gameserver/models"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	// 键为WebSocket连接，值为用户id
	activeWebSockets map[*websocket.Conn]string
	messages         []models.ClientMsg
	webSocketMutex   sync.Mutex
	messagesMutex    sync.Mutex
}

var (
	wsmInstanceOnce sync.Once
	wsmInstance     *WebSocketManager
)

func init() {
	getWebSocketManager().activeWebSockets = make(map[*websocket.Conn]string)
}

func getWebSocketManager() *WebSocketManager {
	// Double Check Lock，当对象创建后不再加锁
	//if instance == nil {
	//	instanceMutex.Lock()
	//	defer instanceMutex.Unlock()
	//	if instance == nil {
	//		instance = new(WebSocketManager)
	//	}
	//}
	// 更好的方法是在程序运行过程中只运行一次创建对象
	wsmInstanceOnce.Do(func() {
		wsmInstance = new(WebSocketManager)
	})

	return wsmInstance
}

func insertWebSocket(ws *websocket.Conn) {
	getWebSocketManager().webSocketMutex.Lock()
	defer getWebSocketManager().webSocketMutex.Unlock()
	getWebSocketManager().activeWebSockets[ws] = ""
}

func deleteWebSocket(ws *websocket.Conn) {
	getWebSocketManager().webSocketMutex.Lock()
	defer getWebSocketManager().webSocketMutex.Unlock()
	delete(getWebSocketManager().activeWebSockets, ws)
}

func insertMessage(msg models.ClientMsg) {
	getWebSocketManager().messagesMutex.Lock()
	defer getWebSocketManager().messagesMutex.Unlock()
	getWebSocketManager().messages = append(getWebSocketManager().messages, msg)
}

func cleanMessage() {
	getWebSocketManager().messagesMutex.Lock()
	defer getWebSocketManager().messagesMutex.Unlock()
	getWebSocketManager().messages = []models.ClientMsg{}
}

func getWebSocketNum() int {
	return len(getWebSocketManager().activeWebSockets)
}

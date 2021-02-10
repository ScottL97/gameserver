package controllers

import (
	"fmt"
	"gameserver/models"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

var (
	upg                 = websocket.Upgrader{}
	serverInfoBroadcast = make(chan models.ClientMsg)
)

type WebSocketController struct {
	beego.Controller
}

func init() {
	go sendServerInfo()
	go sendMessages()
}

func (c *WebSocketController) Get() {
	// todo: 身份校验
	ws, err := upg.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		log.Fatalln(err)
	}
	insertWebSocket(ws)
	handleMessages(ws)
}

func sendServerInfo() {
	// 每隔60ms向客户端广播服务器信息，如果发送信息失败，则客户端已断开连接，从活动websocket中删除
	for {
		time.Sleep(time.Millisecond * 60)
		// fmt.Println(getUserManager().loginUsers)
		serverInfo := models.ServerInfo{Users: getUserManager().loginUsers}
		//fmt.Println("send: ", serverInfo)

		for client := range getWebSocketManager().activeWebSockets {
			err := client.WriteJSON(serverInfo)
			if err != nil {
				log.Printf("[sendServerInfo]client.WriteJSON error: %v", err)
				deleteUser(getWebSocketManager().activeWebSockets[client])
				deleteWebSocket(client)
				err = client.Close()
				if err != nil {
					log.Printf("[sendServerInfo]client.Close error: %v", err)
				}
			}
		}
	}
}

func sendMessages() {
	// 当消息不为空时，向客户端发送信息
	for {
		message := <-serverInfoBroadcast
		// messages := getWebSocketManager().messages
		// cleanMessage()
		fmt.Println("send: ", message)

		for client := range getWebSocketManager().activeWebSockets {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("[sendMessages]client.WriteJSON error: %v", err)
			}
		}
	}
}

func handleMessages(ws *websocket.Conn) {
	for {
		var clientMsg models.ClientMsg
		err := ws.ReadJSON(&clientMsg)
		if err != nil {
			log.Printf("[handleMessages]ws.ReadJSON error: %v", err)
			break
		} else {
			fmt.Println("[handleMessages]receive msg:", clientMsg.Msg, clientMsg.UserName, clientMsg.Id)
			if len(clientMsg.Msg) == 0 {
				getWebSocketManager().activeWebSockets[ws] = clientMsg.UserName
				continue
			}
			insertMessage(clientMsg)
			serverInfoBroadcast <- clientMsg
		}
	}
}

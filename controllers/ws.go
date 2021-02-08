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
	// 每隔100ms向客户端广播服务器信息
	for {
		time.Sleep(time.Millisecond * 100)
		// fmt.Println(getUserManager().loginUsers)
		serverInfo := models.ServerInfo{Users: getUserManager().loginUsers}
		//fmt.Println("send: ", serverInfo)

		for client := range getWebSocketManager().activeWebSockets {
			err := client.WriteJSON(serverInfo)
			if err != nil {
				log.Printf("client.WriteJSON error: %v", err)
				err = client.Close()
				if err != nil {
					log.Printf("client.Close error: %v", err)
				}
				deleteWebSocket(client)
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
				log.Printf("client.WriteJSON error: %v", err)
				err = client.Close()
				if err != nil {
					log.Printf("client.Close error: %v", err)
				}
				deleteWebSocket(client)
			}
		}
	}
}

func handleMessages(ws *websocket.Conn) {
	fmt.Println("-------------------------------------")
	for {
		var clientMsg models.ClientMsg
		err := ws.ReadJSON(&clientMsg)
		if err != nil {
			log.Printf("ws.ReadJSON error: %v", err)
			// fmt.Println(getWebSocketManager().activeWebSockets)
			// fmt.Println("", ws)
			// deleteUser(getWebSocketManager().activeWebSockets[ws])
			// deleteWebSocket(ws)
			break
		} else {
			fmt.Println("receive msg:", clientMsg.Msg, clientMsg.Id)
			if len(clientMsg.Msg) == 0 {
				getWebSocketManager().activeWebSockets[ws] = clientMsg.Id
				continue
			}
			insertMessage(clientMsg)
			serverInfoBroadcast <- clientMsg
		}
	}
	deleteUser(getWebSocketManager().activeWebSockets[ws])
	fmt.Println("-------------------------------------")
}

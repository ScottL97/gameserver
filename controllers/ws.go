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
	c.TplName = "login.tpl"
}

func sendServerInfo() {
	// 每隔60ms向客户端广播服务器信息，如果发送信息失败，则客户端已断开连接，从活动websocket中删除
	for {
		time.Sleep(time.Millisecond * 60)
		// fmt.Println(getUserManager().loginUsers)
		serverInfo := models.ServerInfo{Users: getUserManager().loginUsers}
		//fmt.Println("send: ", serverInfo)
		for client := range getWebSocketManager().activeWebSockets {
			getWebSocketManager().webSocketMutex.Lock()
			err := client.WriteJSON(serverInfo)
			getWebSocketManager().webSocketMutex.Unlock()
			if err != nil {
				log.Printf("[sendServerInfo]client.WriteJSON error: %v", err)
				deleteUser(getWebSocketManager().activeWebSockets[client])
				deletePlayer(getWebSocketManager().activeWebSockets[client])
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
	// 当有消息时，向客户端发送websocket消息
	for {
		message := <-serverInfoBroadcast
		// fmt.Println("send: ", message)
		err := sendJSON(message)
		if err != nil {
			log.Printf("[sendMessages]sendJSON error: %v", err)
		}
	}
}

func sendJSON(i interface{}) error {
	// switch t := i.(type) {
	// default:
	// 	{
	// 		fmt.Printf("type %T", t)
	for client := range getWebSocketManager().activeWebSockets {
		getWebSocketManager().webSocketMutex.Lock()
		err := client.WriteJSON(i)
		getWebSocketManager().webSocketMutex.Unlock()
		if err != nil {
			return err
		}
	}
	// 	}
	// }
	return nil
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

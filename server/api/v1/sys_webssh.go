package v1

import (
	"bytes"
	"gin-vue-admin/model/request"
	"gin-vue-admin/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsSsh(c *gin.Context) {
	var sshRequest request.SshRequestStruct

	username := c.Query("username")
	password := c.Query("password")
	ip := c.Query("host")
	port := c.Query("port")

	intPost, _ := strconv.Atoi(port)
	sshRequest.Port = intPost
	sshRequest.Username = username
	sshRequest.Ip = ip
	sshRequest.Password = password

	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if utils.HandleError(c, err) {
		return
	}
	defer wsConn.Close()

	cols, err := strconv.Atoi(c.DefaultQuery("cols", "150"))
	if utils.WshandleError(wsConn, err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "32"))
	if utils.WshandleError(wsConn, err) {
		return
	}

	client, err := utils.NewSshClient(sshRequest.Ip, sshRequest.Port,
		sshRequest.Username, sshRequest.Password)

	if utils.WshandleError(wsConn, err) {
		return
	}
	defer client.Close()
	//startTime := time.Now()
	ssConn, err := utils.NewSshConn(cols, rows, client)

	if utils.WshandleError(wsConn, err) {
		return
	}
	defer ssConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan

}

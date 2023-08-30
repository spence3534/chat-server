package models

import (
	"chat-server/models/common"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	common.GlobalModel
	ChatType int    // 消息类型群聊 私聊 广播
	TargetId uint   // 接收者
	FormId   uint   // 发送者
	MsgType  string // 消息类型 文字 图片 音频
	Content  string // 内容
	Desc     string // 描述
	Amount   int    // 其他数字统计
}

func (m *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系

var clientMap map[uint]*Node = make(map[uint]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(res http.ResponseWriter, req *http.Request) {
	// 校验参数
	query := req.URL.Query()
	// token := query.Get("token")
	id := query.Get("userId")
	aId, _ := strconv.Atoi(id)
	userId := uint(aId)
	//targetId := query.Get("targetId")
	//content := query.Get("content")
	//chatType := query.Get("chatType")
	isValida := true // checkToken()
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(res, req, nil)
	if err != nil {
		fmt.Println("err")
		return
	}
	// 2. 获取连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// 3. 用户关系
	// 4. userid跟node绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送逻辑
	go sendProc(node)
	// 6.完成接收逻辑
	go recvProc(node)
}

// 把数据队列中的消息取出来发送给前端，这里才是真正的发送
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws]sendMsg >>>>>> data", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 接收前端的消息，把消息存到udpsendChan中
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data)
		broadMsg(data)
		fmt.Println("[ws]recvProc <<<<", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine")
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpSendProc >>>> data:", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()

	if err != nil {
		fmt.Println(err)
	}

	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc >>>> data:", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

// 调度逻辑处理 发送消息给指定的目标用户
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
	}
	switch msg.ChatType {
	case 1:
		fmt.Println("dispatch >>>> data:", string(data))
		sendMsg(msg.TargetId, data)
		//case 2: sendGroupMsg()
		//case 3: sendAllMsg()
	}
}

// 把消息存到node.DataQueue的队列中去
func sendMsg(userId uint, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

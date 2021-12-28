package ipc

import (
	"bilibili-live-communication/bilibili/biliAPI"
	"bilibili-live-communication/bilibili/biliBinConv"
	"bilibili-live-communication/bilibili/biliJsonConv"
	"bilibili-live-communication/bilibili/biliWebsocket"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net"
)

const debug = true

type EncapInterface struct {
	// 筛选
	filtrate []biliJsonConv.CmdType

	// room id
	roomId      string
	trustRoomId string

	// conn, 与外部通信
	conn net.Conn

	// ReceiveCh
	receiveCh chan []byte
}

func (ei *EncapInterface) IPC(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	ei.conn = conn

	// 获取配置数据
	headerJson, err := ei.ReadBytes()
	if err != nil {
		return
	}
	// 解析JSON
	ei.roomId = jsoniter.Get(headerJson, "room_id").ToString()
	if debug {
		log.Println("Room ID: " + ei.roomId)
	}

	// 创建 websocket
	ei.trustRoomId, err = biliAPI.GetRoomInfo(ei.roomId)
	if err != nil {
		log.Println(err)
		return
	}
	hosts, token, err := biliAPI.GetDanmuInfo(ei.trustRoomId)
	if err != nil {
		log.Println(err)
		return
	}
	wsHost := fmt.Sprintf("wss://%s/sub", hosts[0])
	log.Println(wsHost)
	c, _, err := websocket.DefaultDialer.Dial(wsHost, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}(c)

	sendChannel := make(chan biliBinConv.EncodeArgs)
	receiveChannel := make(chan []byte)
	blw := biliWebsocket.New(c, &receiveChannel, &sendChannel, func(blw *biliWebsocket.BiliWebsocket) biliBinConv.EncodeArgs {
		data := fmt.Sprintf("{\"uid\":0,\"roomid\":%s,\"protover\":3,\"platform\":\"web\",\"type\":2,\"key\":\"%s\"}", ei.trustRoomId, token)
		return biliBinConv.EncodeArgs{
			Version: 0x0001, Operation: 0x00000007, Sequence: 0x00000001, Body: []byte(data),
		}
	})

	for {
		select {
		case <-blw.IsClosedCh:
			return
		case b, ok := <-receiveChannel:
			if !ok {
				return
			}
			switch biliJsonConv.GetCmdType(b) {
			case biliJsonConv.DanMuMsgType:
				log.Println(string(b))
				_, err := ei.Write(b)
				if err != nil {
					return
				}
			}
		}
	}
}

func (ei *EncapInterface) Write(b []byte) (n int, err error) {
	// 发送内容长度
	headerBytes := make([]byte, 8)
	headerLength := len(headerBytes)
	count := 0
	for count < headerLength {
		binary.LittleEndian.PutUint64(headerBytes, uint64(len(b)))
		n, err = ei.conn.Write(headerBytes)
		if err != nil {
			return count, err
		} else if n != len(headerBytes) {
			return count, errors.New("Write Length: 0")
		}
		count += n
	}

	// 发送数据
	bodyLength := len(b)
	count = 0
	for count < bodyLength {
		n, err = ei.conn.Write(b[count:])
		if err != nil {
			return count, err
		} else if n == 0 {
			return count, err
		}
		count += n
	}
	return count, nil
}

func (ei *EncapInterface) ReadBytes() ([]byte, error) {
	var b []byte
	// 读取header
	lengthBytes, err := ei.ReadCount(8)
	if err != nil {
		return b, err
	}
	log.Println(lengthBytes)
	length := binary.LittleEndian.Uint64(lengthBytes)
	if debug {
		log.Println("Length: ", length)
	}

	// 读取主体数据
	b, err = ei.ReadCount(int(length))
	log.Println(string(b))

	return b, err
}

func (ei *EncapInterface) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (ei *EncapInterface) ReadCount(count int) ([]byte, error) {
	var readCount = 0
	var b []byte

	for readCount < count {
		b_ := make([]byte, count-readCount)
		n, err := ei.conn.Read(b_)
		if err != nil {
			return b, err
		} else if n == 0 {
			return b, errors.New("Read Length: 0")
		}
		readCount += n
		b = append(b, b_[:n]...)
	}

	return b, nil
}

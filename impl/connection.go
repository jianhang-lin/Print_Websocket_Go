package impl

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"sync"
	"unsafe"
)

type Connection struct {
	wsConnect *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte

	mutex    sync.Mutex // 对closeChan关闭上锁
	isClosed bool       // 防止closeChan被关闭多次
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConnect: wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()
	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
		// doPrintLabel(data)
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func doPrintLabel(data []byte) {
	labelPrinter := LabelPrinter{}
	err := json.Unmarshal(data, &labelPrinter)
	//fmt.Println(labelPrinter.LabelData)
	if err != nil {
		// fmt.Println(err.Error())
		Error.Printf("json failed convert to labelPrinter: %s", err.Error())
	}
	printLabel(labelPrinter)
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *Connection) Close() {
	// 线程安全，可多次调用
	_ = conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// 内部实现
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}
		// fmt.Printf("readLoop = %s\n", String(data))
		Info.Printf("readLoop = %s", String(data))
		select {
		case conn.inChan <- data:
			doPrintLabel(data)
		case <-conn.closeChan: // closeChan 感知 conn断开
			goto ERR
		}

	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	select {
	case data = <-conn.outChan:
	case <-conn.closeChan:
		goto ERR
	}
	for {
		data = <-conn.outChan
		// fmt.Printf("writeLoop = %s\n", String(data))
		Info.Printf("writeLoop = %s", String(data))
		if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}

func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

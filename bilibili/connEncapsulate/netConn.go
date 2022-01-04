package connEncapsulate

import (
	"encoding/binary"
	"net"
	"time"
)

//  net.Conn
type netConn struct {
	conn    *net.Conn
	bufSize int
}

func NetConn(conn *net.Conn) Conn {
	var nC = netConn{conn: conn, bufSize: 1024}
	return &nC
}

func (nc *netConn) SetDeadline(t time.Time) error {
	return (*nc.conn).SetDeadline(t)
}

func (nc *netConn) SetWriteDeadline(t time.Time) error {
	return (*nc.conn).SetWriteDeadline(t)
}

func (nc *netConn) SetReadDeadline(t time.Time) error {
	return (*nc.conn).SetReadDeadline(t)
}

func (nc *netConn) Close() error {
	return (*nc.conn).Close()
}

func (nc *netConn) RemoteAddr() net.Addr {
	return (*nc.conn).RemoteAddr()
}

func (nc *netConn) LocalAddr() net.Addr {
	return (*nc.conn).LocalAddr()
}

// WriteMessage 经过封装的 WriteMessage, 发送数据前会发送8字节的长度信息
func (nc *netConn) WriteMessage(_ int, data []byte) error {
	lengthArray := make([]byte, 8)

	// 序列化长度信息
	binary.LittleEndian.PutUint64(lengthArray, uint64(len(data)))
	// 发送长度信息
	if err := nc.WriteFixedLength(len(lengthArray), lengthArray); err != nil {
		return err
	}
	// 发送主体数据
	return nc.WriteFixedLength(len(data), data)
}

// WriteFixedLength 写入定长数据
func (nc *netConn) WriteFixedLength(length int, data []byte) error {
	for count := 0; count < length; {
		n, err := nc.Write(data[count:])
		count += n
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadMessage 经过封装的 ReadMessage, 会去除头部8字节的长度信息部分
func (nc *netConn) ReadMessage() (int, []byte, error) {
	var b []byte
	var err error = nil
	var length uint64
	var lengthArray []byte

	// 读取8字节
	lengthArray, err = nc.ReadFixedLength(8)
	if err != nil {
		return 0, nil, err
	}
	// 反序列化
	length = binary.LittleEndian.Uint64(lengthArray)

	// 读取 length 长度的字节
	b, err = nc.ReadFixedLength(int(length))

	return len(b), b, err
}

// ReadFixedLength 读取定长数据
func (nc *netConn) ReadFixedLength(length int) ([]byte, error) {
	var b []byte
	var err error = nil
	var readCount = 0
	var n int

	for readCount < length {
		_b_ := make([]byte, length-readCount)
		n, err = nc.Read(_b_)
		b = append(b, _b_[:n]...)
		if err != nil {
			break
		}
		readCount += n
	}

	return b, err
}

func (nc *netConn) Write(b []byte) (int, error) {
	return (*nc.conn).Write(b)
}

func (nc *netConn) Read(b []byte) (int, error) {
	return (*nc.conn).Read(b)
}

package ipc

import (
	"bilibili-live-communication/bilibili/connEncapsulate"
	"encoding/binary"
	"net"
)

type Conn struct {
	conn net.Conn
}

func NetConnEncapsulate(conn net.Conn) connEncapsulate.ReadWriteCloser {
	c := new(Conn)
	c.conn = conn

	return c
}

func (self *Conn) Close() error {
	return self.conn.Close()
}

func (self *Conn) WriteMessage(_ int, data []byte) error {
	lengthBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lengthBytes, uint64(len(data)))
	if err := self.writeCount(lengthBytes); err != nil {
		return err
	}

	return self.writeCount(data)
}

func (self *Conn) writeCount(data []byte) error {
	count := len(data)
	for i := 0; i < count; {
		n, err := self.conn.Write(data[i:])
		if err != nil {
			return err
		}
		i += n
	}
	return nil
}

func (self *Conn) ReadMessage() (int, []byte, error) {
	lengthBytes, err := self.readCount(8)
	if err != nil {
		return 0, nil, err
	}
	length := binary.LittleEndian.Uint64(lengthBytes)

	data, err := self.readCount(int(length))

	return len(data), data, err
}

func (self *Conn) readCount(count int) ([]byte, error) {
	var data []byte
	for i := 0; i < count; {
		tmp := make([]byte, count-i)
		n, err := self.conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		data = append(data, tmp[:n]...)
		i += n
	}
	return data, nil
}

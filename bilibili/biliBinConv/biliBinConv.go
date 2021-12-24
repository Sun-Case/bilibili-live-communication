package biliBinConv

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/brotli/go/cbrotli"
	"io"
	"log"
	"unsafe"
)

/*!
Encode:
	将 JSON数据 编码成 二进制数据
Decode:
	将 原始二进制数据 解码成多个 JSON数据（特殊情况可能是其他数据类型）
*/

/*!
注意使用了 unsafe.Sizeof(BaseHeader{})
这是一个坑，只是现在没问题
如果成员发生改变，可能会因对齐原因，大小有偏差导致数据处理错误
*/

const debug = false

type BaseHeader struct {
	PackLen    uint32
	HeaderSize uint16
	Version    uint16
	Operation  uint32
	Sequence   uint32
}

type BiliConvStruct struct {
	Header BaseHeader
	Body   []uint8

	//Bin []uint8
}

type EncodeArgs struct {
	Version   uint16
	Operation uint32
	Sequence  uint32
	Body      []uint8
}

func Encode(ver uint16, operation uint32, sequence uint32, body []uint8) ([]byte, error) {
	bh := BaseHeader{
		PackLen:    uint32(unsafe.Sizeof(BaseHeader{})) + uint32(len(body)),
		HeaderSize: uint16(unsafe.Sizeof(BaseHeader{})),
		Version:    ver,
		Operation:  operation,
		Sequence:   sequence}
	var buf bytes.Buffer
	var err error

	if err = binary.Write(&buf, binary.BigEndian, &bh); err == nil {
		err = binary.Write(&buf, binary.BigEndian, &body)
	}

	return buf.Bytes(), err
}

func Decode(b []byte) ([]BiliConvStruct, error) {
	var err error
	var buf = bytes.NewReader(b)
	var bcsSlice = make([]BiliConvStruct, 0)
	var bh = BaseHeader{}
	var bcs = BiliConvStruct{}
	var unZlib io.ReadCloser
	var unBrotli io.Reader

	// 解析 Header
	if buf.Len() < int(unsafe.Sizeof(BaseHeader{})) {
		err = errors.New(fmt.Sprintf("#1 Len error, buf len = %d, but need %d", buf.Len(), int(unsafe.Sizeof(BaseHeader{}))))
		goto Label
	}
	if err = binary.Read(buf, binary.BigEndian, &bh); err != nil {
		goto Label
	}

	// 判断 Version ，分别处理
	switch bh.Version {
	case 0:
		fallthrough
	case 1:
		if bh.Operation == 3 { // 心跳包回复
			if debug {
				log.Println("Heartbeat Receive")
			}
			goto Label
		}

		for {
			if int(bh.PackLen)-int(bh.HeaderSize) > buf.Len() {
				err = errors.New(fmt.Sprintf("#2 Len error, PackLen: %d, HeaderSize: %d, but BodyLen: %d\n", bh.PackLen, bh.HeaderSize, buf.Len()))
				goto Label
			}

			bcs.Header = bh
			bcs.Body = make([]byte, int(bh.PackLen)-int(bh.HeaderSize))
			err = binary.Read(buf, binary.BigEndian, &bcs.Body)

			// 添加到切片
			bcsSlice = append(bcsSlice, bcs)

			if buf.Len() == 0 {
				goto Label
			}
			if err = binary.Read(buf, binary.BigEndian, &bh); err != nil {
				goto Label
			}
		}

	case 2: // zlib
		if int(bh.PackLen)-int(bh.HeaderSize) != buf.Len() {
			err = errors.New(fmt.Sprintf("#3 Len error, PackLen: %d, HeaderSize: %d, but BodyLen: %d\n", bh.PackLen, bh.HeaderSize, buf.Len()))
			goto Label
		}
		bcs.Header = bh

		// 解压缩数据
		if unZlib, err = zlib.NewReader(buf); err != nil {
			goto Label
		}

		b := bytes.Buffer{}
		if _, err = b.ReadFrom(unZlib); err != nil {
			goto Label
		}
		bcs.Body = make([]byte, b.Len())
		copy(bcs.Body, b.Bytes())

		// 添加到切片
		bcsSlice = append(bcsSlice, bcs)

	case 3:
		for {
			if int(bh.PackLen)-int(bh.HeaderSize) <= buf.Len() {
				// 分割数据
				body := bytes.Buffer{}
				if _, err = io.CopyN(&body, buf, int64(bh.PackLen)-int64(bh.HeaderSize)); err != nil {
					goto Label
				}

				// 解压缩
				unBrotli = cbrotli.NewReader(&body)
				b := bytes.Buffer{}
				if _, err = b.ReadFrom(unBrotli); err != nil {
					goto Label
				}
				// 再次解析数据
				var sec []BiliConvStruct
				if sec, err = Decode(b.Bytes()); err != nil {
					goto Label
				}
				bcsSlice = append(bcsSlice, sec...)

			} else {
				err = errors.New(fmt.Sprintf("#4 Len error, PackLen: %d, HeaderSize: %d, but BodyLen: %d\n", bh.PackLen, bh.HeaderSize, buf.Len()))
				goto Label
			}

			if buf.Len() == 0 {
				goto Label
			} else if buf.Len() < int(unsafe.Sizeof(BaseHeader{})) {
				err = errors.New(fmt.Sprintf("#5 Len error, buf len = %d, but need %d", buf.Len(), int(unsafe.Sizeof(BaseHeader{}))))
				goto Label
			}
			if err = binary.Read(buf, binary.BigEndian, &bh); err != nil {
				goto Label
			}
		}
	}

Label:
	return bcsSlice, err
}

package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	Xlog "github.com/cihub/seelog"
	"time"
)

// send writes a generic object as JSON to the writer.
func SendJson(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

// send writes a generic object as JSON to the writer.
func Send(w io.Writer, data []byte) error {
	total := 0
	length := len(data)
	n,err := w.Write(data)		// 后续要使用for循环确保数据都正常发送
	if length != n {
		total += n
		writeCount := 0
		for total < length {
			n,err = w.Write(data[total:])		// 要使用for循环确保数据都正常发送
			if err != nil {
				Xlog.Errorf("Write data failed:%s", err.Error())
				return err
			}
			total += n
			writeCount++
			if writeCount > 100 {		// 长时间没有发送完数据则报错
				Xlog.Errorf("Write %d bytes, but it expect %d bytes", total, length)
				return errors.New(fmt.Sprintf("Write %d bytes, but it expect %d bytes", total, length))
			}
			time.Sleep(time.Millisecond*30)	// 休眠下再继续发送数据
		}
	}

	return err
}

func SendEx(w io.Writer, data []byte) error {
	length := len(data)
	n,err := w.Write(data)		// 后续要使用for循环确保数据都正常发送
	if (length != n) {
		Xlog.Errorf("wirte data failed:%s, epxectd is %d, but write %d",
			err.Error(),
			length,
			n)
	}
	return err
}
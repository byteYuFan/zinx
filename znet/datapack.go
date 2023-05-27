package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/byteYuFan/zinx/utils"
	"github.com/byteYuFan/zinx/zinterfance"
)

// DataPack  定义封包，拆包的实体
type DataPack struct {
}

// NewDataPack 返回一个实例
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头长度
func (dp *DataPack) GetHeadLen() uint32 {
	// 封装长度已经固定了
	// ID uint32 4byte
	// DataLen uint32 4byte
	return uint32(8)
}

// Pack 封包
func (dp *DataPack) Pack(msg zinterfance.IMessage) ([]byte, error) {
	// 创建一个byte字节流的缓存
	dataBuff := bytes.NewBuffer([]byte{})
	// 将dataLen写入dataBuff里面
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgID写入dataBuff里面
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	//将MsgData写入dataBuff里面
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// Unpack 拆包 只需要将包的head信息读出来，之后在根据head信息里的data的长度在进行一次读取
func (dp *DataPack) Unpack(binaryData []byte) (zinterfance.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 只解压head信息，得到dateLen和MsgID
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	// 判断是否超过最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg date receive")
	}
	return msg, nil
}

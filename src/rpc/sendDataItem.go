package rpc

// 发送的数据项
type sendDataItem struct {
	// Id
	id int32

	// 内容
	data []byte
}

// 创建发送的数据项对象
func newSendDataItem(_id int32, _data []byte) *sendDataItem {
	return &sendDataItem{
		id:   _id,
		data: _data,
	}
}

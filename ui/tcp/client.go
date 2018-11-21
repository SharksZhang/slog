package tcp

import (
	"fmt"
	"net"
	"os"
)

//默认的服务器地址
var (
	server = "127.0.0.1:8080"
)

type TcpClient struct {
	connection *net.TCPConn
	hawkServer *net.TCPAddr
	stopChan   chan struct{}
}

func (client *TcpClient) SendData(bytes []byte) {
	client.connection.Write(bytes)
	//fmt.Printf("Send log [%v]\n", string(bytes))

}

func NewTcpClient(hawkServer string) (*TcpClient, error) {
	client := &TcpClient{}
	var err error
	client.hawkServer, err = net.ResolveTCPAddr("tcp", hawkServer)
	if err != nil {
		fmt.Printf("hawk server [%s] resolve error: [%s]", server, err.Error())
		return nil, err
	}
	client.connection, err = net.DialTCP("tcp", nil, client.hawkServer)
	if err != nil {
		fmt.Printf("connect to hawk server error: [%s]", err.Error())
		os.Exit(1)
		return nil, err
	}
	client.stopChan = make(chan struct{})
	return client, nil
}

func (client *TcpClient) Close() {
	if client.connection != nil {
		client.connection.Close()
	}

}

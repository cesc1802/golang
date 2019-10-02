package server

import (
	"golang/iso8583"
	"golang/utils"
	"container/list"
	"net"
	"os"
	"sync"
)

// Server basic properties for server
type Server struct {
	sync.Mutex
	sync.WaitGroup
	mServerType  string
	mIPAddr      string
	mPort        string
	mListenner   net.Listener
	mConnections *list.List
	mIsRunning   bool
	mLengthType  iso8583.MessageLengthType
}

// NewServer init param for server, need to be read from config file
func NewServer() *Server {
	return &Server{
		mServerType:  os.Getenv("SERVER_TYPE"),
		mIPAddr:      os.Getenv("SERVER_ADDR"),
		mPort:        os.Getenv("SERVER_PORT"),
		mLengthType:  iso8583.ToMessageLengthType(os.Getenv("SERVER_LENGTH_TYPE")),
		mConnections: list.New()}
}

// Start start server
func (s *Server) Start() {
	if s.mIsRunning {
		return
	}

	var err error
	CStr := s.mIPAddr + ":" + s.mPort

	utils.GetLog().Info("Local Address: ", s.mIPAddr, " Port: ", s.mPort)
	s.mListenner, err = net.Listen(s.mServerType, CStr)
	if err != nil {
		utils.GetLog().Error("IPAddr or Port is not valid")
		panic(err)
	}
	s.mIsRunning = true
	go s.doAccept()
}

// DoAccept accept new client connected to server
func (s *Server) doAccept() {
	var err error
	var streamer net.Conn
	for s.mIsRunning {
		streamer, err = s.mListenner.Accept()
		if err != nil {
			utils.GetLog().Error("Server can't accept new client")
			panic(err)
		}
		utils.GetLog().Info("A remote client connected from IP: ", streamer.RemoteAddr())
		client := &ISO8583Client{
			mClientCon: streamer,
			mServer:    s,
		}
		client.Add(2)
		go client.Listen()
		go client.ProcessMessage()
		client.Wait()
	}
}
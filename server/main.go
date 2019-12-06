package main

import (
	"github/wziww/medusa"
	"github/wziww/medusa/config"
	"github/wziww/medusa/encrpt"
	"github/wziww/medusa/log"
	"github/wziww/medusa/stream"
	"net"
	"os"
	"strconv"
)

func main() {
	addr, resoveErr := net.ResolveTCPAddr("tcp", "0.0.0.0:"+strconv.Itoa(config.C.Server.Port))
	if resoveErr != nil {
		log.FMTLog(log.LOGERROR, resoveErr)
		os.Exit(0)
	}
	config.C.Base.Client = false
	log.FMTLog(log.LOGINFO, "server start")
	stream.APIServerInit()
	password := []byte(config.C.Base.Password)
	encryptor := encrpt.InitEncrypto(&password, config.C.Base.Crypto)
	if encryptor == nil {
		log.FMTLog(log.LOGERROR, "unsupport encrypto:", config.C.Base.Crypto)
		os.Exit(0)
	}
	listener, listenErr := net.ListenTCP("tcp", addr)
	if listenErr != nil {
		log.FMTLog(log.LOGERROR, listenErr)
		os.Exit(0)
	}
	log.FMTLog(log.LOGINFO, "service start listen at", addr)
	defer listener.Close()
	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			continue
		}
		// log.FMTLog(log.LOGINFO, localConn.RemoteAddr(), "connected")
		// localConn被关闭时直接清除所有数据 不管没有发送的数据
		localConn.SetLinger(0)
		go handleConn(&medusa.TCPConn{
			ReadWriteCloser: localConn,
			Encryptor:       encryptor,
		})
	}
}

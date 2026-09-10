package kv

import (
	"io"
	"log"
	"net"
	"runtime"
	"unsafe"
)

const ClientBufLength = 1024
const ClientHeaderLength = unsafe.Sizeof(uint32(42))

type Server struct {
	m *Map

	ln net.Listener

	msgCh    chan Msg
	shutdown chan struct{}
}

func NewServer(m *Map) *Server {
	return &Server{
		m: m,

		msgCh:    make(chan Msg),
		shutdown: make(chan struct{}, 1),
	}
}

func (s *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.ln = ln
	go s.acceptor()
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.processMsgs()
	}

	return nil
}

func (s *Server) acceptor() {
	for {
		c, err := s.ln.Accept()
		if err != nil {
			log.Println("failed to accept:", err)
			continue
		}

		go s.client(c)
	}
}

func (s *Server) client(c net.Conn) {
	defer c.Close()

	hd := make([]byte, ClientHeaderLength)

	for {
		_, err := io.ReadFull(c, hd)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("header read error:", err)
			break
		}

		plen := MsgLen(hd)
		pbuf := make([]byte, plen)

		_, err = io.ReadFull(c, pbuf)
		if err != nil {
			log.Println("payload read error:", err)
			break
		}

		msg := Msg{
			Data: pbuf,
			Conn: c,
		}

		s.msgCh <- msg
	}
}

func (s *Server) processMsgs() {
	for msg := range s.msgCh {
		s.onMsg(msg)
	}
}

func (s *Server) onMsg(msg Msg) {
	args := ReadPayload(msg.Data)
	log.Println(args)
}

func (s *Server) Shutdown() {
	close(s.shutdown)
}

func (s *Server) WaitForShutdown() {
	<-s.shutdown
}

package kv

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"runtime"
	"unsafe"
)

const ClientBufLength = 1024
const ClientHeaderLength = unsafe.Sizeof(uint32(42))

type CmdHandler func(c net.Conn, args []string) bool

type Server struct {
	m *Map

	ln net.Listener

	msgCh    chan Msg
	shutdown chan struct{}

	cmdTable map[string][]CmdHandler
}

func NewServer(m *Map) *Server {
	s := &Server{
		m: m,

		msgCh:    make(chan Msg),
		shutdown: make(chan struct{}, 1),
	}

	s.cmdTable = map[string][]CmdHandler{
		"SET":    {s.requireArgs(2), s.put},
		"GET":    {s.requireArgs(1), s.get},
		"EXISTS": {s.requireArgs(1), s.exists},
		"DEL":    {s.requireArgs(1), s.del},
	}

	return s
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

	if len(args) == 0 {
		s.sendMsg(msg.Conn, []string{"ERR", "missing command name"})
		return
	}

	cmd := args[0]
	chain, ok := s.cmdTable[cmd]
	if !ok {
		s.sendMsg(msg.Conn, []string{"ERR", "unknown command"})
		return
	}

	for _, c := range chain {
		if !c(msg.Conn, args[1:]) {
			break
		}
	}
}

func (s *Server) sendMsg(c net.Conn, args []string) (int, error) {
	b := new(bytes.Buffer)
	for i := 0; i < 4; i++ {
		b.WriteByte(0)
	}

	for _, s := range args {
		b.WriteString(s)
		b.WriteByte(0)
	}

	buf := b.Bytes()
	binary.LittleEndian.PutUint32(buf, uint32(b.Len()-4))

	return c.Write(buf)
}

func (s *Server) requireArgs(n int) CmdHandler {
	return func(c net.Conn, args []string) bool {
		if len(args) < n {
			return false
		}
		return true
	}
}

func (s *Server) put(c net.Conn, args []string) bool {
	s.m.Put(args[0], args[1])
	s.sendMsg(c, []string{"OK"})
	return true
}

func (s *Server) get(c net.Conn, args []string) bool {
	v, ok := s.m.Get(args[0]).(string)
	if !ok {
		v = "nil"
	}

	s.sendMsg(c, []string{v})

	return true
}

func (s *Server) exists(c net.Conn, args []string) bool {
	exists := s.m.Has(args[0])
	if exists {
		s.sendMsg(c, []string{"YES"})
	} else {
		s.sendMsg(c, []string{"NO"})
	}

	return true
}

func (s *Server) del(c net.Conn, args []string) bool {
	s.m.Delete(args[0])
	s.sendMsg(c, []string{"OK"})
	return true
}

func (s *Server) Shutdown() {
	close(s.shutdown)
}

func (s *Server) WaitForShutdown() {
	<-s.shutdown
}

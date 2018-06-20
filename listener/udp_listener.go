package listener

import (
	"github.com/myzhan/goreplay-udp/proto"
	"log"
	"strconv"
)

type UDPListener struct {
	// IP to listen
	addr string
	// Port to listen
	port uint16

	messagesChan chan *proto.UDPMessage

	underlying *IPListener
}

func NewUDPListener(addr string, port string, trackResponse bool) (l *UDPListener) {
	l = &UDPListener{}
	l.messagesChan = make(chan *proto.UDPMessage, 10000)
	l.addr = addr
	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Invaild Port: %s, %v\n", port, err)
	}
	l.port = uint16(intPort)

	l.underlying = NewIPListener(addr, l.port, trackResponse)

	if l.underlying.IsReady() {
		go l.recv()
	} else {
		log.Fatalln("IP Listener is not ready after 5 seconds")
	}

	return
}

func (l *UDPListener) parseUDPPacket(packet *ipPacket) (message *proto.UDPMessage) {
	data := packet.payload
	message = proto.NewUDPMessage(data, false)
	if message.DstPort == l.port {
		message.IsIncoming = true
	}
	message.Start = packet.timestamp
	return
}

func (l *UDPListener) recv() {
	for {
		ipPacketsChan := l.underlying.Receiver()
		select {
		case packet := <-ipPacketsChan:
			message := l.parseUDPPacket(packet)
			l.messagesChan <- message
		}
	}
}

func (l *UDPListener) Receiver() chan *proto.UDPMessage {
	return l.messagesChan
}

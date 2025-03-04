package socket

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
)

// Listener ...
// Capture traffic from socket using RAW_SOCKET's
// http://en.wikipedia.org/wiki/Raw_socket
//
// RAW_SOCKET allow you listen for traffic on any port (e.g. sniffing) because they operate on IP level.
// Ports is TCP feature, same as flow control, reliable transmission and etc.
// Since we can't use default TCP libraries RAWTCPLitener implements own TCP layer
// TCP packets is parsed using tcp_packet.go, and flow control is managed by tcp_message.go
type Listener struct {
	messagesWaiting map[string]*TCPMessage // buffer of TCPMessages waiting to be send

	packets       chan *TCPPacket
	messagesReady chan *TCPMessage // Messages ready to be send to client

	messageDone chan *TCPMessage // Used for notifications about completed or expired messages

	addr string // IP to listen
	port int    // Port to listen
}

// NewListener ...
// RAWTCPListen creates a listener to capture traffic from RAW_SOCKET
func NewListener(addr string, port string) (rawListener *Listener) {
	rawListener = &Listener{}

	rawListener.packets = make(chan *TCPPacket, 100)
	rawListener.messagesReady = make(chan *TCPMessage, 100)
	rawListener.messageDone = make(chan *TCPMessage, 100)
	rawListener.messagesWaiting = make(map[string]*TCPMessage)

	rawListener.addr = addr
	rawListener.port, _ = strconv.Atoi(port)

	go rawListener.listen()
	go rawListener.readRAWSocket()

	return
}

func (t *Listener) listen() {
	for {
		select {
		// If message ready for deletion it means that its also complete or expired by timeout
		case message := <-t.messageDone:
			t.messagesReady <- message
			delete(t.messagesWaiting, message.ID)

		// We need to use channels to process each packet to avoid data races
		case packet := <-t.packets:
			t.processTCPPacket(packet)
		}
	}
}

func (t *Listener) readRAWSocket() {
	conn, e := net.ListenPacket("ip4:tcp", t.addr)
	if e != nil {
		log.Fatal(e)
	}
	defer conn.Close()

	buf := make([]byte, 4096*2)

	for {
		// Note: ReadFrom receive messages without IP header
		n, addr, err := conn.ReadFrom(buf)

		if err != nil {
			log.Println("Error:", err)
			continue
		}

		if n > 0 {
			t.parsePacket(addr, buf[:n])
		}
	}
}

func (t *Listener) parsePacket(addr net.Addr, buf []byte) {
	if t.isIncomingDataPacket(buf) {
		newBuf := make([]byte, len(buf))
		copy(newBuf, buf)

		t.packets <- ParseTCPPacket(addr, newBuf)
	}
}

func (t *Listener) isIncomingDataPacket(buf []byte) bool {
	// To avoid full packet parsing every time, we manually parsing values needed for packet filtering
	// http://en.wikipedia.org/wiki/Transmission_Control_Protocol
	destPort := binary.BigEndian.Uint16(buf[2:4])

	// Because RAW_SOCKET can't be bound to port, we have to control it by ourself
	if int(destPort) == t.port {
		// Get the 'data offset' (size of the TCP header in 32-bit words)
		dataOffset := (buf[12] & 0xF0) >> 4

		// We need only packets with data inside
		// Check that the buffer is larger than the size of the TCP header
		if len(buf) > int(dataOffset*4) {
			// We should create new buffer because go slices is pointers. So buffer data shoud be immutable.
			return true
		}
	}

	return false
}

// Trying to add packet to existing message or creating new message
//
// For TCP message unique id is Acknowledgment number (see tcp_packet.go)
func (t *Listener) processTCPPacket(packet *TCPPacket) {
	defer func() { recover() }()

	var message *TCPMessage
	messageID := packet.Addr.String() + strconv.Itoa(int(packet.Ack))

	message, ok := t.messagesWaiting[messageID]

	if !ok {
		// We sending c_del_message channel, so message object can communicate with Listener and notify it if message completed
		message = NewTCPMessage(messageID, t.messageDone)
		t.messagesWaiting[messageID] = message
	}

	// Adding packet to message
	message.cPackets <- packet
}

// Receive TCP messages from the listener channel
func (t *Listener) Receive() *TCPMessage {
	return <-t.messagesReady
}

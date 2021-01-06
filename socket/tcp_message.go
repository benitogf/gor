package socket

import (
	"log"
	"sort"
	"time"
)

// MsgExpiration ...
const MsgExpiration = 2000 * time.Millisecond

// TCPMessage ensure that all TCP packets for given request is received, and processed in right sequence
// Its needed because all TCP message can be fragmented or re-transmitted
//
// Each TCP Packet have 2 ids: acknowledgment - message_id, and sequence - packet_id
// Message can be compiled from unique packets with same message_id which sorted by sequence
// Message is received if we didn't receive any packets for 2000ms
type TCPMessage struct {
	ID      string // Message ID
	packets []*TCPPacket

	timer *time.Timer // Used for expire check

	cPackets chan *TCPPacket

	cDelMessage chan *TCPMessage
}

// NewTCPMessage pointer created from a Acknowledgment number and a channel of messages readuy to be deleted
func NewTCPMessage(ID string, cDel chan *TCPMessage) (msg *TCPMessage) {
	msg = &TCPMessage{ID: ID}

	msg.cPackets = make(chan *TCPPacket)
	msg.cDelMessage = cDel // used for notifying that message completed or expired

	// Every time we receive packet we reset this timer
	msg.timer = time.AfterFunc(MsgExpiration, msg.Timeout)

	go msg.listen()

	return
}

func (t *TCPMessage) listen() {
	for {
		select {
		case packet, more := <-t.cPackets:
			if more {
				t.AddPacket(packet)
			} else {
				// Stop loop if channel closed
				return
			}
		}
	}
}

// Timeout notifies message to stop listening, close channel and message ready to be sent
func (t *TCPMessage) Timeout() {
	close(t.cPackets)  // Notify to stop listen loop and close channel
	t.cDelMessage <- t // Notify RAWListener that message is ready to be send to replay server
}

// Bytes sorts packets in right orders and return message content
func (t *TCPMessage) Bytes() (output []byte) {
	sort.Sort(SortBySeq(t.packets))

	for _, v := range t.packets {
		output = append(output, v.Data...)
	}

	return
}

// AddPacket to the message and ensure packet uniqueness
// TCP allows that packet can be re-send multiple times
func (t *TCPMessage) AddPacket(packet *TCPPacket) {
	packetFound := false

	for _, pkt := range t.packets {
		if packet.Seq == pkt.Seq {
			packetFound = true
			break
		}
	}

	if packetFound {
		log.Println("Received packet with same sequence")
	} else {
		t.packets = append(t.packets, packet)
	}

	// Reset message timeout timer
	t.timer.Reset(MsgExpiration)
}

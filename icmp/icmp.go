package icmp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	TargetIP    *net.IPAddr
	ListenIP    string
	BlockSize   = 1200
	VerboseMode bool
)

type MessageData struct {
	Type     int
	Code     int
	Checksum int
	Id       int
	Seq      int
	Data     []byte
	DataLen  int
	IP       string
}

// Send message to TargetIP using ICMP.
func MessageSend(c *icmp.PacketConn, message *string) {
	var totalBlocks int
	messageLen := len(*message)
	n := float32(messageLen) / float32(BlockSize)
	if n > float32(int(n)) {
		totalBlocks = int(n) + 1
	} else {
		totalBlocks = int(n)
	}
	totalBlocksStr := strconv.Itoa(totalBlocks) + " " + strconv.Itoa(messageLen)
	if c.IPv6PacketConn() == nil {
		RawMessageSend(c, &totalBlocksStr, 0, 0, ipv4.ICMPTypeEcho, 0)
	} else {
		RawMessageSend(c, &totalBlocksStr, 0, 0, ipv6.ICMPTypeEchoRequest, 0)
	}
	if VerboseMode {
		fmt.Printf("Sending %v bytes to %v in %v blocks\n", messageLen, TargetIP, totalBlocks)
	}
	var tempInt = 1
	for i := 0; i < messageLen; i += BlockSize {
		end := i + BlockSize
		if end > messageLen {
			end = messageLen
		}
		block := (*message)[i:end]
		lenBlock := len(block)
		var id, seq uint16
		if lenBlock >= 1 {
			if lenBlock == 1 {
				id = CombineInt8(block[0], 0)
			} else {
				id = CombineInt8(block[0], block[1])
			}
		}
		if lenBlock >= 3 {
			if lenBlock == 3 {
				seq = CombineInt8(block[2], 0)
			} else {
				seq = CombineInt8(block[2], block[3])
			}
		}
		if lenBlock > 4 {
			block = block[4:]
		} else {
			block = ""
		}
		if VerboseMode {
			fmt.Printf("Block %v: %v bytes sent\n", tempInt, lenBlock)
		}
		tempInt++
		if c.IPv6PacketConn() == nil {
			RawMessageSend(c, &block, id, seq, ipv4.ICMPTypeEcho, 0)
		} else {
			RawMessageSend(c, &block, id, seq, ipv6.ICMPTypeEchoRequest, 0)
		}
	}
}

// Read message from ICMP.
func MessageRead(c *icmp.PacketConn) (*string, error) {
	var message string
	var receivedBytes, messageLen int
	MessageData := RawMessageRead(c)
	x := string(MessageData.Data)
	x2 := strings.Split(x, " ")
	num, err := strconv.Atoi(x2[0])
	messageLen, err2 := strconv.Atoi(x2[1])
	if err != nil && err2 != nil {
		return &message, errors.New("invalid Block length received")
	}
	if VerboseMode {
		fmt.Printf("Receiving %v bytes to %v in %v blocks\n", messageLen, MessageData.IP, num)
	}
	for tempInt := range num {
		messageData := RawMessageRead(c)
		msgLen := messageData.DataLen + 4
		if VerboseMode {
			fmt.Printf("Block %v: %v bytes received\n", tempInt+1, msgLen)
		}
		tempMessage := make([]byte, msgLen)
		receivedBytes += msgLen
		tempMessage[0], tempMessage[1] = SplitInt16(uint16(messageData.Id))
		tempMessage[2], tempMessage[3] = SplitInt16(uint16(messageData.Seq))
		for i := range messageData.DataLen {
			tempMessage[i+4] = messageData.Data[i]
		}
		message += string(tempMessage)
	}
	message = message[:messageLen]
	return &message, nil
}

func RawMessageSend(c *icmp.PacketConn, message *string, id uint16, seq uint16, mtype icmp.Type, code int) {
	wm := icmp.Message{
		Type: mtype, Code: code,
		Body: &icmp.Echo{
			ID: int(id), Seq: int(seq),
			Data: []byte(*message),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := c.WriteTo(wb, TargetIP); err != nil {
		log.Fatalf("WriteTo err, %s", err)
	}
}

func RawMessageRead(c *icmp.PacketConn) MessageData {
	rawBytes := make([]byte, 65507)
	n, peer, err := c.ReadFrom(rawBytes)
	if err != nil {
		log.Fatalf("ReadFrom err, %s", err)
	}
	for {
		if peer.String() == ListenIP || ListenIP == "" {
			break
		} else {
			n, peer, err = c.ReadFrom(rawBytes)
		}
	}
	if err != nil {
		log.Fatalf("ReadFrom err, %s", err)
	}

	MessageData := MessageData{
		Type:     int(rawBytes[0]),
		Code:     int(rawBytes[1]),
		Checksum: int(CombineInt8(rawBytes[2], rawBytes[3])),
		Id:       int(CombineInt8(rawBytes[4], rawBytes[5])),
		Seq:      int(CombineInt8(rawBytes[6], rawBytes[7])),
		Data:     rawBytes[8:n],
		DataLen:  n - 8,
		IP:       peer.String(),
	}
	return MessageData
}

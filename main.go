package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	icmpa "github.com/aftab700/pingo/icmp"
	"golang.org/x/net/icmp"
)

var messagePtr *string
var targetPtr *string
var listeningStr stringFlag
var helpPtr *bool
var quietPtr *bool
var myFlagSet = flag.NewFlagSet("Options", flag.ContinueOnError)
var c *icmp.PacketConn
var err error
var isIPv6 bool

func printUsage() {
	fmt.Println("Usage:\n  ./icmp-go [flags]\n\nFlags:")
	myFlagSet.PrintDefaults()
}

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}
func (sf *stringFlag) String() string {
	return sf.value
}

func parseFlags() {
	messagePtr = myFlagSet.String("m", "Hello", "Message string to send")
	targetPtr = myFlagSet.String("t", "", "IP address of the target or receiver")
	myFlagSet.Var(&listeningStr, "l", "Listen for incoming ICMP packets\nProvide an IP address to Receive ICMP packets from provided IP address only")
	myFlagSet.IntVar(&icmpa.BlockSize, "s", 1200, "ICMP packet Data block size in bytes\nAllowed size is between 5 and 65495 bytes")
	helpPtr = myFlagSet.Bool("h", false, "Print help")
	quietPtr = myFlagSet.Bool("q", false, "Quiet mode")
	myFlagSet.BoolVar(&icmpa.VerboseMode, "v", false, "Verbose mode")
	myFlagSet.SetOutput(io.Discard)
	// -l is string so if there is flag after -l[-l -h] will not consider -h as flag but value to -l
	var Args []string
	var temp = len(os.Args)
	for i, v := range os.Args[1:] {
		Args = append(Args, v)
		if v == "-l" && (i+2) < temp && os.Args[i+2][0] == 45 {
			Args = append(Args, "")
		}
	}

	err := myFlagSet.Parse(Args)
	myFlagSet.SetOutput(nil)

	if err != nil {
		if err.Error() == "flag needs an argument: -l" {
			listeningStr.set = true
		} else {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}

func init() {
	parseFlags()
	if *helpPtr {
		printUsage()
		os.Exit(0)
	}
	if !listeningStr.set && *targetPtr == "" {
		fmt.Println("Expected valid argument for the target or receiver")
		fmt.Println("EXAMPLE:\n  ./icmp-go -t 127.0.0.1")
		os.Exit(0)
	}
	if *targetPtr == "" {
		*targetPtr = "127.0.0.1"
	}
	icmpa.SetTargetIP(*targetPtr)
	icmpa.SetListenIP(listeningStr.value)
	var maxBlockSize = 65495
	if icmpa.BlockSize > maxBlockSize || icmpa.BlockSize < 5 {
		log.Fatal("Invalid ICMP packet Data block size, Allowed size is between 5 and 65495 bytes")
	}
}

func main() {
	if icmpa.IsIPv6(*targetPtr) {
		isIPv6 = true
	}
	if listeningStr.set && icmpa.ListenIP != "" && icmpa.IsIPv6(icmpa.ListenIP) {
		isIPv6 = true
	}
	if isIPv6 {
		c, err = icmp.ListenPacket("ip6:ipv6-icmp", "::")
		if err != nil {
			log.Fatalf("listen err, %s", err)
		}
		defer c.Close()
	} else {
		c, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			log.Fatalf("listen err, %s", err)
		}
		defer c.Close()
	}

	if listeningStr.set {
		if isIPv6 && !*quietPtr {
			fmt.Println("Listening Using IPv6 Packet Connection...")
		} else if !*quietPtr {
			fmt.Println("Listening Using IPv4 Packet Connection...")
		}
		message, err := icmpa.MessageRead(c)
		if err != nil {
			log.Fatalf("MessageRead err: %s", err)
		} else {
			if !*quietPtr {
				fmt.Print("Received Data: ")
			}
			fmt.Println(*message)
		}
	} else {
		if isIPv6 && !*quietPtr {
			fmt.Println("Sending Using IPv6 Packet Connection...")
		} else if !*quietPtr {
			fmt.Println("Sending Using IPv4 Packet Connection...")
		}
		icmpa.MessageSend(c, messagePtr)
		if !*quietPtr {
			fmt.Println("Message sent.")
		}
	}
}

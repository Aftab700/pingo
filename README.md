# icmp-go

# Overview

`icmp-go` is a Tool written in [Go](https://go.dev/) for sending data over the network using the ICMP protocol.\
It supports IPv4 and IPv6.

# Installation

[Download](https://github.com/Aftab700/icmp-go/releases/latest) prebuilt ready-to-run binary from the [releases page](https://github.com/Aftab700/icmp-go/releases/latest) or install using GO.

## Build from source

Prerequisites to build `icmp-go` from source:

- Go 1.20 or later

Run the following command to install the command-line tool:
```
go install github.com/aftab700/icmp-go@latest
```

or run the following commands to build from repo:
```
git clone https://github.com/Aftab700/icmp-go.git
cd icmp-go
go build
```

# Usage

```
icmp-go -h
```

This will display help for the tool. Here are all the switches it supports.
```
Usage:
  ./icmp-go [flags]

Flags:
  -h    Print help
  -l value
        Listen for incoming ICMP packets
        Provide an IP address to Receive ICMP packets from provided IP address only
  -m string
        Message string to send (default "Hello")
  -q    Quiet mode
  -s int
        ICMP packet Data block size in bytes
        Allowed size is between 5 and 65495 bytes (default 1200)
  -t string
        IP address of the target or receiver
  -v    Verbose mode
```

## Example usage:

> [!IMPORTANT]
> To successfully send the data, you must first run the listen function, and then send the data.

Listen for incomming packates:
```
$ icmp-go -l
```
Send data:
```
$ icmp-go -t 127.0.0.1 -m 'Hello There'
```

Output of the listener:
```
$ icmp-go -l
Listening Using IPv4 Packet Connection...
Received Data: Hello There
```

Output of the sender:
```
$ icmp-go -t 127.0.0.1 -m 'Hello There'
Sending Using IPv4 Packet Connection...
Message sent.
```

# Todo
- Add the option for sending the file and save the output to a file.
- Add an interactive chat feature.
- Add a feature to Encrypt or decrypt data before sending or receiving(AES, XOR, etc). Encode and decode using Base64 and more.

# License

ffuf is released under MIT license. See [LICENSE](https://github.com/Aftab700/icmp-go/blob/main/LICENSE).
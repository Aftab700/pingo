# ping Go

# Overview

`pingo` is a Tool written in [Go](https://go.dev/) for sending data over the network using the ICMP protocol.\
It supports IPv4 and IPv6.

# Installation

## Download binary

[Download](https://github.com/Aftab700/pingo/releases/latest) prebuilt ready-to-run binary from the [releases page](https://github.com/Aftab700/pingo/releases/latest) or install using GO.

## Build from source

Prerequisites to build `pingo` from source:

- Go 1.20 or later

Run the following command to install the command-line tool:
```
go install github.com/aftab700/pingo@latest
```

or run the following commands to build from repo:
```
git clone https://github.com/Aftab700/pingo.git
cd pingo
go build
```

# Usage

## Help flag

```
pingo -h
```

This will display help for the tool. Here are all the switches it supports.
```
Usage:
  ./pingo [flags]

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
$ pingo -l
```
Send data:
```
$ pingo -t 127.0.0.1 -m 'Hello There'
```

Output of the listener:
```
$ pingo -l
Listening Using IPv4 Packet Connection...
Received Data: Hello There
```

Output of the sender:
```
$ pingo -t 127.0.0.1 -m 'Hello There'
Sending Using IPv4 Packet Connection...
Message sent.
```

### Video demo:

https://github.com/Aftab700/pingo/assets/79740895/9aceb324-1718-4c52-8b45-88627418017d


# Todo
- Add the option for sending the file and save the output to a file.
- Add an interactive chat feature.
- Add a feature to Encrypt or decrypt data before sending or receiving(AES, XOR, etc). Encode and decode using Base64 and more.

# License

ffuf is released under MIT license. See [LICENSE](https://github.com/Aftab700/pingo/blob/main/LICENSE).

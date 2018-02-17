# Dolon

## Background
Dolon is one of the most notorious [spies in Greek mythology](https://en.wikipedia.org/wiki/Dolon_(mythology)), and that's what this project is all about. Using gopacket the dolon binary captures a complete HTTP stream and then it can either send the complete HttpStream data as a protobuf to clients via GRPC or it can print the request portion as a curl command.

## Usage

Running `go get -u github.com/kai5263499/dolon/...` downloads and builds the dolon and grpc client commands.

## References
* Inspired by [this python project](https://github.com/jullrich/pcap2curl)
* Indebted to [this example](https://github.com/google/gopacket/blob/master/examples/bidirectional/main.go) of a custom stream assembly implementation. 
* [This project](https://github.com/hsiafan/httpparse) was also useful
package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/kai5263499/dolon/generated"

	"github.com/sirupsen/logrus"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var err error
	var conn *grpc.ClientConn
	if true {
		conn, err = grpc.Dial("localhost:9001", grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial("localhost:9001",
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				// remove the following line if the server certificate is signed by a certificate authority
				InsecureSkipVerify: true,
			})))
	}

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := generated.NewTrafficClient(conn)

	ctx := context.Background()
	stream, err := c.GetHttpStream(ctx, &google_protobuf.Empty{})
	if err != nil {
		logrus.Fatalf("error creating content submit stream client err=%#v", err)
	}

	for {
		msg, _ := stream.Recv()
		if err != nil {
			if err == io.EOF {
				os.Exit(-1)
			}

			continue
		}

		if err == nil && msg != nil {
			spew.Dump(msg)
		}
	}
}

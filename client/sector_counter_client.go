package client

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"os"

	"github.com/zhouwu5222/sector-storage-counter/proto"
)

// Client ..
type Client struct {
	DialAddr string
}

// NewClient ..
func NewClient() *Client {
	rpcAddr, ok := os.LookupEnv("SC_LISTEN")
	if !ok {
		log.Println("NO SC_LISTEN ENV")
	}

	return &Client{
		DialAddr: rpcAddr,
	}
}

func (c *Client) connect() (proto.SectorIdGrpcClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(c.DialAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := proto.NewSectorIdGrpcClient(conn)
	return client, conn, nil
}

// GetSectorID
func (c *Client) GetSectorID(ctx context.Context, param string) (uint64, error) {
	client, conn, err := c.connect()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	req := new(proto.SectorIDRequest)
	req.Question = param

	resp, err := client.GetSectorID(ctx, req)
	if err != nil {
		return 0, err
	}
	return resp.Answer, nil
}

package cmd

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"testing"

	"github.com/zhouwu5222/sector-storage-counter/proto"
)

func TestClient(t *testing.T) {

	conn, err := grpc.Dial("127.0.0.1:10086", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	client := proto.NewSectorIdGrpcClient(conn)
	reqBody := new(proto.SectorIDRequest)
	response, err := client.GetSectorID(context.Background(), reqBody)
	fmt.Printf("id:%d\n", response.Answer)

}

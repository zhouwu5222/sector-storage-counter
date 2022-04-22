package cmd

import (
	"testing"

	"github.com/zhouwu5222/sector-storage-counter/server"
)

func TestServer(t *testing.T) {

	server.Run("/root/project3/sectorid/sectorid.data")

}

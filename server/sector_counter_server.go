package server

import (
	"context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/zhouwu5222/sector-storage-counter/proto"
)

// Service
type Service struct {
	SectorIDLk sync.RWMutex
	SectorID   uint64
	SCFilePath string
}

// GetSectorID 实现 GetSectorID 方法
func (s *Service) GetSectorID(ctx context.Context, req *proto.SectorIDRequest) (*proto.SectorIDResponse, error) {
	s.SectorIDLk.Lock()
	defer s.SectorIDLk.Unlock()
	s.SectorID++
	s.WriteSectorID()
	return &proto.SectorIDResponse{Answer: s.SectorID}, nil
}

// WriteSectorID 实现 WriteSectorID 方法
func (s *Service) WriteSectorID() {
	f, err := os.OpenFile(s.SCFilePath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	strID := strconv.FormatUint(s.SectorID, 10)
	_, _ = f.Write([]byte(strID))
}

// Read SectorId
func readFileSid(filePath string) (uint64, error) {
	// check file exists, if not created.
	if _, err := os.Stat(filePath); err != nil {
		f, err := os.Create(filePath)
		if err != nil {
			return 0, err
		}
		_, _ = f.Write([]byte("0"))
		f.Close()
		return 0, nil
	}

	// 存在历史文件
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	byteID, err := ioutil.ReadAll(f)
	if err != nil {
		return 0, err
	}

	stringID := strings.Replace(string(byteID), "\n", "", -1)   // 将最后的\n去掉
	sectorID, err := strconv.ParseUint(string(stringID), 0, 64) // 将字符型数字转化为uint64类型
	if err != nil {
		return 0, err
	}

	return sectorID, nil
}

// Run Server
func Run(scFilePath string) {
	rpcAddr, ok := os.LookupEnv("SC_LISTEN")
	if !ok {
		log.Println("NO SC_LISTEN ENV")
		return
	}

	sectorID, err := readFileSid(scFilePath)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("currn sectorid: ", sectorID)

	listener, err := net.Listen("tcp", rpcAddr) // 监听本地端口
	if err != nil {
		log.Println(err)
	}
	log.Println("grpc server Listing on", rpcAddr)

	grpcServer := grpc.NewServer()
	server := &Service{
		SectorID:   sectorID,
		SCFilePath: scFilePath,
	}
	proto.RegisterSectorIdGrpcServer(grpcServer, server)

	if err = grpcServer.Serve(listener); err != nil {
		log.Println(err)
	}
}

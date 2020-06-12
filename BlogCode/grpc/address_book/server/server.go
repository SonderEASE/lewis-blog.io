package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/address_book/protos"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"
)

var (
	port = flag.Int("port", 8888, "The server port")
	filepath = flag.String("path", "", "文件路径")
)


type addressBookServer struct {
	sync.RWMutex
	address_book pb.AddressBook
}

func (s *addressBookServer) Register (ctx context.Context, request *pb.RegisterReq) (*pb.RegisterRep, error) {
	for _, client := range s.address_book.People {
		if client.Name == request.PersonalInfo.Name {
			return &pb.RegisterRep{Result: "Username has been used"}, nil
		}
	}
	s.Lock()
	defer s.Unlock()
	s.address_book.People = append(s.address_book.People, request.PersonalInfo)
	return &pb.RegisterRep{Result: "Register Success"}, nil
}

func (s *addressBookServer) GetAddressByName(ctx context.Context, request *pb.GetAddressByNameReq) (*pb.GetAddressByNameRep, error) {
	for _, client := range s.address_book.People {
		if client.Name == request.Name {
			return &pb.GetAddressByNameRep{
				PersonalInfo: &pb.PersonalInfo{Name: client.Name, PhoneNumber: client.PhoneNumber, Address: client.Address},
				Result: true,
				}, nil
		}
	}
	return &pb.GetAddressByNameRep{Result: false}, nil
}

func newServer() *addressBookServer {
	s := &addressBookServer{address_book: pb.AddressBook{}}
	s.loadFeatures(*filepath)
	return s
}

// loadFeatures loads features from a JSON file.
func (s *addressBookServer) loadFeatures(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default features: %v", err)
		}
	} else {
		data = LocalData
	}
	if err := json.Unmarshal(data, &s.address_book); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}


func main () {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Start listen at ", *port)
	grpcServer := grpc.NewServer()
	address_server := newServer()
	pb.RegisterAddressBookServiceServer(grpcServer, address_server)
	go func() {
		for {
			address_server.RLock()
			log.Println("当前注册人数: ", len(address_server.address_book.People))
			log.Printf("注册信息: %v\n", address_server.address_book.People)
			address_server.RUnlock()
			time.Sleep(20 * time.Second)
		}
	}()
	_ = grpcServer.Serve(lis)
}


var LocalData = []byte(`{
	"people": [{"name": "Manual customer service", 
	"phone_number": "8888",
	"address": "Beijing"}, {
	"name": "After-sales service",
	"phone_number": "6666",
	"address": "Beijing"}
	]
}`)

package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/address_book/protos"
	"log"
	"time"
)

type ClientInfo struct {
	name string
	phone string
	address string
}


var (
	serverAddr  = flag.String("server-address", "localhost:8888", "服务器的地址")
	clientName  = flag.String("name", "lewis", "姓名")
	clientPhone = flag.String("tel", "12345678910", "电话")
	clientAddress = flag.String("address", "Universe Center", "地址")
)

var (
	clientInfo *ClientInfo
)

// Register
func Register(client pb.AddressBookServiceClient, info ClientInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	register_reply, err := client.Register(ctx, &pb.RegisterReq{
		PersonalInfo: &pb.PersonalInfo{
			Name: info.name,
			PhoneNumber: info.phone,
			Address: info.address}})

	if err != nil {
		log.Fatalf("%v.RegisterAndUpdate(_) = _, %v: ", client, err)
	}
	log.Println(register_reply)

}

// GetAddressByName
func GetAddressByName(client pb.AddressBookServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	address_book, err := client.GetAddressByName(ctx, &pb.GetAddressByNameReq{Name: "lewis"})
	if err != nil {
		log.Fatalf("%v.GetAddressBook(_) = _, %v: ", client, err)
	}
	log.Println(address_book)
}

func main () {
	flag.Parse()

	clientInfo = &ClientInfo{
		name:    *clientName,
		phone:   *clientPhone,
		address: *clientAddress,
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()


	client := pb.NewAddressBookServiceClient(conn)

	Register(client, *clientInfo)

}

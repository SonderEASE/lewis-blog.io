// Code generated by protoc-gen-go. DO NOT EDIT.
// source: address_book.proto

package address_book

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetAddressByNameReq struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Number               string   `protobuf:"bytes,2,opt,name=number,proto3" json:"number,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAddressByNameReq) Reset()         { *m = GetAddressByNameReq{} }
func (m *GetAddressByNameReq) String() string { return proto.CompactTextString(m) }
func (*GetAddressByNameReq) ProtoMessage()    {}
func (*GetAddressByNameReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_351d41b0f01c9adf, []int{0}
}

func (m *GetAddressByNameReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAddressByNameReq.Unmarshal(m, b)
}
func (m *GetAddressByNameReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAddressByNameReq.Marshal(b, m, deterministic)
}
func (m *GetAddressByNameReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAddressByNameReq.Merge(m, src)
}
func (m *GetAddressByNameReq) XXX_Size() int {
	return xxx_messageInfo_GetAddressByNameReq.Size(m)
}
func (m *GetAddressByNameReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAddressByNameReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetAddressByNameReq proto.InternalMessageInfo

func (m *GetAddressByNameReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GetAddressByNameReq) GetNumber() string {
	if m != nil {
		return m.Number
	}
	return ""
}

type GetAddressByNameRep struct {
	PersonalInfo         *PersonalInfo `protobuf:"bytes,1,opt,name=personal_info,json=personalInfo,proto3" json:"personal_info,omitempty"`
	Result               bool          `protobuf:"varint,2,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GetAddressByNameRep) Reset()         { *m = GetAddressByNameRep{} }
func (m *GetAddressByNameRep) String() string { return proto.CompactTextString(m) }
func (*GetAddressByNameRep) ProtoMessage()    {}
func (*GetAddressByNameRep) Descriptor() ([]byte, []int) {
	return fileDescriptor_351d41b0f01c9adf, []int{1}
}

func (m *GetAddressByNameRep) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAddressByNameRep.Unmarshal(m, b)
}
func (m *GetAddressByNameRep) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAddressByNameRep.Marshal(b, m, deterministic)
}
func (m *GetAddressByNameRep) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAddressByNameRep.Merge(m, src)
}
func (m *GetAddressByNameRep) XXX_Size() int {
	return xxx_messageInfo_GetAddressByNameRep.Size(m)
}
func (m *GetAddressByNameRep) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAddressByNameRep.DiscardUnknown(m)
}

var xxx_messageInfo_GetAddressByNameRep proto.InternalMessageInfo

func (m *GetAddressByNameRep) GetPersonalInfo() *PersonalInfo {
	if m != nil {
		return m.PersonalInfo
	}
	return nil
}

func (m *GetAddressByNameRep) GetResult() bool {
	if m != nil {
		return m.Result
	}
	return false
}

type RegisterRep struct {
	Result               string   `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRep) Reset()         { *m = RegisterRep{} }
func (m *RegisterRep) String() string { return proto.CompactTextString(m) }
func (*RegisterRep) ProtoMessage()    {}
func (*RegisterRep) Descriptor() ([]byte, []int) {
	return fileDescriptor_351d41b0f01c9adf, []int{2}
}

func (m *RegisterRep) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRep.Unmarshal(m, b)
}
func (m *RegisterRep) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRep.Marshal(b, m, deterministic)
}
func (m *RegisterRep) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRep.Merge(m, src)
}
func (m *RegisterRep) XXX_Size() int {
	return xxx_messageInfo_RegisterRep.Size(m)
}
func (m *RegisterRep) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRep.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRep proto.InternalMessageInfo

func (m *RegisterRep) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

type RegisterReq struct {
	PersonalInfo         *PersonalInfo `protobuf:"bytes,1,opt,name=personal_info,json=personalInfo,proto3" json:"personal_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RegisterReq) Reset()         { *m = RegisterReq{} }
func (m *RegisterReq) String() string { return proto.CompactTextString(m) }
func (*RegisterReq) ProtoMessage()    {}
func (*RegisterReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_351d41b0f01c9adf, []int{3}
}

func (m *RegisterReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterReq.Unmarshal(m, b)
}
func (m *RegisterReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterReq.Marshal(b, m, deterministic)
}
func (m *RegisterReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterReq.Merge(m, src)
}
func (m *RegisterReq) XXX_Size() int {
	return xxx_messageInfo_RegisterReq.Size(m)
}
func (m *RegisterReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterReq.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterReq proto.InternalMessageInfo

func (m *RegisterReq) GetPersonalInfo() *PersonalInfo {
	if m != nil {
		return m.PersonalInfo
	}
	return nil
}

type PersonalInfo struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	PhoneNumber          string   `protobuf:"bytes,2,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	Address              string   `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PersonalInfo) Reset()         { *m = PersonalInfo{} }
func (m *PersonalInfo) String() string { return proto.CompactTextString(m) }
func (*PersonalInfo) ProtoMessage()    {}
func (*PersonalInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_351d41b0f01c9adf, []int{4}
}

func (m *PersonalInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PersonalInfo.Unmarshal(m, b)
}
func (m *PersonalInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PersonalInfo.Marshal(b, m, deterministic)
}
func (m *PersonalInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PersonalInfo.Merge(m, src)
}
func (m *PersonalInfo) XXX_Size() int {
	return xxx_messageInfo_PersonalInfo.Size(m)
}
func (m *PersonalInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PersonalInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PersonalInfo proto.InternalMessageInfo

func (m *PersonalInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PersonalInfo) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

func (m *PersonalInfo) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type AddressBook struct {
	People               []*PersonalInfo `protobuf:"bytes,1,rep,name=people,proto3" json:"people,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AddressBook) Reset()         { *m = AddressBook{} }
func (m *AddressBook) String() string { return proto.CompactTextString(m) }
func (*AddressBook) ProtoMessage()    {}
func (*AddressBook) Descriptor() ([]byte, []int) {
	return fileDescriptor_351d41b0f01c9adf, []int{5}
}

func (m *AddressBook) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddressBook.Unmarshal(m, b)
}
func (m *AddressBook) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddressBook.Marshal(b, m, deterministic)
}
func (m *AddressBook) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddressBook.Merge(m, src)
}
func (m *AddressBook) XXX_Size() int {
	return xxx_messageInfo_AddressBook.Size(m)
}
func (m *AddressBook) XXX_DiscardUnknown() {
	xxx_messageInfo_AddressBook.DiscardUnknown(m)
}

var xxx_messageInfo_AddressBook proto.InternalMessageInfo

func (m *AddressBook) GetPeople() []*PersonalInfo {
	if m != nil {
		return m.People
	}
	return nil
}

func init() {
	proto.RegisterType((*GetAddressByNameReq)(nil), "address_book.GetAddressByNameReq")
	proto.RegisterType((*GetAddressByNameRep)(nil), "address_book.GetAddressByNameRep")
	proto.RegisterType((*RegisterRep)(nil), "address_book.RegisterRep")
	proto.RegisterType((*RegisterReq)(nil), "address_book.RegisterReq")
	proto.RegisterType((*PersonalInfo)(nil), "address_book.PersonalInfo")
	proto.RegisterType((*AddressBook)(nil), "address_book.AddressBook")
}

func init() {
	proto.RegisterFile("address_book.proto", fileDescriptor_351d41b0f01c9adf)
}

var fileDescriptor_351d41b0f01c9adf = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x6d, 0xad, 0xd4, 0x3a, 0x89, 0x20, 0x23, 0x48, 0xcc, 0xc9, 0x2e, 0x08, 0x9e, 0x7a, 0x88,
	0x1f, 0x20, 0xe9, 0x45, 0xbc, 0x04, 0x89, 0x37, 0x2f, 0x21, 0xb1, 0x53, 0x0d, 0x4d, 0x76, 0xb6,
	0x9b, 0x54, 0xf0, 0xbf, 0xfc, 0x40, 0xc9, 0x36, 0xc1, 0x8d, 0x44, 0x7a, 0xf0, 0x96, 0x79, 0xf3,
	0xf2, 0xde, 0xe3, 0xed, 0x00, 0xa6, 0xab, 0x95, 0xa6, 0xaa, 0x4a, 0x32, 0xe6, 0xcd, 0x42, 0x69,
	0xae, 0x19, 0x5d, 0x1b, 0x13, 0x21, 0x5c, 0x3c, 0x50, 0x1d, 0xee, 0xa1, 0xe5, 0x67, 0x94, 0x96,
	0x14, 0xd3, 0x16, 0x11, 0x8e, 0x65, 0x5a, 0x92, 0x37, 0xbe, 0x1e, 0xdf, 0x9e, 0xc6, 0xe6, 0x1b,
	0x2f, 0x61, 0x2a, 0x77, 0x65, 0x46, 0xda, 0x3b, 0x32, 0x68, 0x3b, 0x09, 0x39, 0x24, 0xa1, 0xf0,
	0x1e, 0xce, 0x14, 0xe9, 0x8a, 0x65, 0x5a, 0x24, 0xb9, 0x5c, 0xb3, 0xd1, 0x72, 0x02, 0x7f, 0xd1,
	0xcb, 0xf4, 0xd4, 0x52, 0x1e, 0xe5, 0x9a, 0x63, 0x57, 0x59, 0x53, 0xe3, 0xa7, 0xa9, 0xda, 0x15,
	0xb5, 0xf1, 0x9b, 0xc5, 0xed, 0x24, 0x6e, 0xc0, 0x89, 0xe9, 0x2d, 0xaf, 0x6a, 0xd2, 0x8d, 0xcf,
	0x0f, 0x6d, 0x1f, 0xb6, 0xa3, 0x45, 0x36, 0x6d, 0xfb, 0xef, 0x38, 0x22, 0x01, 0xd7, 0xde, 0x0e,
	0x56, 0x34, 0x07, 0x57, 0xbd, 0xb3, 0xa4, 0xa4, 0x57, 0x94, 0x63, 0xb0, 0xc8, 0x40, 0xe8, 0xc1,
	0x49, 0xeb, 0xe8, 0x4d, 0xcc, 0xb6, 0x1b, 0x45, 0x08, 0x4e, 0x57, 0x22, 0xf3, 0x06, 0x03, 0x98,
	0x2a, 0x62, 0x55, 0x34, 0x0e, 0x93, 0x03, 0x49, 0x5b, 0x66, 0xf0, 0x35, 0x06, 0xb4, 0x34, 0x9e,
	0x49, 0x7f, 0xe4, 0xaf, 0x84, 0x2f, 0x70, 0xfe, 0xfb, 0x85, 0x70, 0xde, 0x97, 0x1b, 0x38, 0x02,
	0xff, 0x20, 0x45, 0x89, 0x11, 0x2e, 0x61, 0xd6, 0xd5, 0x8c, 0x57, 0xfd, 0x1f, 0xac, 0xfa, 0xfd,
	0x3f, 0x57, 0x4a, 0x8c, 0xb2, 0xa9, 0xb9, 0xcc, 0xbb, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xce,
	0xc6, 0xd0, 0x10, 0xaf, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AddressBookServiceClient is the client API for AddressBookService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AddressBookServiceClient interface {
	GetAddressByName(ctx context.Context, in *GetAddressByNameReq, opts ...grpc.CallOption) (*GetAddressByNameRep, error)
	Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRep, error)
}

type addressBookServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAddressBookServiceClient(cc grpc.ClientConnInterface) AddressBookServiceClient {
	return &addressBookServiceClient{cc}
}

func (c *addressBookServiceClient) GetAddressByName(ctx context.Context, in *GetAddressByNameReq, opts ...grpc.CallOption) (*GetAddressByNameRep, error) {
	out := new(GetAddressByNameRep)
	err := c.cc.Invoke(ctx, "/address_book.AddressBookService/GetAddressByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressBookServiceClient) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRep, error) {
	out := new(RegisterRep)
	err := c.cc.Invoke(ctx, "/address_book.AddressBookService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AddressBookServiceServer is the server API for AddressBookService service.
type AddressBookServiceServer interface {
	GetAddressByName(context.Context, *GetAddressByNameReq) (*GetAddressByNameRep, error)
	Register(context.Context, *RegisterReq) (*RegisterRep, error)
}

// UnimplementedAddressBookServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAddressBookServiceServer struct {
}

func (*UnimplementedAddressBookServiceServer) GetAddressByName(ctx context.Context, req *GetAddressByNameReq) (*GetAddressByNameRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddressByName not implemented")
}
func (*UnimplementedAddressBookServiceServer) Register(ctx context.Context, req *RegisterReq) (*RegisterRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func RegisterAddressBookServiceServer(s *grpc.Server, srv AddressBookServiceServer) {
	s.RegisterService(&_AddressBookService_serviceDesc, srv)
}

func _AddressBookService_GetAddressByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAddressByNameReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddressBookServiceServer).GetAddressByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/address_book.AddressBookService/GetAddressByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddressBookServiceServer).GetAddressByName(ctx, req.(*GetAddressByNameReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddressBookService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddressBookServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/address_book.AddressBookService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddressBookServiceServer).Register(ctx, req.(*RegisterReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _AddressBookService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "address_book.AddressBookService",
	HandlerType: (*AddressBookServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAddressByName",
			Handler:    _AddressBookService_GetAddressByName_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _AddressBookService_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "address_book.proto",
}

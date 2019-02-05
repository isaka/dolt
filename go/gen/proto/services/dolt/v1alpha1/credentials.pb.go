// Code generated by protoc-gen-go. DO NOT EDIT.
// source: services/dolt/v1alpha1/credentials.proto

package v1alpha1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type WhoAmIRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WhoAmIRequest) Reset()         { *m = WhoAmIRequest{} }
func (m *WhoAmIRequest) String() string { return proto.CompactTextString(m) }
func (*WhoAmIRequest) ProtoMessage()    {}
func (*WhoAmIRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_credentials_e7445091f28022eb, []int{0}
}
func (m *WhoAmIRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WhoAmIRequest.Unmarshal(m, b)
}
func (m *WhoAmIRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WhoAmIRequest.Marshal(b, m, deterministic)
}
func (dst *WhoAmIRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WhoAmIRequest.Merge(dst, src)
}
func (m *WhoAmIRequest) XXX_Size() int {
	return xxx_messageInfo_WhoAmIRequest.Size(m)
}
func (m *WhoAmIRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WhoAmIRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WhoAmIRequest proto.InternalMessageInfo

type WhoAmIResponse struct {
	// Ex: "bheni"
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// Ex: "Brian Hendriks"
	DisplayName string `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// Ex: "brian@liquidata.co"
	EmailAddress         string   `protobuf:"bytes,3,opt,name=email_address,json=emailAddress,proto3" json:"email_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WhoAmIResponse) Reset()         { *m = WhoAmIResponse{} }
func (m *WhoAmIResponse) String() string { return proto.CompactTextString(m) }
func (*WhoAmIResponse) ProtoMessage()    {}
func (*WhoAmIResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_credentials_e7445091f28022eb, []int{1}
}
func (m *WhoAmIResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WhoAmIResponse.Unmarshal(m, b)
}
func (m *WhoAmIResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WhoAmIResponse.Marshal(b, m, deterministic)
}
func (dst *WhoAmIResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WhoAmIResponse.Merge(dst, src)
}
func (m *WhoAmIResponse) XXX_Size() int {
	return xxx_messageInfo_WhoAmIResponse.Size(m)
}
func (m *WhoAmIResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WhoAmIResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WhoAmIResponse proto.InternalMessageInfo

func (m *WhoAmIResponse) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *WhoAmIResponse) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *WhoAmIResponse) GetEmailAddress() string {
	if m != nil {
		return m.EmailAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*WhoAmIRequest)(nil), "services.dolt.v1alpha1.WhoAmIRequest")
	proto.RegisterType((*WhoAmIResponse)(nil), "services.dolt.v1alpha1.WhoAmIResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CredentialsServiceClient is the client API for CredentialsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CredentialsServiceClient interface {
	WhoAmI(ctx context.Context, in *WhoAmIRequest, opts ...grpc.CallOption) (*WhoAmIResponse, error)
}

type credentialsServiceClient struct {
	cc *grpc.ClientConn
}

func NewCredentialsServiceClient(cc *grpc.ClientConn) CredentialsServiceClient {
	return &credentialsServiceClient{cc}
}

func (c *credentialsServiceClient) WhoAmI(ctx context.Context, in *WhoAmIRequest, opts ...grpc.CallOption) (*WhoAmIResponse, error) {
	out := new(WhoAmIResponse)
	err := c.cc.Invoke(ctx, "/services.dolt.v1alpha1.CredentialsService/WhoAmI", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CredentialsServiceServer is the server API for CredentialsService service.
type CredentialsServiceServer interface {
	WhoAmI(context.Context, *WhoAmIRequest) (*WhoAmIResponse, error)
}

func RegisterCredentialsServiceServer(s *grpc.Server, srv CredentialsServiceServer) {
	s.RegisterService(&_CredentialsService_serviceDesc, srv)
}

func _CredentialsService_WhoAmI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WhoAmIRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CredentialsServiceServer).WhoAmI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.dolt.v1alpha1.CredentialsService/WhoAmI",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CredentialsServiceServer).WhoAmI(ctx, req.(*WhoAmIRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CredentialsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "services.dolt.v1alpha1.CredentialsService",
	HandlerType: (*CredentialsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WhoAmI",
			Handler:    _CredentialsService_WhoAmI_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/dolt/v1alpha1/credentials.proto",
}

func init() {
	proto.RegisterFile("services/dolt/v1alpha1/credentials.proto", fileDescriptor_credentials_e7445091f28022eb)
}

var fileDescriptor_credentials_e7445091f28022eb = []byte{
	// 217 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x28, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0x2d, 0xd6, 0x4f, 0xc9, 0xcf, 0x29, 0xd1, 0x2f, 0x33, 0x4c, 0xcc, 0x29, 0xc8,
	0x48, 0x34, 0xd4, 0x4f, 0x2e, 0x4a, 0x4d, 0x49, 0xcd, 0x2b, 0xc9, 0x4c, 0xcc, 0x29, 0xd6, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x83, 0xa9, 0xd4, 0x03, 0xa9, 0xd4, 0x83, 0xa9, 0x54, 0xe2,
	0xe7, 0xe2, 0x0d, 0xcf, 0xc8, 0x77, 0xcc, 0xf5, 0x0c, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x51,
	0x2a, 0xe1, 0xe2, 0x83, 0x09, 0x14, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0x49, 0x71, 0x71, 0x94,
	0x16, 0xa7, 0x16, 0xe5, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xf9,
	0x42, 0x8a, 0x5c, 0x3c, 0x29, 0x99, 0xc5, 0x05, 0x39, 0x89, 0x95, 0xf1, 0x60, 0x79, 0x26, 0xb0,
	0x3c, 0x37, 0x54, 0xcc, 0x0f, 0xa4, 0x44, 0x99, 0x8b, 0x37, 0x35, 0x37, 0x31, 0x33, 0x27, 0x3e,
	0x31, 0x25, 0xa5, 0x28, 0xb5, 0xb8, 0x58, 0x82, 0x19, 0xac, 0x86, 0x07, 0x2c, 0xe8, 0x08, 0x11,
	0x33, 0xca, 0xe5, 0x12, 0x72, 0x46, 0xb8, 0x39, 0x18, 0xe2, 0x56, 0xa1, 0x70, 0x2e, 0x36, 0x88,
	0x5b, 0x84, 0x54, 0xf5, 0xb0, 0xbb, 0x5f, 0x0f, 0xc5, 0xf1, 0x52, 0x6a, 0x84, 0x94, 0x41, 0xbc,
	0xe4, 0xc4, 0x15, 0xc5, 0x01, 0x93, 0x4a, 0x62, 0x03, 0x07, 0x90, 0x31, 0x20, 0x00, 0x00, 0xff,
	0xff, 0x6f, 0x23, 0x23, 0xd0, 0x4c, 0x01, 0x00, 0x00,
}

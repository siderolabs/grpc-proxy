// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

package talos_testproto

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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type PingRequest struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{1}
}

func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (m *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(m, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

func (m *PingRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type PingResponse struct {
	Value                string   `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	Counter              int32    `protobuf:"varint,2,opt,name=counter,proto3" json:"counter,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingResponse) Reset()         { *m = PingResponse{} }
func (m *PingResponse) String() string { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()    {}
func (*PingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{2}
}

func (m *PingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingResponse.Unmarshal(m, b)
}
func (m *PingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingResponse.Marshal(b, m, deterministic)
}
func (m *PingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingResponse.Merge(m, src)
}
func (m *PingResponse) XXX_Size() int {
	return xxx_messageInfo_PingResponse.Size(m)
}
func (m *PingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PingResponse proto.InternalMessageInfo

func (m *PingResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *PingResponse) GetCounter() int32 {
	if m != nil {
		return m.Counter
	}
	return 0
}

type ResponseMetadata struct {
	Hostname             string   `protobuf:"bytes,99,opt,name=hostname,proto3" json:"hostname,omitempty"`
	UpstreamError        string   `protobuf:"bytes,100,opt,name=upstream_error,json=upstreamError,proto3" json:"upstream_error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseMetadata) Reset()         { *m = ResponseMetadata{} }
func (m *ResponseMetadata) String() string { return proto.CompactTextString(m) }
func (*ResponseMetadata) ProtoMessage()    {}
func (*ResponseMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{3}
}

func (m *ResponseMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseMetadata.Unmarshal(m, b)
}
func (m *ResponseMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseMetadata.Marshal(b, m, deterministic)
}
func (m *ResponseMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseMetadata.Merge(m, src)
}
func (m *ResponseMetadata) XXX_Size() int {
	return xxx_messageInfo_ResponseMetadata.Size(m)
}
func (m *ResponseMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseMetadata proto.InternalMessageInfo

func (m *ResponseMetadata) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *ResponseMetadata) GetUpstreamError() string {
	if m != nil {
		return m.UpstreamError
	}
	return ""
}

type ResponseMetadataPrepender struct {
	Metadata             *ResponseMetadata `protobuf:"bytes,99,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ResponseMetadataPrepender) Reset()         { *m = ResponseMetadataPrepender{} }
func (m *ResponseMetadataPrepender) String() string { return proto.CompactTextString(m) }
func (*ResponseMetadataPrepender) ProtoMessage()    {}
func (*ResponseMetadataPrepender) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{4}
}

func (m *ResponseMetadataPrepender) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseMetadataPrepender.Unmarshal(m, b)
}
func (m *ResponseMetadataPrepender) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseMetadataPrepender.Marshal(b, m, deterministic)
}
func (m *ResponseMetadataPrepender) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseMetadataPrepender.Merge(m, src)
}
func (m *ResponseMetadataPrepender) XXX_Size() int {
	return xxx_messageInfo_ResponseMetadataPrepender.Size(m)
}
func (m *ResponseMetadataPrepender) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseMetadataPrepender.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseMetadataPrepender proto.InternalMessageInfo

func (m *ResponseMetadataPrepender) GetMetadata() *ResponseMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type MultiPingResponse struct {
	Metadata             *ResponseMetadata `protobuf:"bytes,99,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Value                string            `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	Counter              int32             `protobuf:"varint,2,opt,name=counter,proto3" json:"counter,omitempty"`
	Server               string            `protobuf:"bytes,3,opt,name=server,proto3" json:"server,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MultiPingResponse) Reset()         { *m = MultiPingResponse{} }
func (m *MultiPingResponse) String() string { return proto.CompactTextString(m) }
func (*MultiPingResponse) ProtoMessage()    {}
func (*MultiPingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{5}
}

func (m *MultiPingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MultiPingResponse.Unmarshal(m, b)
}
func (m *MultiPingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MultiPingResponse.Marshal(b, m, deterministic)
}
func (m *MultiPingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiPingResponse.Merge(m, src)
}
func (m *MultiPingResponse) XXX_Size() int {
	return xxx_messageInfo_MultiPingResponse.Size(m)
}
func (m *MultiPingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiPingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MultiPingResponse proto.InternalMessageInfo

func (m *MultiPingResponse) GetMetadata() *ResponseMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *MultiPingResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *MultiPingResponse) GetCounter() int32 {
	if m != nil {
		return m.Counter
	}
	return 0
}

func (m *MultiPingResponse) GetServer() string {
	if m != nil {
		return m.Server
	}
	return ""
}

type MultiPingReply struct {
	Response             []*MultiPingResponse `protobuf:"bytes,1,rep,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *MultiPingReply) Reset()         { *m = MultiPingReply{} }
func (m *MultiPingReply) String() string { return proto.CompactTextString(m) }
func (*MultiPingReply) ProtoMessage()    {}
func (*MultiPingReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{6}
}

func (m *MultiPingReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MultiPingReply.Unmarshal(m, b)
}
func (m *MultiPingReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MultiPingReply.Marshal(b, m, deterministic)
}
func (m *MultiPingReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiPingReply.Merge(m, src)
}
func (m *MultiPingReply) XXX_Size() int {
	return xxx_messageInfo_MultiPingReply.Size(m)
}
func (m *MultiPingReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiPingReply.DiscardUnknown(m)
}

var xxx_messageInfo_MultiPingReply proto.InternalMessageInfo

func (m *MultiPingReply) GetResponse() []*MultiPingResponse {
	if m != nil {
		return m.Response
	}
	return nil
}

type EmptyReply struct {
	Response             []*EmptyResponse `protobuf:"bytes,1,rep,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *EmptyReply) Reset()         { *m = EmptyReply{} }
func (m *EmptyReply) String() string { return proto.CompactTextString(m) }
func (*EmptyReply) ProtoMessage()    {}
func (*EmptyReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{7}
}

func (m *EmptyReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyReply.Unmarshal(m, b)
}
func (m *EmptyReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyReply.Marshal(b, m, deterministic)
}
func (m *EmptyReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyReply.Merge(m, src)
}
func (m *EmptyReply) XXX_Size() int {
	return xxx_messageInfo_EmptyReply.Size(m)
}
func (m *EmptyReply) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyReply.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyReply proto.InternalMessageInfo

func (m *EmptyReply) GetResponse() []*EmptyResponse {
	if m != nil {
		return m.Response
	}
	return nil
}

type EmptyResponse struct {
	Metadata             *ResponseMetadata `protobuf:"bytes,99,opt,name=metadata,proto3" json:"metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *EmptyResponse) Reset()         { *m = EmptyResponse{} }
func (m *EmptyResponse) String() string { return proto.CompactTextString(m) }
func (*EmptyResponse) ProtoMessage()    {}
func (*EmptyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{8}
}

func (m *EmptyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyResponse.Unmarshal(m, b)
}
func (m *EmptyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyResponse.Marshal(b, m, deterministic)
}
func (m *EmptyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyResponse.Merge(m, src)
}
func (m *EmptyResponse) XXX_Size() int {
	return xxx_messageInfo_EmptyResponse.Size(m)
}
func (m *EmptyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyResponse proto.InternalMessageInfo

func (m *EmptyResponse) GetMetadata() *ResponseMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "talos.testproto.Empty")
	proto.RegisterType((*PingRequest)(nil), "talos.testproto.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "talos.testproto.PingResponse")
	proto.RegisterType((*ResponseMetadata)(nil), "talos.testproto.ResponseMetadata")
	proto.RegisterType((*ResponseMetadataPrepender)(nil), "talos.testproto.ResponseMetadataPrepender")
	proto.RegisterType((*MultiPingResponse)(nil), "talos.testproto.MultiPingResponse")
	proto.RegisterType((*MultiPingReply)(nil), "talos.testproto.MultiPingReply")
	proto.RegisterType((*EmptyReply)(nil), "talos.testproto.EmptyReply")
	proto.RegisterType((*EmptyResponse)(nil), "talos.testproto.EmptyResponse")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x93, 0x5d, 0xcf, 0xd2, 0x30,
	0x14, 0xc7, 0x1d, 0x08, 0x8c, 0xc3, 0x9b, 0x36, 0x86, 0x4c, 0x7c, 0xaf, 0x31, 0xe1, 0x6a, 0x21,
	0x78, 0x67, 0x22, 0x37, 0x8a, 0x92, 0x28, 0xba, 0x0c, 0x34, 0xd1, 0x1b, 0x33, 0xa1, 0xd1, 0x25,
	0x7b, 0xb3, 0xeb, 0x48, 0xf8, 0x2a, 0x7e, 0x46, 0xbf, 0x82, 0x89, 0x6d, 0xb7, 0x21, 0xb0, 0x87,
	0x87, 0x92, 0x5d, 0x9e, 0x7f, 0xcf, 0xf9, 0xf5, 0x9c, 0xfe, 0x4f, 0x01, 0x18, 0x89, 0x99, 0x19,
	0xd1, 0x90, 0x85, 0xa8, 0xc7, 0x1c, 0x2f, 0x8c, 0x4d, 0xa1, 0x48, 0x01, 0x37, 0xa0, 0x36, 0xf5,
	0x23, 0xb6, 0xc5, 0x4f, 0xa1, 0x65, 0xb9, 0xc1, 0x0f, 0x9b, 0xfc, 0x4a, 0xf8, 0x21, 0xba, 0x03,
	0xb5, 0x8d, 0xe3, 0x25, 0xc4, 0xd0, 0x1e, 0x6b, 0xc3, 0xa6, 0x9d, 0x06, 0x78, 0x02, 0xed, 0x34,
	0x29, 0x8e, 0xc2, 0x20, 0x26, 0x22, 0xeb, 0xf3, 0x7e, 0x96, 0x0c, 0x90, 0x01, 0x8d, 0x55, 0x98,
	0x04, 0x8c, 0x50, 0xa3, 0xc2, 0xf5, 0x9a, 0x9d, 0x87, 0xf8, 0x13, 0xdc, 0xca, 0x6b, 0xe7, 0x84,
	0x39, 0x6b, 0x87, 0x39, 0x68, 0x00, 0xfa, 0xcf, 0x30, 0x66, 0x81, 0xe3, 0x13, 0x63, 0x25, 0x31,
	0xbb, 0x18, 0x3d, 0x83, 0x6e, 0x12, 0xc5, 0x8c, 0x12, 0xc7, 0xff, 0x46, 0x28, 0x0d, 0xa9, 0xb1,
	0x96, 0x19, 0x9d, 0x5c, 0x9d, 0x0a, 0x11, 0x7f, 0x85, 0xbb, 0xc7, 0x58, 0x8b, 0x92, 0x88, 0x04,
	0x6b, 0x42, 0xd1, 0x4b, 0xd0, 0xfd, 0x4c, 0x94, 0xfc, 0xd6, 0xf8, 0x89, 0x79, 0xf4, 0x0a, 0xe6,
	0x71, 0xb5, 0xbd, 0x2b, 0xc1, 0xbf, 0x35, 0xb8, 0x3d, 0x4f, 0x3c, 0xe6, 0x1e, 0x0c, 0x5e, 0x0e,
	0x7a, 0xe9, 0xbb, 0xa1, 0x3e, 0xd4, 0x63, 0x42, 0x37, 0xfc, 0xa0, 0x2a, 0x0b, 0xb2, 0x08, 0x5b,
	0xd0, 0xdd, 0xeb, 0x2d, 0xf2, 0xb6, 0x68, 0x02, 0x3a, 0xcd, 0xee, 0xe5, 0xf0, 0x2a, 0x6f, 0x0c,
	0x17, 0x1a, 0x2b, 0x8c, 0x63, 0xef, 0x6a, 0xf0, 0x0c, 0x40, 0xee, 0x43, 0x4a, 0x7b, 0x51, 0xa0,
	0x3d, 0x2c, 0xd0, 0xb2, 0xf4, 0x02, 0xe9, 0x03, 0x74, 0x0e, 0x8e, 0x4a, 0xbe, 0xd9, 0xf8, 0x6f,
	0x05, 0x5a, 0x4b, 0x9e, 0xb8, 0xe0, 0xa3, 0xbb, 0x2b, 0x82, 0x5e, 0x43, 0x53, 0xcc, 0x20, 0xef,
	0x40, 0xfd, 0xab, 0xdb, 0x1a, 0x3c, 0x28, 0xe8, 0xfb, 0x73, 0xe3, 0x1b, 0x68, 0x0a, 0x37, 0x85,
	0x82, 0xee, 0x9f, 0x48, 0x94, 0xbf, 0xe1, 0x3c, 0xe6, 0x55, 0xd6, 0x8c, 0x58, 0xc7, 0x33, 0xac,
	0x13, 0xad, 0x72, 0xc8, 0x3b, 0xd0, 0x45, 0xe2, 0x7b, 0x97, 0xff, 0xbf, 0x72, 0xfd, 0x8c, 0x34,
	0xf4, 0x11, 0x40, 0x68, 0x0b, 0xf9, 0x4d, 0x4a, 0xe2, 0x86, 0xda, 0x48, 0x1b, 0xff, 0xa9, 0x42,
	0x5b, 0x6e, 0x4e, 0x6e, 0xc0, 0x1b, 0x15, 0x03, 0x1e, 0x5d, 0xb7, 0x7d, 0x7c, 0xc5, 0xf8, 0xd8,
	0x6f, 0x95, 0x2c, 0x50, 0x00, 0xcd, 0xd4, 0x4d, 0xb8, 0x77, 0x6a, 0x8d, 0x53, 0x92, 0xa5, 0xec,
	0x84, 0xc2, 0xef, 0x92, 0x76, 0x2c, 0x2f, 0xb0, 0x43, 0x89, 0x29, 0x3c, 0x41, 0x5f, 0xa0, 0xf7,
	0x9f, 0xaa, 0x32, 0xb7, 0x32, 0xfa, 0x7b, 0x5d, 0x1e, 0x3f, 0xff, 0x17, 0x00, 0x00, 0xff, 0xff,
	0x8f, 0xae, 0x5f, 0xeb, 0x3e, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TestServiceClient is the client API for TestService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TestServiceClient interface {
	PingEmpty(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PingResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	PingError(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*Empty, error)
	PingList(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (TestService_PingListClient, error)
	PingStream(ctx context.Context, opts ...grpc.CallOption) (TestService_PingStreamClient, error)
}

type testServiceClient struct {
	cc *grpc.ClientConn
}

func NewTestServiceClient(cc *grpc.ClientConn) TestServiceClient {
	return &testServiceClient{cc}
}

func (c *testServiceClient) PingEmpty(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/talos.testproto.TestService/PingEmpty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/talos.testproto.TestService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) PingError(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/talos.testproto.TestService/PingError", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) PingList(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (TestService_PingListClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TestService_serviceDesc.Streams[0], "/talos.testproto.TestService/PingList", opts...)
	if err != nil {
		return nil, err
	}
	x := &testServicePingListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TestService_PingListClient interface {
	Recv() (*PingResponse, error)
	grpc.ClientStream
}

type testServicePingListClient struct {
	grpc.ClientStream
}

func (x *testServicePingListClient) Recv() (*PingResponse, error) {
	m := new(PingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testServiceClient) PingStream(ctx context.Context, opts ...grpc.CallOption) (TestService_PingStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TestService_serviceDesc.Streams[1], "/talos.testproto.TestService/PingStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &testServicePingStreamClient{stream}
	return x, nil
}

type TestService_PingStreamClient interface {
	Send(*PingRequest) error
	Recv() (*PingResponse, error)
	grpc.ClientStream
}

type testServicePingStreamClient struct {
	grpc.ClientStream
}

func (x *testServicePingStreamClient) Send(m *PingRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testServicePingStreamClient) Recv() (*PingResponse, error) {
	m := new(PingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestServiceServer is the server API for TestService service.
type TestServiceServer interface {
	PingEmpty(context.Context, *Empty) (*PingResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	PingError(context.Context, *PingRequest) (*Empty, error)
	PingList(*PingRequest, TestService_PingListServer) error
	PingStream(TestService_PingStreamServer) error
}

// UnimplementedTestServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTestServiceServer struct {
}

func (*UnimplementedTestServiceServer) PingEmpty(ctx context.Context, req *Empty) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingEmpty not implemented")
}
func (*UnimplementedTestServiceServer) Ping(ctx context.Context, req *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedTestServiceServer) PingError(ctx context.Context, req *PingRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingError not implemented")
}
func (*UnimplementedTestServiceServer) PingList(req *PingRequest, srv TestService_PingListServer) error {
	return status.Errorf(codes.Unimplemented, "method PingList not implemented")
}
func (*UnimplementedTestServiceServer) PingStream(srv TestService_PingStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method PingStream not implemented")
}

func RegisterTestServiceServer(s *grpc.Server, srv TestServiceServer) {
	s.RegisterService(&_TestService_serviceDesc, srv)
}

func _TestService_PingEmpty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).PingEmpty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talos.testproto.TestService/PingEmpty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).PingEmpty(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talos.testproto.TestService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_PingError_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).PingError(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talos.testproto.TestService/PingError",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).PingError(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_PingList_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PingRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestServiceServer).PingList(m, &testServicePingListServer{stream})
}

type TestService_PingListServer interface {
	Send(*PingResponse) error
	grpc.ServerStream
}

type testServicePingListServer struct {
	grpc.ServerStream
}

func (x *testServicePingListServer) Send(m *PingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _TestService_PingStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestServiceServer).PingStream(&testServicePingStreamServer{stream})
}

type TestService_PingStreamServer interface {
	Send(*PingResponse) error
	Recv() (*PingRequest, error)
	grpc.ServerStream
}

type testServicePingStreamServer struct {
	grpc.ServerStream
}

func (x *testServicePingStreamServer) Send(m *PingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testServicePingStreamServer) Recv() (*PingRequest, error) {
	m := new(PingRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _TestService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "talos.testproto.TestService",
	HandlerType: (*TestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PingEmpty",
			Handler:    _TestService_PingEmpty_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _TestService_Ping_Handler,
		},
		{
			MethodName: "PingError",
			Handler:    _TestService_PingError_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PingList",
			Handler:       _TestService_PingList_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PingStream",
			Handler:       _TestService_PingStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "test.proto",
}

// MultiServiceClient is the client API for MultiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MultiServiceClient interface {
	PingEmpty(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MultiPingReply, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*MultiPingReply, error)
	PingError(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*EmptyReply, error)
	PingList(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (MultiService_PingListClient, error)
	PingStream(ctx context.Context, opts ...grpc.CallOption) (MultiService_PingStreamClient, error)
	PingStreamError(ctx context.Context, opts ...grpc.CallOption) (MultiService_PingStreamErrorClient, error)
}

type multiServiceClient struct {
	cc *grpc.ClientConn
}

func NewMultiServiceClient(cc *grpc.ClientConn) MultiServiceClient {
	return &multiServiceClient{cc}
}

func (c *multiServiceClient) PingEmpty(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*MultiPingReply, error) {
	out := new(MultiPingReply)
	err := c.cc.Invoke(ctx, "/talos.testproto.MultiService/PingEmpty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*MultiPingReply, error) {
	out := new(MultiPingReply)
	err := c.cc.Invoke(ctx, "/talos.testproto.MultiService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiServiceClient) PingError(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*EmptyReply, error) {
	out := new(EmptyReply)
	err := c.cc.Invoke(ctx, "/talos.testproto.MultiService/PingError", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiServiceClient) PingList(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (MultiService_PingListClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MultiService_serviceDesc.Streams[0], "/talos.testproto.MultiService/PingList", opts...)
	if err != nil {
		return nil, err
	}
	x := &multiServicePingListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MultiService_PingListClient interface {
	Recv() (*MultiPingResponse, error)
	grpc.ClientStream
}

type multiServicePingListClient struct {
	grpc.ClientStream
}

func (x *multiServicePingListClient) Recv() (*MultiPingResponse, error) {
	m := new(MultiPingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *multiServiceClient) PingStream(ctx context.Context, opts ...grpc.CallOption) (MultiService_PingStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MultiService_serviceDesc.Streams[1], "/talos.testproto.MultiService/PingStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &multiServicePingStreamClient{stream}
	return x, nil
}

type MultiService_PingStreamClient interface {
	Send(*PingRequest) error
	Recv() (*MultiPingResponse, error)
	grpc.ClientStream
}

type multiServicePingStreamClient struct {
	grpc.ClientStream
}

func (x *multiServicePingStreamClient) Send(m *PingRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *multiServicePingStreamClient) Recv() (*MultiPingResponse, error) {
	m := new(MultiPingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *multiServiceClient) PingStreamError(ctx context.Context, opts ...grpc.CallOption) (MultiService_PingStreamErrorClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MultiService_serviceDesc.Streams[2], "/talos.testproto.MultiService/PingStreamError", opts...)
	if err != nil {
		return nil, err
	}
	x := &multiServicePingStreamErrorClient{stream}
	return x, nil
}

type MultiService_PingStreamErrorClient interface {
	Send(*PingRequest) error
	Recv() (*MultiPingResponse, error)
	grpc.ClientStream
}

type multiServicePingStreamErrorClient struct {
	grpc.ClientStream
}

func (x *multiServicePingStreamErrorClient) Send(m *PingRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *multiServicePingStreamErrorClient) Recv() (*MultiPingResponse, error) {
	m := new(MultiPingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MultiServiceServer is the server API for MultiService service.
type MultiServiceServer interface {
	PingEmpty(context.Context, *Empty) (*MultiPingReply, error)
	Ping(context.Context, *PingRequest) (*MultiPingReply, error)
	PingError(context.Context, *PingRequest) (*EmptyReply, error)
	PingList(*PingRequest, MultiService_PingListServer) error
	PingStream(MultiService_PingStreamServer) error
	PingStreamError(MultiService_PingStreamErrorServer) error
}

// UnimplementedMultiServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMultiServiceServer struct {
}

func (*UnimplementedMultiServiceServer) PingEmpty(ctx context.Context, req *Empty) (*MultiPingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingEmpty not implemented")
}
func (*UnimplementedMultiServiceServer) Ping(ctx context.Context, req *PingRequest) (*MultiPingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedMultiServiceServer) PingError(ctx context.Context, req *PingRequest) (*EmptyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingError not implemented")
}
func (*UnimplementedMultiServiceServer) PingList(req *PingRequest, srv MultiService_PingListServer) error {
	return status.Errorf(codes.Unimplemented, "method PingList not implemented")
}
func (*UnimplementedMultiServiceServer) PingStream(srv MultiService_PingStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method PingStream not implemented")
}
func (*UnimplementedMultiServiceServer) PingStreamError(srv MultiService_PingStreamErrorServer) error {
	return status.Errorf(codes.Unimplemented, "method PingStreamError not implemented")
}

func RegisterMultiServiceServer(s *grpc.Server, srv MultiServiceServer) {
	s.RegisterService(&_MultiService_serviceDesc, srv)
}

func _MultiService_PingEmpty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiServiceServer).PingEmpty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talos.testproto.MultiService/PingEmpty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiServiceServer).PingEmpty(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talos.testproto.MultiService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiService_PingError_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiServiceServer).PingError(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talos.testproto.MultiService/PingError",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiServiceServer).PingError(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiService_PingList_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PingRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MultiServiceServer).PingList(m, &multiServicePingListServer{stream})
}

type MultiService_PingListServer interface {
	Send(*MultiPingResponse) error
	grpc.ServerStream
}

type multiServicePingListServer struct {
	grpc.ServerStream
}

func (x *multiServicePingListServer) Send(m *MultiPingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MultiService_PingStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MultiServiceServer).PingStream(&multiServicePingStreamServer{stream})
}

type MultiService_PingStreamServer interface {
	Send(*MultiPingResponse) error
	Recv() (*PingRequest, error)
	grpc.ServerStream
}

type multiServicePingStreamServer struct {
	grpc.ServerStream
}

func (x *multiServicePingStreamServer) Send(m *MultiPingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *multiServicePingStreamServer) Recv() (*PingRequest, error) {
	m := new(PingRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _MultiService_PingStreamError_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MultiServiceServer).PingStreamError(&multiServicePingStreamErrorServer{stream})
}

type MultiService_PingStreamErrorServer interface {
	Send(*MultiPingResponse) error
	Recv() (*PingRequest, error)
	grpc.ServerStream
}

type multiServicePingStreamErrorServer struct {
	grpc.ServerStream
}

func (x *multiServicePingStreamErrorServer) Send(m *MultiPingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *multiServicePingStreamErrorServer) Recv() (*PingRequest, error) {
	m := new(PingRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _MultiService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "talos.testproto.MultiService",
	HandlerType: (*MultiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PingEmpty",
			Handler:    _MultiService_PingEmpty_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _MultiService_Ping_Handler,
		},
		{
			MethodName: "PingError",
			Handler:    _MultiService_PingError_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PingList",
			Handler:       _MultiService_PingList_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PingStream",
			Handler:       _MultiService_PingStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "PingStreamError",
			Handler:       _MultiService_PingStreamError_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "test.proto",
}

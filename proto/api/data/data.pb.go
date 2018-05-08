// Code generated by protoc-gen-go. DO NOT EDIT.
// source: data.proto

/*
Package data is a generated protocol buffer package.

It is generated from these files:
	data.proto

It has these top-level messages:
	DumpChunk
	Empty
*/
package data

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_api1 "google.golang.org/genproto/googleapis/api/annotations"

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

// http from public import google/api/annotations.proto
var E_Http = google_api1.E_Http

// A type that represents one chunk of a data's dump.
type DumpChunk struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *DumpChunk) Reset()                    { *m = DumpChunk{} }
func (m *DumpChunk) String() string            { return proto.CompactTextString(m) }
func (*DumpChunk) ProtoMessage()               {}
func (*DumpChunk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DumpChunk) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*DumpChunk)(nil), "spire.api.data.DumpChunk")
	proto.RegisterType((*Empty)(nil), "spire.api.data.Empty")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Data service

type DataClient interface {
	// Dump retrieves all the data from the datastore chunk by chunk.
	Dump(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Data_DumpClient, error)
	// Replays puts all the DumpChunks of a dump into datastore.
	Replay(ctx context.Context, opts ...grpc.CallOption) (Data_ReplayClient, error)
}

type dataClient struct {
	cc *grpc.ClientConn
}

func NewDataClient(cc *grpc.ClientConn) DataClient {
	return &dataClient{cc}
}

func (c *dataClient) Dump(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Data_DumpClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Data_serviceDesc.Streams[0], c.cc, "/spire.api.data.Data/Dump", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataDumpClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Data_DumpClient interface {
	Recv() (*DumpChunk, error)
	grpc.ClientStream
}

type dataDumpClient struct {
	grpc.ClientStream
}

func (x *dataDumpClient) Recv() (*DumpChunk, error) {
	m := new(DumpChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dataClient) Replay(ctx context.Context, opts ...grpc.CallOption) (Data_ReplayClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Data_serviceDesc.Streams[1], c.cc, "/spire.api.data.Data/Replay", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataReplayClient{stream}
	return x, nil
}

type Data_ReplayClient interface {
	Send(*DumpChunk) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type dataReplayClient struct {
	grpc.ClientStream
}

func (x *dataReplayClient) Send(m *DumpChunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dataReplayClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Data service

type DataServer interface {
	// Dump retrieves all the data from the datastore chunk by chunk.
	Dump(*Empty, Data_DumpServer) error
	// Replays puts all the DumpChunks of a dump into datastore.
	Replay(Data_ReplayServer) error
}

func RegisterDataServer(s *grpc.Server, srv DataServer) {
	s.RegisterService(&_Data_serviceDesc, srv)
}

func _Data_Dump_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DataServer).Dump(m, &dataDumpServer{stream})
}

type Data_DumpServer interface {
	Send(*DumpChunk) error
	grpc.ServerStream
}

type dataDumpServer struct {
	grpc.ServerStream
}

func (x *dataDumpServer) Send(m *DumpChunk) error {
	return x.ServerStream.SendMsg(m)
}

func _Data_Replay_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataServer).Replay(&dataReplayServer{stream})
}

type Data_ReplayServer interface {
	SendAndClose(*Empty) error
	Recv() (*DumpChunk, error)
	grpc.ServerStream
}

type dataReplayServer struct {
	grpc.ServerStream
}

func (x *dataReplayServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dataReplayServer) Recv() (*DumpChunk, error) {
	m := new(DumpChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Data_serviceDesc = grpc.ServiceDesc{
	ServiceName: "spire.api.data.Data",
	HandlerType: (*DataServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Dump",
			Handler:       _Data_Dump_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Replay",
			Handler:       _Data_Replay_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "data.proto",
}

func init() { proto.RegisterFile("data.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x49, 0x2c, 0x49,
	0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2b, 0x2e, 0xc8, 0x2c, 0x4a, 0xd5, 0x4b, 0x2c,
	0xc8, 0xd4, 0x03, 0x89, 0x4a, 0xc9, 0xa4, 0xe7, 0xe7, 0xa7, 0xe7, 0xa4, 0xea, 0x27, 0x16, 0x64,
	0xea, 0x27, 0xe6, 0xe5, 0xe5, 0x97, 0x24, 0x96, 0x64, 0xe6, 0xe7, 0x15, 0x43, 0x54, 0x2b, 0xc9,
	0x73, 0x71, 0xba, 0x94, 0xe6, 0x16, 0x38, 0x67, 0x94, 0xe6, 0x65, 0x0b, 0x09, 0x71, 0xb1, 0x80,
	0xb4, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x04, 0x81, 0xd9, 0x4a, 0xec, 0x5c, 0xac, 0xae, 0xb9,
	0x05, 0x25, 0x95, 0x46, 0x0d, 0x8c, 0x5c, 0x2c, 0x2e, 0x89, 0x25, 0x89, 0x42, 0x56, 0x5c, 0x2c,
	0x20, 0x2d, 0x42, 0xa2, 0x7a, 0xa8, 0x36, 0xe9, 0x81, 0xd5, 0x49, 0x49, 0xa2, 0x0b, 0xc3, 0xcd,
	0x37, 0x60, 0x14, 0xb2, 0xe1, 0x62, 0x0b, 0x4a, 0x2d, 0xc8, 0x49, 0xac, 0x14, 0xc2, 0xad, 0x4c,
	0x0a, 0xbb, 0xc1, 0x1a, 0x8c, 0x4e, 0x6c, 0x51, 0x60, 0x37, 0x05, 0x30, 0x24, 0xb1, 0x81, 0x5d,
	0x6f, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x14, 0x50, 0x40, 0xab, 0xf9, 0x00, 0x00, 0x00,
}
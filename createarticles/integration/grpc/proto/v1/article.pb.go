// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/v1/article.proto

package v1

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

type Article struct {
	Api                  string   `protobuf:"bytes,1,opt,name=api,proto3" json:"api,omitempty"`
	Id                   string   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Title                string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Date                 string   `protobuf:"bytes,4,opt,name=date,proto3" json:"date,omitempty"`
	Body                 string   `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	Tags                 []string `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Article) Reset()         { *m = Article{} }
func (m *Article) String() string { return proto.CompactTextString(m) }
func (*Article) ProtoMessage()    {}
func (*Article) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d61e4162546950e, []int{0}
}

func (m *Article) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Article.Unmarshal(m, b)
}
func (m *Article) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Article.Marshal(b, m, deterministic)
}
func (m *Article) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Article.Merge(m, src)
}
func (m *Article) XXX_Size() int {
	return xxx_messageInfo_Article.Size(m)
}
func (m *Article) XXX_DiscardUnknown() {
	xxx_messageInfo_Article.DiscardUnknown(m)
}

var xxx_messageInfo_Article proto.InternalMessageInfo

func (m *Article) GetApi() string {
	if m != nil {
		return m.Api
	}
	return ""
}

func (m *Article) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Article) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Article) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

func (m *Article) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *Article) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

type Acknowledgement struct {
	Api                  string   `protobuf:"bytes,1,opt,name=api,proto3" json:"api,omitempty"`
	Id                   string   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Errmessage           string   `protobuf:"bytes,3,opt,name=errmessage,proto3" json:"errmessage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Acknowledgement) Reset()         { *m = Acknowledgement{} }
func (m *Acknowledgement) String() string { return proto.CompactTextString(m) }
func (*Acknowledgement) ProtoMessage()    {}
func (*Acknowledgement) Descriptor() ([]byte, []int) {
	return fileDescriptor_1d61e4162546950e, []int{1}
}

func (m *Acknowledgement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Acknowledgement.Unmarshal(m, b)
}
func (m *Acknowledgement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Acknowledgement.Marshal(b, m, deterministic)
}
func (m *Acknowledgement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Acknowledgement.Merge(m, src)
}
func (m *Acknowledgement) XXX_Size() int {
	return xxx_messageInfo_Acknowledgement.Size(m)
}
func (m *Acknowledgement) XXX_DiscardUnknown() {
	xxx_messageInfo_Acknowledgement.DiscardUnknown(m)
}

var xxx_messageInfo_Acknowledgement proto.InternalMessageInfo

func (m *Acknowledgement) GetApi() string {
	if m != nil {
		return m.Api
	}
	return ""
}

func (m *Acknowledgement) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Acknowledgement) GetErrmessage() string {
	if m != nil {
		return m.Errmessage
	}
	return ""
}

func init() {
	proto.RegisterType((*Article)(nil), "v1.Article")
	proto.RegisterType((*Acknowledgement)(nil), "v1.Acknowledgement")
}

func init() { proto.RegisterFile("proto/v1/article.proto", fileDescriptor_1d61e4162546950e) }

var fileDescriptor_1d61e4162546950e = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0xd0, 0xc1, 0x4a, 0x85, 0x40,
	0x14, 0x06, 0xe0, 0xee, 0x78, 0xaf, 0xd1, 0x89, 0x2c, 0xa6, 0x88, 0xa1, 0x45, 0x88, 0xab, 0x56,
	0x8a, 0xf9, 0x04, 0x12, 0xbd, 0x40, 0x3e, 0xc1, 0xe8, 0x1c, 0x64, 0x48, 0x1d, 0x19, 0x0f, 0x13,
	0x2d, 0x7a, 0xf7, 0x98, 0x51, 0xa1, 0x76, 0x77, 0xf7, 0xff, 0xdf, 0xea, 0x3f, 0x07, 0x1e, 0x67,
	0x6b, 0xc8, 0x14, 0xae, 0x2c, 0xa4, 0x25, 0xdd, 0x0d, 0x98, 0x07, 0xe0, 0xcc, 0x95, 0xd9, 0x0f,
	0x5c, 0xd6, 0x2b, 0xf2, 0x3b, 0x88, 0xe4, 0xac, 0xc5, 0x21, 0x3d, 0xbc, 0x5c, 0x7d, 0xf8, 0xc8,
	0x13, 0x60, 0x5a, 0x09, 0x16, 0x80, 0x69, 0xc5, 0x1f, 0xe0, 0x44, 0x9a, 0x06, 0x14, 0x51, 0xa0,
	0xb5, 0x70, 0x0e, 0x47, 0x25, 0x09, 0xc5, 0x31, 0x60, 0xc8, 0xde, 0x5a, 0xa3, 0xbe, 0xc5, 0x69,
	0x35, 0x9f, 0xbd, 0x91, 0xec, 0x17, 0x11, 0xa7, 0x91, 0x37, 0x9f, 0xb3, 0x06, 0x6e, 0xeb, 0xee,
	0x73, 0x32, 0x5f, 0x03, 0xaa, 0x1e, 0x47, 0x9c, 0xe8, 0x8c, 0x19, 0xcf, 0x00, 0x68, 0xed, 0x88,
	0xcb, 0x22, 0xfb, 0x7d, 0xcb, 0x1f, 0x79, 0x7d, 0x87, 0x64, 0xbb, 0xa9, 0x41, 0xeb, 0x74, 0x87,
	0xbc, 0x82, 0x9b, 0x37, 0x8b, 0x92, 0x70, 0xbf, 0xf5, 0x3a, 0x77, 0x65, 0xbe, 0x95, 0xa7, 0xfb,
	0x50, 0xfe, 0xcf, 0xc8, 0x2e, 0xda, 0x38, 0x7c, 0xa9, 0xfa, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xe3,
	0x17, 0xf5, 0xcb, 0x3f, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ArticleServiceClient is the client API for ArticleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ArticleServiceClient interface {
	CreateArticle(ctx context.Context, in *Article, opts ...grpc.CallOption) (*Acknowledgement, error)
}

type articleServiceClient struct {
	cc *grpc.ClientConn
}

func NewArticleServiceClient(cc *grpc.ClientConn) ArticleServiceClient {
	return &articleServiceClient{cc}
}

func (c *articleServiceClient) CreateArticle(ctx context.Context, in *Article, opts ...grpc.CallOption) (*Acknowledgement, error) {
	out := new(Acknowledgement)
	err := c.cc.Invoke(ctx, "/v1.ArticleService/CreateArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ArticleServiceServer is the server API for ArticleService service.
type ArticleServiceServer interface {
	CreateArticle(context.Context, *Article) (*Acknowledgement, error)
}

// UnimplementedArticleServiceServer can be embedded to have forward compatible implementations.
type UnimplementedArticleServiceServer struct {
}

func (*UnimplementedArticleServiceServer) CreateArticle(ctx context.Context, req *Article) (*Acknowledgement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateArticle not implemented")
}

func RegisterArticleServiceServer(s *grpc.Server, srv ArticleServiceServer) {
	s.RegisterService(&_ArticleService_serviceDesc, srv)
}

func _ArticleService_CreateArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Article)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticleServiceServer).CreateArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.ArticleService/CreateArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticleServiceServer).CreateArticle(ctx, req.(*Article))
	}
	return interceptor(ctx, in, info, handler)
}

var _ArticleService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.ArticleService",
	HandlerType: (*ArticleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateArticle",
			Handler:    _ArticleService_CreateArticle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/article.proto",
}

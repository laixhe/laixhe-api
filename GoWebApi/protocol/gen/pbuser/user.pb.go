// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: pbuser/user.proto

package pbuser

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type User struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 用户ID
	Uid uint64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid"`
	// 用户名
	Uname string `protobuf:"bytes,2,opt,name=uname,proto3" json:"uname"`
	// 用户邮箱
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email"`
	// 创建时间
	CreatedAt     string `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_pbuser_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_pbuser_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_pbuser_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetUid() uint64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *User) GetUname() string {
	if x != nil {
		return x.Uname
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

var File_pbuser_user_proto protoreflect.FileDescriptor

const file_pbuser_user_proto_rawDesc = "" +
	"\n" +
	"\x11pbuser/user.proto\x12\x06pbuser\"c\n" +
	"\x04User\x12\x10\n" +
	"\x03uid\x18\x01 \x01(\x04R\x03uid\x12\x14\n" +
	"\x05uname\x18\x02 \x01(\tR\x05uname\x12\x14\n" +
	"\x05email\x18\x03 \x01(\tR\x05email\x12\x1d\n" +
	"\n" +
	"created_at\x18\x04 \x01(\tR\tcreatedAtB#Z!webapi/protocol/gen/pbuser;pbuserb\x06proto3"

var (
	file_pbuser_user_proto_rawDescOnce sync.Once
	file_pbuser_user_proto_rawDescData []byte
)

func file_pbuser_user_proto_rawDescGZIP() []byte {
	file_pbuser_user_proto_rawDescOnce.Do(func() {
		file_pbuser_user_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pbuser_user_proto_rawDesc), len(file_pbuser_user_proto_rawDesc)))
	})
	return file_pbuser_user_proto_rawDescData
}

var file_pbuser_user_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pbuser_user_proto_goTypes = []any{
	(*User)(nil), // 0: pbuser.User
}
var file_pbuser_user_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pbuser_user_proto_init() }
func file_pbuser_user_proto_init() {
	if File_pbuser_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pbuser_user_proto_rawDesc), len(file_pbuser_user_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pbuser_user_proto_goTypes,
		DependencyIndexes: file_pbuser_user_proto_depIdxs,
		MessageInfos:      file_pbuser_user_proto_msgTypes,
	}.Build()
	File_pbuser_user_proto = out.File
	file_pbuser_user_proto_goTypes = nil
	file_pbuser_user_proto_depIdxs = nil
}

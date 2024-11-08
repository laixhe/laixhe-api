// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: proto/code/code.proto

package code

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 错误状态码
type Code int32

const (
	Code_Success       Code = 0   // 成功
	Code_AuthUserError Code = 101 // 用户或密码错误
	Code_UserExist     Code = 102 // 用户已存在
	Code_UserNotExist  Code = 103 // 用户不存在
	Code_EmailExist    Code = 104 // 邮箱已存在
	Code_EmailNotExist Code = 105 // 邮箱不存在
	Code_PhoneExist    Code = 106 // 手机号码已存在
	Code_PhoneNotExist Code = 107 // 手机号码不存在
)

// Enum value maps for Code.
var (
	Code_name = map[int32]string{
		0:   "Success",
		101: "AuthUserError",
		102: "UserExist",
		103: "UserNotExist",
		104: "EmailExist",
		105: "EmailNotExist",
		106: "PhoneExist",
		107: "PhoneNotExist",
	}
	Code_value = map[string]int32{
		"Success":       0,
		"AuthUserError": 101,
		"UserExist":     102,
		"UserNotExist":  103,
		"EmailExist":    104,
		"EmailNotExist": 105,
		"PhoneExist":    106,
		"PhoneNotExist": 107,
	}
)

func (x Code) Enum() *Code {
	p := new(Code)
	*p = x
	return p
}

func (x Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Code) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_code_code_proto_enumTypes[0].Descriptor()
}

func (Code) Type() protoreflect.EnumType {
	return &file_proto_code_code_proto_enumTypes[0]
}

func (x Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Code.Descriptor instead.
func (Code) EnumDescriptor() ([]byte, []int) {
	return file_proto_code_code_proto_rawDescGZIP(), []int{0}
}

var File_proto_code_code_proto protoreflect.FileDescriptor

var file_proto_code_code_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x63, 0x6f, 0x64,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x2a, 0x8d, 0x01,
	0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x55, 0x73, 0x65, 0x72, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x10, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x45, 0x78,
	0x69, 0x73, 0x74, 0x10, 0x66, 0x12, 0x10, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x6f, 0x74,
	0x45, 0x78, 0x69, 0x73, 0x74, 0x10, 0x67, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x45, 0x78, 0x69, 0x73, 0x74, 0x10, 0x68, 0x12, 0x11, 0x0a, 0x0d, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x4e, 0x6f, 0x74, 0x45, 0x78, 0x69, 0x73, 0x74, 0x10, 0x69, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x68,
	0x6f, 0x6e, 0x65, 0x45, 0x78, 0x69, 0x73, 0x74, 0x10, 0x6a, 0x12, 0x11, 0x0a, 0x0d, 0x50, 0x68,
	0x6f, 0x6e, 0x65, 0x4e, 0x6f, 0x74, 0x45, 0x78, 0x69, 0x73, 0x74, 0x10, 0x6b, 0x42, 0x1a, 0x5a,
	0x18, 0x77, 0x65, 0x62, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x63, 0x6f, 0x64, 0x65, 0x3b, 0x63, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_code_code_proto_rawDescOnce sync.Once
	file_proto_code_code_proto_rawDescData = file_proto_code_code_proto_rawDesc
)

func file_proto_code_code_proto_rawDescGZIP() []byte {
	file_proto_code_code_proto_rawDescOnce.Do(func() {
		file_proto_code_code_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_code_code_proto_rawDescData)
	})
	return file_proto_code_code_proto_rawDescData
}

var file_proto_code_code_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_code_code_proto_goTypes = []any{
	(Code)(0), // 0: code.Code
}
var file_proto_code_code_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_code_code_proto_init() }
func file_proto_code_code_proto_init() {
	if File_proto_code_code_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_code_code_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_code_code_proto_goTypes,
		DependencyIndexes: file_proto_code_code_proto_depIdxs,
		EnumInfos:         file_proto_code_code_proto_enumTypes,
	}.Build()
	File_proto_code_code_proto = out.File
	file_proto_code_code_proto_rawDesc = nil
	file_proto_code_code_proto_goTypes = nil
	file_proto_code_code_proto_depIdxs = nil
}

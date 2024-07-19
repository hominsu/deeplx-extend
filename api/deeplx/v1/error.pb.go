// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: deeplx/v1/error.proto

package v1

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type PallasErrorReason int32

const (
	PallasErrorReason_UNKNOWN   PallasErrorReason = 0
	PallasErrorReason_INTERNAL  PallasErrorReason = 1
	PallasErrorReason_NOT_FOUND PallasErrorReason = 2
	PallasErrorReason_CONFLICT  PallasErrorReason = 3
)

// Enum value maps for PallasErrorReason.
var (
	PallasErrorReason_name = map[int32]string{
		0: "UNKNOWN",
		1: "INTERNAL",
		2: "NOT_FOUND",
		3: "CONFLICT",
	}
	PallasErrorReason_value = map[string]int32{
		"UNKNOWN":   0,
		"INTERNAL":  1,
		"NOT_FOUND": 2,
		"CONFLICT":  3,
	}
)

func (x PallasErrorReason) Enum() *PallasErrorReason {
	p := new(PallasErrorReason)
	*p = x
	return p
}

func (x PallasErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PallasErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_deeplx_v1_error_proto_enumTypes[0].Descriptor()
}

func (PallasErrorReason) Type() protoreflect.EnumType {
	return &file_deeplx_v1_error_proto_enumTypes[0]
}

func (x PallasErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PallasErrorReason.Descriptor instead.
func (PallasErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_deeplx_v1_error_proto_rawDescGZIP(), []int{0}
}

var File_deeplx_v1_error_proto protoreflect.FileDescriptor

var file_deeplx_v1_error_proto_rawDesc = []byte{
	0x0a, 0x15, 0x64, 0x65, 0x65, 0x70, 0x6c, 0x78, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x64, 0x65, 0x65, 0x70, 0x6c, 0x78, 0x2e,
	0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x5d, 0x0a, 0x11, 0x50, 0x61, 0x6c, 0x6c, 0x61,
	0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x0b, 0x0a, 0x07,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x4e, 0x54,
	0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f, 0x46,
	0x4f, 0x55, 0x4e, 0x44, 0x10, 0x02, 0x1a, 0x04, 0xa8, 0x45, 0x94, 0x03, 0x12, 0x12, 0x0a, 0x08,
	0x43, 0x4f, 0x4e, 0x46, 0x4c, 0x49, 0x43, 0x54, 0x10, 0x03, 0x1a, 0x04, 0xa8, 0x45, 0x99, 0x03,
	0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x39, 0x50, 0x01, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x69, 0x6f, 0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x2f, 0x64, 0x65, 0x65, 0x70, 0x6c, 0x78, 0x2d, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x65, 0x70, 0x6c, 0x78, 0x2f, 0x76, 0x31, 0x3b, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_deeplx_v1_error_proto_rawDescOnce sync.Once
	file_deeplx_v1_error_proto_rawDescData = file_deeplx_v1_error_proto_rawDesc
)

func file_deeplx_v1_error_proto_rawDescGZIP() []byte {
	file_deeplx_v1_error_proto_rawDescOnce.Do(func() {
		file_deeplx_v1_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_deeplx_v1_error_proto_rawDescData)
	})
	return file_deeplx_v1_error_proto_rawDescData
}

var file_deeplx_v1_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_deeplx_v1_error_proto_goTypes = []any{
	(PallasErrorReason)(0), // 0: deeplx.v1.PallasErrorReason
}
var file_deeplx_v1_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_deeplx_v1_error_proto_init() }
func file_deeplx_v1_error_proto_init() {
	if File_deeplx_v1_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_deeplx_v1_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_deeplx_v1_error_proto_goTypes,
		DependencyIndexes: file_deeplx_v1_error_proto_depIdxs,
		EnumInfos:         file_deeplx_v1_error_proto_enumTypes,
	}.Build()
	File_deeplx_v1_error_proto = out.File
	file_deeplx_v1_error_proto_rawDesc = nil
	file_deeplx_v1_error_proto_goTypes = nil
	file_deeplx_v1_error_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: user.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6e, 0x63,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x0b, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xa6, 0x02, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72,
	0x12, 0x37, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x6e, 0x63, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x1a, 0x13, 0x2e, 0x6e, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x34, 0x0a, 0x06, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x12, 0x15, 0x2e, 0x6e, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x1a, 0x13, 0x2e, 0x6e, 0x63, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x3c, 0x0a, 0x09, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x15, 0x2e, 0x6e,
	0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x4b, 0x65, 0x79, 0x1a, 0x18, 0x2e, 0x6e, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x36, 0x0a,
	0x07, 0x41, 0x64, 0x64, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x6e, 0x63, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65,
	0x1a, 0x13, 0x2e, 0x6e, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x39, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x6f, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x6e, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x1a, 0x13, 0x2e, 0x6e, 0x63,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e,
	0x2d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x76, 0x65, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f,
	0x61, 0x70, 0x69, 0x2d, 0x72, 0x62, 0x61, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_user_proto_goTypes = []interface{}{
	(*UserEntity)(nil), // 0: ncs.protobuf.userEntity
	(*UserKey)(nil),    // 1: ncs.protobuf.userKey
	(*UserRole)(nil),   // 2: ncs.protobuf.userRole
	(*Empty)(nil),      // 3: ncs.protobuf.empty
}
var file_user_proto_depIdxs = []int32{
	0, // 0: ncs.protobuf.User.Create:input_type -> ncs.protobuf.userEntity
	1, // 1: ncs.protobuf.User.Delete:input_type -> ncs.protobuf.userKey
	1, // 2: ncs.protobuf.User.FindByKey:input_type -> ncs.protobuf.userKey
	2, // 3: ncs.protobuf.User.AddRole:input_type -> ncs.protobuf.userRole
	2, // 4: ncs.protobuf.User.DeleteRole:input_type -> ncs.protobuf.userRole
	3, // 5: ncs.protobuf.User.Create:output_type -> ncs.protobuf.empty
	3, // 6: ncs.protobuf.User.Delete:output_type -> ncs.protobuf.empty
	0, // 7: ncs.protobuf.User.FindByKey:output_type -> ncs.protobuf.userEntity
	3, // 8: ncs.protobuf.User.AddRole:output_type -> ncs.protobuf.empty
	3, // 9: ncs.protobuf.User.DeleteRole:output_type -> ncs.protobuf.empty
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	file_types_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}

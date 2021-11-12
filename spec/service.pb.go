// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: spec/service.proto

package spec

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Type int32

const (
	Type_OPEN Type = 0
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "OPEN",
	}
	Type_value = map[string]int32{
		"OPEN": 0,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_spec_service_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_spec_service_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_spec_service_proto_rawDescGZIP(), []int{0}
}

type PingRequestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	N int32 `protobuf:"varint,1,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *PingRequestResponse) Reset() {
	*x = PingRequestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spec_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequestResponse) ProtoMessage() {}

func (x *PingRequestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_spec_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequestResponse.ProtoReflect.Descriptor instead.
func (*PingRequestResponse) Descriptor() ([]byte, []int) {
	return file_spec_service_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequestResponse) GetN() int32 {
	if x != nil {
		return x.N
	}
	return 0
}

type ConnectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token       string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ServiceName string `protobuf:"bytes,2,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Private     bool   `protobuf:"varint,3,opt,name=private,proto3" json:"private,omitempty"`
}

func (x *ConnectRequest) Reset() {
	*x = ConnectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spec_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectRequest) ProtoMessage() {}

func (x *ConnectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_spec_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectRequest.ProtoReflect.Descriptor instead.
func (*ConnectRequest) Descriptor() ([]byte, []int) {
	return file_spec_service_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ConnectRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *ConnectRequest) GetPrivate() bool {
	if x != nil {
		return x.Private
	}
	return false
}

type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind Type `protobuf:"varint,1,opt,name=kind,proto3,enum=Type" json:"kind,omitempty"`
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_spec_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_spec_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_spec_service_proto_rawDescGZIP(), []int{2}
}

func (x *Command) GetKind() Type {
	if x != nil {
		return x.Kind
	}
	return Type_OPEN
}

var File_spec_service_proto protoreflect.FileDescriptor

var file_spec_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x23, 0x0a, 0x13, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x6e, 0x22, 0x63, 0x0a, 0x0e, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x22, 0x24,
	0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x19, 0x0a, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x05, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x6b, 0x69, 0x6e, 0x64, 0x2a, 0x10, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04,
	0x4f, 0x50, 0x45, 0x4e, 0x10, 0x00, 0x32, 0x69, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x28, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x0f, 0x2e, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x00, 0x30, 0x01, 0x12, 0x34, 0x0a, 0x04, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x14, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x14, 0x2e, 0x50, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_spec_service_proto_rawDescOnce sync.Once
	file_spec_service_proto_rawDescData = file_spec_service_proto_rawDesc
)

func file_spec_service_proto_rawDescGZIP() []byte {
	file_spec_service_proto_rawDescOnce.Do(func() {
		file_spec_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_spec_service_proto_rawDescData)
	})
	return file_spec_service_proto_rawDescData
}

var file_spec_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_spec_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_spec_service_proto_goTypes = []interface{}{
	(Type)(0),                   // 0: Type
	(*PingRequestResponse)(nil), // 1: PingRequestResponse
	(*ConnectRequest)(nil),      // 2: ConnectRequest
	(*Command)(nil),             // 3: Command
}
var file_spec_service_proto_depIdxs = []int32{
	0, // 0: Command.kind:type_name -> Type
	2, // 1: Service.Connect:input_type -> ConnectRequest
	1, // 2: Service.Ping:input_type -> PingRequestResponse
	3, // 3: Service.Connect:output_type -> Command
	1, // 4: Service.Ping:output_type -> PingRequestResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_spec_service_proto_init() }
func file_spec_service_proto_init() {
	if File_spec_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_spec_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequestResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_spec_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_spec_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_spec_service_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_spec_service_proto_goTypes,
		DependencyIndexes: file_spec_service_proto_depIdxs,
		EnumInfos:         file_spec_service_proto_enumTypes,
		MessageInfos:      file_spec_service_proto_msgTypes,
	}.Build()
	File_spec_service_proto = out.File
	file_spec_service_proto_rawDesc = nil
	file_spec_service_proto_goTypes = nil
	file_spec_service_proto_depIdxs = nil
}

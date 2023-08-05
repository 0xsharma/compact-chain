// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: protos/p2p.proto

package protos

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

type LatestBlockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LatestBlockRequest) Reset() {
	*x = LatestBlockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_p2p_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LatestBlockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatestBlockRequest) ProtoMessage() {}

func (x *LatestBlockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_p2p_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatestBlockRequest.ProtoReflect.Descriptor instead.
func (*LatestBlockRequest) Descriptor() ([]byte, []int) {
	return file_protos_p2p_proto_rawDescGZIP(), []int{0}
}

type LatestBlockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height       uint64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	EncodedBlock []byte `protobuf:"bytes,2,opt,name=encodedBlock,proto3" json:"encodedBlock,omitempty"`
}

func (x *LatestBlockResponse) Reset() {
	*x = LatestBlockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_p2p_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LatestBlockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatestBlockResponse) ProtoMessage() {}

func (x *LatestBlockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_p2p_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatestBlockResponse.ProtoReflect.Descriptor instead.
func (*LatestBlockResponse) Descriptor() ([]byte, []int) {
	return file_protos_p2p_proto_rawDescGZIP(), []int{1}
}

func (x *LatestBlockResponse) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *LatestBlockResponse) GetEncodedBlock() []byte {
	if x != nil {
		return x.EncodedBlock
	}
	return nil
}

type TxpoolPendingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TxpoolPendingRequest) Reset() {
	*x = TxpoolPendingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_p2p_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxpoolPendingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxpoolPendingRequest) ProtoMessage() {}

func (x *TxpoolPendingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_p2p_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxpoolPendingRequest.ProtoReflect.Descriptor instead.
func (*TxpoolPendingRequest) Descriptor() ([]byte, []int) {
	return file_protos_p2p_proto_rawDescGZIP(), []int{2}
}

type TxpoolPendingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncodedTxs [][]byte `protobuf:"bytes,1,rep,name=encodedTxs,proto3" json:"encodedTxs,omitempty"`
}

func (x *TxpoolPendingResponse) Reset() {
	*x = TxpoolPendingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_p2p_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxpoolPendingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxpoolPendingResponse) ProtoMessage() {}

func (x *TxpoolPendingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_p2p_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxpoolPendingResponse.ProtoReflect.Descriptor instead.
func (*TxpoolPendingResponse) Descriptor() ([]byte, []int) {
	return file_protos_p2p_proto_rawDescGZIP(), []int{3}
}

func (x *TxpoolPendingResponse) GetEncodedTxs() [][]byte {
	if x != nil {
		return x.EncodedTxs
	}
	return nil
}

var File_protos_p2p_proto protoreflect.FileDescriptor

var file_protos_p2p_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x70, 0x32, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x22, 0x14, 0x0a, 0x12, 0x4c, 0x61,
	0x74, 0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x51, 0x0a, 0x13, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12,
	0x22, 0x0a, 0x0c, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x22, 0x16, 0x0a, 0x14, 0x54, 0x78, 0x70, 0x6f, 0x6f, 0x6c, 0x50, 0x65, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x37, 0x0a, 0x15, 0x54,
	0x78, 0x70, 0x6f, 0x6f, 0x6c, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x54,
	0x78, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0a, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65,
	0x64, 0x54, 0x78, 0x73, 0x32, 0x9b, 0x01, 0x0a, 0x03, 0x50, 0x32, 0x50, 0x12, 0x46, 0x0a, 0x0b,
	0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2e, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0d, 0x54, 0x78, 0x50, 0x6f, 0x6f, 0x6c, 0x50, 0x65,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x54,
	0x78, 0x70, 0x6f, 0x6f, 0x6c, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x54, 0x78, 0x70,
	0x6f, 0x6f, 0x6c, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_p2p_proto_rawDescOnce sync.Once
	file_protos_p2p_proto_rawDescData = file_protos_p2p_proto_rawDesc
)

func file_protos_p2p_proto_rawDescGZIP() []byte {
	file_protos_p2p_proto_rawDescOnce.Do(func() {
		file_protos_p2p_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_p2p_proto_rawDescData)
	})
	return file_protos_p2p_proto_rawDescData
}

var file_protos_p2p_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protos_p2p_proto_goTypes = []interface{}{
	(*LatestBlockRequest)(nil),    // 0: protos.LatestBlockRequest
	(*LatestBlockResponse)(nil),   // 1: protos.LatestBlockResponse
	(*TxpoolPendingRequest)(nil),  // 2: protos.TxpoolPendingRequest
	(*TxpoolPendingResponse)(nil), // 3: protos.TxpoolPendingResponse
}
var file_protos_p2p_proto_depIdxs = []int32{
	0, // 0: protos.P2P.LatestBlock:input_type -> protos.LatestBlockRequest
	2, // 1: protos.P2P.TxPoolPending:input_type -> protos.TxpoolPendingRequest
	1, // 2: protos.P2P.LatestBlock:output_type -> protos.LatestBlockResponse
	3, // 3: protos.P2P.TxPoolPending:output_type -> protos.TxpoolPendingResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_p2p_proto_init() }
func file_protos_p2p_proto_init() {
	if File_protos_p2p_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_p2p_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LatestBlockRequest); i {
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
		file_protos_p2p_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LatestBlockResponse); i {
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
		file_protos_p2p_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxpoolPendingRequest); i {
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
		file_protos_p2p_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxpoolPendingResponse); i {
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
			RawDescriptor: file_protos_p2p_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_p2p_proto_goTypes,
		DependencyIndexes: file_protos_p2p_proto_depIdxs,
		MessageInfos:      file_protos_p2p_proto_msgTypes,
	}.Build()
	File_protos_p2p_proto = out.File
	file_protos_p2p_proto_rawDesc = nil
	file_protos_p2p_proto_goTypes = nil
	file_protos_p2p_proto_depIdxs = nil
}

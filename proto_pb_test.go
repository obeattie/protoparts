// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.1
// source: proto_test.proto

package protoparts

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

type Person_MaritalStatus int32

const (
	Person_PREFER_NOT_TO_SAY Person_MaritalStatus = 0
	Person_SINGLE            Person_MaritalStatus = 1
	Person_DIVORCED          Person_MaritalStatus = 2
	Person_WIDOWED           Person_MaritalStatus = 3
)

// Enum value maps for Person_MaritalStatus.
var (
	Person_MaritalStatus_name = map[int32]string{
		0: "PREFER_NOT_TO_SAY",
		1: "SINGLE",
		2: "DIVORCED",
		3: "WIDOWED",
	}
	Person_MaritalStatus_value = map[string]int32{
		"PREFER_NOT_TO_SAY": 0,
		"SINGLE":            1,
		"DIVORCED":          2,
		"WIDOWED":           3,
	}
)

func (x Person_MaritalStatus) Enum() *Person_MaritalStatus {
	p := new(Person_MaritalStatus)
	*p = x
	return p
}

func (x Person_MaritalStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Person_MaritalStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_test_proto_enumTypes[0].Descriptor()
}

func (Person_MaritalStatus) Type() protoreflect.EnumType {
	return &file_proto_test_proto_enumTypes[0]
}

func (x Person_MaritalStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Person_MaritalStatus.Descriptor instead.
func (Person_MaritalStatus) EnumDescriptor() ([]byte, []int) {
	return file_proto_test_proto_rawDescGZIP(), []int{2, 0}
}

type LatLng struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *LatLng) Reset() {
	*x = LatLng{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LatLng) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatLng) ProtoMessage() {}

func (x *LatLng) ProtoReflect() protoreflect.Message {
	mi := &file_proto_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatLng.ProtoReflect.Descriptor instead.
func (*LatLng) Descriptor() ([]byte, []int) {
	return file_proto_test_proto_rawDescGZIP(), []int{0}
}

func (x *LatLng) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *LatLng) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreetAddress string  `protobuf:"bytes,1,opt,name=street_address,json=streetAddress,proto3" json:"street_address,omitempty"`
	City          string  `protobuf:"bytes,2,opt,name=city,proto3" json:"city,omitempty"`
	LatLng        *LatLng `protobuf:"bytes,3,opt,name=lat_lng,json=latLng,proto3" json:"lat_lng,omitempty"`
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_proto_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_proto_test_proto_rawDescGZIP(), []int{1}
}

func (x *Address) GetStreetAddress() string {
	if x != nil {
		return x.StreetAddress
	}
	return ""
}

func (x *Address) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *Address) GetLatLng() *LatLng {
	if x != nil {
		return x.LatLng
	}
	return nil
}

type Person struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name            *string              `protobuf:"bytes,1,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Address         *Address             `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	MoarAddresses   []*Address           `protobuf:"bytes,3,rep,name=moar_addresses,json=moarAddresses,proto3" json:"moar_addresses,omitempty"`
	Tags            []string             `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	Boop            [][]byte             `protobuf:"bytes,5,rep,name=boop,proto3" json:"boop,omitempty"`
	MapStringLatlng map[string]*LatLng   `protobuf:"bytes,6,rep,name=map_string_latlng,json=mapStringLatlng,proto3" json:"map_string_latlng,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	MaritalStatus   Person_MaritalStatus `protobuf:"varint,7,opt,name=marital_status,json=maritalStatus,proto3,enum=protoparts.Person_MaritalStatus" json:"marital_status,omitempty"`
	// Types that are assignable to StringOrLatlng:
	//	*Person_MaybeString
	//	*Person_MaybeLatlng
	StringOrLatlng  isPerson_StringOrLatlng `protobuf_oneof:"string_or_latlng"`
	MapStringString map[string]string       `protobuf:"bytes,10,rep,name=map_string_string,json=mapStringString,proto3" json:"map_string_string,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Person) Reset() {
	*x = Person{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_test_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Person) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Person) ProtoMessage() {}

func (x *Person) ProtoReflect() protoreflect.Message {
	mi := &file_proto_test_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Person.ProtoReflect.Descriptor instead.
func (*Person) Descriptor() ([]byte, []int) {
	return file_proto_test_proto_rawDescGZIP(), []int{2}
}

func (x *Person) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Person) GetAddress() *Address {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *Person) GetMoarAddresses() []*Address {
	if x != nil {
		return x.MoarAddresses
	}
	return nil
}

func (x *Person) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Person) GetBoop() [][]byte {
	if x != nil {
		return x.Boop
	}
	return nil
}

func (x *Person) GetMapStringLatlng() map[string]*LatLng {
	if x != nil {
		return x.MapStringLatlng
	}
	return nil
}

func (x *Person) GetMaritalStatus() Person_MaritalStatus {
	if x != nil {
		return x.MaritalStatus
	}
	return Person_PREFER_NOT_TO_SAY
}

func (m *Person) GetStringOrLatlng() isPerson_StringOrLatlng {
	if m != nil {
		return m.StringOrLatlng
	}
	return nil
}

func (x *Person) GetMaybeString() string {
	if x, ok := x.GetStringOrLatlng().(*Person_MaybeString); ok {
		return x.MaybeString
	}
	return ""
}

func (x *Person) GetMaybeLatlng() *LatLng {
	if x, ok := x.GetStringOrLatlng().(*Person_MaybeLatlng); ok {
		return x.MaybeLatlng
	}
	return nil
}

func (x *Person) GetMapStringString() map[string]string {
	if x != nil {
		return x.MapStringString
	}
	return nil
}

type isPerson_StringOrLatlng interface {
	isPerson_StringOrLatlng()
}

type Person_MaybeString struct {
	MaybeString string `protobuf:"bytes,8,opt,name=maybe_string,json=maybeString,proto3,oneof"`
}

type Person_MaybeLatlng struct {
	MaybeLatlng *LatLng `protobuf:"bytes,9,opt,name=maybe_latlng,json=maybeLatlng,proto3,oneof"`
}

func (*Person_MaybeString) isPerson_StringOrLatlng() {}

func (*Person_MaybeLatlng) isPerson_StringOrLatlng() {}

var File_proto_test_proto protoreflect.FileDescriptor

var file_proto_test_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x70, 0x61, 0x72, 0x74, 0x73, 0x22, 0x42,
	0x0a, 0x06, 0x4c, 0x61, 0x74, 0x4c, 0x6e, 0x67, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75,
	0x64, 0x65, 0x22, 0x71, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x25, 0x0a,
	0x0e, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x63, 0x69, 0x74, 0x79, 0x12, 0x2b, 0x0a, 0x07, 0x6c, 0x61, 0x74, 0x5f,
	0x6c, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x70, 0x61, 0x72, 0x74, 0x73, 0x2e, 0x4c, 0x61, 0x74, 0x4c, 0x6e, 0x67, 0x52, 0x06, 0x6c,
	0x61, 0x74, 0x4c, 0x6e, 0x67, 0x22, 0x8d, 0x06, 0x0a, 0x06, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e,
	0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x70, 0x61, 0x72, 0x74, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3a, 0x0a, 0x0e, 0x6d, 0x6f, 0x61, 0x72,
	0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x70, 0x61, 0x72, 0x74, 0x73, 0x2e, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x0d, 0x6d, 0x6f, 0x61, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x6f, 0x70,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x70, 0x12, 0x53, 0x0a, 0x11,
	0x6d, 0x61, 0x70, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x6c, 0x61, 0x74, 0x6c, 0x6e,
	0x67, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x70,
	0x61, 0x72, 0x74, 0x73, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x70, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x4c, 0x61, 0x74, 0x6c, 0x6e, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0f, 0x6d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4c, 0x61, 0x74, 0x6c, 0x6e,
	0x67, 0x12, 0x47, 0x0a, 0x0e, 0x6d, 0x61, 0x72, 0x69, 0x74, 0x61, 0x6c, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x70, 0x61, 0x72, 0x74, 0x73, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x4d, 0x61,
	0x72, 0x69, 0x74, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0d, 0x6d, 0x61, 0x72,
	0x69, 0x74, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x23, 0x0a, 0x0c, 0x6d, 0x61,
	0x79, 0x62, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x0b, 0x6d, 0x61, 0x79, 0x62, 0x65, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12,
	0x37, 0x0a, 0x0c, 0x6d, 0x61, 0x79, 0x62, 0x65, 0x5f, 0x6c, 0x61, 0x74, 0x6c, 0x6e, 0x67, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x70, 0x61, 0x72,
	0x74, 0x73, 0x2e, 0x4c, 0x61, 0x74, 0x4c, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x0b, 0x6d, 0x61, 0x79,
	0x62, 0x65, 0x4c, 0x61, 0x74, 0x6c, 0x6e, 0x67, 0x12, 0x53, 0x0a, 0x11, 0x6d, 0x61, 0x70, 0x5f,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x0a, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x70, 0x61, 0x72, 0x74, 0x73,
	0x2e, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0f, 0x6d, 0x61,
	0x70, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x56, 0x0a,
	0x14, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4c, 0x61, 0x74, 0x6c, 0x6e, 0x67,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x70, 0x61,
	0x72, 0x74, 0x73, 0x2e, 0x4c, 0x61, 0x74, 0x4c, 0x6e, 0x67, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x42, 0x0a, 0x14, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x4d, 0x0a, 0x0d, 0x4d, 0x61, 0x72,
	0x69, 0x74, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x52,
	0x45, 0x46, 0x45, 0x52, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x54, 0x4f, 0x5f, 0x53, 0x41, 0x59, 0x10,
	0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x49, 0x4e, 0x47, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x0c, 0x0a,
	0x08, 0x44, 0x49, 0x56, 0x4f, 0x52, 0x43, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x57,
	0x49, 0x44, 0x4f, 0x57, 0x45, 0x44, 0x10, 0x03, 0x42, 0x12, 0x0a, 0x10, 0x73, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x5f, 0x6f, 0x72, 0x5f, 0x6c, 0x61, 0x74, 0x6c, 0x6e, 0x67, 0x42, 0x07, 0x0a, 0x05,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x62, 0x65, 0x61, 0x74, 0x74, 0x69, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x70, 0x61, 0x72, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_test_proto_rawDescOnce sync.Once
	file_proto_test_proto_rawDescData = file_proto_test_proto_rawDesc
)

func file_proto_test_proto_rawDescGZIP() []byte {
	file_proto_test_proto_rawDescOnce.Do(func() {
		file_proto_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_test_proto_rawDescData)
	})
	return file_proto_test_proto_rawDescData
}

var file_proto_test_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_test_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_test_proto_goTypes = []interface{}{
	(Person_MaritalStatus)(0), // 0: protoparts.Person.MaritalStatus
	(*LatLng)(nil),            // 1: protoparts.LatLng
	(*Address)(nil),           // 2: protoparts.Address
	(*Person)(nil),            // 3: protoparts.Person
	nil,                       // 4: protoparts.Person.MapStringLatlngEntry
	nil,                       // 5: protoparts.Person.MapStringStringEntry
}
var file_proto_test_proto_depIdxs = []int32{
	1, // 0: protoparts.Address.lat_lng:type_name -> protoparts.LatLng
	2, // 1: protoparts.Person.address:type_name -> protoparts.Address
	2, // 2: protoparts.Person.moar_addresses:type_name -> protoparts.Address
	4, // 3: protoparts.Person.map_string_latlng:type_name -> protoparts.Person.MapStringLatlngEntry
	0, // 4: protoparts.Person.marital_status:type_name -> protoparts.Person.MaritalStatus
	1, // 5: protoparts.Person.maybe_latlng:type_name -> protoparts.LatLng
	5, // 6: protoparts.Person.map_string_string:type_name -> protoparts.Person.MapStringStringEntry
	1, // 7: protoparts.Person.MapStringLatlngEntry.value:type_name -> protoparts.LatLng
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_proto_test_proto_init() }
func file_proto_test_proto_init() {
	if File_proto_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LatLng); i {
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
		file_proto_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
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
		file_proto_test_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Person); i {
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
	file_proto_test_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Person_MaybeString)(nil),
		(*Person_MaybeLatlng)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_test_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_test_proto_goTypes,
		DependencyIndexes: file_proto_test_proto_depIdxs,
		EnumInfos:         file_proto_test_proto_enumTypes,
		MessageInfos:      file_proto_test_proto_msgTypes,
	}.Build()
	File_proto_test_proto = out.File
	file_proto_test_proto_rawDesc = nil
	file_proto_test_proto_goTypes = nil
	file_proto_test_proto_depIdxs = nil
}
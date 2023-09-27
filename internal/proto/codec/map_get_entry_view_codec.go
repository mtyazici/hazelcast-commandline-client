/*
* Copyright (c) 2008-2023, Hazelcast, Inc. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License")
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package codec

import (
	control "github.com/hazelcast/hazelcast-commandline-client/internal/proto/codec/control"
	iserialization "github.com/hazelcast/hazelcast-go-client"
	proto "github.com/hazelcast/hazelcast-go-client"
)

const (
	MapGetEntryViewCodecRequestMessageType  = int32(0x011D00)
	MapGetEntryViewCodecResponseMessageType = int32(0x011D01)

	MapGetEntryViewCodecRequestThreadIdOffset   = proto.PartitionIDOffset + proto.IntSizeInBytes
	MapGetEntryViewCodecRequestInitialFrameSize = MapGetEntryViewCodecRequestThreadIdOffset + proto.LongSizeInBytes

	MapGetEntryViewResponseMaxIdleOffset = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
)

// Returns the EntryView for the specified key.
// This method returns a clone of original mapping, modifying the returned value does not change the actual value
// in the map. One should put modified value back to make changes visible to all nodes.

func EncodeMapGetEntryViewRequest(name string, key iserialization.Data, threadId int64) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(true)

	initialFrame := proto.NewFrameWith(make([]byte, MapGetEntryViewCodecRequestInitialFrameSize), proto.UnfragmentedMessage)
	EncodeLong(initialFrame.Content, MapGetEntryViewCodecRequestThreadIdOffset, threadId)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(MapGetEntryViewCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	EncodeString(clientMessage, name)
	EncodeData(clientMessage, key)

	return clientMessage
}

func DecodeMapGetEntryViewResponse(clientMessage *proto.ClientMessage) (response control.SimpleEntryView, maxIdle int64) {
	frameIterator := clientMessage.FrameIterator()
	initialFrame := frameIterator.Next()

	maxIdle = DecodeLong(initialFrame.Content, MapGetEntryViewResponseMaxIdleOffset)
	response = DecodeNullableForSimpleEntryView(frameIterator)

	return response, maxIdle
}

func DecodeNullableForSimpleEntryView(frameIterator *proto.ForwardFrameIterator) control.SimpleEntryView {
	if NextFrameIsNullFrame(frameIterator) {
		return control.SimpleEntryView{}
	}
	return DecodeSimpleEntryView(frameIterator)
}

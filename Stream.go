package main

import "github.com/pion/webrtc/v4"

type Stream struct {
	id   string
	kind webrtc.RTPCodecType
}

func newStream(id string, kind webrtc.RTPCodecType) *Stream {
	stream := Stream{
		id:   id,
		kind: kind,
	}
	return &stream
}

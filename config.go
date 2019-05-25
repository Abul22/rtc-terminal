package main

import (
	"github.com/pion/webrtc/v2"
)
	

var (
	configRTC = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
				},
			},
		},
	}
	
)

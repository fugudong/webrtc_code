package model

type IceServer  struct {
	Credential string   `json:"credential"`
	Urls       []string `json:"urls"`
	Username   string   `json:"username"`
}

type RtcConfiguration struct {
	BundlePolicy string `json:"bundlePolicy"`	// "balanced" | "max-compat" | "max-bundle";
	IceServers []IceServer `json:"iceServers"`
	IceTransportPolicy string `json:"iceTransportPolicy"`	// relay, all
	RtcpMuxPolicy      string `json:"rtcpMuxPolicy"` // "negotiate" | "require"
}

package types

import "github.com/prometheus/client_golang/prometheus"

type Config struct {
	ChainId             string
	BinaryPath          string
	HomePath            string
	PoolId              int
	Seeds               string
	FallbackEndpoints   string
	StateRequests       bool
	Interval            int
	HeightDifferenceMax int
	HeightDifferenceMin int
	Metrics             bool
}

type HeightResponse struct {
	Result struct {
		Response struct {
			LastBlockHeight string `json:"last_block_height"`
		} `json:"response"`
	} `json:"result"`
}

type Metrics struct {
	PoolHeight  prometheus.Gauge
	NodeHeight  prometheus.Gauge
	MaxHeight   prometheus.Gauge
	MinHeight   prometheus.Gauge
	DataDirSize prometheus.Gauge
}

type PoolSettingsType struct {
	MaxBundleSize  int
	UploadInterval int
}

type ProcessType struct {
	Id        int
	GhostMode bool
}

type SettingsResponse struct {
	Pool struct {
		Data struct {
			StartKey       string `json:"start_key"`
			CurrentKey     string `json:"current_key"`
			UploadInterval string `json:"upload_interval"`
			MaxBundleSize  string `json:"max_bundle_size"`
		} `json:"data"`
	} `json:"pool"`
}

type SettingsType struct {
	MaxDifference int
	Seeds         string
	Interval      int
	KeepEvery     int
	KeepRecent    int
}

package types

var (
	KaonEndpoints = []string{
		"https://api.kaon.kyve.network",
	}
	KorelliaEndpoints = []string{
		"https://api.korellia.kyve.network",
	}
	MainnetEndpoints = []string{
		"https://api.kyve.network",
	}
)

var Version string

const (
	BackoffMaxRetries = 15
	SegmentKey        = "oLhjq9j6pOrIB7TjNHxWWB1ILhK5Fwn6"
)

const (
	BACKUP  = "BACKUP"
	INIT    = "INIT"
	PRUNE   = "PRUNE"
	START   = "START"
	VERSION = "VERSION"
)

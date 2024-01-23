package types

var (
	KaonEndpoints = []string{
		"https://api-eu-1.kaon.kyve.network",
		"https://api-us-1.kaon.kyve.network",
	}
	KorelliaEndpoints = []string{
		"https://api.korellia.kyve.network",
		"https://api-eu-1.korellia.kyve.network",
		"https://api-explorer.korellia.kyve.network",
	}
	MainnetEndpoints = []string{
		"https://api-eu-1.kyve.network",
		"https://api-us-1.kyve.network",
	}
)

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

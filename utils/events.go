package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/KYVENetwork/supervysor/cmd/supervysor/commands/helpers"

	"github.com/KYVENetwork/supervysor/types"
	"github.com/google/uuid"
	"github.com/segmentio/analytics-go"
)

var (
	startId = uuid.New().String()
	client  = analytics.New(types.SegmentKey)
)

func getContext() *analytics.Context {
	version := "local"
	build, _ := debug.ReadBuildInfo()

	if strings.TrimSpace(build.Main.Version) != "" {
		version = strings.TrimSpace(build.Main.Version)
	}

	timezone, _ := time.Now().Zone()
	locale := os.Getenv("LANG")

	return &analytics.Context{
		App: analytics.AppInfo{
			Name:    "supervysor",
			Version: version,
		},
		Location: analytics.LocationInfo{},
		OS: analytics.OSInfo{
			Name: fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH),
		},
		Locale:   locale,
		Timezone: timezone,
	}
}

func getUserId() (string, error) {
	supervysorDir, err := helpers.GetSupervysorDir()
	if err != nil {
		return "", fmt.Errorf("could not find .supervysor directory: %s", err)
	}

	userId := uuid.New().String()

	idFile := filepath.Join(supervysorDir, "id")
	if _, err = os.Stat(idFile); os.IsNotExist(err) {
		if err := os.WriteFile(idFile, []byte(userId), 0o755); err != nil {
			return "", err
		}
	} else {
		data, err := os.ReadFile(idFile)
		if err != nil {
			return "", err
		}
		userId = string(data)
	}

	return userId, nil
}

func TrackBackupEvent(optOut bool) {
	if optOut {
		return
	}

	userId, err := getUserId()
	if err != nil {
		return
	}

	err = client.Enqueue(analytics.Track{
		UserId:  userId,
		Event:   types.BACKUP,
		Context: getContext(),
	})

	if err != nil {
		return
	}

	err = client.Close()
	_ = err
}

func TrackInitEvent(chainId string, optOut bool) {
	if optOut {
		return
	}

	userId, err := getUserId()
	if err != nil {
		return
	}

	err = client.Enqueue(analytics.Track{
		UserId:     userId,
		Event:      types.INIT,
		Properties: analytics.NewProperties().Set("chain_id", chainId),
		Context:    getContext(),
	})

	if err != nil {
		return
	}

	err = client.Close()
	_ = err
}

func TrackPruneEvent(optOut bool) {
	if optOut {
		return
	}

	userId, err := getUserId()
	if err != nil {
		return
	}

	err = client.Enqueue(analytics.Track{
		UserId:  userId,
		Event:   types.PRUNE,
		Context: getContext(),
	})

	if err != nil {
		return
	}

	err = client.Close()
	_ = err
}

func TrackStartEvent(chainId string, optOut bool) {
	if optOut {
		return
	}

	userId, err := getUserId()
	if err != nil {
		return
	}

	err = client.Enqueue(analytics.Track{
		UserId:     userId,
		Event:      types.START,
		Properties: analytics.NewProperties().Set("chain_id", chainId).Set("start_id", startId),
		Context:    getContext(),
	})

	if err != nil {
		return
	}

	err = client.Close()
	_ = err
}

func TrackVersionEvent(optOut bool) {
	if optOut {
		return
	}

	userId, err := getUserId()
	if err != nil {
		return
	}

	err = client.Enqueue(analytics.Track{
		UserId:  userId,
		Event:   types.VERSION,
		Context: getContext(),
	})

	if err != nil {
		return
	}

	err = client.Close()
	_ = err
}

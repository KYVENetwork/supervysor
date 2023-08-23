package helpers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/KYVENetwork/supervysor/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cfg "github.com/tendermint/tendermint/config"
)

func GetDirectorySize(dirPath string) (float64, error) {
	var s int64
	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			s += info.Size()
		}
		return err
	})

	return float64(s), err
}

func GetLogsDir() (string, error) {
	supervysorDir, err := GetSupervysorDir()
	if err != nil {
		return "", fmt.Errorf("could not find .supervysor directory: %s", err)
	}

	logsDir := filepath.Join(supervysorDir, "logs")

	if _, err = os.Stat(logsDir); os.IsNotExist(err) {
		err = os.Mkdir(logsDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("could not create logs directory: %s", err)
		}
	}

	return logsDir, nil
}

func GetSupervysorDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %s", err)
	}

	supervysorDir := filepath.Join(home, ".supervysor")

	if _, err = os.Stat(supervysorDir); os.IsNotExist(err) {
		err = os.Mkdir(supervysorDir, 0o755)
		if err != nil {
			return "", err
		}
	}

	return supervysorDir, nil
}

func LoadConfig(homeDir string) (config *cfg.Config, err error) {
	config = cfg.DefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(homeDir)
	viper.AddConfigPath(filepath.Join(homeDir, "config"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	config.SetRoot(homeDir)

	return config, nil
}

func NewMetrics(reg prometheus.Registerer) *types.Metrics {
	m := &types.Metrics{
		PoolHeight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "supervysor",
			Name:      "pool_height",
			Help:      "Height of the specified KYVE data pool.",
		}),
		NodeHeight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "supervysor",
			Name:      "node_height",
			Help:      "Height of the running data source node.",
		}),
		MaxHeight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "supervysor",
			Name:      "max_height",
			Help:      "Maximum height of node until Ghost Mode enabling.",
		}),
		MinHeight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "supervysor",
			Name:      "min_height",
			Help:      "Minimum height of node until Normal Mode enabling.",
		}),
		DataDirSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "supervysor",
			Name:      "data_dir_size",
			Help:      "Size of data dir in --home dir.",
		}),
	}
	reg.MustRegister(m.PoolHeight, m.NodeHeight, m.MaxHeight, m.MinHeight, m.DataDirSize)
	return m
}

func StartMetricsServer(reg *prometheus.Registry) error {
	// Create metrics endpoint
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	http.Handle("/metrics", promHandler)
	err := http.ListenAndServe(":26660", nil)
	if err != nil {
		return err
	}
	return nil
}

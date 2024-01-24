package helpers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"cosmossdk.io/log"

	"github.com/spf13/viper"

	"github.com/KYVENetwork/supervysor/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cfg "github.com/tendermint/tendermint/config"
)

var logger = log.NewLogger(os.Stdout)

func CreateDestPath(backupDir string, latestHeight int64) (string, error) {
	if err := os.Mkdir(filepath.Join(backupDir, strconv.FormatInt(latestHeight, 10)), 0o755); err != nil {
		return "", fmt.Errorf("error creating backup directory: %v", err)
	}
	fmt.Println(filepath.Join(backupDir, strconv.FormatInt(latestHeight, 10)))
	if err := os.Mkdir(filepath.Join(backupDir, strconv.FormatInt(latestHeight, 10), "data"), 0o755); err != nil {
		return "", fmt.Errorf("error creating data backup directory: %v", err)
	}
	return filepath.Join(backupDir, strconv.FormatInt(latestHeight, 10), "data"), nil
}

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

func GetBackupDir() (string, error) {
	home, err := GetSupervysorDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %s", err)
	}

	backupDir := filepath.Join(home, "backups")
	if _, err = os.Stat(backupDir); os.IsNotExist(err) {
		err = os.Mkdir(backupDir, 0o755)
		if err != nil {
			return "", err
		}
	}

	return backupDir, nil
}

func GetHomePathFromBinary(binaryPath string) string {
	cmdPath, err := exec.LookPath(binaryPath)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to lookup binary path: %s", err.Error()))
		os.Exit(1)
	}

	startArgs := make([]string, 0)

	// if we run with cosmovisor we start with the cosmovisor run command
	if strings.HasSuffix(binaryPath, "cosmovisor") {
		startArgs = append(startArgs, "run")
	}

	out, err := exec.Command(cmdPath, startArgs...).Output()
	if err != nil {
		logger.Error("failed to get output of binary")
		os.Exit(1)
	}

	// here we search for a specific line in the binary output when simply
	// executed without arguments. In the output, the default home path
	// is printed, which is parsed and used by KSYNC
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "--home") {
			if strings.Count(line, "\"") != 2 {
				logger.Error(fmt.Sprintf("did not found default home path in help line: %s", line))
				os.Exit(1)
			}

			return strings.Split(line, "\"")[1]
		}
	}

	logger.Error("did not found default home path in entire binary output")
	os.Exit(1)
	return ""
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

func StartMetricsServer(reg *prometheus.Registry, port int) error {
	// Create metrics endpoint
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	http.Handle("/metrics", promHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		return err
	}
	return nil
}

func ValidatePaths(srcPath, destPath string) error {
	pathInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if !pathInfo.IsDir() {
		return err
	}
	pathInfo, err = os.Stat(destPath)
	if err != nil {
		return err
	}
	if !pathInfo.IsDir() {
		return err
	}

	return nil
}

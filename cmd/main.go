package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"bitbucket.org/publimonitor/pmwatcher/models"
	"github.com/coreos/pkg/capnslog"
	"github.com/sandovalrr/galley/config"
	"gopkg.in/yaml.v2"
)

var log = capnslog.NewPackageLogger("github.com/sandovalrr/goowatcher", "main")
var homeDir, _ = homedir.Dir()

func main() {
	flagLogLevel := flag.String("log-level", "info", "Define the logging level.")
	flagConfigPath := flag.String("config", "/etc/galley/config.yaml", "Load configuration from the specified file.")
	flag.Parse()

	config, err := config.Load(*flagConfigPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	// Initialize logging system
	logLevel, err := capnslog.ParseLevel(strings.ToUpper(*flagLogLevel))
	if err != nil {
		log.Errorf("Error on: %v", err)
	}
	capnslog.SetGlobalLogLevel(logLevel)
	capnslog.SetFormatter(capnslog.NewPrettyFormatter(os.Stdout, false))

}

// DefaultConfig is a configuration that can be used as a fallback value.
func DefaultConfig() models.Config {
	return models.Config{
		Watcher: models.WatcherConfig{
			Dir:  []string{homeDir + "/.galley/watcher"},
			Ext:  []string{"mp3", "mp4", "avi"},
			Wait: 3600,
		},
		Database: &models.DatabaseConfig{
			Port:     5432,
			Host:     "127.0.0.1",
			User:     "postgres",
			Password: "postgres",
			Database: "postgres",
		},
	}
}

// Load is a shortcut to open a file, read it, and generate a Config.
// It supports relative and absolute paths. Given "", it returns DefaultConfig.
func Load(path string) (config *models.Config, err error) {
	var cfgFile models.File

	cfgFile.Galley = DefaultConfig()
	if path == "" {
		return &cfgFile.Galley, nil
	}

	f, err := os.Open(os.ExpandEnv(path))
	if err != nil {
		return
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(d, &cfgFile)
	if err != nil {
		return
	}
	config = &cfgFile.Galley

	return
}

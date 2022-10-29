package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/jessevdk/go-flags"
	"github.com/lkolacz/OrderBook/rest"
	"github.com/lkolacz/OrderBook/rest/config"
	"go.uber.org/zap"
)

var (
	BuildType     = ""
	BuildVersion  = ""
	BuildDateTime = ""
)

func introductionText() {
	fmt.Printf("App service\nbuild mode type: %v\nbuild version: %v\nbuild datetime: %v\n\n", BuildType, BuildVersion, BuildDateTime)
}

type versionCmd struct{}

func (c *versionCmd) Execute([]string) error {
	return nil
}

type startCmd struct {
	ConfigFile string `short:"c" long:"config" description:"path to configuration file" default:"config.yaml"`
}

func (c *startCmd) Execute([]string) error {
	var cfg config.Config
	cfg.Default()

	err := cfg.Load(c.ConfigFile)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	setCtrlC()
	log := logger(&cfg.Logging)
	log.Info("Loaded config file from " + c.ConfigFile)

	ver := rest.DefaultVersion()
	if len(BuildType) > 0 {
		ver.BuildType = BuildType
	}
	if len(BuildVersion) > 0 {
		ver.BuildVersion = BuildVersion
	}
	if len(BuildDateTime) > 0 {
		ver.BuildDateTime = BuildDateTime
	}

	router := rest.NewQRouter(log, &cfg, ver)

	if err = router.Start(); err != nil {
		log.Error("HTTP Listener error", "err", err)
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	return nil
}

type initCmd struct {
	FileName string `short:"f" long:"file-name" description:"output configuration filename" default:"config.yaml"`
}

func (c *initCmd) Execute([]string) error {
	var cfg config.Config
	cfg.Default()
	if err := cfg.Save(c.FileName); err != nil {
		return err
	}

	fmt.Printf("written filename is %s\n\n", c.FileName)
	return nil
}

func main() {
	introductionText()

	var parser = flags.NewParser(nil, flags.Default)

	parser.AddCommand("init", "init config", "write default config", &initCmd{})
	parser.AddCommand("start", "start service", "start application service", &startCmd{})
	parser.AddCommand("version", "print version", "print service version and quit", &versionCmd{})

	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}

func logger(cfg *config.Logging) *zap.SugaredLogger {
	logConfig := zap.NewProductionConfig()

	if cfg.Format == "text" {
		logConfig = zap.NewDevelopmentConfig()
	}

	logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	if cfg.Level == "info" {
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logConfig.DisableStacktrace = true
	l, _ := logConfig.Build()

	return l.Sugar()
}

func setCtrlC() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			os.Exit(0)
		}
	}()
}

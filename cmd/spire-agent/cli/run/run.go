package run

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/hashicorp/hcl"
	"github.com/spiffe/spire/pkg/agent"
	"github.com/spiffe/spire/pkg/common/catalog"
	"github.com/spiffe/spire/pkg/common/log"
)

const (
	defaultConfigPath = "conf/agent/agent.conf"

	defaultSocketPath = "./spire_api"

	// TODO: Make my defaults sane
	defaultDataDir  = "."
	defaultLogLevel = "INFO"
	defaultUmask    = 0077
)

// RunConfig represents the available configurables for file
// and CLI options
type runConfig struct {
	AgentConfig   agentConfig             `hcl:"agent"`
	PluginConfigs catalog.PluginConfigMap `hcl:"plugins"`
}

type agentConfig struct {
	ServerAddress   string `hcl:"server_address"`
	ServerPort      int    `hcl:"server_port"`
	TrustDomain     string `hcl:"trust_domain"`
	TrustBundlePath string `hcl:"trust_bundle_path"`
	JoinToken       string `hcl:"join_token"`

	SocketPath string `hcl:"socket_path"`
	DataDir    string `hcl:"data_dir"`
	LogFile    string `hcl:"log_file"`
	LogLevel   string `hcl:"log_level"`

	ConfigPath string
	Umask      string `hcl:"umask"`

	ProfilingEnabled string   `hcl:"profiling_enabled"`
	ProfilingPort    string   `hcl:"profiling_port"`
	ProfilingFreq    string   `hcl:"profiling_freq"`
	ProfilingNames   []string `hcl:"profiling_names"`
}

type RunCLI struct {
}

func (*RunCLI) Help() string {
	_, err := parseFlags([]string{"-h"})
	return err.Error()
}

func (*RunCLI) Run(args []string) int {
	cliConfig, err := parseFlags(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	fileConfig, err := parseFile(cliConfig.AgentConfig.ConfigPath)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	c := newDefaultConfig()

	// Get the plugin configurations from the file
	c.PluginConfigs = fileConfig.PluginConfigs

	err = mergeConfigs(c, fileConfig, cliConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = validateConfig(c)
	if err != nil {
		fmt.Println(err.Error())
	}

	agt := agent.New(c)
	signalListener(agt)

	err = agt.Run()
	if err != nil {
		c.Log.Errorf("agent crashed: %v", err)
		return 1
	}

	c.Log.Infof("Agent stopped gracefully")
	return 0
}

func (*RunCLI) Synopsis() string {
	return "Runs the agent"
}

func parseFile(filePath string) (*runConfig, error) {
	c := &runConfig{}

	// Return a friendly error if the file is missing
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		msg := "could not find config file %s: please use the -config flag"
		p, err := filepath.Abs(filePath)
		if err != nil {
			p = filePath
			msg = "could not determine CWD; config file not found at %s: use -config"
		}
		return nil, fmt.Errorf(msg, p)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	hclTree, err := hcl.Parse(string(data))
	if err != nil {
		return nil, err
	}
	if err := hcl.DecodeObject(&c, hclTree); err != nil {
		return nil, err
	}

	return c, nil
}

func parseFlags(args []string) (*runConfig, error) {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	c := &runConfig{}

	flags.StringVar(&c.AgentConfig.ServerAddress, "serverAddress", "", "IP address or DNS name of the SPIRE server")
	flags.IntVar(&c.AgentConfig.ServerPort, "serverPort", 0, "Port number of the SPIRE server")
	flags.StringVar(&c.AgentConfig.TrustDomain, "trustDomain", "", "The trust domain that this agent belongs to")
	flags.StringVar(&c.AgentConfig.TrustBundlePath, "trustBundle", "", "Path to the SPIRE server CA bundle")
	flags.StringVar(&c.AgentConfig.JoinToken, "joinToken", "", "An optional token which has been generated by the SPIRE server")
	flags.StringVar(&c.AgentConfig.SocketPath, "socketPath", "", "Location to bind the workload API socket")
	flags.StringVar(&c.AgentConfig.DataDir, "dataDir", "", "A directory the agent can use for its runtime data")
	flags.StringVar(&c.AgentConfig.LogFile, "logFile", "", "File to write logs to")
	flags.StringVar(&c.AgentConfig.LogLevel, "logLevel", "", "DEBUG, INFO, WARN or ERROR")

	flags.StringVar(&c.AgentConfig.ConfigPath, "config", defaultConfigPath, "Path to a SPIRE config file")
	flags.StringVar(&c.AgentConfig.Umask, "umask", "", "Umask value to use for new files")

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func mergeConfigs(c *agent.Config, fileConfig, cliConfig *runConfig) error {
	// CLI > File, merge fileConfig first
	err := mergeConfig(c, fileConfig)
	if err != nil {
		return err
	}

	return mergeConfig(c, cliConfig)
}

func mergeConfig(orig *agent.Config, cmd *runConfig) error {
	// Parse server address
	if cmd.AgentConfig.ServerAddress != "" {
		ips, err := net.LookupIP(cmd.AgentConfig.ServerAddress)
		if err != nil {
			return err
		}

		if len(ips) == 0 {
			return fmt.Errorf("Could not resolve ServerAddress %s", cmd.AgentConfig.ServerAddress)
		}
		serverAddress := ips[0]

		orig.ServerAddress.IP = serverAddress
	}

	if cmd.AgentConfig.ServerPort != 0 {
		orig.ServerAddress.Port = cmd.AgentConfig.ServerPort
	}

	if cmd.AgentConfig.TrustDomain != "" {
		trustDomain := url.URL{
			Scheme: "spiffe",
			Host:   cmd.AgentConfig.TrustDomain,
		}

		orig.TrustDomain = trustDomain
	}

	// Parse trust bundle
	if cmd.AgentConfig.TrustBundlePath != "" {
		bundle, err := parseTrustBundle(cmd.AgentConfig.TrustBundlePath)
		if err != nil {
			return fmt.Errorf("Error parsing trust bundle: %s", err)
		}

		orig.TrustBundle = bundle
	}

	if cmd.AgentConfig.JoinToken != "" {
		orig.JoinToken = cmd.AgentConfig.JoinToken
	}

	if cmd.AgentConfig.SocketPath != "" {
		orig.BindAddress.Name = cmd.AgentConfig.SocketPath
	}

	if cmd.AgentConfig.DataDir != "" {
		orig.DataDir = cmd.AgentConfig.DataDir
	}

	// Handle log file and level
	if cmd.AgentConfig.LogFile != "" || cmd.AgentConfig.LogLevel != "" {
		logLevel := defaultLogLevel
		if cmd.AgentConfig.LogLevel != "" {
			logLevel = cmd.AgentConfig.LogLevel
		}

		logger, err := log.NewLogger(logLevel, cmd.AgentConfig.LogFile)
		if err != nil {
			return fmt.Errorf("Could not open log file %s: %s", cmd.AgentConfig.LogFile, err)
		}

		orig.Log = logger
	}

	if cmd.AgentConfig.Umask != "" {
		umask, err := strconv.ParseInt(cmd.AgentConfig.Umask, 0, 0)
		if err != nil {
			return fmt.Errorf("Could not parse umask %s: %s", cmd.AgentConfig.Umask, err)
		}
		orig.Umask = int(umask)
	}

	if cmd.AgentConfig.ProfilingEnabled != "" {
		value, err := strconv.ParseBool(cmd.AgentConfig.ProfilingEnabled)
		if err != nil {
			return fmt.Errorf("Could not parse profiling_enabled %s: %s", cmd.AgentConfig.ProfilingEnabled, err)
		}
		orig.ProfilingEnabled = value
	}

	if orig.ProfilingEnabled {
		if cmd.AgentConfig.ProfilingPort != "" {
			value, err := strconv.ParseInt(cmd.AgentConfig.ProfilingPort, 0, 0)
			if err != nil {
				if orig.Log != nil {
					orig.Log.Warnf("Could not parse profiling_port %s: %s. pprof web server would not be run", cmd.AgentConfig.ProfilingPort, err)
				}
			} else {
				orig.ProfilingPort = int(value)
			}
		}

		if cmd.AgentConfig.ProfilingFreq != "" {
			value, err := strconv.ParseInt(cmd.AgentConfig.ProfilingFreq, 0, 0)
			if err != nil {
				if orig.Log != nil {
					orig.Log.Warnf("Could not parse profiling_freq %s: %s. Profiling data would not be generated", cmd.AgentConfig.ProfilingFreq, err)
				}
			} else {
				orig.ProfilingFreq = int(value)
			}
		}

		if len(cmd.AgentConfig.ProfilingNames) > 0 {
			orig.ProfilingNames = cmd.AgentConfig.ProfilingNames
		}
	}
	return nil
}

func validateConfig(c *agent.Config) error {
	if c.ServerAddress.IP == nil || c.ServerAddress.Port == 0 {
		return errors.New("ServerAddress and ServerPort are required")
	}

	if c.TrustDomain.String() == "" {
		return errors.New("TrustDomain is required")
	}

	if c.TrustBundle == nil {
		return errors.New("TrustBundle is required")
	}

	return nil
}

func newDefaultConfig() *agent.Config {
	bindAddr := &net.UnixAddr{Name: defaultSocketPath, Net: "unix"}

	// log.NewLogger() cannot return error when using STDOUT
	logger, _ := log.NewLogger(defaultLogLevel, "")
	serverAddress := &net.TCPAddr{}

	return &agent.Config{
		BindAddress:   bindAddr,
		DataDir:       defaultDataDir,
		Log:           logger,
		ServerAddress: serverAddress,
		Umask:         defaultUmask,
	}
}

func parseTrustBundle(path string) ([]*x509.Certificate, error) {
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data []byte
	for len(pemData) > 1 {
		var block *pem.Block
		block, pemData = pem.Decode(pemData)
		if block == nil && len(data) < 1 {
			return nil, errors.New("no certificates found")
		}

		if block == nil {
			return nil, errors.New("encountered unknown data in trust bundle")
		}

		if block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("non-certificate type %v found in trust bundle", block.Type)
		}

		data = append(data, block.Bytes...)
	}

	bundle, err := x509.ParseCertificates(data)
	if err != nil {
		return nil, fmt.Errorf("parse certificates from %v, %v", path, err)
	}

	return bundle, nil
}

func stringDefault(option string, defaultValue string) string {
	if option == "" {
		return defaultValue
	}

	return option
}

func signalListener(agt *agent.Agent) {
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalCh:
			agt.Shutdown()
		}
	}()
	return
}

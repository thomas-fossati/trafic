package config

// CommonConfig implements the Configurer interface for configuration keys that
// are common to iperf3 clients and servers
type CommonConfig struct {
	BindInterface  string  `yaml:"bind"`
	Debug          bool    `yaml:"debug"`
	ExchangedFile  string  `yaml:"file"`
	ForceFlush     bool    `yaml:"forceflush"`
	JSON           bool    `yaml:"json"`
	LogFile        string  `yaml:"logfile"`
	ReportFormat   string  `yaml:"format"`
	ReportInterval float32 `yaml:"report-interval-s"`
	ServerPort     uint16  `yaml:"server-port"`
	Verbose        bool    `yaml:"verbose"`
}

func (cfg *CommonConfig) ToArgs(args []string) ([]string, error) {
	args = AppendKeyVal(args, "--bind", cfg.BindInterface)
	args = AppendKey(args, "--debug", cfg.Debug)
	args = AppendKeyVal(args, "--file", cfg.ExchangedFile)
	args = AppendKey(args, "--forceflush", cfg.ForceFlush)
	args = AppendKey(args, "--json", cfg.JSON)
	args = AppendKeyVal(args, "--logfile", cfg.LogFile)
	args = AppendKeyVal(args, "--format", cfg.ReportFormat)
	args = AppendKeyVal(args, "--interval", cfg.ReportInterval)
	args = AppendKeyVal(args, "--port", cfg.ServerPort)
	args = AppendKey(args, "--verbose", cfg.Verbose)

	return args, nil
}

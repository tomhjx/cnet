package options

import (
	"encoding/json"
	"time"

	"github.com/spf13/pflag"
	"github.com/tomhjx/cnet/pkg/config"
	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/field"
	"github.com/tomhjx/cnet/pkg/flow"
	"github.com/tomhjx/cnet/pkg/metric"
	"github.com/tomhjx/xlog"
)

type ConfigDataItem struct {
	ClientID        string `mapstructure:"client_id"`
	ADDR            string
	Queue           string
	QueueExchange   string
	QueueRoutingKey string
	Method          string
	Data            string
	Interval        time.Duration
	Includes        []string
	Sinks           []flow.SinkConfig
}

type ConfigData struct {
	ConfigDataItem `mapstructure:",squash"`
	Items          []ConfigDataItem
}

type Config struct {
	Content            ConfigData
	ConfigPath         string
	ConfigPollInterval time.Duration
	Version            bool
	OnLoaded           func(*Config)
	LogFile            string
	Verbosity          int
	MetricsServerPort  int
	ProfileServerPort  int
	LogLevel           string
}

func NewConfig() *Config {
	c := &Config{
		Content: ConfigData{
			ConfigDataItem: ConfigDataItem{
				Interval: 10 * time.Second,
				Sinks:    []flow.SinkConfig{},
			},
		},
	}
	return c
}

func (c *Config) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&c.Content.ClientID, "cid", c.Content.ClientID, "Used to differentiate reporting clients")
	flags.StringVar(&c.Content.ADDR, "addr", c.Content.ADDR, "Address to work with")
	flags.StringVarP(&c.ConfigPath, "config", "K", "", "Specify which config file to read")
	flags.DurationVar(&c.ConfigPollInterval, "config-poll-interval", time.Minute, "Control the interval duration for automatic polling of remote configuration files.")
	flags.StringVarP(&c.Content.Method, "request", "X", c.Content.Method, "Specify request command to use")
	flags.StringVarP(&c.Content.Data, "data", "d", "", "HTTP POST data (H)")
	flags.DurationVar(&c.Content.Interval, "interval", c.Content.Interval, "Control the interval duration between each request.")
	flags.BoolVarP(&c.Version, "version", "V", c.Version, "Show version number and quit")
	flags.StringArrayVarP(&c.Content.Includes, "include", "i", c.Content.Includes, "Include protocol fields (header,body) in the output")
	flags.StringVar(&c.LogFile, "log-file", c.LogFile, "Logging output file path")
	flags.IntVarP(&c.Verbosity, "verbosity", "v", c.Verbosity, "Number for the log level verbosity")
	flags.IntVar(&c.MetricsServerPort, "metrics-server-port", c.MetricsServerPort, "Prometheus metrics http server port")
	flags.IntVar(&c.ProfileServerPort, "profile-server-port", c.ProfileServerPort, "Profile http server port")
	flags.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Name for the log severity level:info,warning,error,fatal.")

	sinks := []string{}
	flags.StringArrayVar(&sinks, "sink", sinks, "SINK to work with")
	for _, v := range sinks {
		vv := flow.SinkConfig{}
		json.Unmarshal([]byte(v), &vv)
		c.Content.Sinks = append(c.Content.Sinks, vv)
	}
}

func (c *Config) Init(onInited func(), onLoaded func(ConfigData)) {
	xlog.SetVerbosity(c.Verbosity)
	xlog.SetSeverityName(c.LogLevel)
	xlog.SetFile(c.LogFile)
	if c.MetricsServerPort <= 0 {
		metric.Disable()
	}
	onInited()

	if c.ConfigPath == "" {
		return
	}
	b, err := config.NewBuilder(config.InputOption{Path: c.ConfigPath, PollInterval: c.ConfigPollInterval}, c.Content)
	if err != nil {
		xlog.Fatalln(err)
	}
	b.OnLoad(func(d any) {
		onLoaded(d.(ConfigData))
	})
	if err := b.Load(); err != nil {
		xlog.Fatalln(err)
	}
	b.Watch()
}

func (c ConfigData) Complete() []*CompletedConfig {

	cls := []*CompletedConfig{}
	sinks := c.Sinks
	if len(sinks) == 0 {
		sinks = append(c.Sinks, flow.SinkConfig{Name: flow.StdOutSinkName})
	}

	for _, v := range c.Items {
		cc := &CompletedConfig{ConfigDataItem: v, IncludeFields: []field.Field{}}
		if len(cc.Includes) == 0 {
			cc.Includes = c.Includes
		}
		if len(v.Sinks) == 0 {
			cc.Sinks = sinks
		}
		if cc.ClientID == "" {
			cc.ClientID = c.ClientID
		}
		if cc.Interval == 0 {
			cc.Interval = c.Interval
		}
		if cc.Method == "" {
			cc.Method = c.Method
		}
		if cc.ADDR == "" {
			cc.ADDR = c.ADDR
		}
		if cc.QueueRoutingKey == "" {
			cc.QueueRoutingKey = cc.Queue
		}
		for _, f := range cc.Includes {
			cc.IncludeFields = append(cc.IncludeFields, field.Field(f))
		}
		cls = append(cls, cc)
	}

	return cls
}

type CompletedConfig struct {
	ConfigDataItem
	IncludeFields []field.Field
}

func (c *CompletedConfig) CreateRequest() *core.Request {
	req := &core.Request{}
	req.ClientID = c.ClientID
	req.ADDR = c.ADDR
	req.Method = c.Method
	req.Queue = c.Queue
	req.QueueExchange = c.QueueExchange
	req.QueueRoutingKey = c.QueueRoutingKey
	return req
}

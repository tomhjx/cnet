package options

import (
	"encoding/json"
	"log"
	"time"

	"github.com/spf13/pflag"
	"github.com/tomhjx/cnet/pkg/config"
	"github.com/tomhjx/cnet/pkg/core"
	"github.com/tomhjx/cnet/pkg/field"
	"github.com/tomhjx/cnet/pkg/flow"
)

type ConfigDataItem struct {
	ClientID string `mapstructure:"client_id"`
	URL      string
	Method   string
	Data     string
	Interval time.Duration
	Includes []string
	Sinks    []flow.SinkConfig
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
	flags.StringVar(&c.Content.URL, "url", c.Content.URL, "URL to work with")
	flags.StringVarP(&c.ConfigPath, "config", "K", "", "Specify which config file to read")
	flags.DurationVar(&c.ConfigPollInterval, "config-poll-interval", time.Minute, "Control the interval duration for automatic polling of remote configuration files.")
	flags.StringVarP(&c.Content.Method, "request", "X", c.Content.Method, "Specify request command to use")
	flags.StringVarP(&c.Content.Data, "data", "d", "", "HTTP POST data (H)")
	flags.DurationVar(&c.Content.Interval, "interval", c.Content.Interval, "Control the interval duration between each request.")
	flags.BoolVarP(&c.Version, "version", "V", c.Version, "Show version number and quit")
	flags.StringArrayVarP(&c.Content.Includes, "include", "i", c.Content.Includes, "Include protocol fields (header,body) in the output")

	sinks := []string{}
	flags.StringArrayVar(&sinks, "sink", sinks, "SINK to work with")
	for _, v := range sinks {
		vv := flow.SinkConfig{}
		json.Unmarshal([]byte(v), &vv)
		c.Content.Sinks = append(c.Content.Sinks, vv)
	}
}

func (c *Config) Init(onLoaded func(ConfigData)) {
	if c.ConfigPath == "" {
		return
	}
	b, err := config.NewBuilder(config.InputOption{Path: c.ConfigPath, PollInterval: c.ConfigPollInterval}, c.Content)
	if err != nil {
		log.Fatalln(err)
	}
	b.OnLoad(func(d any) {
		onLoaded(d.(ConfigData))
	})
	b.Load()
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
		if cc.URL == "" {
			cc.URL = c.URL
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
	req.RawURL = c.URL
	req.Method = c.Method
	return req
}

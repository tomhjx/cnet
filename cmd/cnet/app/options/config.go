package options

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	Content      ConfigData
	ConfigPath   string
	Version      bool
	loadInitOnce sync.Once
	OnChanged    func(*Config)
	changed      chan bool
}

func NewConfig() *Config {
	c := &Config{
		Content: ConfigData{
			ConfigDataItem: ConfigDataItem{
				Interval: 10 * time.Second,
				Sinks:    []flow.SinkConfig{},
			},
		},
		changed: make(chan bool),
	}
	go func() {
		for {
			<-c.changed
			fmt.Println("config changed.")
			c.OnChanged(c)
		}
	}()
	return c
}

func (c *Config) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&c.Content.ClientID, "cid", c.Content.ClientID, "Used to differentiate reporting clients")
	flags.StringVar(&c.Content.URL, "url", c.Content.URL, "URL to work with")
	flags.StringVarP(&c.ConfigPath, "config", "K", "", "Specify which config file to read")
	flags.StringVarP(&c.Content.Method, "request", "X", c.Content.Method, "Specify request command to use")
	flags.StringVarP(&c.Content.Data, "data", "d", "", "HTTP POST data (H)")
	flags.DurationVar(&c.Content.Interval, "interval", c.Content.Interval, "Make a request every N seconds, where the configuration declares N.")
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

func (c *Config) Load() {
	if c.ConfigPath == "" {
		return
	}

	c.loadInitOnce.Do(func() {
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
			c.changed <- true
		})
		viper.SetConfigFile(c.ConfigPath)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Read config [%s] failed, %s", c.ConfigPath, err)
		}
		viper.WatchConfig()
	})

	log.Println("load.Unmarshal.start")
	if err := viper.Unmarshal(&c.Content); err != nil {
		log.Fatalf("Mapping config [%s] failed, %s", c.ConfigPath, err)
	}
	log.Println("load.Unmarshal.value:", c.Content)

}

func (c *Config) Complete() []*CompletedConfig {
	log.Println("config.Complete.start")
	c.Load()

	cls := []*CompletedConfig{}
	if len(c.Content.Sinks) == 0 {
		c.Content.Sinks = append(c.Content.Sinks, flow.SinkConfig{Name: flow.StdOutSinkName})
	}

	log.Println("config.Complete.root.value:", c.Content)
	for _, v := range c.Content.Items {
		cc := &CompletedConfig{ConfigDataItem: v, IncludeFields: []field.Field{}}
		if len(cc.Includes) == 0 {
			cc.Includes = c.Content.Includes
		}
		if len(v.Sinks) == 0 {
			cc.Sinks = c.Content.Sinks
		}
		if cc.ClientID == "" {
			cc.ClientID = c.Content.ClientID
		}
		if cc.Interval == 0 {
			cc.Interval = c.Content.Interval
		}
		if cc.Method == "" {
			cc.Method = c.Content.Method
		}
		if cc.URL == "" {
			cc.URL = c.Content.URL
		}
		for _, f := range cc.Includes {
			cc.IncludeFields = append(cc.IncludeFields, field.Field(f))
		}
		log.Println("config.Complete.item.value.ClientID:", cc.ClientID)
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

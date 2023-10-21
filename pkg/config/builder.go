package config

import (
	"bytes"
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Builder struct {
	input        Inputer
	out          any
	raw          Content
	watchOnce    sync.Once
	viper        *viper.Viper
	onLoadedFunc func(any)
}

func NewBuilder(inputOption InputOption, out any) (*Builder, error) {
	input, err := inputOption.NewInputer()
	if err != nil {
		return nil, err
	}
	v := viper.New()
	v.SetConfigType("json")

	b := &Builder{
		input: input,
		out:   out,
		viper: v,
	}
	return b, nil
}

func (b *Builder) OnLoad(f func(any)) {
	b.onLoadedFunc = f
}

func (b *Builder) Load() error {
	c, err := b.input.Read()
	if err != nil {
		return err
	}
	if c.ModTime.Before(b.raw.ModTime) {
		log.Println("Old configuration.")
		return nil
	}
	if bytes.Equal(c.Body, b.raw.Body) {
		log.Println("No changes in configuration.")
		return nil
	}
	if err := b.viper.ReadConfig(bytes.NewReader(c.Body)); err != nil {
		return err
	}
	o := b.out
	if err := b.viper.Unmarshal(&o); err != nil {
		return err
	}
	b.raw = c
	b.onLoadedFunc(o)
	return nil
}

func (b *Builder) Watch() error {
	b.watchOnce.Do(func() {
		go func() {
			defer b.input.Close()
			b.input.Watch(func(i Inputer) {
				b.Load()
			})
		}()
	})
	return nil
}

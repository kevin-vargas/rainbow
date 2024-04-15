package wiz

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type Params struct {
	Color
	PhoneMac string        `json:"phoneMac,omitempty"`
	Register *bool         `json:"register,omitempty"`
	PhoneIP  string        `json:"phoneIp,omitempty"`
	ID       string        `json:"id,omitempty"`
	Duration time.Duration `json:"-"`
	Dimming  uint          `json:"dimming,omitempty"`
	SceneID  uint          `json:"sceneId,omitempty"`
	Delta    int           `json:"delta,omitempty"`
	State    *bool         `json:"state,omitempty"`
}

type Result struct {
	Mac     string `json:"mac"`
	Success bool   `json:"success"`
}

type Message struct {
	Method string  `json:"method"`
	Params *Params `json:"params,omitempty"`
	Env    string  `json:"env,omitempty"`
	Result *Result `json:"result,omitempty"`
}

type Device struct {
	s    *sync.Mutex
	Addr string
	Mac  string
}

type Color struct {
	Red   uint8 `json:"r,omitempty"`
	Blue  uint8 `json:"b,omitempty"`
	Green uint8 `json:"g,omitempty"`
}

func (d *Device) Pulse(t time.Duration, ops ...Option) error {
	err := d.Set(ops...)
	if err != nil {
		return err
	}
	time.Sleep(t)
	err = d.Set(WithState(false))
	if err != nil {
		return nil
	}
	return nil
}

type Option func(*Message)

func WithColor(c Color) Option {
	return func(m *Message) {
		if m.Params == nil {
			m.Params = &Params{}
		}
		m.Params.Color = c
	}
}

func WithDimming(i uint) Option {
	return func(m *Message) {
		if m.Params == nil {
			m.Params = &Params{}
		}
		m.Params.Dimming = i
	}
}

func WithSceneID(id uint) Option {
	return func(m *Message) {
		if m.Params == nil {
			m.Params = &Params{}
		}
		m.Params.SceneID = id
	}
}

func WithDuration(t time.Duration) Option {
	return func(m *Message) {
		m.Params.Duration = t
	}
}

func WithState(s bool) Option {
	return func(m *Message) {
		if m.Params == nil {
			m.Params = &Params{}
		}
		m.Params.State = &s
	}
}

func (d *Device) Set(ops ...Option) error {
	msg := Message{
		Env:    "pro",
		Method: "setPilot",
	}
	for _, o := range ops {
		o(&msg)
	}
	s, err := net.ResolveUDPAddr("udp4", d.Addr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		return err
	}
	defer conn.Close()
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	d.s.Lock()
	defer d.s.Unlock()
	fmt.Println(string(raw))
	_, err = conn.Write(raw)
	if err != nil {
		return err
	}
	return nil
}

func newDevice(addr string, mac string) *Device {
	return &Device{
		Addr: addr,
		Mac:  mac,
		s:    &sync.Mutex{},
	}
}

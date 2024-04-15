package wiz

import (
	"rainbow/starter"
)

func New(addr, mac string, messages <-chan []Option) starter.Starter {
	d := newDevice(addr, mac)

	return &manager{
		d:  d,
		in: messages,
	}
}

type manager struct {
	d  *Device
	in <-chan []Option
}

func (m *manager) Start(args ...interface{}) (err error) {
	for ops := range m.in {
		err := m.d.Set(ops...)
		if err != nil {
			return err
		}
	}
	return nil
}

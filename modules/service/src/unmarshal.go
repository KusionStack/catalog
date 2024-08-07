package main

import (
	"encoding/json"
	"errors"
)

// UnmarshalJSON implements the json.Unmarshaller interface for ProbeHandler.
func (p *ProbeHandler) UnmarshalJSON(data []byte) error {
	var probeType TypeWrapper
	err := json.Unmarshal(data, &probeType)
	if err != nil {
		return err
	}

	p.Type = probeType.Type
	switch p.Type {
	case TypeHTTP:
		handler := &HTTPGetAction{}
		err = json.Unmarshal(data, handler)
		p.HTTPGetAction = handler
	case TypeExec:
		handler := &ExecAction{}
		err = json.Unmarshal(data, handler)
		p.ExecAction = handler
	case TypeTCP:
		handler := &TCPSocketAction{}
		err = json.Unmarshal(data, handler)
		p.TCPSocketAction = handler
	default:
		return errors.New("unrecognized probe handler type")
	}

	return err
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for ProbeHandler.
func (p *ProbeHandler) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var probeType TypeWrapper
	err := unmarshal(&probeType)
	if err != nil {
		return err
	}

	p.Type = probeType.Type
	switch p.Type {
	case TypeHTTP:
		handler := &HTTPGetAction{}
		err = unmarshal(handler)
		p.HTTPGetAction = handler
	case TypeExec:
		handler := &ExecAction{}
		err = unmarshal(handler)
		p.ExecAction = handler
	case TypeTCP:
		handler := &TCPSocketAction{}
		err = unmarshal(handler)
		p.TCPSocketAction = handler
	default:
		return errors.New("unrecognized probe handler type")
	}

	return err
}

// UnmarshalJSON implements the json.Unmarshaller interface for LifecycleHandler.
func (l *LifecycleHandler) UnmarshalJSON(data []byte) error {
	var handlerType TypeWrapper
	err := json.Unmarshal(data, &handlerType)
	if err != nil {
		return err
	}

	l.Type = handlerType.Type
	switch l.Type {
	case TypeHTTP:
		handler := &HTTPGetAction{}
		err = json.Unmarshal(data, handler)
		l.HTTPGetAction = handler
	case TypeExec:
		handler := &ExecAction{}
		err = json.Unmarshal(data, handler)
		l.ExecAction = handler
	default:
		return errors.New("unrecognized lifecycle handler type")
	}

	return err
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for LifecycleHandler.
func (l *LifecycleHandler) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var handlerType TypeWrapper
	err := unmarshal(&handlerType)
	if err != nil {
		return err
	}

	l.Type = handlerType.Type
	switch l.Type {
	case TypeHTTP:
		handler := &HTTPGetAction{}
		err = unmarshal(handler)
		l.HTTPGetAction = handler
	case TypeExec:
		handler := &ExecAction{}
		err = unmarshal(handler)
		l.ExecAction = handler
	default:
		return errors.New("unrecognized lifecycle handler type")
	}

	return err
}

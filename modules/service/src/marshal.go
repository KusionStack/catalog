package main

import (
	"encoding/json"
	"errors"
)

// MarshalJSON implements the json.Marshaler interface for ProbeHandler.
func (p *ProbeHandler) MarshalJSON() ([]byte, error) {
	switch p.Type {
	case TypeHTTP:
		return json.Marshal(struct {
			TypeWrapper    `json:",inline"`
			*HTTPGetAction `json:",inline"`
		}{
			TypeWrapper:   TypeWrapper{p.Type},
			HTTPGetAction: p.HTTPGetAction,
		})
	case TypeExec:
		return json.Marshal(struct {
			TypeWrapper `json:",inline"`
			*ExecAction `json:",inline"`
		}{
			TypeWrapper: TypeWrapper{p.Type},
			ExecAction:  p.ExecAction,
		})
	case TypeTCP:
		return json.Marshal(struct {
			TypeWrapper      `json:",inline"`
			*TCPSocketAction `json:",inline"`
		}{
			TypeWrapper:     TypeWrapper{p.Type},
			TCPSocketAction: p.TCPSocketAction,
		})
	default:
		return nil, errors.New("unrecognized probe handler type")
	}
}

// MarshalYAML implements the yaml.Marshaler interface for ProbeHandler.
func (p *ProbeHandler) MarshalYAML() (interface{}, error) {
	switch p.Type {
	case TypeHTTP:
		return struct {
			TypeWrapper   `yaml:",inline" json:",inline"`
			HTTPGetAction `yaml:",inline" json:",inline"`
		}{
			TypeWrapper:   TypeWrapper{Type: p.Type},
			HTTPGetAction: *p.HTTPGetAction,
		}, nil
	case TypeExec:
		return struct {
			TypeWrapper `yaml:",inline" json:",inline"`
			ExecAction  `yaml:",inline" json:",inline"`
		}{
			TypeWrapper: TypeWrapper{Type: p.Type},
			ExecAction:  *p.ExecAction,
		}, nil
	case TypeTCP:
		return struct {
			TypeWrapper     `yaml:",inline" json:",inline"`
			TCPSocketAction `yaml:",inline" json:",inline"`
		}{
			TypeWrapper:     TypeWrapper{Type: p.Type},
			TCPSocketAction: *p.TCPSocketAction,
		}, nil
	}

	return nil, nil
}

// MarshalJSON implements the json.Marshaler interface for LifecycleHandler.
func (l *LifecycleHandler) MarshalJSON() ([]byte, error) {
	switch l.Type {
	case TypeHTTP:
		return json.Marshal(struct {
			TypeWrapper    `json:",inline"`
			*HTTPGetAction `json:",inline"`
		}{
			TypeWrapper:   TypeWrapper{l.Type},
			HTTPGetAction: l.HTTPGetAction,
		})
	case TypeExec:
		return json.Marshal(struct {
			TypeWrapper `json:",inline"`
			*ExecAction `json:",inline"`
		}{
			TypeWrapper: TypeWrapper{l.Type},
			ExecAction:  l.ExecAction,
		})
	default:
		return nil, errors.New("unrecognized lifecycle handler type")
	}
}

// MarshalYAML implements the yaml.Marshaler interface for LifecycleHandler.
func (l *LifecycleHandler) MarshalYAML() (interface{}, error) {
	switch l.Type {
	case TypeHTTP:
		return struct {
			TypeWrapper   `yaml:",inline" json:",inline"`
			HTTPGetAction `yaml:",inline" json:",inline"`
		}{
			TypeWrapper:   TypeWrapper{Type: l.Type},
			HTTPGetAction: *l.HTTPGetAction,
		}, nil
	case TypeExec:
		return struct {
			TypeWrapper `yaml:",inline" json:",inline"`
			ExecAction  `yaml:",inline" json:",inline"`
		}{
			TypeWrapper: TypeWrapper{Type: l.Type},
			ExecAction:  *l.ExecAction,
		}, nil
	default:
		return nil, errors.New("unrecognized lifecycle handler type")
	}
}

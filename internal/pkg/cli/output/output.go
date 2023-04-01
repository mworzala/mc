package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type FormatType string

const (
	Default  FormatType = "default"
	JSON     FormatType = "json"
	YAML     FormatType = "yaml"
	Template FormatType = "template"
)

type Format struct {
	Type  FormatType
	Value string
}

var ErrInvalidFormat = errors.New("invalid output format")

// All commands get a cli.App instance in their new command function

func ParseFormat(s string) (Format, error) {
	eqChar := strings.IndexRune(s, '=')
	if eqChar == -1 {
		switch strings.ToLower(s) {
		case string(Default), "":
			return Format{Default, ""}, nil
		case string(JSON):
			return Format{JSON, ""}, nil
		case string(YAML):
			return Format{YAML, ""}, nil
		case string(Template):
			return Format{Template, "missing template"}, nil
		default:
			return Format{}, ErrInvalidFormat
		}
	}

	// Attempt to parse as template
	name := s[0:eqChar]
	if name != "template" {
		return Format{}, ErrInvalidFormat
	}
	format := strings.Trim(s[eqChar+1:], `"`)
	return Format{Template, format}, nil
}

func (f Format) Write(out io.Writer, entry interface{}) (err error) {
	switch f.Type {
	case Default:
		stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
		if reflect.TypeOf(entry).Implements(stringerType) {
			_, err = out.Write([]byte(entry.(fmt.Stringer).String()))
		} else {
			_, err = out.Write([]byte(fmt.Sprintf("%#v", entry)))
		}
	case JSON:
		if err = json.NewEncoder(out).Encode(entry); err != nil {
			err = fmt.Errorf("unable to write JSON output: %w", err)
		}
	case YAML:
		if err = yaml.NewEncoder(out).Encode(entry); err != nil {
			err = fmt.Errorf("unable to write YAML output: %w", err)
		}
	case Template:
		t, err := template.New("_").Parse(f.Value)
		if err != nil {
			return fmt.Errorf("unable to parse template: %w", err)
		}
		return t.Execute(out, entry)
	default:
		return ErrInvalidFormat
	}
	return
}

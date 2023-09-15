package rule

import (
	"encoding/json"
	"fmt"
	"regexp"
	"runtime"

	"github.com/mworzala/mc/internal/pkg/platform"
)

//todo test this package a bunch, i don't feel super confident

type Action bool

const (
	Allow Action = true
	Deny  Action = false
)

func (a *Action) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*a = value == "allow"
	return nil
}

func (a *Action) MarshalJSON() ([]byte, error) {
	value := "deny"
	if *a {
		value = "allow"
	}
	return json.Marshal(value)
}

type Rule struct {
	Action Action `json:"action" mapstructure:"action"`
	OS     struct {
		Name    string `json:"name" mapstructure:"name"`
		Arch    string `json:"arch" mapstructure:"arch"`
		Version string `json:"version" mapstructure:"version"`
	}
	Features map[string]bool `json:"features" mapstructure:"features"`
}

// Evaluator handles evaluating vanilla rule entries found in libraries and
// argument validation for installation and launch.
type Evaluator struct {
	os       string
	arch     string
	version  string
	features map[string]bool
}

func NewEvaluator(features ...string) *Evaluator {
	featureMap := map[string]bool{}
	for _, feat := range features {
		featureMap[feat] = true
	}

	return &Evaluator{
		os:       determineOS(),
		arch:     determineArch(),
		version:  determineVersion(),
		features: featureMap,
	}
}

func (e *Evaluator) Eval(rules []*Rule) (action Action) {
	if len(rules) == 0 {
		return Allow
	}

	// If any rules evaluate to a deny, return deny
	for _, rule := range rules {
		if e.evalSingle(rule) == Deny {
			return Deny
		}
	}

	return Allow
}

func (e *Evaluator) evalSingle(rule *Rule) Action {
	action := rule.Action

	// If any rules do not match, return the opposite of the given action
	if expected := rule.OS.Name; expected != "" && e.os != expected {
		return !action
	}
	if expected := rule.OS.Arch; expected != "" && e.arch != expected {
		return !action
	}
	if expected := rule.OS.Version; expected != "" {
		re, err := regexp.Compile(expected)
		if err != nil && e.version != expected {
			return !action
		}
		if !re.MatchString(e.version) {
			return !action
		}
	}
	if rule.Features != nil {
		for feat, expected := range rule.Features {
			if expected != e.features[feat] {
				return !action
			}
		}
	}

	return action
}

func determineOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "osx"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		panic(fmt.Sprintf("unsupported os/arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	}
}

func determineArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"
	case "386":
		return "x86"
	case "arm64":
		return "arm64"
	default:
		panic(fmt.Sprintf("unsupported os/arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	}
}

func determineVersion() string {
	v, err := platform.GetVersion()
	if err != nil {
		panic(fmt.Errorf("unable to get os version: %w", err))
	}
	return v
}

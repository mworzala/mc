package launch

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/mworzala/mc/internal/pkg/platform"

	"github.com/mworzala/mc/internal/pkg/game/rule"

	"github.com/mworzala/mc/internal/pkg/account"
	gameModel "github.com/mworzala/mc/internal/pkg/game/model"
	"github.com/mworzala/mc/internal/pkg/java"
	"github.com/mworzala/mc/internal/pkg/profile"
	"github.com/mworzala/mc/internal/pkg/util"
)

// todo need to rewrite this whole thing... it's a mess
func LaunchProfile(
	dataDir string,
	p *profile.Profile,
	acc *account.Account,
	accessToken string,
	javaInstall *java.Installation,
	tail bool,
	quickPlay *QuickPlay,
) error {
	var spec gameModel.VersionSpec

	versionSpecPath := path.Join(dataDir, "versions", p.Version, fmt.Sprintf("%s.json", p.Version))
	if err := util.ReadFile(versionSpecPath, &spec); err != nil {
		return err
	}
	if spec.InheritsFrom != "" {

		//todo should move away from merging specs, it creates weird edge cases like the one below to choose the client jar
		var inheritedSpec gameModel.VersionSpec
		inheritedVersionSpecPath := path.Join(dataDir, "versions", spec.InheritsFrom, fmt.Sprintf("%s.json", spec.InheritsFrom))
		if err := util.ReadFile(inheritedVersionSpecPath, &inheritedSpec); err != nil {
			return err
		}

		spec = *mergeSpec(&spec, &inheritedSpec)
	}

	vars := map[string]string{
		// jvm
		"natives_directory": ".",
		"launcher_name":     "mc",
		"launcher_version":  "0.0.1",
		// game
		"version_name":      p.Version,
		"game_directory":    p.Directory,
		"assets_root":       path.Join(dataDir, "assets"),
		"assets_index_name": spec.Assets,
		"auth_player_name":  acc.Profile.Username,
		"auth_uuid":         util.TrimUUID(acc.UUID),
		"auth_access_token": accessToken,
		// Clientid seems to be the mso client id, without dashes, base64 encoded. Should try it with my own client id to see if that works
		"clientid":          "MTMwQUU2ODYwQUE1NDUwNkIyNUZCMzZBNjFCNjc3M0Q=",
		"user_type":         "msa",
		"version_type":      "release", //todo this needs to be release/snapshot
		"resolution_width":  "1920",
		"resolution_height": "1080",
	}

	var features []string
	if quickPlay != nil {
		//todo need to check game version for this
		switch quickPlay.Type {
		case QuickPlaySingleplayer:
			features = append(features, "is_quick_play_singleplayer")
			vars["quickPlaySingleplayer"] = quickPlay.Id
		case QuickPlayMultiplayer:
			features = append(features, "is_quick_play_multiplayer")
			vars["quickPlayMultiplayer"] = quickPlay.Id
		case QuickPlayRealms:
			features = append(features, "is_quick_play_realms")
			vars["quickPlayRealms"] = quickPlay.Id
		}
	}
	rules := rule.NewEvaluator(features...)

	// Build classpath
	classpath := strings.Builder{}
	librariesPath := path.Join(dataDir, "libraries")

	for _, lib := range spec.Libraries {
		if rules.Eval(lib.Rules) == rule.Deny {
			continue
		}

		if lib.Downloads != nil { // Vanilla-type library
			libPath := path.Join(librariesPath, lib.Downloads.Artifact.Path)
			classpath.WriteString(libPath)
		} else if lib.Url != "" { // Direct maven library
			parts := strings.Split(lib.Name, ":")
			groupId := parts[0]
			artifactName := parts[1]
			version := parts[2]

			artifactPath := fmt.Sprintf("%s/%s/%s/%s-%s.jar", strings.ReplaceAll(groupId, ".", "/"), artifactName, version, artifactName, version)
			classpath.WriteString(path.Join(librariesPath, artifactPath))
		}
		classpath.WriteString(platform.ClasspathSeparator)
	}

	if spec.InheritsFrom != "" {
		classpath.WriteString(path.Join(dataDir, "versions", spec.InheritsFrom, fmt.Sprintf("%s.jar", spec.InheritsFrom)))
	} else {
		classpath.WriteString(path.Join(dataDir, "versions", p.Version, fmt.Sprintf("%s.jar", p.Version)))
	}

	vars["classpath"] = classpath.String()

	if msoTokenData, ok := acc.Source.(*account.MicrosoftTokenData); ok {
		vars["auth_xuid"] = msoTokenData.UserHash
	}

	replaceVars := func(s string) string {
		for k, v := range vars {
			s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
		}
		return s
	}

	var args []string

	for _, arg := range spec.Arguments.JVM {
		if s, ok := arg.(string); ok {
			args = append(args, replaceVars(s))
		} else if m, ok := arg.(map[string]interface{}); ok {
			var ruleDef []*rule.Rule

			md, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				DecodeHook: func(from, to reflect.Type, value interface{}) (interface{}, error) {
					if from.Kind() == reflect.String && to == reflect.TypeOf(rule.Allow) {
						if value.(string) == "allow" {
							return rule.Allow, nil
						}
						return rule.Deny, nil
					}
					return value, nil
				},
				Result: &ruleDef,
			})
			if err != nil {
				panic("todo")
			}
			if err := md.Decode(m["rules"]); err != nil {
				panic(fmt.Errorf("invalid rule: %w", err)) //todo better error handling. Should print about this and add an option to ignore unknown rules
			}

			if rules.Eval(ruleDef) == rule.Deny {
				continue
			}

			// Add the rules
			switch value := m["value"].(type) {
			case string:
				args = append(args, replaceVars(value))
			case []interface{}:
				for _, v := range value {
					if s, ok := v.(string); ok {
						args = append(args, replaceVars(s))
					} else {
						panic(fmt.Sprintf("unknown inner value type: %T", v))
					}
				}
			default:
				panic(fmt.Sprintf("unknown value type: %T", value))
			}
		} else {
			panic("unknown arg type")
		}
	}

	args = append(args, spec.MainClass)

	for _, arg := range spec.Arguments.Game {
		if s, ok := arg.(string); ok {
			args = append(args, replaceVars(s))
		} else if m, ok := arg.(map[string]interface{}); ok {
			//todo duplicated above
			var ruleDef []*rule.Rule

			md, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				DecodeHook: func(from, to reflect.Type, value interface{}) (interface{}, error) {
					if from.Kind() == reflect.String && to == reflect.TypeOf(rule.Allow) {
						if value.(string) == "allow" {
							return rule.Allow, nil
						}
						return rule.Deny, nil
					}
					return value, nil
				},
				Result: &ruleDef,
			})
			if err != nil {
				panic("todo")
			}
			if err := md.Decode(m["rules"]); err != nil {
				panic(fmt.Errorf("invalid rule: %w", err)) //todo better error handling. Should print about this and add an option to ignore unknown rules
			}

			if rules.Eval(ruleDef) == rule.Deny {
				continue
			}

			// Add the rules
			switch value := m["value"].(type) {
			case string:
				args = append(args, replaceVars(value))
			case []interface{}:
				for _, v := range value {
					if s, ok := v.(string); ok {
						args = append(args, replaceVars(s))
					} else {
						panic(fmt.Sprintf("unknown inner value type: %T", v))
					}
				}
			default:
				panic(fmt.Sprintf("unknown value type: %T", value))
			}
		} else {
			panic("unknown arg type")
		}
	}

	cmd := exec.Command(javaInstall.Path, args...)
	cmd.Dir = p.Directory

	if tail {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
	} else {
		cmd.Stdout = io.Discard
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if tail {
		if err := cmd.Wait(); err != nil {
			panic(err)
		}
	}

	return nil
}

func mergeSpec(spec, base *gameModel.VersionSpec) *gameModel.VersionSpec {
	var result gameModel.VersionSpec

	if spec.Id != "" {
		result.Id = spec.Id
	} else {
		result.Id = base.Id
	}

	if spec.InheritsFrom != "" {
		result.InheritsFrom = spec.InheritsFrom
	} else {
		result.InheritsFrom = base.InheritsFrom
	}

	if spec.MinimumLauncherVersion != 0 {
		result.MinimumLauncherVersion = spec.MinimumLauncherVersion
	} else {
		result.MinimumLauncherVersion = base.MinimumLauncherVersion
	}

	if spec.ComplianceLevel != 0 {
		result.ComplianceLevel = spec.ComplianceLevel
	} else {
		result.ComplianceLevel = base.ComplianceLevel
	}

	if spec.Downloads != nil {
		result.Downloads = spec.Downloads
	} else {
		result.Downloads = base.Downloads
	}

	result.Libraries = append(spec.Libraries, base.Libraries...)

	if spec.AssetIndex != nil {
		result.AssetIndex = spec.AssetIndex
	} else {
		result.AssetIndex = base.AssetIndex
	}

	if spec.Logging != nil {
		result.Logging = spec.Logging
	} else {
		result.Logging = base.Logging
	}

	if spec.JavaVersion != nil {
		result.JavaVersion = spec.JavaVersion
	} else {
		result.JavaVersion = base.JavaVersion
	}

	if spec.MainClass != "" {
		result.MainClass = spec.MainClass
	} else {
		result.MainClass = base.MainClass
	}

	if spec.Assets != "" {
		result.Assets = spec.Assets
	} else {
		result.Assets = base.Assets
	}

	result.Arguments.JVM = append(spec.Arguments.JVM, base.Arguments.JVM...)
	result.Arguments.Game = append(spec.Arguments.Game, base.Arguments.Game...)

	return &result
}

package launch

import (
	"github.com/mworzala/mc-cli/internal/pkg/account"
	"github.com/mworzala/mc-cli/internal/pkg/java"
	"github.com/mworzala/mc-cli/internal/pkg/profile"
)

func LaunchProfile(dataDir string, p *profile.Profile, acc *account.Account, javaInstall *java.Installation) error {
	return nil
	// var spec VersionSpec
	// versionSpecPath := path.Join(dataDir, "versions", p.Version, fmt.Sprintf("%s.json", p.Version))
	//
	//	if err := readFile(versionSpecPath, &spec); err != nil {
	//		return err
	//	}
	//
	// // Build classpath
	// classpath := strings.Builder{}
	// librariesPath := path.Join(dataDir, "libraries")
	//
	//	for _, lib := range spec.Libraries {
	//		libPath := path.Join(librariesPath, lib.Downloads.Artifact.Path)
	//		classpath.WriteString(libPath)
	//		classpath.WriteString(":")
	//	}
	//
	// classpath.WriteString(path.Join(dataDir, "versions", p.Version, fmt.Sprintf("%s.jar", p.Version)))
	//
	//	vars := map[string]string{
	//		// jvm
	//		"natives_directory": ".",
	//		"launcher_name":     "mc-cli",
	//		"launcher_version":  "0.0.1",
	//		"classpath":         classpath.String(),
	//		// game
	//		"version_name":      p.Version,
	//		"game_directory":    p.Directory,
	//		"assets_root":       path.Join(dataDir, "assets"),
	//		"assets_index_name": spec.Assets,
	//		"auth_player_name":  acc.Profile.Username,
	//		"auth_uuid":         util.TrimUUID(acc.UUID),
	//		"auth_access_token": acc.AccessToken,
	//		// Clientid seems to be the mso client id, without dashes, base64 encoded. Should try it with my own client id to see if that works
	//		"clientid":          "MTMwQUU2ODYwQUE1NDUwNkIyNUZCMzZBNjFCNjc3M0Q=",
	//		"user_type":         "msa",
	//		"version_type":      "release", //todo this needs to be release/snapshot
	//		"resolution_width":  "1920",
	//		"resolution_height": "1080",
	//	}
	//
	//	if msoTokenData, ok := acc.Source.(*account.MicrosoftTokenData); ok {
	//		vars["auth_xuid"] = msoTokenData.UserHash
	//	}
	//
	//	replaceVars := func(s string) string {
	//		for k, v := range vars {
	//			s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	//		}
	//		return s
	//	}
	//
	// var args []string
	// args = append(args, "-XstartOnFirstThread")
	//
	//	for _, arg := range spec.Arguments.JVM {
	//		if s, ok := arg.(string); ok {
	//			args = append(args, replaceVars(s))
	//		} else if m, ok := arg.(map[string]interface{}); ok {
	//			_ = m
	//			//value := m["value"]
	//			//if s, ok := value.(string); ok {
	//			//	args = append(args, replaceVars(s))
	//			//} else if a, ok := value.([]interface{}); ok {
	//			//	for _, v := range a {
	//			//		if s, ok := v.(string); ok {
	//			//			args = append(args, replaceVars(s))
	//			//		}
	//			//	}
	//			//} else {
	//			//	panic(fmt.Sprintf("unknown type: %T", value))
	//			//}
	//		} else {
	//			panic("unknown arg type")
	//		}
	//	}
	//
	// args = append(args, spec.MainClass)
	//
	//	for _, arg := range spec.Arguments.Game {
	//		if s, ok := arg.(string); ok {
	//			args = append(args, replaceVars(s))
	//		} else if m, ok := arg.(map[string]interface{}); ok {
	//			_ = m
	//			//value := m["value"]
	//			//if s, ok := value.(string); ok {
	//			//	args = append(args, replaceVars(s))
	//			//} else if a, ok := value.([]interface{}); ok {
	//			//	for _, v := range a {
	//			//		if s, ok := v.(string); ok {
	//			//			args = append(args, replaceVars(s))
	//			//		}
	//			//	}
	//			//} else {
	//			//	panic(fmt.Sprintf("unknown type: %T", value))
	//			//}
	//		} else {
	//			panic("unknown arg type")
	//		}
	//	}
	//
	// cmd := exec.Command(javaInstall.Path, args...)
	// cmd.Dir = p.Directory
	//
	// tail := false
	//
	//	if tail {
	//		cmd.Stdout = os.Stdout
	//	} else {
	//
	//		cmd.Stdout = io.Discard
	//	}
	//
	//	if err := cmd.Start(); err != nil {
	//		panic(err)
	//	}
	//
	//	if tail {
	//		if err := cmd.Wait(); err != nil {
	//			panic(err)
	//		}
	//	}
	//
	// return nil
}

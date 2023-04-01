package model

import (
	"github.com/mworzala/mc-cli/internal/pkg/game/rule"
	"github.com/mworzala/mc-cli/internal/pkg/util"
)

const (
	MojangObjectBaseUrl = "https://resources.download.minecraft.net"
)

type VersionSpec struct {
	// General
	Id                     string `json:"id"`
	MinimumLauncherVersion int    `json:"minimumLauncherVersion"`
	InheritsFrom           string `json:"inheritsFrom"`
	// ComplianceLevel of 1 indicates that it supports "new security features"
	// Its used for that warning in the launcher.
	ComplianceLevel int `json:"complianceLevel"`
	//	ReleaseTime     time.Time `json:"releaseTime"`
	//	Time            time.Time `json:"time"`

	// Installation
	Downloads *struct {
		Client         *util.FileDownload `json:"client"`
		ClientMappings *util.FileDownload `json:"client_mappings"`
		Server         *util.FileDownload `json:"server"`
		ServerMappings *util.FileDownload `json:"server_mappings"`
	} `json:"downloads"`
	Libraries  []*Library `json:"libraries"`
	AssetIndex *struct {
		Id        string `json:"id"`
		TotalSize int64  `json:"totalSize"`
		util.FileDownload
	} `json:"assetIndex"`
	Logging *struct {
		Client struct {
			Argument string `json:"argument"`
			File     struct {
				util.FileDownload
				Id string `json:"id"`
			} `json:"file"`
			Type string `json:"type"`
		} `json:"client"`
	} `json:"logging"`

	// Launch
	JavaVersion *struct {
		Component    string `json:"component"`
		MajorVersion int    `json:"majorVersion"`
	} `json:"javaVersion"`
	MainClass string `json:"mainClass"`
	Assets    string `json:"assets"`
	Arguments *struct {
		Game []interface{} `json:"game"`
		JVM  []interface{} `json:"jvm"`
	} `json:"arguments"`
}

// Library is blah blah blah
//
// There are two unique types of library being expressed here: vanilla and direct maven
//   - Vanilla libraries include a Downloads field, and never a Url.
//   - Direct maven libraries never include Downloads, and always have a Url.
//     The URL in these dependencies is the base url of the maven repo, so it must be merged with
//     the name, which is given in dependency form, eg `net.fabricmc:access-widener:2.1.0`
type Library struct {
	Name  string       `json:"name"`
	Rules []*rule.Rule `json:"rules"`

	// Vanilla
	Downloads *struct {
		Artifact struct {
			Path string `json:"path"`
			util.FileDownload
		} `json:"artifact"`
	} `json:"downloads"`

	// Direct maven
	Url string `json:"url"`
}

type AssetIndex struct {
	Objects map[string]*AssetObject `json:"objects"`
}

type AssetObject struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

//func MergeSpec(spec, base *VersionSpec) *VersionSpec {
//	var result VersionSpec
//	if err := copier.Copy(&result, spec); err != nil {
//		panic(err)
//	}
//
//	if result.Id == "" {
//		result.Id = base.Id
//	}
//
//	if result.MinimumLauncherVersion == 0 {
//		result.MinimumLauncherVersion = base.MinimumLauncherVersion
//	}
//
//	if result.ComplianceLevel == 0 {
//		result.ComplianceLevel = base.ComplianceLevel
//	}
//
//	var zeroTime time.Time
//	if result.ReleaseTime == zeroTime {
//		result.ReleaseTime = base.ReleaseTime
//	}
//	if result.Time == zeroTime {
//		result.Time = base.Time
//	}
//
//	var emptyDownload util.FileDownload
//	if result.Downloads.Client == emptyDownload {
//		result.Downloads.Client = base.Downloads.Client
//	}
//	if result.Downloads.ClientMappings == emptyDownload {
//		result.Downloads.ClientMappings = base.Downloads.ClientMappings
//	}
//	if result.Downloads.Server == emptyDownload {
//		result.Downloads.Server = base.Downloads.Server
//	}
//	if result.Downloads.ServerMappings == emptyDownload {
//		result.Downloads.ServerMappings = base.Downloads.ServerMappings
//	}
//
//	result.Libraries = append(result.Libraries, base.Libraries...)
//
//	if result.AssetIndex.Id == "" {
//		result.AssetIndex = base.AssetIndex
//	}
//
//	if result.Logging.Client.File.Id == "" {
//		result.Logging = base.Logging
//	}
//
//	if result.JavaVersion.MajorVersion == 0 {
//		result.JavaVersion = base.JavaVersion
//	}
//
//	if result.MainClass == "" {
//		result.MainClass = base.MainClass
//	}
//
//	if result.Assets == "" {
//		result.Assets = base.Assets
//	}
//
//	result.Arguments.Game = append(result.Arguments.Game, base.Arguments.Game...)
//	result.Arguments.JVM = append(result.Arguments.JVM, base.Arguments.JVM...)
//
//	return &result
//}

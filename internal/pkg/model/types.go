package model

import "time"

type Manifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	}
	Versions []*ManifestVersion
}

type ManifestVersion struct {
	Id          string    `json:"id"`
	ReleaseTime time.Time `json:"releaseTime"`
	Time        time.Time `json:"time"`
	Type        string    `json:"type"`
	Url         string    `json:"url"`
}

// Version JSON

type Version struct {
	Arguments struct {
		Game []interface{} `json:"game"`
		JVM  []interface{} `json:"jvm"`
	} `json:"arguments"`
	AssetIndex struct {
		Id        string `json:"id"`
		TotalSize int64  `json:"totalSize"`
		Download
	} `json:"assetIndex"`
	Assets          string `json:"assets"`
	ComplianceLevel int    `json:"complianceLevel"` //todo enum?
	Downloads       struct {
		Client         Download `json:"client"`
		ClientMappings Download `json:"client_mappings"`
		Server         Download `json:"server"`
		ServerMappings Download `json:"server_mappings"`
	} `json:"downloads"`
	Id          string `json:"id"`
	JavaVersion struct {
		Component    string `json:"component"`
		MajorVersion int    `json:"majorVersion"`
	} `json:"javaVersion"`
	Libraries []*Library `json:"libraries"`
	Logging   struct {
		Client struct {
			Argument string `json:"argument"`
			File     struct {
				Id string `json:"id"`
				Download
			} `json:"file"`
			Type string `json:"type"`
		} `json:"client"`
	} `json:"logging"` //todo
	MainClass              string    `json:"mainClass"`
	MinimumLauncherVersion int       `json:"minimumLauncherVersion"`
	ReleaseTime            time.Time `json:"releaseTime"`
	Time                   time.Time `json:"time"`
	Type                   string    `json:"type"`
}

type Library struct {
	Name      string `json:"name"`
	Downloads struct {
		Artifact struct {
			Path string `json:"path"`
			Download
		} `json:"artifact"`
	} `json:"downloads"`
	Rules []map[string]interface{} `json:"rules"` //todo
}

type Download struct {
	Sha1 string `json:"sha1"`
	Size int64  `json:"size"`
	Url  string `json:"url"`
}

/*
Known rules
- os.name == osx|windows|linux
- os.version == regex(?)
- os.arch == x86
- features.is_demo_user == bool
- features.has_custom_resolution == bool
*/

type AssetIndex struct {
	Objects map[string]struct {
		Hash string `json:"hash"`
		Size int64  `json:"size"`
	} `json:"objects"`
}

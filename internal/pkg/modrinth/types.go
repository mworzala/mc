package modrinth

import "time"

type VersionType string

const (
	Alpha   VersionType = "alpha"
	Beta    VersionType = "beta"
	Release VersionType = "release"
)

type Version struct {
	VersionType   VersionType `json:"version_type"`
	DatePublished time.Time   `json:"date_published"`
	Files         []*struct {
		Hashes struct {
			Sha1 string `json:"sha1"`
		}
		Url      string `json:"url"`
		Filename string `json:"filename"`
		Primary  bool   `json:"primary"`
		Size     int64  `json:"size"`
	} `json:"files"`
}

package util

import "fmt"

const userAgentFormat = "mworzala/mc/%s"

func MakeUserAgent(idVersion string) string {
	return fmt.Sprintf(userAgentFormat, idVersion)
}

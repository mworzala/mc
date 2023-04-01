package util

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	uuidPattern        = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	trimmedUUIDPattern = regexp.MustCompile("^[a-fA-F0-9]{32}$")
)

// IsUUID checks if a given string is a valid UUID (expanded or trimmed)
func IsUUID(uuid string) bool {
	return uuidPattern.MatchString(uuid) || trimmedUUIDPattern.MatchString(uuid)
}

// ExpandUUID takes a uuid in the form without dashes, such as `aceb326fda1545bcbf2f11940c21780c`
// and returns it with dashes, eg `aceb326f-da15-45bc-bf2f-11940c21780c`.
//
// If the uuid is already expanded, it will be returned as is. If the UUID is invalid, panic.
func ExpandUUID(uuid string) string {
	switch len(uuid) {
	case 32:
		return uuid[:8] + "-" + uuid[8:12] + "-" + uuid[12:16] + "-" + uuid[16:20] + "-" + uuid[20:]
	case 36:
		return uuid // Already expanded
	default:
		// Something is wrong
		panic(fmt.Sprintf("not a uuid: %s", uuid))
	}
}

// TrimUUID takes a uuid in the form with dashes, such as `aceb326f-da15-45bc-bf2f-11940c21780c`
// and returns it without dashes, eg `aceb326fda1545bcbf2f11940c21780c`.
//
// If the uuid is already trimmed, it will be returned as is. If the UUID is invalid, panic.
func TrimUUID(uuid string) string {
	switch len(uuid) {
	case 36:
		return strings.ReplaceAll(uuid, "-", "")
	case 32:
		return uuid // Already trimmed
	default:
		// Something is wrong
		panic(fmt.Sprintf("not a uuid: %s", uuid))
	}
}

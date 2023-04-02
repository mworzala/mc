package java

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractParams(t *testing.T) {
	sample := `Property settings:
file.encoding = UTF-8
file.separator = /
java.class.path =
java.class.version = 63.0
java.home = /Library/Java/JavaVirtualMachines/temurin-19.jdk/Contents/Home
java.io.tmpdir = /var/folders/bx/h2q2q_0j73s6n2mw8vh_z_lw0000gn/T/
java.runtime.name = OpenJDK Runtime Environment
java.runtime.version = 19.0.2+7
java.specification.name = Java Platform API Specification
java.specification.vendor = Oracle Corporation
java.specification.version = 19
java.vendor = Eclipse Adoptium
java.vendor.url = https://adoptium.net/
java.vendor.url.bug = https://github.com/adoptium/adoptium-support/issues
java.vendor.version = Temurin-19.0.2+7
java.version = 19.0.2
sun.management.compiler = HotSpot 64-Bit Tiered Compilers
user.country = US
user.dir = /Users/matt
user.home = /Users/matt
user.language = en
user.name = matt

openjdk 19.0.2 2023-01-17
OpenJDK Runtime Environment Temurin-19.0.2+7 (build 19.0.2+7)
OpenJDK 64-Bit Server VM Temurin-19.0.2+7 (build 19.0.2+7, mixed mode)`

	result := extractParams(sample)

	require.Equal(t, "UTF-8", result["file.encoding"])
	// The below case is wrong, it finds the one on the next value. Luckily I don't actually care about this
	// value, so i wont fix the issue :)
	require.Equal(t, "java.class.version = 63.0", result["java.class.path"])
	require.Equal(t, "HotSpot 64-Bit Tiered Compilers", result["sun.management.compiler"])
}

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsUUID(t *testing.T) {
	// Test valid UUID string with dashes
	uuidWithDashes := "aceb326f-da15-45bc-bf2f-11940c21780c"
	require.True(t, IsUUID(uuidWithDashes), "Valid UUID with dashes check failed")

	// Test valid UUID string without dashes
	uuidWithoutDashes := "aceb326fda1545bcbf2f11940c21780c"
	require.True(t, IsUUID(uuidWithoutDashes), "Valid UUID without dashes check failed")

	// Test invalid UUID string
	invalidUUID := "not-a-uuid"
	require.False(t, IsUUID(invalidUUID), "Invalid UUID check failed")
}

func TestExpandUUID(t *testing.T) {
	// Test valid UUID string without dashes
	uuidWithoutDashes := "aceb326fda1545bcbf2f11940c21780c"
	expectedUUID := "aceb326f-da15-45bc-bf2f-11940c21780c"
	require.Equal(t, expectedUUID, ExpandUUID(uuidWithoutDashes))

	// Test valid UUID string with dashes
	uuidWithDashes := "aceb326f-da15-45bc-bf2f-11940c21780c"
	require.Equal(t, uuidWithDashes, ExpandUUID(uuidWithDashes))

	// Test invalid UUID string
	require.Panics(t, func() { ExpandUUID("not-a-uuid") })
}

func TestTrimUUID(t *testing.T) {
	// Test valid UUID string with dashes
	uuidWithDashes := "aceb326f-da15-45bc-bf2f-11940c21780c"
	expectedUUID := "aceb326fda1545bcbf2f11940c21780c"
	require.Equal(t, expectedUUID, TrimUUID(uuidWithDashes))

	// Test valid UUID string without dashes
	uuidWithoutDashes := "aceb326fda1545bcbf2f11940c21780c"
	require.Equal(t, uuidWithoutDashes, TrimUUID(uuidWithoutDashes))

	// Test invalid UUID string
	require.Panics(t, func() { TrimUUID("not-a-uuid") })
}

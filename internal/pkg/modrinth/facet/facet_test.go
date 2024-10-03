package facet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		name string
		in   Root
		out  string
	}{
		{"empty and", And{}, ""},
		{"empty or", Or{}, ""},
		{"or no and", Or{Eq{ProjectType, "mod"}}, `[["project_type=mod"]]`},
		{"and no or", And{Eq{ProjectType, "mod"}}, `[["project_type=mod"]]`},
		{"no and or", Eq{ProjectType, "mod"}, `[["project_type=mod"]]`},
		{"and or", And{Or{Eq{ProjectType, "mod"}}}, `[["project_type=mod"]]`},
		{"and multi", And{Eq{ProjectType, "mod"}, NEq{ProjectType, "modpack"}}, `[["project_type=mod"],["project_type!=modpack"]]`},
		{"or multi", Or{Eq{ProjectType, "mod"}, NEq{ProjectType, "modpack"}}, `[["project_type=mod","project_type!=modpack"]]`},

		{"eq", Eq{ProjectType, "mod"}, `[["project_type=mod"]]`},
		{"neq", NEq{ProjectType, "mod"}, `[["project_type!=mod"]]`},
		{"gt", Gt{ProjectType, "mod"}, `[["project_type>mod"]]`},
		{"gteq", GtEq{ProjectType, "mod"}, `[["project_type>=mod"]]`},
		{"lt", Lt{ProjectType, "mod"}, `[["project_type<mod"]]`},
		{"lteq", LtEq{ProjectType, "mod"}, `[["project_type<=mod"]]`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.in.asString(false, false)
			require.NoError(t, err)
			require.Equal(t, test.out, res)
		})
	}
}

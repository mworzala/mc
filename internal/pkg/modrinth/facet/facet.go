package facet

import (
	"fmt"
	"strconv"
	"time"
)

func ToString(facet Root) (string, error) {
	if facet == nil {
		return "", nil
	}
	return facet.asString(false, false)
}

type And []OrFilter

type Or []Filter

type Eq struct {
	Type  Type
	Value interface{}
}

type NEq struct {
	Type  Type
	Value interface{}
}

type Gt struct {
	Type  Type
	Value interface{}
}

type GtEq struct {
	Type  Type
	Value interface{}
}

type Lt struct {
	Type  Type
	Value interface{}
}

type LtEq struct {
	Type  Type
	Value interface{}
}

type Type string

const (
	ProjectType       Type = "project_type"
	Categories        Type = "categories"
	Versions          Type = "versions"
	ClientSide        Type = "client_side"
	ServerSide        Type = "server_side"
	OpenSource        Type = "open_source"
	Title             Type = "title"
	Author            Type = "author"
	Follows           Type = "follows"
	ProjectID         Type = "project_id"
	License           Type = "license"
	Downloads         Type = "downloads"
	Color             Type = "color"
	CreatedTimestamp  Type = "created_timestamp"
	ModifiedTimestamp Type = "modified_timestamp"
)

type (
	anyFilter interface {
		asString(hasAnd, hasOr bool) (string, error)
	}
	Root interface {
		anyFilter
		facetRoot()
	}
	OrFilter interface {
		anyFilter
		orFilter()
	}
	Filter interface {
		anyFilter
		facet()
	}
)

func (f And) asString(_, _ bool) (string, error) {
	if len(f) == 0 {
		return "", nil
	}

	out := "["
	for i, child := range f {
		childStr, err := child.asString(true, false)
		if err != nil {
			return "", err
		}
		out += childStr
		if i != len(f)-1 {
			out += ","
		}
	}

	return out + "]", nil
}

func (f Or) asString(hasAnd, _ bool) (string, error) {
	if len(f) == 0 {
		return "", nil
	}

	out := "["
	for i, child := range f {
		childStr, err := child.asString(true, true)
		if err != nil {
			return "", err
		}
		out += childStr
		if i != len(f)-1 {
			out += ","
		}
	}
	out += "]"

	if !hasAnd {
		out = "[" + out + "]"
	}

	return out, nil
}

func (f Eq) asString(hasAnd, hasOr bool) (string, error) {
	return serializeOperation(hasAnd, hasOr, "=", f.Type, f.Value)
}
func (f NEq) asString(hasAnd, hasOr bool) (string, error) {
	return serializeOperation(hasAnd, hasOr, "!=", f.Type, f.Value)
}
func (f Gt) asString(hasAnd, hasOr bool) (string, error) {
	return serializeOperation(hasAnd, hasOr, ">", f.Type, f.Value)
}
func (f GtEq) asString(hasAnd, hasOr bool) (string, error) {
	return serializeOperation(hasAnd, hasOr, ">=", f.Type, f.Value)
}
func (f Lt) asString(hasAnd, hasOr bool) (string, error) {
	return serializeOperation(hasAnd, hasOr, "<", f.Type, f.Value)
}
func (f LtEq) asString(hasAnd, hasOr bool) (string, error) {
	return serializeOperation(hasAnd, hasOr, "<=", f.Type, f.Value)
}

func (And) facetRoot()  {}
func (Or) facetRoot()   {}
func (Eq) facetRoot()   {}
func (NEq) facetRoot()  {}
func (Gt) facetRoot()   {}
func (GtEq) facetRoot() {}
func (Lt) facetRoot()   {}
func (LtEq) facetRoot() {}

func (Or) orFilter()   {}
func (Eq) orFilter()   {}
func (NEq) orFilter()  {}
func (Gt) orFilter()   {}
func (GtEq) orFilter() {}
func (Lt) orFilter()   {}
func (LtEq) orFilter() {}

func (Or) facet()   {}
func (Eq) facet()   {}
func (NEq) facet()  {}
func (Gt) facet()   {}
func (GtEq) facet() {}
func (Lt) facet()   {}
func (LtEq) facet() {}

func serializeOperation(hasAnd, hasOr bool, op string, _type Type, value interface{}) (string, error) {
	valueStr, err := serializeValue(value)
	if err != nil {
		return "", err
	}

	out := `"` + string(_type) + op + valueStr + `"`

	if !hasAnd {
		out = "[" + out + "]"
	}
	if !hasOr {
		out = "[" + out + "]"
	}
	return out, nil
}

func serializeValue(value interface{}) (string, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return "true", nil
		} else {
			return "false", nil
		}
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64), nil
	case time.Time:
		return v.Format(time.RFC3339), nil
	default:
		return "", fmt.Errorf("unexpected facet value type %T", v)
	}
}

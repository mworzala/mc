package skin

import (
	"time"
)

type Skin struct {
	Name      string    `json:"name"`
	Variant   string    `json:"variant"`
	Skin      string    `json:"skin"`
	Cape      string    `json:"cape"`
	AddedDate time.Time `json:"added_date"`
}

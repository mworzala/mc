package mojang

type ProfileInformationResponse struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Skins          []ProfileSkin `json:"skins"`
	Capes          []ProfileCape `json:"capes"`
	ProfileActions struct {
	} `json:"profileActions"`
}

type ProfileSkin struct {
	ID         string `json:"id"`
	State      string `json:"state"`
	URL        string `json:"url"`
	TextureKey string `json:"textureKey"`
	Variant    string `json:"variant"`
}

type ProfileCape struct {
	ID    string `json:"id"`
	State string `json:"state"`
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

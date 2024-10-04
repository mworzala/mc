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

type UsernameToUuidResponse struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type UuidToProfileResponse struct {
	Id         string              `json:"id"`
	Name       string              `json:"name"`
	Properties []ProfileProperties `json:"properties"`
	Legacy     bool                `json:"legacy"`
}

type ProfileProperties struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

type TextureInformation struct {
	Timestamp   int      `json:"timestamp"`
	ProfileId   string   `json:"profileId"`
	ProfileName string   `json:"profileName"`
	Textures    Textures `json:"textures"`
}

type Textures struct {
	Skin struct {
		Url      string `json:"url"`
		Metadata struct {
			Model string `json:"model"`
		} `json:"metadata"`
	} `json:"SKIN"`
	Cape struct {
		Url string `json:"url"`
	} `json:"CAPE"`
}

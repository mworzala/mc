package launch

type QuickPlayType int

const (
	QuickPlaySingleplayer QuickPlayType = iota
	QuickPlayMultiplayer
	QuickPlayRealms
)

type QuickPlay struct {
	Type QuickPlayType
	Id   string
}

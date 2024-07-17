package chat

type participant struct {
	role string
}

var (
	Customer  = participant{"customer"}
	Assistant = participant{"assistant"}
)

func (p participant) Role() string {
	return p.role
}

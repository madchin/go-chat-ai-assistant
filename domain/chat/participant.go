package chat

type Participant struct {
	role string
}

var (
	Customer  = Participant{"customer"}
	Assistant = Participant{"assistant"}
)

func (p Participant) Role() string {
	return p.role
}

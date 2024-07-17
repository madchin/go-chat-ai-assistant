package chat

type participant struct {
	role string
}

var (
	Customer  = participant{"customer"}
	Assistant = participant{"assistant"}
)

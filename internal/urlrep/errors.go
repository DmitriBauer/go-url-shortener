package urlrep

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrDuplicateURL Error = "duplicate url"
)

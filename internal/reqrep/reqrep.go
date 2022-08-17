package reqrep

type ReqRepo interface {
	DataBySessionID(sessionID string) ([]byte, error)

	Save(req Req) error
}

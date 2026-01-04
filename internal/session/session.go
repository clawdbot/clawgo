package session

type ID string

type Context struct {
	SessionKey string
	RunID      string
}

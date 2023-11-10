package rand

const (
	SessionTokenBytes = 32
)

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}

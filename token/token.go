package token

type Tokenizer interface {
	CreateToken(username string) (string, error)

	VerifyToken(token string) error
}

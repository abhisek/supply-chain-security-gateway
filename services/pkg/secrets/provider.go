package secrets

type SecretProvider interface {
	GetSecret(string) (string, error)
}

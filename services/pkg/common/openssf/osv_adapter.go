package openssf

type OsvServiceAdapter struct{}

func NewOsvServiceAdapter() (*OsvServiceAdapter, error) {
	return &OsvServiceAdapter{}, nil
}

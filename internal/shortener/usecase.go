package shortener

type (
	ShortenUseCase struct {
		storage WriteOnlyStorage
	}

	GenerateShortenUseCase struct{}
)

func NewShortenUseCase(storage WriteOnlyStorage) ShortenUseCase {
	return ShortenUseCase{
		storage: storage,
	}
}

func (uc ShortenUseCase) Handle(original string) error {
	lnk, err := newLink(original)
	if err != nil {
		return err
	}

	return uc.storage.Store(lnk)
}

func NewGenerateShortenUseCase() GenerateShortenUseCase {
	return GenerateShortenUseCase{}
}

func (uc GenerateShortenUseCase) Handle(host string, original string) (DestinationURL, error) {
	lnk, err := newLink(original)
	if err != nil {
		return DestinationURL{}, err
	}

	return newFullURL(host, lnk.short), nil
}

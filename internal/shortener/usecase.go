package shortener

type (
	ShortenUseCase struct {
		storage WriteOnlyStorage
	}

	GenerateShortenQuery struct {
		isSSL    bool
		host     string
		original string
	}

	GenerateShortenUseCase struct{}

	GetOriginalUseCase struct {
		storage ReadOnlyStorage
	}
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

func NewGenerateShortenQuery(isSSL bool, host string, original string) GenerateShortenQuery {
	return GenerateShortenQuery{
		isSSL:    isSSL,
		host:     host,
		original: original,
	}
}

func NewGenerateShortenUseCase() GenerateShortenUseCase {
	return GenerateShortenUseCase{}
}

func (uc GenerateShortenUseCase) Handle(query GenerateShortenQuery) (DestinationURL, error) {
	lnk, err := newLink(query.original)
	if err != nil {
		return DestinationURL{}, err
	}

	return newDestinationURL(query.isSSL, query.host, lnk.shortID), nil
}

func NewGetOriginalUseCase(storage ReadOnlyStorage) GetOriginalUseCase {
	return GetOriginalUseCase{
		storage: storage,
	}
}

func (uc GetOriginalUseCase) Handle(short string) (string, error) {
	original, err := uc.storage.GetOriginalURL(shortID{encoded: short})
	if err != nil {
		return "", err
	}

	return original, nil
}

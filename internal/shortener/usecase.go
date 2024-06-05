package shortener

type (
	ShortenUseCase struct {
		storage WriteOnlyStorage
	}

	GetShortQuery struct {
		isSSL    bool
		host     string
		original string
	}

	GetShortUseCase struct {
		storage ReadOnlyStorage
	}

	GetOriginalUseCase struct {
		storage ReadOnlyStorage
	}

	GetOriginalByIDUseCase struct {
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

func NewGenerateShortenQuery(isSSL bool, host string, original string) GetShortQuery {
	return GetShortQuery{
		isSSL:    isSSL,
		host:     host,
		original: original,
	}
}

func NewGenerateShortenUseCase(storage ReadOnlyStorage) GetShortUseCase {
	return GetShortUseCase{storage: storage}
}

func (uc GetShortUseCase) Handle(query GetShortQuery) (DestinationURL, error) {
	lnk, err := newLink(query.original)
	if err != nil {
		return DestinationURL{}, err
	}

	_, err = uc.storage.GetOriginalURL(lnk.ShortID())
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

func (uc GetOriginalUseCase) Handle(shortURL string) (string, error) {
	short, err := newShortLink(shortURL)
	if err != nil {
		return "", err
	}

	original, err := uc.storage.GetOriginalURL(short.shortID())
	if err != nil {
		return "", err
	}

	return original, nil
}

func NewGetOriginalForRedirectUseCase(storage ReadOnlyStorage) GetOriginalByIDUseCase {
	return GetOriginalByIDUseCase{
		storage: storage,
	}
}

func (uc GetOriginalByIDUseCase) Handle(id string) (string, error) {
	short := shortID{encoded: id}

	original, err := uc.storage.GetOriginalURL(short)
	if err != nil {
		return "", err
	}

	return original, nil
}

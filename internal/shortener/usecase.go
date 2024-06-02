package shortener

type (
	ShortenUseCase struct {
		storage WriteOnlyStorage
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

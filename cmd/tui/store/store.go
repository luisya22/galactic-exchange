package store

type Store struct {
	ContentHeight int
}

func NewStore(contentHeight int) *Store {
	return &Store{
		ContentHeight: contentHeight,
	}
}

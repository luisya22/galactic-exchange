package store

type Store struct {
	NavBarHeight int
}

func NewStore(navHeight int) *Store {
	return &Store{
		NavBarHeight: navHeight,
	}
}

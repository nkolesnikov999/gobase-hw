package bins

import "time"

// Repository определяет интерфейс для работы с хранилищем bin-ов
type Repository interface {
	Save(binList BinList, filename string) error
	Load(filename string) (BinList, error)
}

// Bin represents a bin entity with metadata
type Bin struct {
	ID        string    `json:"id"`
	Private   bool      `json:"private"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}

// BinList represents a collection of bins
type BinList []Bin

// NewBin creates a new Bin instance with the provided parameters
func NewBin(id, name string, private bool) *Bin {
	return &Bin{
		ID:        id,
		Private:   private,
		CreatedAt: time.Now(),
		Name:      name,
	}
}

// NewBinList creates a new empty BinList
func NewBinList() BinList {
	return make(BinList, 0)
}

// AddBin adds a new bin to the BinList
func (bl *BinList) AddBin(bin *Bin) {
	*bl = append(*bl, *bin)
}

// Service представляет бизнес-логику для работы с bin-ами
type Service struct {
	repo Repository
}

// NewService создает новый экземпляр Service с переданным репозиторием
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// SaveBins сохраняет список bin-ов используя подключенный репозиторий
func (s *Service) SaveBins(binList BinList, filename string) error {
	return s.repo.Save(binList, filename)
}

// LoadBins загружает список bin-ов используя подключенный репозиторий
func (s *Service) LoadBins(filename string) (BinList, error) {
	return s.repo.Load(filename)
}

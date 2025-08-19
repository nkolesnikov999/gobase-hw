package bins

import "time"

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

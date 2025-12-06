package literaladapter

import (
	"context"
)

// Adapter implements ports.LiteralPort using the local literal package.
type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) CurrentlyReading(ctx context.Context, count int) ([]Book, error) {
	books, err := CurrentlyReading()
	if err != nil {
		return nil, err
	}
	if len(books) > count {
		return books[:count], nil
	}
	return books, nil
}

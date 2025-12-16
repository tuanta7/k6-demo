package location

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
)

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (uc *UseCase) UpdateLatestLocation(ctx context.Context, location *domain.Location) error {
	return nil
}

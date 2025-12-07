package location

import "github.com/tuanta7/k6-demo/services/location/internal/domain"

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (uc *UseCase) UpdateLocation(location domain.Location) error {
	return nil
}

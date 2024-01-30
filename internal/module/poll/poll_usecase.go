package poll

import (
	"context"
	"fmt"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"
	"github.com/google/uuid"
)

// Usecase contains the function of business logic of domain poll
type Usecase interface {
	GetCurrentPoll(ctx context.Context, presentationID uuid.UUID) (model.CurrentPoll, error)
	UpdateCurrentPoll(ctx context.Context, presentationID uuid.UUID) (model.CurrentPoll, error)
}

type service struct {
	store db.Store
}

// Initialize takes all necessary service for domain poll to run the business logic of domain poll
func Initialize(
	store db.Store,
) Usecase {
	return &service{store: store}
}

func (s service) GetCurrentPoll(ctx context.Context, presentationID uuid.UUID) (model.CurrentPoll, error) {
	presentation, err := s.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		return model.CurrentPoll{}, fmt.Errorf("no presentation found")
	}
	if presentation.Currentpollindex.Int32 == 0 {
		return model.CurrentPoll{}, model.ErrNoPollDisplayed
	}
	currentPoll, err := s.store.UpdateCurrentPollToForwardTx(context.Background(), presentationID)
	if err != nil {
		return model.CurrentPoll{}, model.ErrRunOutOfIndex
	}
	return currentPoll, nil

}
func (s service) UpdateCurrentPoll(ctx context.Context, presentationID uuid.UUID) (model.CurrentPoll, error) {
	_, err := s.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		return model.CurrentPoll{}, model.ErrNotFound
	}

	currentPoll, err := s.store.GetCurrentPoll(context.Background(), presentationID)
	if err != nil {
		return model.CurrentPoll{}, model.ErrRunOutOfIndex
	}
	return currentPoll, nil
}

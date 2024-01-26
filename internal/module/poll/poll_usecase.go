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

// Initialize takes all necessary service for domain presentation to run the business logic of domain presentation
func Initialize(
	store db.Store,
) Usecase {
	return &service{store: store}
}

func (s service) GetCurrentPoll(ctx context.Context, presentationID uuid.UUID) (model.CurrentPoll, error) {
	presentation, err := s.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "No presentation found"})
		return model.CurrentPoll{}, fmt.Errorf("no presentation found")
	}
	if presentation.Currentpollindex.Int32 == 0 {
		// c.JSON(http.StatusConflict, "There are no polls currently displayed")
		return model.CurrentPoll{}, fmt.Errorf("there are no polls currently displayed")
	}
	currentPoll, err := s.store.UpdateCurrentPollToForwardTx(context.Background(), presentationID)
	if err != nil {
		// c.JSON(http.StatusConflict, gin.H{"error": "The presentation ran out of polls."})
		return model.CurrentPoll{}, fmt.Errorf("the presentation ran out of polls")
	}
	return currentPoll, nil

}
func (s service) UpdateCurrentPoll(ctx context.Context, presentationID uuid.UUID) (model.CurrentPoll, error) {
	_, err := s.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		// c.JSON(http.StatusNotFound, "There is no presentation with the provided `presentation_id`")
		return model.CurrentPoll{}, fmt.Errorf("there is no presentation with the provided `presentation_id`")
	}

	currentPoll, err := s.store.GetCurrentPoll(context.Background(), presentationID)
	if err != nil {
		// c.JSON(http.StatusConflict, gin.H{"error": "The presentation ran out of polls."})
		return model.CurrentPoll{}, fmt.Errorf("the presentation ran out of polls")
	}
	return currentPoll, nil
}

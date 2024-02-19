package presentation

import (
	"context"
	"encoding/json"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"
	"github.com/google/uuid"
)

// Usecase contains the function of business logic of domain presentation
type Usecase interface {
	CreatePresentation(ctx context.Context, presenation *model.CreatePresentionRequestDTO) (uuid.UUID, error)
	GetPresentation(ctx context.Context, presenationID uuid.UUID) (model.PresentionResponseDTO, error)
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

func (s service) CreatePresentation(ctx context.Context, presenation *model.CreatePresentionRequestDTO) (uuid.UUID, error) {
	jsonb, err := json.Marshal(presenation.Polls)
	if err != nil {
		return uuid.Nil, err
	}

	presID, err := s.store.CreatePresentationAndPolls(context.Background(), jsonb)
	if err != nil {
		return uuid.Nil, err
	}

	return presID, err
}

func (s service) GetPresentation(ctx context.Context, presentationID uuid.UUID) (model.PresentionResponseDTO, error) {
	presentation, err := s.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		return model.PresentionResponseDTO{}, model.ErrNotFound
	}
	
	if presentation.Currentpollindex.Int32 == 0 {
		return model.PresentionResponseDTO{}, model.ErrNoPollDisplayed
	}

	polls, err := s.store.ListPolls(context.Background(), presentationID)
	if err != nil {
		return model.PresentionResponseDTO{}, err
	}

	return model.PresentionResponseDTO{
		CurrentPollIndex: presentation.Currentpollindex.Int32,
		Polls: func() []model.PresentationResponsePoll {
			pollList := make([]model.PresentationResponsePoll, 0)
			for _, p := range polls {
				pollList = append(pollList, model.PresentationResponsePoll{
					Question: p.Question,
					Options:  p.Options,
				})
			}
			return pollList
		}(),
	}, nil
}


func (s service) GetCurrentPoll(ctx context.Context, presentationID uuid.UUID) (model.CurrentPoll, error) {
	presentation, err := s.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		return model.CurrentPoll{}, model.ErrNotFound
	}
	if presentation.Currentpollindex.Int32 == 0 {
		return model.CurrentPoll{}, model.ErrNoPollDisplayed
	}
	currentPoll, err := s.store.GetCurrentPoll(context.Background(), presentationID)
	if err != nil {
		return model.CurrentPoll{}, err
	}
	return currentPoll, nil

}
func (s service) UpdateCurrentPoll(ctx context.Context, presentationID uuid.UUID) (model.CurrentPoll, error) {
	_, err := s.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		return model.CurrentPoll{}, model.ErrNotFound
	}

	currentPoll, err := s.store.UpdateCurrentPollToForwardTx(context.Background(), presentationID)
	if err != nil {
		return model.CurrentPoll{}, model.ErrRunOutOfIndex
	}
	return currentPoll, nil
}

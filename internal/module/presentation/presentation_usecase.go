package presentation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"
	"github.com/google/uuid"
)

// Usecase contains the function of business logic of domain presentation
type Usecase interface {
	CreatePresentation(ctx context.Context, presenation *model.CreatePresentionRequestDTO) (uuid.UUID, error)
	GetPresentation(ctx context.Context, presenationID uuid.UUID) (model.PresentionResponseDTO, error)
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
		// c.AbortWithStatus(http.StatusInternalServerError)
		return uuid.Nil, err
	}

	presID, err := s.store.CreatePresentationAndPolls(context.Background(), []byte(jsonb))
	if err != nil {
		// c.AbortWithStatus(http.StatusBadRequest)
		return uuid.Nil, err
	}

	return presID, err
}

func (s service) GetPresentation(ctx context.Context, presentationID uuid.UUID) (model.PresentionResponseDTO, error) {
	presentation, err := s.store.GetPresentation(context.Background(), presentationID)
	if err != nil {
		// c.JSON(http.StatusNotFound, "There is no presentation with the provided `presentation_id`")
		return model.PresentionResponseDTO{}, fmt.Errorf("there is no presentation with the provided `presentation_id`")
	}
	if presentation.Currentpollindex.Int32 == 0 {
		// c.JSON(http.StatusConflict, "There are no polls currently displayed")
		return model.PresentionResponseDTO{}, fmt.Errorf("there are no polls currently displayed")
	}

	polls, err := s.store.ListPolls(context.Background(), presentationID)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, err.Error())
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

	// c.JSON(http.StatusOK, struct {
	// 	CurrentPollIndex int32             `json:"current_poll_index"`
	// 	Polls            []db.ListPollsRow `json:"polls"`
	// }{
	// 	CurrentPollIndex: presentation.Currentpollindex.Int32,
	// 	Polls:            polls,
	// })
}

package vote

import (
	"context"
	"fmt"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"
	"github.com/google/uuid"
)

// Usecase contains the function of business logic of domain poll
type Usecase interface {
	CreateVote(ctx context.Context, presentationID uuid.UUID, param model.CreateVoteRequestDTO) error
	GetPollVotes(ctx context.Context, presentationID, pollID uuid.UUID) ([]db.GetVotesRow, error)
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

func (s service) CreateVote(ctx context.Context, presentationID uuid.UUID, param model.CreateVoteRequestDTO) error {
	err := s.store.VoteCurrentPollTx(context.Background(), db.VoteParams{
		PresentationID: presentationID,
		Pollid:         param.PollID,
		Optionkey:      param.Key,
		Clientid:       param.ClientId,
	})
	return err
}

func (s service) GetPollVotes(ctx context.Context, presentationID, pollID uuid.UUID) ([]db.GetVotesRow, error) {
	result, err := s.store.GetPresentationAndPoll(context.Background(), db.GetPresentationAndPollParams{
		PresentationID: presentationID,
		PollID:         pollID,
	})

	if err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Either `presentation_id` or `poll_id` not found"})
		return []db.GetVotesRow{}, fmt.Errorf("either `presentation_id` or `poll_id` not found")
	}
	if result.Currentpollindex.Int32 != result.Pollindex {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "The `poll_id` in the request body doesn't match the current poll."})
		return []db.GetVotesRow{}, fmt.Errorf("the `poll_id` in the request body doesn't match the current poll")
	}
	votes, err := s.store.GetVotes(context.Background(), db.GetVotesParams{
		PresentationID: presentationID,
		PollID:         pollID,
	})
	return votes, err
}

package vote

import (
	"context"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	db "github.com/aleale2121/interactive-presentation/internal/storage/persistence"
	"github.com/google/uuid"
)

// Usecase contains the function of business logic of domain voter
type Usecase interface {
	CreateVote(ctx context.Context, presentationID uuid.UUID, param model.CreateVoteRequestDTO) error
	GetPollVotes(ctx context.Context, presentationID, pollID uuid.UUID) ([]db.GetVotesRow, error)
}

type service struct {
	store db.Store
}

// Initialize takes all necessary service for domain vote to run the business logic of domain vote
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
		return []db.GetVotesRow{}, model.ErrNotFound
	}
	if result.Currentpollindex.Int32 != result.Pollindex {
		return []db.GetVotesRow{}, model.ErrIDMismatch
	}
	votes, err := s.store.GetVotes(context.Background(), db.GetVotesParams{
		PresentationID: presentationID,
		PollID:         pollID,
	})
	return votes, err
}

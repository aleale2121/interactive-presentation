package db

import (
	"context"

	"github.com/aleale2121/interactive-presentation/models"
)

func (store *SQLStore) VoteCurrentPollTx(ctx context.Context, arg VoteParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		result, err := store.Queries.GetPresentationAndPoll(context.Background(), GetPresentationAndPollParams{
			PresentationID: arg.PresentationID,
			PollID:         arg.Pollid,
		})

		if err != nil {
			return models.ErrNotFound
		}
		if result.Currentpollindex.Int32 != result.Pollindex {
			return models.ErrConflict
		}

		err = q.CreateVote(ctx, CreateVoteParams{
			Pollid:    arg.Pollid,
			Optionkey: arg.Optionkey,
			Clientid:  arg.Clientid,
		})
		if err != nil {
			return err
		}
		return err
	})

	return err
}

package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (store *SQLStore) GetCurrentPoll(ctx context.Context, id uuid.UUID) (CurrentPoll, error) {
	var result CurrentPoll
	err := store.execTx(ctx, func(q *Queries) error {

		poll, err := q.GetPoll(ctx, id)
		if err != nil {
			return err
		}
		result.ID = poll.ID.String()
		result.Question = poll.Question
		var optionRows []PollOption
		err = json.Unmarshal(poll.Options, &optionRows)
		if err != nil {
			return err
		}
		fmt.Println(optionRows)
		result.Options = optionRows
		return err
	})
	return result, err
}

func (store *SQLStore) UpdateCurrentPollToForwardTx(ctx context.Context, id uuid.UUID) (CurrentPoll, error) {
	var result CurrentPoll
	err := store.execTx(ctx, func(q *Queries) error {

		// Update the current poll index for the presentation.
		nextPoll, err := q.GetNextPoll(ctx, id)
		if err != nil {
			return err
		}
		result.ID = nextPoll.ID.String()
		result.Question = nextPoll.Question
		err = json.Unmarshal(nextPoll.Options, &result.Options)
		if err != nil {
			return err
		}

		return err
	})

	return result, err

}

func (store *SQLStore) UpdateCurrentPollToPreviousTx(ctx context.Context, id uuid.UUID) (CurrentPoll, error) {
	var result CurrentPoll
	err := store.execTx(ctx, func(q *Queries) error {

		// Update the current poll index for the presentation.
		nextPoll, err := q.GetPreviousPoll(ctx, id)
		if err != nil {
			return err
		}
		result.ID = nextPoll.ID.String()
		result.Question = nextPoll.Question
		err = json.Unmarshal(nextPoll.Options, &result.Options)
		if err != nil {
			return err
		}

		return err
	})

	return result, err

}

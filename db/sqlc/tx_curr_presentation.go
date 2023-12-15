package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

)

func (store *SQLStore) UpdateCurrentPollTx(ctx context.Context, arg PollIndexParams) (CurrPollIndexResult, error) {
	var result CurrPollIndexResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var pres Presentation
		prID,err:=uuid.Parse(arg.ID)
		if err != nil {
			return err
		}
		pres, err = q.GetPresentation(context.Background(), prID)
		if err != nil {
			return err
		}
		
		var polls []GetPresentationPollsRow
		// Get the number of polls for the presentation.
		polls, err = q.GetPresentationPolls(ctx, prID)
		if err != nil {
			return err
		}

		newPollIndex := pres.Currentpollindex.Int32 + 1

		// Check the poll index is not less than 0 and greater than the presentation polls length.
		if newPollIndex < 0 {
			newPollIndex = 0
		}
		if len(polls) > 0 && newPollIndex >= int32(len(polls)) {
			newPollIndex = 1
		}

		// Update the current poll index for the presentation.
		_, err = q.UpdateCurrPollIndex(ctx, UpdateCurrPollIndexParams{
			Currentpollindex: sql.NullInt32{
				Int32: newPollIndex,
				Valid: true,
			},
			ID: prID,
		})
		if err != nil {
			return err
		}

		result.ID = polls[newPollIndex].ID.String()
		result.Question = polls[newPollIndex].Question

		// Get the poll at the new poll index.
		resPID,err:=uuid.Parse(result.ID)
		if err != nil {
			return err
		}
		result.Options, err = q.GetPollOptions(ctx, resPID)
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}

func (store *SQLStore) GetCurrentPoll(ctx context.Context, arg PollIndexParams) (CurrPollIndexResult, error) {
	var result CurrPollIndexResult
	prID,err:=uuid.Parse(arg.ID)
	if err != nil {
		return CurrPollIndexResult{},err
	}
	pres, err := store.GetPresentation(context.Background(), prID)
	if err != nil {
		return result, err
	}
	// Get the number of polls for the presentation.
	polls, err := store.GetPresentationPolls(ctx, prID)
	if err != nil {
		return result, err
	}

	// Check the poll index is not less than 0 and greater than the presentation polls length.
	if pres.Currentpollindex.Int32 < 1 || pres.Currentpollindex.Int32 > int32(len(polls)) {
		return result, fmt.Errorf("invalid poll index")
	}

	result.ID = polls[pres.Currentpollindex.Int32].ID.String()
	result.Question = polls[pres.Currentpollindex.Int32].Question

	// Get the poll at the new poll index.
	result.Options, err = store.GetPollOptions(ctx, polls[pres.Currentpollindex.Int32].ID)
	if err != nil {
		return result, err
	}

	return result, err
}

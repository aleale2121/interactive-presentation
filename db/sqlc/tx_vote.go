package db

import (
	"context"

	"github.com/aleale2121/interactive-presentation/util"
	"github.com/google/uuid"
)

func (store *SQLStore) Vote(ctx context.Context, arg VoteParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		pollID,err:=uuid.Parse(arg.Pollid)
		if err != nil {
			return err
		}
		pID,err:=uuid.Parse(arg.PresentationID)
		if err != nil {
			return err
		}
		_, err = store.GetPollByPID(context.Background(), GetPollByPIDParams{
			ID:            pollID,
			Presentationid: pID,
		})
		if err != nil {
			return err
		}
	
		err = q.CreateVote(ctx, CreateVoteParams{
			ID:        util.RandomUUID(),
			Pollid:    pollID,
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

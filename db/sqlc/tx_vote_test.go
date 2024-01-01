package db

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aleale2121/interactive-presentation/util"
	"github.com/stretchr/testify/require"
)

func TestVoteCurrentPollTx(t *testing.T) {
	store := NewStore(testDB)
	createdPresenationID := createRandomPresentationWithPolls(t, 2)
	polls, err := store.GetPollsByPresentationID(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 2)
	var options []Option

	err = json.Unmarshal(polls[0].Options, &options)
	require.NoError(t, err)

	// initially no polls displayed for vote should return error
	err = store.VoteCurrentPollTx(context.Background(), VoteParams{
		PresentationID: createdPresenationID,
		Pollid:         polls[0].PollID,
		Optionkey:      options[0].Optionkey,
		Clientid:       util.RandomUUID().String(),
	})
	require.Error(t, err)

	// slide to the first presentation
	currPoll, err := testQueries.MoveForwardToNextPoll(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, currPoll)

	// create vote with existing keys
	for i := 0; i < len(options); i++ {
		err = store.VoteCurrentPollTx(context.Background(), VoteParams{
			PresentationID: createdPresenationID,
			Pollid:         currPoll.ID,
			Optionkey:      options[i].Optionkey,
			Clientid:       util.RandomUUID().String(),
		})
		require.NoError(t, err)
	}
	votes, err := testQueries.GetPollVotes(context.Background(), GetPollVotesParams{
		PresentationID: createdPresenationID,
		PollID:         currPoll.ID,
	})
	require.NoError(t, err)
	require.Len(t, votes, len(options))
}

func TestVoteDeadlock(t *testing.T) {
	store := NewStore(testDB)
	n := 100

	createdPresenationID := createRandomPresentationWithPolls(t, 2)
	polls, err := store.GetPollsByPresentationID(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 2)
	var options []Option

	err = json.Unmarshal(polls[0].Options, &options)
	require.NoError(t, err)

	// slide to the first presentation
	currPoll, err := testQueries.MoveForwardToNextPoll(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, currPoll)

	errs := make(chan error)
	l:=len(options)
	
	// run n concurrent vote transaction
	for i := 0; i < n; i++ {
		go func(j int) {
			err = store.VoteCurrentPollTx(context.Background(), VoteParams{
				PresentationID: createdPresenationID,
				Pollid:         currPoll.ID,
				Optionkey:      options[j].Optionkey,
				Clientid:       util.RandomUUID().String(),
			})
			errs <- err
		}(i%l)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated votes
	votes, err := testQueries.GetPollVotes(context.Background(), GetPollVotesParams{
		PresentationID: createdPresenationID,
		PollID:         currPoll.ID,
	})
	require.NoError(t, err)
	require.Len(t, votes, n)
}

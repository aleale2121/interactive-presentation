package db

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aleale2121/interactive-presentation/models"
	"github.com/aleale2121/interactive-presentation/util"
	"github.com/stretchr/testify/require"
)

func createRandomVote(t *testing.T) models.Vote {
	createdPresenationID := createRandomPresentationWithPolls(t, 2)

	polls, err := testQueries.GetPollsByPresentationID(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 2)
	var options []Option

	err = json.Unmarshal(polls[0].Options, &options)
	require.NoError(t, err)

	// create vote with existing keys
	for i := 0; i < len(options); i++ {
		err = testQueries.CreateVote(context.Background(), CreateVoteParams{
			Pollid:    polls[0].PollID,
			Optionkey: options[i].Optionkey,
			Clientid:  util.RandomUUID().String(),
		})
		require.NoError(t, err)
	}

	// create vote with non-existing option key
	err = testQueries.CreateVote(context.Background(), CreateVoteParams{
		Pollid:    polls[0].PollID,
		Optionkey: util.RandomUUID().String(),
		Clientid:  util.RandomUUID().String(),
	})
	require.Error(t, err)

	// create vote with non-existing poll id
	err = testQueries.CreateVote(context.Background(), CreateVoteParams{
		Pollid:    util.RandomUUID(),
		Optionkey: options[0].Optionkey,
		Clientid:  util.RandomUUID().String(),
	})
	require.Error(t, err)

	return models.Vote{
		PollID:         polls[0].PollID,
		PresentationID: createdPresenationID,
		VoteCount:      len(options),
	}
}

func TestCreateVote(t *testing.T) {
	createRandomVote(t)
}

func TestGetPollVotes(t *testing.T) {
	vote := createRandomVote(t)

	votes, err := testQueries.GetPollVotes(context.Background(), GetPollVotesParams{
		PresentationID: vote.PresentationID,
		PollID:         vote.PollID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, votes)
	require.Len(t, votes, vote.VoteCount)
	for _, vote := range votes {
		require.NotEmpty(t, vote)
	}
}

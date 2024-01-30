package persistence

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aleale2121/interactive-presentation/internal/constant/model"
	"github.com/aleale2121/interactive-presentation/pkg/random"
	"github.com/stretchr/testify/require"
)

func createRandomVote(t *testing.T) model.Vote {
	createdPresenationID := createRandomPresentationWithPolls(t, 2)

	polls, err := testQueries.ListPolls(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, polls)
	require.Len(t, polls, 2)
	var options []Option

	err = json.Unmarshal(polls[0].Options, &options)
	require.NoError(t, err)

	// create vote with existing keys
	for i := 0; i < len(options); i++ {
		err = testQueries.CreateVote(context.Background(), CreateVoteParams{
			Pollid:    polls[0].ID,
			Optionkey: options[i].Optionkey,
			Clientid:  random.RandomUUID().String(),
		})
		require.NoError(t, err)
	}

	// create vote with non-existing option key
	err = testQueries.CreateVote(context.Background(), CreateVoteParams{
		Pollid:    polls[0].ID,
		Optionkey: random.RandomUUID().String(),
		Clientid:  random.RandomUUID().String(),
	})
	require.Error(t, err)

	// create vote with non-existing poll id
	err = testQueries.CreateVote(context.Background(), CreateVoteParams{
		Pollid:    random.RandomUUID(),
		Optionkey: options[0].Optionkey,
		Clientid:  random.RandomUUID().String(),
	})
	require.Error(t, err)

	return model.Vote{
		PollID:         polls[0].ID,
		PresentationID: createdPresenationID,
		VoteCount:      len(options),
	}
}

func TestCreateVote(t *testing.T) {
	createRandomVote(t)
}

func TestGetVotes(t *testing.T) {
	vote := createRandomVote(t)

	votes, err := testQueries.GetVotes(context.Background(), GetVotesParams{
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

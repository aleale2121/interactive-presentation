package db

// import (
// 	"context"
// 	"testing"

// 	"github.com/aleale2121/interactive-presentation/util"
// 	"github.com/stretchr/testify/require"
// 	"github.com/google/uuid"

// )

// func createRandomVote(t *testing.T) string {
// 	key := util.RandomOptionKey()
// 	presentation := createRandomPresentation(t)
// 	poll := createPresentationPoll(t, presentation)

// 	arg := CreateOptionParams{
// 		Pollid:      poll.ID,
// 		Optionkey:   key,
// 		Optionvalue: util.RandomOptionValue(),
// 	}

// 	err := testQueries.CreateOption(context.Background(), arg)
// 	require.NoError(t, err)

// 	err = testQueries.CreateVote(context.Background(), CreateVoteParams{
// 		ID:        util.RandomUUID(),
// 		Pollid:    poll.ID,
// 		Optionkey: key,
// 		Clientid:  util.RandomUUID().String(),
// 	})

// 	require.NoError(t, err)
// 	return key
// }

// func createRandomPollVote(t *testing.T, pollID uuid.UUID, optionKey string) string {

// 	err := testQueries.CreateVote(context.Background(), CreateVoteParams{
// 		ID:        util.RandomUUID(),
// 		Pollid:    pollID,
// 		Optionkey: optionKey,
// 		Clientid:  util.RandomUUID().String(),
// 	})

// 	require.NoError(t, err)
// 	return optionKey
// }

// func TestCreateVote(t *testing.T) {
// 	createRandomVote(t)
// }

// func TestGetVote(t *testing.T) {
// 	pres := createRandomPresentation(t)
// 	poll := createPresentationPoll(t, pres)
// 	for i := 0; i < 10; i++ {
// 		optionKey := createRandomPollOption(t, poll)
// 		createRandomPollVote(t, poll.ID, optionKey)
// 	}

// 	votes, err := testQueries.GetVote(context.Background(), poll.ID)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, votes)
// 	require.Len(t, votes, 10)
// 	for _, vote := range votes {
// 		require.NotEmpty(t, vote)
// 	}
// }

// func TestGetPollVotes(t *testing.T) {
// 	presID := createRandomPresentation(t)
// 	poll := createPresentationPoll(t, presID)
// 	for i := 0; i < 10; i++ {
// 		optionKey := createRandomPollOption(t, poll)
// 		createRandomPollVote(t, poll.ID, optionKey)
// 	}

// 	votes, err := testQueries.GetPollVotes(context.Background(), GetPollVotesParams{
// 		ID:   presID,
// 		ID_2: poll.ID,
// 	})

// 	require.NoError(t, err)
// 	require.NotEmpty(t, votes)
// 	require.Len(t, votes, 10)
// 	for _, vote := range votes {
// 		require.NotEmpty(t, vote)
// 	}
// }

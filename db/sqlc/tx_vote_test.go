package db

// import (
// 	"context"
// 	"testing"

// 	"github.com/aleale2121/interactive-presentation/util"
// 	"github.com/stretchr/testify/require"
// )


// func TestVote(t *testing.T) {
// 	store := NewStore(testDB)
// 	n := 10

// 	presID := createRandomPresentation(t)
// 	poll := createPresentationPoll(t, presID)
// 	pollOptions := make([]string, 0, 10)
// 	for i := 0; i < n; i++ {
// 		optionKey := createRandomPollOption(t, poll)
// 		pollOptions = append(pollOptions, optionKey)
// 	}

// 	for i := 0; i < n; i++ {
// 		err := store.VoteCurrentPollTx(context.Background(), VoteParams{
// 			PresentationID: presID.String(),
// 			Pollid:         poll.ID.String(),
// 			Optionkey:      pollOptions[i%2],
// 			Clientid:       util.RandomUUID().String(),
// 		})
// 		require.NoError(t, err)
// 	}

// 	// check the final updated votes
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

// func TestVoteDeadlock(t *testing.T) {
// 	store := NewStore(testDB)
// 	n := 100

// 	presID := createRandomPresentation(t)
// 	poll := createPresentationPoll(t, presID)
// 	pollOptions := make([]string, 0, 10)
// 	for i := 0; i < n; i++ {
// 		optionKey := createRandomPollOption(t, poll)
// 		pollOptions = append(pollOptions, optionKey)
// 	}

// 	errs := make(chan error)

// 	// run n concurrent vote transaction
// 	for i := 0; i < n; i++ {
// 		go func(j int) {
// 			err := store.VoteCurrentPollTx(context.Background(), VoteParams{
// 				PresentationID: presID.String(),
// 				Pollid:         poll.ID.String(),
// 				Optionkey:      pollOptions[j],
// 				Clientid:       util.RandomUUID().String(),
// 			})

// 			errs <- err
// 		}(i / 2)
// 	}

// 	for i := 0; i < n; i++ {
// 		err := <-errs
// 		require.NoError(t, err)
// 	}
// 	// check the final updated votes
// 	votes, err := testQueries.GetPollVotes(context.Background(), GetPollVotesParams{
// 		ID:   presID,
// 		ID_2: poll.ID,
// 	})

// 	require.NoError(t, err)
// 	require.NotEmpty(t, votes)
// 	require.Len(t, votes, n)
// 	for _, vote := range votes {
// 		require.NotEmpty(t, vote)
// 	}

// }

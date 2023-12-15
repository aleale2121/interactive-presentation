package db

import (
	"context"
	"testing"

	"github.com/aleale2121/interactive-presentation/util"
	"github.com/stretchr/testify/require"
)

func TestUpdateCurrentPollTx(t *testing.T) {

	store := NewStore(testDB)
	n := 5

	presID := createRandomPresentation(t)
	for i := 0; i < n; i++ {
		createPresentationPoll(t, presID)
	}

	for i := 0; i < n-1; i++ {
		currPoll, err := store.UpdateCurrentPollTx(context.Background(), PollIndexParams{
			ID: presID.String(),
		})

		require.NoError(t, err)
		require.NotEmpty(t, currPoll)
	}

	// check the final updated presentation
	presentation, err := testQueries.GetPresentation(context.Background(), presID)

	require.NoError(t, err)
	require.NotEmpty(t, presentation)
	require.Equal(t, presentation.Currentpollindex.Int32, int32(n-1))
}

func TestVote(t *testing.T) {
	store := NewStore(testDB)
	n := 10

	presID := createRandomPresentation(t)
	poll := createPresentationPoll(t, presID)
	pollOptions := make([]string,0,10)
	for i := 0; i < n; i++ {
		optionKey := createRandomPollOption(t, poll)
		pollOptions = append(pollOptions, optionKey)
	}

	for i := 0; i < n; i++ {
		err := store.Vote(context.Background(), VoteParams{
			PresentationID: presID.String(),
			Pollid:         poll.ID.String(),
			Optionkey:      pollOptions[i%2],
			Clientid:       util.RandomUUID().String(),
		})
		require.NoError(t, err)
	}

	// check the final updated votes
	votes, err := testQueries.GetPollVotes(context.Background(), GetPollVotesParams{
		ID:   presID,
		ID_2: poll.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, votes)
	require.Len(t, votes, 10)
	for _, vote := range votes {
		require.NotEmpty(t, vote)
	}
}

func TestUpdateCurrentPollTxDeadlock(t *testing.T) {
		store := NewStore(testDB)
		n := 15

		presID := createRandomPresentation(t)
		for i := 0; i < n; i++ {
			createPresentationPoll(t, presID)
		}

		errs := make(chan error)
		results := make(chan CurrPollIndexResult)

		// run n concurrent vote transaction
		for i := 0; i < n-1; i++ {
			go func() {
				currPoll, err := store.UpdateCurrentPollTx(context.Background(), PollIndexParams{
					ID: presID.String(),
				})
				errs <- err
				results <- currPoll

			}()
		}

		for i := 0; i < n-1; i++ {
			err := <-errs
				require.NoError(t, err)
			currPoll := <-results
				require.NotEmpty(t, currPoll)
		}

		close(errs)
		close(results)
		// check the final updated presentation
		presentation, err := testQueries.GetPresentation(context.Background(), presID)

		require.NoError(t, err)
		require.NotEmpty(t, presentation)
		require.Equal(t, presentation.Currentpollindex.Int32, int32(n-1))
}

func TestVoteDeadlock(t *testing.T) {
		store := NewStore(testDB)
		n := 100

		presID := createRandomPresentation(t)
		poll := createPresentationPoll(t, presID)
		pollOptions := make([]string,0, 10)
		for i := 0; i < n; i++ {
			optionKey := createRandomPollOption(t, poll)
			pollOptions = append(pollOptions, optionKey)
		}

		errs := make(chan error)

		// run n concurrent vote transaction
		for i := 0; i < n; i++ {
			go func(j int) {
				err := store.Vote(context.Background(), VoteParams{
					PresentationID: presID.String(),
					Pollid:         poll.ID.String(),
					Optionkey:      pollOptions[j],
					Clientid:       util.RandomUUID().String(),
				})

				errs <- err
			}(i/2)
		}

		for i := 0; i < n; i++ {
			err := <-errs
			require.NoError(t, err)
		}
		// check the final updated votes
		votes, err := testQueries.GetPollVotes(context.Background(), GetPollVotesParams{
			ID:   presID,
			ID_2: poll.ID,
		})

		require.NoError(t, err)
		require.NotEmpty(t, votes)
		require.Len(t, votes, n)
		for _, vote := range votes {
			require.NotEmpty(t, vote)
		}
	
}

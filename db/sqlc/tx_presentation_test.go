package db

import (
	"context"
	"testing"

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
		currPoll, err := store.UpdateCurrentPollToForwardTx(context.Background(), presID)

		require.NoError(t, err)
		require.NotEmpty(t, currPoll)
	}

	// check the final updated presentation
	presentation, err := testQueries.GetPresentation(context.Background(), presID)

	require.NoError(t, err)
	require.NotEmpty(t, presentation)
	require.Equal(t, presentation.Currentpollindex.Int32, int32(n-1))
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
			currPoll, err := store.UpdateCurrentPollToForwardTx(context.Background(), presID)
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

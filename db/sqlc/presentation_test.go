package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomPresentation(t *testing.T) uuid.UUID {
	arg := sql.NullInt32{
		Int32: 0,
		Valid: true,
	}

	presentation, err := testQueries.CreatePresentation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, presentation)

	return presentation
}

func TestCreatePresentation(t *testing.T) {
	createRandomPresentation(t)
}

func createRandomPresentationWithPolls(t *testing.T) uuid.UUID {
	presentationID, err := testQueries.CreatePresentationAndPolls(context.Background(), []byte(presenatationData8))
	require.NoError(t, err)
	require.NotEmpty(t, presentationID)

	return presentationID
}

func TestCreatePresentationAndPolls(t *testing.T) {
	createRandomPresentationWithPolls(t)
}

func TestGetPresentation(t *testing.T) {
	createdPresenationID := createRandomPresentation(t)
	retrievedPresentation, err := testQueries.GetPresentation(context.Background(), createdPresenationID)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedPresentation)

	require.Equal(t, createdPresenationID, retrievedPresentation.ID)
}

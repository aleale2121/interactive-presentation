// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Querier interface {
	CreatePresentationAndPolls(ctx context.Context, dollar_1 json.RawMessage) (uuid.UUID, error)
	CreateVote(ctx context.Context, arg CreateVoteParams) error
	GetNextPoll(ctx context.Context, id uuid.UUID) (GetNextPollRow, error)
	GetPresentation(ctx context.Context, id uuid.UUID) (Presentation, error)
	GetPresentationAndPoll(ctx context.Context, arg GetPresentationAndPollParams) (GetPresentationAndPollRow, error)
	GetPresentationCurrentPoll(ctx context.Context, presentationid uuid.UUID) (GetPresentationCurrentPollRow, error)
	GetPreviousPoll(ctx context.Context, id uuid.UUID) (GetPreviousPollRow, error)
	GetVotes(ctx context.Context, arg GetVotesParams) ([]GetVotesRow, error)
	ListPolls(ctx context.Context, presentationid uuid.UUID) ([]ListPollsRow, error)
}

var _ Querier = (*Queries)(nil)

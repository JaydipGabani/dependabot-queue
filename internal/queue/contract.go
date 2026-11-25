// Package queue defines the polling contract. The single writer never checks out PR code.
package queue

import (
	"context"
	"time"
)

var Repositories = [...]string{
	"open-policy-agent/gatekeeper", "open-policy-agent/gatekeeper-library",
	"open-policy-agent/cert-controller", "open-policy-agent/frameworks",
}

type Candidate struct {
	Number  int
	Created time.Time
}

// State must be a fresh, internally consistent snapshot for this head and base.
// Policy and provenance failures set Eligible=false; unknowns never imply green.
type State struct {
	Number           int
	Head, Base       string
	Eligible         bool
	Reason           string
	Behind, Conflict bool
	CI               string // green, pending, missing, failed, transient
	RunID            int64
	RunAttempt       int
	ReviewDecision   string // APPROVED, REVIEW_REQUIRED, CHANGES_REQUESTED or empty
	OwnApproval      bool   // currently effective automated review on this head by expected user
	MergeReady       bool   // GitHub reports a clean PR; never bypass rules or a merge queue
}

type GitHub interface {
	List(context.Context, string) ([]Candidate, error)
	Inspect(context.Context, string, int) (State, error)
	Rebase(context.Context, string, State) error
	Retry(context.Context, string, State) error
	Approve(context.Context, string, State) error
	Merge(context.Context, string, State) error
}

package domain

import "errors"

var (
	ErrCannotChangeReviewersOnMerged = errors.New("cannot change reviewers on merged pull request")
	ErrTooManyReviewers              = errors.New("too many reviewers (max 2)")
)

var (
	ErrAlreadyExists = errors.New("entity already exists")
	ErrNotFound      = errors.New("entity not found")
	ErrMerged        = errors.New("pull request is already merged")
	ErrNotAssigned   = errors.New("user is not assigned to the pull request")
	ErrNoCandidate   = errors.New("no candidate available for reassignment")
)

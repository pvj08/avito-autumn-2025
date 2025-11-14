package domain

import "errors"

var (
	ErrCannotChangeReviewersOnMerged = errors.New("cannot change reviewers on merged pull request")
	ErrTooManyReviewers              = errors.New("too many reviewers (max 2)")
)

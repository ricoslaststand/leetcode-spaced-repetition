package utils

import (
	"time"

	"leetcode-spaced-repetition/models"

	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

func stateToFSRS(s models.ProblemCardStateStatus) fsrs.State {
	switch s {
	case models.ProblemCardStateLearning:
		return fsrs.Learning
	case models.ProblemCardStateReview:
		return fsrs.Review
	case models.ProblemCardStateRelearning:
		return fsrs.Relearning
	default:
		return fsrs.New
	}
}

func toCard(state models.ProblemCardState) fsrs.Card {
	return fsrs.Card{
		Due:           state.Due,
		Stability:     state.Stability,
		Difficulty:    state.Difficulty,
		ElapsedDays:   state.ElapsedDays,
		ScheduledDays: state.ScheduledDays,
		Reps:          state.Reps,
		Lapses:        state.Lapses,
		State:         stateToFSRS(state.State),
		LastReview:    state.LastReview,
	}
}

func toRating(confidenceLevel models.ConfidenceLevel) fsrs.Rating {
	switch confidenceLevel {
	case models.AgainConfidence:
		return fsrs.Again
	case models.HardConfidence:
		return fsrs.Hard
	case models.GoodConfidence:
		return fsrs.Good
	case models.EasyConfidence:
		return fsrs.Easy
	default:
		return fsrs.Again
	}
}

// ApplyReview computes the updated FSRS state after a review.
// Pass nil for cardState when recording the first-ever submission for a problem.
// The caller is responsible for setting ID, ProblemID, and UserID on the returned state.
func ApplyReview(
	cardState *models.ProblemCardState,
	confidenceLevel models.ConfidenceLevel,
	now time.Time,
) models.ProblemCardState {
	var card fsrs.Card
	if cardState == nil {
		card = fsrs.NewCard()
	} else {
		card = toCard(*cardState)
	}

	params := fsrs.DefaultParam()
	params.EnableShortTerm = false
	scheduler := params.NewLongTermScheduler(card, now)
	result := scheduler.Review(toRating(confidenceLevel))

	return models.ProblemCardState{
		Due:           result.Card.Due,
		State:         models.ProblemCardStateStatusFromFSRS(result.Card.State),
		Stability:     result.Card.Stability,
		Difficulty:    result.Card.Difficulty,
		ElapsedDays:   result.Card.ElapsedDays,
		ScheduledDays: result.Card.ScheduledDays,
		Reps:          result.Card.Reps,
		Lapses:        result.Card.Lapses,
		LastReview:    now,
	}
}

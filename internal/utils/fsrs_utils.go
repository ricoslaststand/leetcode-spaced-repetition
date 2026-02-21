package utils

import (
	"leetcode-spaced-repetition/models"

	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

func toCard(state models.CardState) fsrs.Card {
	return fsrs.Card{
		Stability:     state.Stability,
		Difficulty:    state.Difficulty,
		ElapsedDays:   state.ElapsedDays,
		ScheduledDays: state.ScheduledDays,
		Reps:          state.Reps,
		Lapses:        state.Lapses,
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

// func ApplyReview(
// 	cardState *models.CardState,
// 	confidenceLevel models.ConfidenceLevel,
// 	now time.Time,
// ) (models.CardState, time.Time) {
// 	var card fsrs.Card
// 	if cardState == nil {
// 		card = fsrs.NewCard()
// 	} else {
// 		card = toCard(*cardState)
// 	}

// 	defaultParams := fsrs.DefaultParam()
// 	defaultParams.EnableShortTerm = false
// 	scheduler := defaultParams.NewLongTermScheduler(card, now)
// 	// scheduler.Review()

// 	// 3. Apply the review
// 	result := scheduler.Review(toRating(confidenceLevel))

// 	// 4. Convert result back into your DB model
// 	updatedState := FSRSState{
// 		Stability:     result.Card.Stability,
// 		Difficulty:    result.Card.Difficulty,
// 		ElapsedDays:   result.Card.ElapsedDays,
// 		ScheduledDays: result.Card.ScheduledDays,
// 		Reps:          result.Card.Reps,
// 		Lapses:        result.Card.Lapses,
// 		LastReview:    now,
// 		NextReview:    nextReview,
// 	}

// 	// 5. Compute next review timestamp
// 	nextReview := now.AddDate(0, 0, result.Card.ScheduledDays)

// 	return updatedState, nextReview
// }

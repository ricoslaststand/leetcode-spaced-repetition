package internal

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// Difficulty represents the difficulty level of a LeetCode problem
type Difficulty int

const (
	Easy Difficulty = iota + 1
	Medium
	Hard
)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "EASY"
	case Medium:
		return "MEDIUM"
	case Hard:
		return "HARD"
	default:
		return "UNKNOWN"
	}
}

// ComfortLevel represents how comfortable the user is with a problem
type ComfortLevel int

const (
	Struggled   ComfortLevel = iota + 1 // Couldn't solve or needed significant help
	Difficult                           // Solved with difficulty, lots of hints
	Moderate                            // Solved with some effort/minor hints
	Comfortable                         // Solved easily with minor issues
	Mastered                            // Solved quickly and confidently
)

func (c ComfortLevel) String() string {
	switch c {
	case Struggled:
		return "STRUGGLED"
	case Difficult:
		return "DIFFICULT"
	case Moderate:
		return "MODERATE"
	case Comfortable:
		return "COMFORTABLE"
	case Mastered:
		return "MASTERED"
	default:
		return "UNKNOWN"
	}
}

// ProblemAttempt represents a single attempt at solving a problem
type ProblemAttempt struct {
	Timestamp        time.Time
	Notes            string
	ComfortLevel     ComfortLevel
	TimeTakenMinutes int
}

// LeetCodeProblem represents a LeetCode problem with spaced repetition data
type LeetCodeProblem struct {
	NextReviewDate  time.Time
	ProblemID       string
	Title           string
	Attempts        []ProblemAttempt
	Difficulty      Difficulty
	EaseFactor      float64
	IntervalDays    float64
	RepetitionCount int
}

// NewLeetCodeProblem creates a new LeetCode problem
func NewLeetCodeProblem(problemID, title string, difficulty Difficulty) *LeetCodeProblem {
	return &LeetCodeProblem{
		ProblemID:       problemID,
		Title:           title,
		Difficulty:      difficulty,
		Attempts:        make([]ProblemAttempt, 0),
		EaseFactor:      2.5,
		IntervalDays:    1.0,
		RepetitionCount: 0,
		NextReviewDate:  time.Now(),
	}
}

// ProblemStats contains statistics about a problem
type ProblemStats struct {
	NextReview       time.Time
	LatestComfort    *ComfortLevel
	ImprovementTrend *float64
	TotalAttempts    int
	AverageComfort   float64
	CurrentInterval  float64
	EaseFactor       float64
	RepetitionCount  int
}

// LeetCodeSpacedRepetition manages the spaced repetition system
type LeetCodeSpacedRepetition struct {
	problems map[string]*LeetCodeProblem

	// Base intervals by difficulty (in days)
	baseIntervals map[Difficulty]float64

	// Comfort level multipliers for interval calculation
	comfortMultipliers map[ComfortLevel]float64

	// Ease factor adjustments
	easeAdjustments map[ComfortLevel]float64
}

// NewLeetCodeSpacedRepetition creates a new spaced repetition system
func NewLeetCodeSpacedRepetition() *LeetCodeSpacedRepetition {
	return &LeetCodeSpacedRepetition{
		problems: make(map[string]*LeetCodeProblem),
		baseIntervals: map[Difficulty]float64{
			Easy:   1.0,
			Medium: 1.5,
			Hard:   2.0,
		},
		comfortMultipliers: map[ComfortLevel]float64{
			Struggled:   0.3,
			Difficult:   0.6,
			Moderate:    1.0,
			Comfortable: 1.3,
			Mastered:    1.6,
		},
		easeAdjustments: map[ComfortLevel]float64{
			Struggled:   -0.20,
			Difficult:   -0.15,
			Moderate:    0.0,
			Comfortable: 0.05,
			Mastered:    0.10,
		},
	}
}

// AddProblem adds a new problem to the system
func (sr *LeetCodeSpacedRepetition) AddProblem(problemID, title string, difficulty Difficulty) *LeetCodeProblem {
	problem := NewLeetCodeProblem(problemID, title, difficulty)
	sr.problems[problemID] = problem
	return problem
}

// RecordAttempt records an attempt at solving a problem
func (sr *LeetCodeSpacedRepetition) RecordAttempt(problemID string, comfortLevel ComfortLevel, timeTakenMinutes int, notes string) error {
	problem, exists := sr.problems[problemID]
	if !exists {
		return fmt.Errorf("problem %s not found", problemID)
	}

	attempt := ProblemAttempt{
		Timestamp:        time.Now(),
		ComfortLevel:     comfortLevel,
		TimeTakenMinutes: timeTakenMinutes,
		Notes:            notes,
	}

	problem.Attempts = append(problem.Attempts, attempt)
	sr.updateSpacedRepetitionParams(problem, comfortLevel)

	return nil
}

// updateSpacedRepetitionParams updates the spaced repetition parameters based on the latest attempt
func (sr *LeetCodeSpacedRepetition) updateSpacedRepetitionParams(problem *LeetCodeProblem, comfortLevel ComfortLevel) {
	// Update ease factor based on comfort level
	easeAdjustment := sr.easeAdjustments[comfortLevel]
	problem.EaseFactor = math.Max(1.3, problem.EaseFactor+easeAdjustment)

	// Calculate historical performance factor
	historicalFactor := sr.calculateHistoricalFactor(problem)

	// Difficulty factor - harder problems get longer initial intervals
	difficultyFactor := sr.baseIntervals[problem.Difficulty]

	// Recent performance factor based on current attempt
	recentFactor := sr.comfortMultipliers[comfortLevel]

	// If user struggled, reset to shorter interval regardless of history
	if comfortLevel == Struggled || comfortLevel == Difficult {
		problem.RepetitionCount = 0
		problem.IntervalDays = difficultyFactor * recentFactor
	} else {
		// Successful attempt - increase repetition count
		problem.RepetitionCount++

		switch problem.RepetitionCount {
		case 1:
			problem.IntervalDays = difficultyFactor
		case 2:
			problem.IntervalDays = difficultyFactor * 2
		default:
			// Use modified SM-2 algorithm with our factors
			problem.IntervalDays = problem.IntervalDays * problem.EaseFactor * historicalFactor * recentFactor
		}
	}

	// Cap the maximum interval to prevent problems from disappearing too long
	maxInterval := sr.getMaxInterval(problem.Difficulty)
	problem.IntervalDays = math.Min(problem.IntervalDays, maxInterval)

	// Set next review date
	problem.NextReviewDate = time.Now().Add(time.Duration(problem.IntervalDays * float64(24*time.Hour)))
}

// calculateHistoricalFactor calculates a factor based on historical performance, emphasizing recent attempts
func (sr *LeetCodeSpacedRepetition) calculateHistoricalFactor(problem *LeetCodeProblem) float64 {
	if len(problem.Attempts) <= 1 {
		return 1.0
	}

	// Take last 5 attempts, weight more recent ones higher
	recentAttempts := problem.Attempts
	if len(recentAttempts) > 5 {
		recentAttempts = problem.Attempts[len(problem.Attempts)-5:]
	}

	weights := []float64{1.0, 1.2, 1.4, 1.7, 2.0} // Most recent gets highest weight

	weightedComfort := 0.0
	totalWeight := 0.0

	for i, attempt := range recentAttempts {
		var weight float64
		if i < len(weights) {
			weight = weights[i]
		} else {
			weight = weights[len(weights)-1]
		}

		weightedComfort += float64(attempt.ComfortLevel) * weight
		totalWeight += weight
	}

	avgComfort := weightedComfort / totalWeight

	// Convert to multiplier (3.0 is neutral, representing MODERATE comfort)
	return 0.5 + (avgComfort/3.0)*0.8 // Range: ~0.67 to 1.33
}

// getMaxInterval gets maximum interval based on difficulty to prevent problems from being forgotten
func (sr *LeetCodeSpacedRepetition) getMaxInterval(difficulty Difficulty) float64 {
	maxIntervals := map[Difficulty]float64{
		Easy:   30, // 1 month max
		Medium: 21, // 3 weeks max
		Hard:   14, // 2 weeks max
	}
	return maxIntervals[difficulty]
}

// GetProblemsDueForReview gets all problems that are due for review, sorted by priority
func (sr *LeetCodeSpacedRepetition) GetProblemsDueForReview(includeNew bool) []*LeetCodeProblem {
	now := time.Now()
	var dueProblems []*LeetCodeProblem

	for _, problem := range sr.problems {
		if problem.NextReviewDate.Before(now) || problem.NextReviewDate.Equal(now) {
			dueProblems = append(dueProblems, problem)
		} else if includeNew && problem.RepetitionCount == 0 {
			dueProblems = append(dueProblems, problem)
		}
	}

	// Sort by priority: overdue problems first, then by difficulty
	sort.Slice(dueProblems, func(i, j int) bool {
		daysOverdueI := int(now.Sub(dueProblems[i].NextReviewDate).Hours() / 24)
		daysOverdueJ := int(now.Sub(dueProblems[j].NextReviewDate).Hours() / 24)

		if daysOverdueI != daysOverdueJ {
			return daysOverdueI > daysOverdueJ // More overdue first
		}

		return dueProblems[i].Difficulty < dueProblems[j].Difficulty // Easier first among equally overdue
	})

	return dueProblems
}

// GetProblemStats gets detailed statistics for a specific problem
func (sr *LeetCodeSpacedRepetition) GetProblemStats(problemID string) (*ProblemStats, error) {
	problem, exists := sr.problems[problemID]
	if !exists {
		return nil, fmt.Errorf("problem %s not found", problemID)
	}

	stats := &ProblemStats{
		TotalAttempts:   len(problem.Attempts),
		NextReview:      problem.NextReviewDate,
		CurrentInterval: problem.IntervalDays,
		EaseFactor:      problem.EaseFactor,
		RepetitionCount: problem.RepetitionCount,
	}

	if len(problem.Attempts) == 0 {
		return stats, nil
	}

	// Calculate average comfort
	totalComfort := 0
	for _, attempt := range problem.Attempts {
		totalComfort += int(attempt.ComfortLevel)
	}
	stats.AverageComfort = float64(totalComfort) / float64(len(problem.Attempts))

	// Get latest comfort
	latestComfort := problem.Attempts[len(problem.Attempts)-1].ComfortLevel
	stats.LatestComfort = &latestComfort

	// Calculate improvement trend (last 3 vs first 3 attempts)
	if len(problem.Attempts) >= 6 {
		earlySum := 0
		recentSum := 0

		for i := 0; i < 3; i++ {
			earlySum += int(problem.Attempts[i].ComfortLevel)
			recentSum += int(problem.Attempts[len(problem.Attempts)-3+i].ComfortLevel)
		}

		earlyAvg := float64(earlySum) / 3.0
		recentAvg := float64(recentSum) / 3.0
		trend := recentAvg - earlyAvg
		stats.ImprovementTrend = &trend
	}

	return stats, nil
}

// GetStudyPlan gets a recommended study plan for today
func (sr *LeetCodeSpacedRepetition) GetStudyPlan(maxProblems int) []*LeetCodeProblem {
	dueProblems := sr.GetProblemsDueForReview(true)
	if len(dueProblems) > maxProblems {
		return dueProblems[:maxProblems]
	}
	return dueProblems
}

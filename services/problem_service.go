package services

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"leetcode-spaced-repetition/internal/utils"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type ProblemService struct {
	problemRepo   repositories.ProblemRepository
	cardStateRepo repositories.ProblemCardStateRepository
	logger        *zap.Logger
}

func NewProblemsService(problemsRepo repositories.ProblemRepository, cardStateRepo repositories.ProblemCardStateRepository, logger *zap.Logger) *ProblemService {
	return &ProblemService{
		problemRepo:   problemsRepo,
		cardStateRepo: cardStateRepo,
		logger:        logger,
	}
}

func (s ProblemService) GetProblems(ctx context.Context, topics, difficulties []string, page, limit int) ([]models.Problem, error) {
	problems, err := s.problemRepo.GetProblems(ctx, topics, difficulties, page, limit)
	if err != nil {
		s.logger.Error("failed to get problems", zap.Error(err))
	}
	return problems, err
}

func (s ProblemService) GetProblemByID(c context.Context, id int) (*models.Problem, error) {
	problem, err := s.problemRepo.GetProblemByID(c, id)
	if err != nil {
		s.logger.Error("failed to get problem by ID", zap.Int("problemID", id), zap.Error(err))
	}
	return problem, err
}

func (s ProblemService) GetProblemSubmissions(c context.Context, problemIDs []int) ([]models.ProblemSubmissionWithDetails, error) {
	submissions, err := s.problemRepo.GetProblemSubmissions(c, problemIDs)
	if err != nil {
		s.logger.Error("failed to get problem submissions", zap.Error(err))
	}
	return submissions, err
}

func (s ProblemService) GetAllProblemTopics(c context.Context) ([]string, error) {
	topics, err := s.problemRepo.GetAllProblemTopics(c)
	if err != nil {
		s.logger.Error("failed to get all problem topics", zap.Error(err))
	}
	return topics, err
}

func (s ProblemService) GetTopicsForProblem(c context.Context, id int) ([]string, error) {
	topics, err := s.problemRepo.GetTopicsForProblem(c, id)
	if err != nil {
		s.logger.Error("failed to get topics for problem", zap.Int("problemID", id), zap.Error(err))
	}
	return topics, err
}

func (s ProblemService) GetDashboard(ctx context.Context, userID uuid.UUID, limit int) (models.DashboardData, error) {
	now := time.Now().UTC()

	due, err := s.cardStateRepo.GetDueReviewItems(ctx, userID, now, limit)
	if err != nil {
		s.logger.Error("failed to get due review items", zap.Error(err))
		return models.DashboardData{}, err
	}

	lowStability, err := s.cardStateRepo.GetLowestStabilityItems(ctx, userID, limit)
	if err != nil {
		s.logger.Error("failed to get lowest stability items", zap.Error(err))
		return models.DashboardData{}, err
	}

	overdueCount, err := s.cardStateRepo.GetOverdueProblemCount(ctx, userID, now)
	if err != nil {
		s.logger.Error("failed to get overdue problem count", zap.Error(err))
		return models.DashboardData{}, err
	}

	return models.DashboardData{
		Due:          due,
		LowStability: lowStability,
		OverdueCount: overdueCount,
	}, nil
}

func (s ProblemService) GetProblemCardState(ctx context.Context, userID uuid.UUID, problemID int) (*models.ProblemCardState, error) {
	return s.cardStateRepo.GetProblemCardStateByProblemID(ctx, userID, problemID)
}

func (s ProblemService) GetAllProblemsPastReviewDate(c context.Context, limit uint) ([]models.Problem, error) {
	problems, err := s.problemRepo.GetAllProblemsPastReviewDate(c, limit)
	if err != nil {
		s.logger.Error("failed to get problems past review date", zap.Error(err))
	}
	return problems, err
}

func (s ProblemService) GetAllSubmissionsForProblem(c context.Context, problemID int) ([]models.ProblemSubmission, error) {
	submissions, err := s.problemRepo.GetSubmissionsByProblemID(c, problemID)
	if err != nil {
		s.logger.Error("failed to get submissions for problem", zap.Int("problemID", problemID), zap.Error(err))
		return []models.ProblemSubmission{}, err
	}

	return submissions, nil
}

// parsedSubmissionRow holds a successfully parsed Excel row before any DB work.
type parsedSubmissionRow struct {
	submissionDate  time.Time
	timeTaken       *time.Duration
	confidenceLevel models.ConfidenceLevel
	rowNumber       int
	problemID       int
}

func (s ProblemService) ImportSubmissions(ctx context.Context, r io.Reader) (models.ImportSubmissionsResult, error) {
	result := models.ImportSubmissionsResult{
		Errors: make([]models.ImportSubmissionRowError, 0),
	}

	f, err := excelize.OpenReader(r)
	if err != nil {
		return result, fmt.Errorf("failed to parse excel file: %w", err)
	}
	defer func() { _ = f.Close() }()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return result, fmt.Errorf("failed to read sheet: %w", err)
	}

	// Phase 1: Parse all rows, collecting errors without touching the DB.
	var validRows []parsedSubmissionRow

	for i, row := range rows {
		if i == 0 {
			continue // skip header row
		}
		rowNumber := i + 1 // 1-based for error messages

		// Skip blank rows
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}

		if len(row) < 4 {
			result.Errors = append(result.Errors, models.ImportSubmissionRowError{
				Row:    rowNumber,
				Reason: "row has fewer than 4 columns",
			})
			continue
		}

		// Column 0: "Problem #"
		problemNumStr := strings.TrimSpace(row[0])
		problemNumber, err := strconv.Atoi(problemNumStr)
		if err != nil || problemNumber < 1 {
			result.Errors = append(result.Errors, models.ImportSubmissionRowError{
				Row:    rowNumber,
				Reason: fmt.Sprintf("%q is not a valid problem number", problemNumStr),
			})
			continue
		}

		// Column 1: "Date" (YYYY-MM-DD)
		dateStr := strings.TrimSpace(row[1])
		submissionDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// Fallback: Excel may serialize dates as serial numbers (float strings)
			if serialFloat, parseErr := strconv.ParseFloat(dateStr, 64); parseErr == nil {
				if t, exErr := excelize.ExcelDateToTime(serialFloat, false); exErr == nil {
					submissionDate = t
				} else {
					result.Errors = append(result.Errors, models.ImportSubmissionRowError{
						Row:    rowNumber,
						Reason: fmt.Sprintf("%q is not a valid date (expected YYYY-MM-DD)", dateStr),
					})
					continue
				}
			} else {
				result.Errors = append(result.Errors, models.ImportSubmissionRowError{
					Row:    rowNumber,
					Reason: fmt.Sprintf("%q is not a valid date (expected YYYY-MM-DD)", dateStr),
				})
				continue
			}
		}
		// Explicitly preserve year/month/day in UTC to prevent timezone drift
		// when the driver encodes the time.Time value for the DATE column.
		submissionDate = time.Date(submissionDate.Year(), submissionDate.Month(), submissionDate.Day(), 0, 0, 0, 0, time.UTC)

		// Column 2: "Time Taken" (e.g. "1m50s", "1 m 36 s") — optional
		var timeTaken *time.Duration
		timeTakenStr := strings.ReplaceAll(strings.TrimSpace(row[2]), " ", "")
		if timeTakenStr != "" {
			d, parseErr := time.ParseDuration(timeTakenStr)
			if parseErr != nil {
				result.Errors = append(result.Errors, models.ImportSubmissionRowError{
					Row:    rowNumber,
					Reason: fmt.Sprintf("%q is not a valid time duration", timeTakenStr),
				})
				continue
			}
			timeTaken = &d
		}

		// Column 3: "Confidence" (1-4)
		confidenceStr := strings.TrimSpace(row[3])
		confidenceLevel, err := models.DetermineConfidenceLevelFromString(confidenceStr)
		if err != nil {
			result.Errors = append(result.Errors, models.ImportSubmissionRowError{
				Row:    rowNumber,
				Reason: fmt.Sprintf("%q is not a valid confidence level (expected 1-4)", confidenceStr),
			})
			continue
		}

		validRows = append(validRows, parsedSubmissionRow{
			rowNumber:       rowNumber,
			problemID:       problemNumber,
			submissionDate:  submissionDate,
			timeTaken:       timeTaken,
			confidenceLevel: confidenceLevel,
		})
	}

	// Phase 2: Sort valid rows by submissionDate ascending for correct FSRS chronological order.
	sort.Slice(validRows, func(i, j int) bool {
		return validRows[i].submissionDate.Before(validRows[j].submissionDate)
	})

	// Phase 3: Save submissions and compute FSRS card states.
	userID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

	for _, row := range validRows {
		inserted, err := s.problemRepo.SaveProblemSubmission(ctx, row.problemID, userID, row.submissionDate, row.timeTaken, row.confidenceLevel, nil)
		if err != nil {
			s.logger.Error("failed to save imported submission",
				zap.Int("row", row.rowNumber),
				zap.Int("problemID", row.problemID),
				zap.Error(err),
			)
			result.Errors = append(result.Errors, models.ImportSubmissionRowError{
				Row:    row.rowNumber,
				Reason: fmt.Sprintf("failed to save: %s", err.Error()),
			})
			continue
		}

		if !inserted {
			result.Skipped++
			continue
		}

		result.Imported++

		// Compute and upsert the FSRS card state for this newly inserted submission.
		currentState, err := s.cardStateRepo.GetProblemCardStateByProblemID(ctx, userID, row.problemID)
		if err != nil {
			s.logger.Error("failed to fetch card state for FSRS computation",
				zap.Int("row", row.rowNumber),
				zap.Int("problemID", row.problemID),
				zap.Error(err),
			)
			continue
		}

		newState := utils.ApplyReview(currentState, row.confidenceLevel, row.submissionDate)
		newState.ProblemID = row.problemID
		newState.UserID = userID

		if err := s.cardStateRepo.UpsertProblemCardState(ctx, &newState); err != nil {
			s.logger.Error("failed to upsert card state",
				zap.Int("row", row.rowNumber),
				zap.Int("problemID", row.problemID),
				zap.Error(err),
			)
		}
	}

	return result, nil
}

func (s ProblemService) SaveProblemSubmission(
	c context.Context,
	problemID int,
	userID uuid.UUID,
	date time.Time,
	timeTaken *time.Duration,
	confidenceLevel models.ConfidenceLevel,
	language *models.SubmissionLanguage,
) error {
	inserted, err := s.problemRepo.SaveProblemSubmission(c, problemID, userID, date, timeTaken, confidenceLevel, language)
	if err != nil {
		s.logger.Error("failed to save problem submission", zap.Int("problemID", problemID), zap.Error(err))
		return err
	}

	if !inserted {
		return nil
	}

	currentState, err := s.cardStateRepo.GetProblemCardStateByProblemID(c, userID, problemID)
	if err != nil {
		s.logger.Error("failed to fetch card state for FSRS computation", zap.Int("problemID", problemID), zap.Error(err))
		return nil
	}

	newState := utils.ApplyReview(currentState, confidenceLevel, date)
	newState.ProblemID = problemID
	newState.UserID = userID

	if err := s.cardStateRepo.UpsertProblemCardState(c, &newState); err != nil {
		s.logger.Error("failed to upsert card state", zap.Int("problemID", problemID), zap.Error(err))
	}

	return nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"log"
	"os"
	"strconv"
	"strings"
)

type LeetcodeProblem struct {
	Name       string
	Difficulty models.QuestionDifficulty
	Slug       string
	Acceptance float64
	Frequency  float64
	TimeTag    string
}

type questionsResponse struct {
	Questions []struct {
		Title       string   `json:"title"`
		ProblemID   string   `json:"problem_id"`
		FrontendID  string   `json:"frontend_id"`
		Difficulty  string   `json:"difficulty"`
		ProblemSlug string   `json:"problem_slug"`
		Topics      []string `json:"topics"`
		Description string   `json:"description"`
		Constraints []string `json:"constraints"`
		FollowUps   []string `json:"follow_ups"`
		Hints       []string `json:"hints"`
		Solution    string   `json:"solution"`
	} `json:"questions"`
}

type TimeTag string

const (
	LastThirtyDays    = "last_thirty_days"
	LastThreeMonths   = "last_three_months"
	LastSixMonths     = "last_six_months"
	MoreThanSixMonths = "more_than_six_months"
)

const (
	DifficultyIdx int = iota
	NameIdx
	FrequencyIdx
	AcceptanceIdx
	LinkIdx
)

func convertFileNameToTimeTag(filename string) (TimeTag, error) {
	lCase := strings.ToLower(filename)

	if strings.Contains(lCase, "more than six months") {
		return MoreThanSixMonths, nil
	} else if strings.Contains(lCase, "six months") {
		return LastSixMonths, nil
	} else if strings.Contains(lCase, "three months") {
		return LastThreeMonths, nil
	} else if strings.Contains(lCase, "thirty days") {
		return LastThirtyDays, nil
	}
	return LastThirtyDays, fmt.Errorf("'%s' is not a valid time tag", filename)
}

func convertStringToDifficulty(diffStr string) (models.QuestionDifficulty, error) {
	lowerDiffStr := strings.ToLower(strings.Trim(diffStr, " `"))

	switch lowerDiffStr {
	case "easy":
		return models.EasyDifficulty, nil
	case "medium":
		return models.MediumDifficulty, nil
	case "hard":
		return models.HardDifficulty, nil
	default:
		return models.EasyDifficulty, fmt.Errorf("'%s' is not a valid difficulty", diffStr)
	}
}

func getSlugFromLink(link string) string {
	parts := strings.Split(link, "/")

	return parts[len(parts)-1]
}

func main() {
	config, err := internal.GetConfig()
	if err != nil {
		panic(err)
	}

	db, err := internal.GetDBConnFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	questionsRepo := repositories.NewQuestionPostgresRepository(db)

	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}

	tagsFileContent, err := os.ReadFile("merged_problems.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var mergedQuestions questionsResponse
	err = json.Unmarshal(tagsFileContent, &mergedQuestions)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, question := range mergedQuestions.Questions {
		questionDifficulty, err := convertStringToDifficulty(question.Difficulty)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		intQuestionID, err := strconv.Atoi(question.ProblemID)
		if err != nil {
			return
		}

		err = questionsRepo.SaveQuestion(context.Background(), &models.Question{
			ID:          intQuestionID,
			Title:       question.Title,
			Description: question.Description,
			Slug:        question.ProblemSlug,
			Difficulty:  questionDifficulty,
			Tags:        question.Topics,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, tag := range question.Topics {
			err = questionsRepo.SaveQuestionTag(context.Background(), intQuestionID, tag)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}

}

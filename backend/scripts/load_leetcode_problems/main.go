package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	config "leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"os"
)

type Difficulty struct {
	Level int `json:"level"`
}

type Stat struct {
	QuestionID              int     `json:"question_id"`
	QuestionArticleLive     *bool   `json:"question__article__live"`
	QuestionArticleSlug     *string `json:"question__article__slug"`
	QuestionArticleHasVideo *bool   `json:"question__article__has_video_solution"`
	QuestionTitle           string  `json:"question__title"`
	QuestionTitleSlug       string  `json:"question__title_slug"`
	QuestionHide            bool    `json:"question__hide"`
	TotalACs                int     `json:"total_acs"`
	TotalSubmitted          int     `json:"total_submitted"`
	FrontendQuestionID      int     `json:"frontend_question_id"`
	IsNewQuestion           bool    `json:"is_new_question"`
}

type StatStatusPair struct {
	Stat       Stat       `json:"stat"`
	Status     *string    `json:"status"` // Could be null (use *string for nullable field)
	Difficulty Difficulty `json:"difficulty"`
	PaidOnly   bool       `json:"paid_only"`
	IsFavor    bool       `json:"is_favor"`
	Frequency  int        `json:"frequency"`
	Progress   int        `json:"progress"`
}

type ProblemData struct {
	UserName        string           `json:"user_name"`
	NumSolved       int              `json:"num_solved"`
	NumTotal        int              `json:"num_total"`
	AcEasy          int              `json:"ac_easy"`
	AcMedium        int              `json:"ac_medium"`
	AcHard          int              `json:"ac_hard"`
	StatStatusPairs []StatStatusPair `json:"stat_status_pairs"`
}

var db *sql.DB

func main() {
	config, err := config.GetConfig()
	if err != nil {

	}

	cfg := post

	fmt.Printf("Loading leetcode problems...")

	filepath := "leetcode_problems.json"
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Println("We are here")

	var responseData ProblemData
	err = json.Unmarshal(fileContent, &responseData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	numOfQuestions := len(responseData.StatStatusPairs)

	questionRepo := repositories.NewQuestionPostgresRepository(db)

	fmt.Printf("We have successfully downloaded %d problems\n", numOfQuestions)

	var questions []models.Question
	for i := 0; i < numOfQuestions; i++ {
		currQuestion := responseData.StatStatusPairs[i]
		questionDifficulty, err := models.DetermineDifficulty(currQuestion.Difficulty.Level)

		if err == nil {
			fmt.Println(err.Error())
			return
		}

		questions = append(questions, models.Question{
			ID:         currQuestion.Stat.QuestionID,
			Title:      currQuestion.Stat.QuestionTitle,
			Difficulty: questionDifficulty,
		})
	}

	for i := 0; i < len(questions); i++ {
		questionRepo.SaveQuestion(questions[i])
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"leetcode-spaced-repetition/internal"
	config "leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
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

type Questions struct {
	Questions []struct {
		ID   string   `json:"problem_id"`
		Tags []string `json:"topics"`
	} `json:"questions"`
}

func main() {
	config, err := config.GetConfig()
	if err != nil {
		return
	}

	db, err := internal.GetDBConnFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		difficultyMap := map[int]models.QuestionDifficulty{
			1: models.EasyDifficulty,
			2: models.MildDifficulty,
			3: models.HardDifficulty,
		}
		questionDifficulty, ok := difficultyMap[currQuestion.Difficulty.Level]
		if !ok {
			fmt.Printf("unrecognized difficulty level: %d\n", currQuestion.Difficulty.Level)
			return
		}

		questions = append(questions, models.Question{
			ID:          currQuestion.Stat.QuestionID,
			Title:       currQuestion.Stat.QuestionTitle,
			Description: "",
			Slug:        currQuestion.Stat.QuestionTitleSlug,
			Difficulty:  questionDifficulty,
		})
	}

	fmt.Println("Saving questions to the database")

	c := context.Background()

	for i := 0; i < len(questions); i++ {
		err = questionRepo.SaveQuestion(c, &questions[i])
		if err != nil {
			panic(err)
		}
	}

	// Adding question tags
	tagsFileContent, err := os.ReadFile("merged_problems.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var mergedQuestions Questions
	err = json.Unmarshal(tagsFileContent, &mergedQuestions)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(mergedQuestions.Questions); i++ {
		question := mergedQuestions.Questions[i]
		questionId, err := strconv.Atoi(question.ID)
		if err != nil {
			fmt.Printf("Invalid questionId\n")
		}

		for j := 0; j < len(question.Tags); j++ {
			fmt.Printf("Processing questionId %d, tag %s\n", questionId, question.Tags[j])
			err = questionRepo.SaveQuestionTag(c, questionId, question.Tags[j])
			if err != nil {
				panic(err)
			}
		}
	}
}

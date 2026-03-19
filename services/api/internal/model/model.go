package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type TaskType string

const (
	TaskCoding   TaskType = "coding"
	TaskLeetcode TaskType = "leetcode"
	TaskContent  TaskType = "content"
)

type Task struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Type       TaskType `json:"type"`
	Date       string   `json:"date"`
	Completed  bool     `json:"completed"`
	CreatedAt  time.Time `json:"created_at"`
}

type IdeaType string
type IdeaStatus string

const (
	IdeaTypeReel     IdeaType = "reel"
	IdeaTypeTweet    IdeaType = "tweet"
	IdeaTypeThread   IdeaType = "thread"
	IdeaTypeLinkedin IdeaType = "linkedin"

	IdeaStatusIdea   IdeaStatus = "idea"
	IdeaStatusReady  IdeaStatus = "ready"
	IdeaStatusPosted IdeaStatus = "posted"
)

type Idea struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Hook      string     `json:"hook"`
	Idea      string     `json:"idea"`
	Type      IdeaType   `json:"type"`
	Status    IdeaStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type LeetCodeLog struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	ProblemName  string     `json:"problem_name"`
	Difficulty   Difficulty `json:"difficulty"`
	Approach     string     `json:"approach"`
	Mistake      string     `json:"mistake"`
	TimeTaken    *int       `json:"time_taken"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CodingLog struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Session struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

type OpportunityType string
type OpportunityStage string

const (
	OppTypeJob      OpportunityType = "job"
	OppTypeFreelance OpportunityType = "freelance"

	OppStageApplied   OpportunityStage = "applied"
	OppStageInterview OpportunityStage = "interview"
	OppStageClosed    OpportunityStage = "closed"
)

type Opportunity struct {
	ID        uuid.UUID         `json:"id"`
	UserID    uuid.UUID         `json:"user_id"`
	Name      string            `json:"name"`
	Type      OpportunityType   `json:"type"`
	Stage     OpportunityStage  `json:"stage"`
	Source    string            `json:"source"`
	Notes     string            `json:"notes"`
	CreatedAt time.Time         `json:"created_at"`
}

type FinanceType string

const (
	FinanceSalary   FinanceType = "salary"
	FinanceFreelance FinanceType = "freelance"
	FinanceOther    FinanceType = "other"
)

type Finance struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	Amount    float64     `json:"amount"`
	Type      FinanceType `json:"type"`
	Note      string      `json:"note"`
	Date      string      `json:"date"`
	CreatedAt time.Time   `json:"created_at"`
}

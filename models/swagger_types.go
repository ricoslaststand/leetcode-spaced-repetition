package models

// Concrete type aliases for the generic Pagaination[T] type.
// swag v1 does not support Go generics, so these aliases are needed
// to generate correct Swagger schema definitions for paginated responses.

type QuestionPage = Pagaination[Question]
type QuestionSubmissionPage = Pagaination[QuestionSubmission]
type QuestionSubmissionWithDetailsPage = Pagaination[QuestionSubmissionWithDetails]

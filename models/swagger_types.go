package models

// Concrete type aliases for the generic Pagaination[T] type.
// swag v1 does not support Go generics, so these aliases are needed
// to generate correct Swagger schema definitions for paginated responses.

type ProblemPage = Pagaination[Problem]
type ProblemSubmissionPage = Pagaination[ProblemSubmission]
type ProblemSubmissionWithDetailsPage = Pagaination[ProblemSubmissionWithDetails]

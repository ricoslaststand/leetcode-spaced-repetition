# Leetcode Buddy

## Architecture
- frontend: React-based single page application
- backend: Go-based backend that interacts with PostgreSQL database

# Technology
- `frontend`
    - `react`
    - `tailwindcss`
    - `shadcn`

## Code Structure
- `frontend`: holds all the frontend code for the application
    - `src`
        - `pages`: each component maps to a unique URL
        - `components`: includes shared components across pages in addition to components for individual pages
            - `ProblemDifficultyTag`: the canonical component for rendering a problem's difficulty badge. Use this everywhere a difficulty label needs to be shown.
- `applications`: holds the different applications for the Go code. For example, it holds the API, one-off scripts for backfilling data, a worker.
- `controllers`: the different HTTP API handlers 
- `respositories`: holds both the interfaces and datastore implementations when interacting with the database

### Naming convention
- Every struct in the `repositories` folder that implements the a repository interface should have the following naming convention `<domain object>_<data store type>_repository.go`. For example, if there was an interface called `TicketRepository` that was in a file called `repositories/ticket_repository` and you wanted to add a MySQL datastore to your project, then the file should be called `repositories/ticket_mysql_repository.go`.


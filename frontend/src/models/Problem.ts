export const ConfidenceLevel = {
    Again: 1,
    Hard: 2,
    Good: 3,
    Easy: 4,
} as const;
export type ConfidenceLevel = typeof ConfidenceLevel[keyof typeof ConfidenceLevel];

const confidenceLevelStrings: Record<ConfidenceLevel, string> = {
    [ConfidenceLevel.Again]: 'again',
    [ConfidenceLevel.Hard]: 'hard',
    [ConfidenceLevel.Good]: 'good',
    [ConfidenceLevel.Easy]: 'easy',
}
export function confidenceLevelToString(level: ConfidenceLevel): string {
    return confidenceLevelStrings[level]
}

export type SubmissionLanguage =
    'cpp' | 'java' | 'python' | 'python3' | 'c' | 'csharp' |
    'javascript' | 'typescript' | 'php' | 'swift' | 'kotlin' |
    'dart' | 'golang' | 'ruby' | 'scala' | 'rust' | 'racket' |
    'erlang' | 'elixir'

export const SUBMISSION_LANGUAGES: { value: SubmissionLanguage; label: string }[] = [
    { value: 'cpp',        label: 'C++' },
    { value: 'java',       label: 'Java' },
    { value: 'python3',    label: 'Python 3' },
    { value: 'python',     label: 'Python 2' },
    { value: 'c',          label: 'C' },
    { value: 'csharp',     label: 'C#' },
    { value: 'javascript', label: 'JavaScript' },
    { value: 'typescript', label: 'TypeScript' },
    { value: 'php',        label: 'PHP' },
    { value: 'swift',      label: 'Swift' },
    { value: 'kotlin',     label: 'Kotlin' },
    { value: 'dart',       label: 'Dart' },
    { value: 'golang',     label: 'Go' },
    { value: 'ruby',       label: 'Ruby' },
    { value: 'scala',      label: 'Scala' },
    { value: 'rust',       label: 'Rust' },
    { value: 'racket',     label: 'Racket' },
    { value: 'erlang',     label: 'Erlang' },
    { value: 'elixir',     label: 'Elixir' },
]

export type ProblemReviewItem = {
    problemId: number;
    title: string;
    slug: string;
    difficulty: string;
    due: string;
    stability: number;
}

export type DashboardData = {
    due: ProblemReviewItem[];
    lowStability: ProblemReviewItem[];
    overdueCount: number;
}

export type ProblemSubmissionWithDetails = {
    id: number;
    problemId: number;
    confidenceLevel: string;
    timeTaken: number | null;
    problem: {
        id: number;
        title: string;
        slug: string;
        difficulty: string;
        topics: string[];
    };
}

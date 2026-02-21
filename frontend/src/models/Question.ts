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

export type QuestionSubmissionWithDetails = {
    id: number;
    questionId: number;
    confidenceLevel: string;
    timeTaken: number | null;
    question: {
        id: number;
        title: string;
        slug: string;
        difficulty: string;
        tags: string[];
    };
}

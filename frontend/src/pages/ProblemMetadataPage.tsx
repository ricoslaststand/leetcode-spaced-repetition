import { type FC } from 'react'

import { useParams } from '@tanstack/react-router'
import { Spinner } from "../components/ui/spinner"

import { useProblem, useProblemSubmissions } from "../hooks/api"
import ProblemSubmissionsTable from '../components/ProblemSubmissionsTable'
import ProblemDifficultyTag from '../components/ProblemDifficultyTag'

const ProblemMetadataPage: FC = () => {
    const problemId = useParams({
        from: '/problems/$problemId/submissions',
        select: (params) => params.problemId
    });

    const { data, isLoading } = useProblemSubmissions(problemId)
    const { data: problem } = useProblem(problemId)

    return (
        <div>
            <div className="mb-6">
                <h1 className="text-2xl font-semibold tracking-tight">
                    {problem ? `${problemId}. ${problem.title}` : `#${problemId}`}
                </h1>
                {problem && (
                    <div className="flex items-center gap-4 mt-2 text-sm text-muted-foreground">
                        <ProblemDifficultyTag difficulty={problem.difficulty} />
                        <span>Stability: {problem.stability != null ? problem.stability.toFixed(2) : '—'}</span>
                        <span># of Submissions: {data?.data?.length ?? 0}</span>
                    </div>
                )}
            </div>
            {isLoading && (
                <div className="flex justify-center py-12">
                    <Spinner />
                </div>
            )}
            {data && <ProblemSubmissionsTable submissions={data.data} />}
        </div>
    )
}

export default ProblemMetadataPage

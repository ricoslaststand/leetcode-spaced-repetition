import { type FC } from 'react'

import { useParams } from '@tanstack/react-router'
import { Spinner } from "../components/ui/spinner"

import { useProblemSubmissions } from "../hooks/api"
import ProblemSubmissionsTable from '../components/ProblemSubmissionsTable'

const ProblemMetadataPage: FC = () => {
    const problemId = useParams({
        from: '/problems/$problemId',
        select: (params) => params.problemId
    });

    const { data, isLoading } = useProblemSubmissions(problemId)

    return (
        <div>
            <h1 className="text-2xl font-semibold tracking-tight mb-6">
                Submission History
                <span className="text-muted-foreground font-normal text-lg ml-2">
                    #{problemId}
                </span>
            </h1>
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

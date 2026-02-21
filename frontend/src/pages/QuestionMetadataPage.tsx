import { type FC } from 'react'

import { useParams } from '@tanstack/react-router'
import { Spinner } from "../components/ui/spinner"

import { useQuestionSubmissions } from "../hooks/api"
import QuestionSubmissionsTable from '../components/QuestionSubmissionsTable'

const QuestionMetadataPage: FC = () => {
    const questionId = useParams({
        from: '/questions/$questionId',
        select: (params) => params.questionId
    });

    const { data, isLoading } = useQuestionSubmissions(questionId)

    return (
        <div>
            <h1 className="text-2xl font-semibold tracking-tight mb-6">
                Submission History
                <span className="text-muted-foreground font-normal text-lg ml-2">
                    #{questionId}
                </span>
            </h1>
            {isLoading && (
                <div className="flex justify-center py-12">
                    <Spinner />
                </div>
            )}
            {data && <QuestionSubmissionsTable submissions={data.data} />}
        </div>
    )
}

export default QuestionMetadataPage

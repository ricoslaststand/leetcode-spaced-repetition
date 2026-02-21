import axios from 'axios';
import qs from 'qs';

const instance = axios.create({
    baseURL: "http://localhost:8000",
    timeout: 5_000
})

export const getAllQuestionTags = async () => {
    const response = await instance.get("/questions/tags")
    return response.data
}

export const getQuestionByID = async (questionID: string) => {
    const response = await instance.get(`/questions/${questionID}`)
    return response.data
}

export const getAllQuestions = async (tags: string[]) => {
    const response = await instance.get("/questions", {
        params: {
            "tags": tags
        },
        paramsSerializer: params => {
            return qs.stringify(params, {
                arrayFormat: "repeat"
            })
        }
    })

    return response.data
}

export const getQuestionSubmissions = async (id: number) => {
    const response = await instance.get(`/questions/${id}/submissions`)

    return response.data
}

export const getQuestionSubmissionsV2 = async (questionIds: number[]) => {
    const response = await instance.get('/questions/submissions',
        {
            params: {
                questionId: questionIds
            },
            paramsSerializer: params => {
                return qs.stringify(params, {
                    arrayFormat: "comma"
                })
            }
        }
    )

    return response.data
}

export const createQuestionSubmission
    = async (questionID: number, confidenceLevel: string, timeTaken?: number) => {
    const response = await instance.post(`/questions/submissions`, {
        questionID,
        confidenceLevel,
        timeTaken
    })
    
    return response.data
}

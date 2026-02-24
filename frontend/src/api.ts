import axios from 'axios';
import qs from 'qs';

const instance = axios.create({
    baseURL: "http://localhost:8000",
    timeout: 5_000
})

export const getAllProblemTopics = async () => {
    const response = await instance.get("/problems/topics")
    return response.data
}

export const getProblemByID = async (problemID: string) => {
    const response = await instance.get(`/problems/${problemID}`)
    return response.data
}

export const getAllProblems = async (topics: string[]) => {
    const response = await instance.get("/problems", {
        params: {
            "topics": topics
        },
        paramsSerializer: params => {
            return qs.stringify(params, {
                arrayFormat: "repeat"
            })
        }
    })

    return response.data
}

export const getProblemSubmissions = async (id: number) => {
    const response = await instance.get(`/problems/${id}/submissions`)

    return response.data
}

export const getProblemSubmissionsV2 = async (problemIds: number[]) => {
    const response = await instance.get('/problems/submissions',
        {
            params: {
                problemId: problemIds
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

export const createProblemSubmission
    = async (problemID: number, confidenceLevel: string, timeTaken?: number) => {
    const response = await instance.post(`/problems/submissions`, {
        problemId: problemID,
        confidenceLevel,
        timeTaken
    })

    return response.data
}

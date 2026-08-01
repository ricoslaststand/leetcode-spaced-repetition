import axios from 'axios';
import qs from 'qs';
import type { DashboardData } from './models/Problem';

// In production the SPA and the API sit behind one Traefik host, so a relative path keeps
// requests same-origin and lets the ForwardAuth session cookie ride along. The Vite dev
// server sets VITE_API_BASE_URL to reach the API container directly.
const instance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL ?? "/api",
    timeout: 5_000,
    withCredentials: true
})

export const getAllProblemTopics = async () => {
    const response = await instance.get("/problems/topics")
    return response.data
}

export const getProblemByID = async (problemID: string) => {
    const response = await instance.get(`/problems/${problemID}`)
    return response.data
}

export const getAllProblems = async (topics: string[], difficulties: string[]) => {
    const response = await instance.get("/problems", {
        params: { topics, difficulties },
        paramsSerializer: params => qs.stringify(params, { arrayFormat: "repeat" })
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
    = async (problemID: number, confidenceLevel: string, timeTaken?: number, language?: string) => {
    const response = await instance.post(`/problems/submissions`, {
        problemId: problemID,
        confidenceLevel,
        timeTaken,
        language
    })

    return response.data
}

export const getDashboard = async (limit: number): Promise<DashboardData> => {
    const response = await instance.get('/dashboard', { params: { limit } })
    return response.data
}

export const importSubmissions = async (file: File) => {
    const formData = new FormData()
    formData.append("file", file)
    const response = await instance.post("/problems/submissions/import", formData, {
        headers: { "Content-Type": "multipart/form-data" },
        timeout: 30_000
    })
    return response.data as { imported: number; skipped: number; errors: { row: number; reason: string }[] }
}

import { useEffect, useState } from "react"
import { getProblemSubmissions, getAllProblemTopics, getProblemByID } from "../api"

export function useProblemSubmissions(problemId?: number) {
  const [data, setData] = useState<any>(null)
  const [error, setError] = useState<Error | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (!problemId) return

    const fetchData = async () => {
      setIsLoading(true)
      setError(null)

      try {
        const result = await getProblemSubmissions(problemId)
        setData(result)
      } catch (err) {
        setError(err instanceof Error ? err : new Error("Unknown error"))
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [problemId])

  return { data, error, isLoading }
}

export function useProblem(problemId?: string) {
  const [data, setData] = useState<any>(null)
  const [error, setError] = useState<Error | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (!problemId) return

    const fetchData = async () => {
      setIsLoading(true)
      setError(null)

      try {
        const result = await getProblemByID(problemId)
        setData(result)
      } catch (err) {
        setError(err instanceof Error ? err : new Error("Unknown error"))
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [problemId])

  return { data, error, isLoading }
}

export function useProblemTopics() {
  const [data, setData] = useState(null)
  const [error, setError] = useState<Error | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true)
      setError(null)

      try {
        const result = await getAllProblemTopics()
        setData(result)
      } catch (err) {
        setError(err instanceof Error ? err : new Error("Unknown error"))
      } finally {
        setIsLoading(false)
      }
    }

    fetchData()
  }, [])

  return { data, error, isLoading }
}

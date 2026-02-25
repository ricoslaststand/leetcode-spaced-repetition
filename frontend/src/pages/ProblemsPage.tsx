import { useEffect, useState } from "react";
import { Link } from "@tanstack/react-router";
import { ChevronDown, ChevronUp, ChevronsUpDown } from "lucide-react";

import { getAllProblems } from "../api";
import { Badge } from "../components/ui/badge";
import ProblemDifficultyTag from "../components/ProblemDifficultyTag";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { cn } from "../lib/utils";
import { Button } from "../components/ui/button";
import { useProblemTopics } from "../hooks/api";
import { generateLinkForLeetcode } from "../lib/leetcodeUtils";

const DIFFICULTY_ORDER: Record<string, number> = { easy: 0, medium: 1, hard: 2 }


const ProblemsPage = () => {
    const [selectedTopics, setSelectedTopics] = useState<Set<string>>(new Set())
    const [selectedDifficulties, setSelectedDifficulties] = useState<Set<string>>(new Set())
    const [problems, setProblems] = useState<any[]>([])
    const [isMetaSelected, setIsCtrlSelected] = useState<boolean>(false)
    const [difficultySort, setDifficultySort] = useState<'asc' | 'desc' | null>(null)

    const sortedProblems = difficultySort === null
        ? problems
        : [...problems].sort((a, b) => {
            const diff = (DIFFICULTY_ORDER[a.difficulty] ?? 0) - (DIFFICULTY_ORDER[b.difficulty] ?? 0)
            return difficultySort === 'asc' ? diff : -diff
        })

    const handleDifficultySort = () => {
        setDifficultySort(prev => prev === null ? 'asc' : prev === 'asc' ? 'desc' : null)
    }

    const { data } = useProblemTopics()

    useEffect(() => {
        (async () => {
            const data = await getAllProblems(Array.from(selectedTopics), Array.from(selectedDifficulties))
            setProblems(data?.data || []);
        })()

        const handleDownPress = (event: KeyboardEvent) => {
            if (event.metaKey) setIsCtrlSelected(true)
        }

        const handleUpPress = (event: KeyboardEvent) => {
            if (event.metaKey) setIsCtrlSelected(true)
            else setIsCtrlSelected(false)
        }

        window.addEventListener("keyup", handleUpPress)
        window.addEventListener("keydown", handleDownPress)

        return () => {
            window.removeEventListener("keyup", handleUpPress)
            window.removeEventListener("keydown", handleDownPress)
        }
    }, [selectedTopics, selectedDifficulties])

    const handleTopicClick = (topic: string) => {
        if (isMetaSelected) {
            if (selectedTopics.has(topic)) {
                setSelectedTopics(new Set([...selectedTopics].filter(t => t !== topic)))
            } else {
                setSelectedTopics(new Set([...selectedTopics, topic]))
            }
        } else {
            setSelectedTopics(new Set([topic]))
        }
    }

    return (
        <div>
            <h1 className="text-2xl font-semibold tracking-tight">Problems</h1>

            <p className="text-sm font-medium text-muted-foreground mt-4">Topics</p>
            <div className="flex flex-wrap gap-2 mt-2">
                {((data as any)?.topics || []).map((topic: string) => (
                    <Badge
                        key={topic}
                        variant={selectedTopics.has(topic) ? "secondary" : "outline"}
                        onClick={() => handleTopicClick(topic)}
                        className={cn(
                            "cursor-pointer select-none",
                            selectedTopics.has(topic) && "bg-primary text-primary-foreground hover:bg-primary/90"
                        )}
                    >
                        {topic}
                    </Badge>
                ))}
            </div>

            <p className="text-sm font-medium text-muted-foreground mt-4">Difficulty</p>
            <div className="flex items-center gap-2 mt-2">
                {["easy", "medium", "hard"].map(d => (
                    <Badge
                        key={d}
                        variant={selectedDifficulties.has(d) ? "secondary" : "outline"}
                        onClick={() => {
                            setSelectedDifficulties(prev => {
                                const next = new Set(prev)
                                next.has(d) ? next.delete(d) : next.add(d)
                                return next
                            })
                        }}
                        className={cn(
                            "cursor-pointer select-none capitalize",
                            selectedDifficulties.has(d) && "bg-primary text-primary-foreground hover:bg-primary/90"
                        )}
                    >
                        {d}
                    </Badge>
                ))}
            </div>

            <div className="flex items-center gap-2 mt-4 mb-6">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setSelectedTopics(new Set([]))}
                    disabled={selectedTopics.size === 0}
                >
                    Clear Topics
                </Button>
            </div>

            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-24">#</TableHead>
                        <TableHead>Title</TableHead>
                        <TableHead onClick={handleDifficultySort} className="cursor-pointer select-none">
                            <span className="flex items-center gap-1">
                                Difficulty
                                {difficultySort === null && <ChevronsUpDown className="h-3 w-3 text-muted-foreground" />}
                                {difficultySort === 'asc' && <ChevronUp className="h-3 w-3" />}
                                {difficultySort === 'desc' && <ChevronDown className="h-3 w-3" />}
                            </span>
                        </TableHead>
                        <TableHead>Submissions</TableHead>
                        <TableHead>LeetCode</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {problems.length === 0 && selectedTopics.size === 0 ? (
                        <TableRow>
                            <TableCell colSpan={5} className="text-center text-muted-foreground py-8">
                                Select a topic above to browse problems.
                            </TableCell>
                        </TableRow>
                    ) : (
                        sortedProblems.map(problem => (
                            <TableRow key={problem.id}>
                                <TableCell className="text-muted-foreground">{problem.id}</TableCell>
                                <TableCell className="font-medium">{problem.title}</TableCell>
                                <TableCell>
                                    <ProblemDifficultyTag difficulty={problem.difficulty} />
                                </TableCell>
                                <TableCell>
                                    <Link
                                        to={`/problems/${problem.id}/submissions`}
                                        className="text-sm text-primary underline-offset-4 hover:underline"
                                    >
                                        View
                                    </Link>
                                </TableCell>
                                <TableCell>
                                    <a
                                        href={generateLinkForLeetcode(problem.slug)}
                                        target="_blank"
                                        rel="noreferrer"
                                        className="text-sm text-primary underline-offset-4 hover:underline"
                                    >
                                        Open
                                    </a>
                                </TableCell>
                            </TableRow>
                        ))
                    )}
                </TableBody>
            </Table>
        </div>
    );
}

export default ProblemsPage;

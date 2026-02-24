import { useEffect, useState } from "react";
import { Link } from "@tanstack/react-router";

import { getAllProblems } from "../api";
import { Badge } from "../components/ui/badge";
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

const difficultyBadgeClass = (d: string) => {
  const normalized = d.toLowerCase()
  if (normalized === "easy") return "bg-green-100 text-green-700 border-green-200"
  if (normalized === "medium") return "bg-yellow-100 text-yellow-700 border-yellow-200"
  return "bg-red-100 text-red-700 border-red-200"
}

const ProblemsPage = () => {
    const [selectedTopics, setSelectedTopics] = useState<Set<string>>(new Set())
    const [problems, setProblems] = useState<any[]>([])
    const [isMetaSelected, setIsCtrlSelected] = useState<boolean>(false)

    const { data } = useProblemTopics()

    useEffect(() => {
        (async () => {
            const data = await getAllProblems(Array.from(selectedTopics))
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
    }, [selectedTopics])

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

            <div className="flex flex-wrap gap-2 mt-4">
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
                        <TableHead>Difficulty</TableHead>
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
                        problems.map(problem => (
                            <TableRow key={problem.id}>
                                <TableCell className="text-muted-foreground">{problem.id}</TableCell>
                                <TableCell className="font-medium">{problem.title}</TableCell>
                                <TableCell>
                                    <Badge className={difficultyBadgeClass(problem.difficulty)}>
                                        {problem.difficulty.charAt(0).toUpperCase() + problem.difficulty.slice(1).toLowerCase()}
                                    </Badge>
                                </TableCell>
                                <TableCell>
                                    <Link
                                        to={`/problems/${problem.id}`}
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

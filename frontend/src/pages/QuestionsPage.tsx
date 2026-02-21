import { useEffect, useState } from "react";
import { Link } from "@tanstack/react-router";

import { getAllQuestions } from "../api";
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
import { useQuestionTags } from "../hooks/api";
import { generateLinkForLeetcode } from "../lib/leetcodeUtils";

const difficultyColor = (d: string) =>
  d === "Easy" ? "text-emerald-600 font-medium" :
  d === "Medium" ? "text-amber-500 font-medium" :
  "text-red-500 font-medium"

const QuestionsPage = () => {
    const [selectedTags, setSelectedTags] = useState<Set<string>>(new Set())
    const [questions, setQuestions] = useState<any[]>([])
    const [isMetaSelected, setIsCtrlSelected] = useState<boolean>(false)

    const { data } = useQuestionTags()

    useEffect(() => {
        (async () => {
            const data = await getAllQuestions(Array.from(selectedTags))
            setQuestions(data?.data || []);
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
    }, [selectedTags])

    const handleTopicClick = (tag: string) => {
        if (isMetaSelected) {
            if (selectedTags.has(tag)) {
                setSelectedTags(new Set([...selectedTags].filter(t => t !== tag)))
            } else {
                setSelectedTags(new Set([...selectedTags, tag]))
            }
        } else {
            setSelectedTags(new Set([tag]))
        }
    }

    return (
        <div>
            <h1 className="text-2xl font-semibold tracking-tight">Questions</h1>

            <div className="flex flex-wrap gap-2 mt-4">
                {((data as any)?.tags || []).map((tag: string) => (
                    <Badge
                        key={tag}
                        variant={selectedTags.has(tag) ? "secondary" : "outline"}
                        onClick={() => handleTopicClick(tag)}
                        className={cn(
                            "cursor-pointer select-none",
                            selectedTags.has(tag) && "bg-primary text-primary-foreground hover:bg-primary/90"
                        )}
                    >
                        {tag}
                    </Badge>
                ))}
            </div>

            <div className="flex items-center gap-2 mt-4 mb-6">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setSelectedTags(new Set([]))}
                    disabled={selectedTags.size === 0}
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
                    {questions.length === 0 && selectedTags.size === 0 ? (
                        <TableRow>
                            <TableCell colSpan={5} className="text-center text-muted-foreground py-8">
                                Select a topic above to browse questions.
                            </TableCell>
                        </TableRow>
                    ) : (
                        questions.map(question => (
                            <TableRow key={question.id}>
                                <TableCell className="text-muted-foreground">{question.id}</TableCell>
                                <TableCell className="font-medium">{question.title}</TableCell>
                                <TableCell className={difficultyColor(question.difficulty)}>
                                    {question.difficulty}
                                </TableCell>
                                <TableCell>
                                    <Link
                                        to={`/questions/${question.id}`}
                                        className="text-sm text-primary underline-offset-4 hover:underline"
                                    >
                                        View
                                    </Link>
                                </TableCell>
                                <TableCell>
                                    <a
                                        href={generateLinkForLeetcode(question.slug)}
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

export default QuestionsPage;

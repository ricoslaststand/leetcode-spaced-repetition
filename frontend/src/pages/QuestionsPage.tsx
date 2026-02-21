import { useEffect, useState } from "react";
import { useNavigate, Link } from "@tanstack/react-router";

import { getAllQuestions } from "../api";
import { Badge } from "../components/ui/badge";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { cn } from "../lib/utils";
import { Button } from "../components/ui/button";
import { useQuestionTags } from "../hooks/api";
import { generateLinkForLeetcode } from "../lib/leetcodeUtils";

const baseBadgeClass = "cursor-pointer"

const QuestionsPage = () => {
    const [selectedTags, setSelectedTags] = useState<Set<string>>(new Set())
    const [questions, setQuestions] = useState<any[]>([])
    const [isMetaSelected, setIsCtrlSelected] = useState<boolean>(false)

    const navigate = useNavigate({ from: '/questions' })

    const { data, isLoading, error } = useQuestionTags()

    useEffect(() => {
        (
            async () => {
                const data = await getAllQuestions(Array.from(selectedTags))
                console.log("data =", data)
                setQuestions(data?.data || []);
            }
        )()

        const handleDownPress = (event: KeyboardEvent) => {
            if (event.metaKey) {
                setIsCtrlSelected(true)
            }
        }

        const handleUpPress = (event: KeyboardEvent) => {
            if (event.metaKey) {
                setIsCtrlSelected(true)
            } else {
                setIsCtrlSelected(false)
            }
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
                setSelectedTags(
                    new Set([...selectedTags].filter(t => t !== tag))
                )
            } else {
                console.log("adding tag:", tag)
                setSelectedTags(
                    new Set([...selectedTags, tag])
                )        
            }
        } else {
            setSelectedTags(new Set([tag]))
        }
    }

    return (
        <div className="absolute inset-0 w-9/10 mx-auto">
            <div>Questions Page</div>
            <div>
                {
                    (data?.tags || []).map(tag => (
                        <Badge
                            key={tag}
                            variant={(selectedTags.has(tag)) ? "secondary" : "outline"}
                            onClick={() => handleTopicClick(tag)}
                            className={selectedTags.has(tag) ? cn("bg-blue-500", "text-white", baseBadgeClass) : baseBadgeClass}>
                                {tag}
                        </Badge>
                    ))
                }
            </div>
            <div>
                <Button variant="outline" onClick={() => navigate({ to: "/" })}>Create Submission</Button>
                <Button className="my-4" variant="outline" onClick={() => setSelectedTags(new Set([]))}>Clear Topics</Button>
            </div>
            <div>
                <Table>
                    {
                        (selectedTags.size  === 0) ?
                        <TableCaption>Please select a topic.</TableCaption> :
                        null
                    }
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-[100px]">Question #</TableHead>
                            <TableHead className="w-[100px] text-left">Title</TableHead>
                            <TableHead className="text-left">Difficulty</TableHead>
                            <TableHead className="text-left">Open Submissions</TableHead>
                            <TableHead className="text-left">Open Link in LeetCode</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {
                            questions.map(question => (
                                <TableRow>
                                    <TableCell className="text-left">{question.id}</TableCell>
                                    <TableCell className="font-medium text-left">{question.title}</TableCell>
                                    <TableCell className="text-left">{question.difficulty}</TableCell>
                                    <TableCell className="text-left">
                                        <Link to={`/questions/${question.id}`}>Open Submissions</Link>
                                    </TableCell>
                                    <TableCell className="text-left">
                                        <a href={generateLinkForLeetcode(question.slug)} target="_blank">Open</a>
                                    </TableCell>
                                </TableRow>
                            ))
                        }
                    </TableBody>
                </Table>
            </div>
        </div>
    );
}

export default QuestionsPage;

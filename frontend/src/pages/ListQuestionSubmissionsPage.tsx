import { useEffect, useState } from "react";
import { Link } from "@tanstack/react-router";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { getQuestionSubmissionsV2 } from "../api";
import { generateLeetcodeURL } from "../lib/leetcodeUtils";

const difficultyColor = (d: string) =>
  d === "Easy" ? "text-emerald-600 font-medium" :
  d === "Medium" ? "text-amber-500 font-medium" :
  "text-red-500 font-medium"

const ListQuestionSubmissionsPage = () => {
    const [questions, setQuestions] = useState<any[]>([])

    useEffect(() => {
        (async () => {
            const data = await getQuestionSubmissionsV2([])
            setQuestions(data?.data || []);
        })()
    }, [])

    return (
        <div>
            <h1 className="text-2xl font-semibold tracking-tight mb-6">All Submissions</h1>
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
                    {questions.length === 0 ? (
                        <TableRow>
                            <TableCell colSpan={5} className="text-center text-muted-foreground py-8">
                                No submissions yet.
                            </TableCell>
                        </TableRow>
                    ) : (
                        questions.map(submission => (
                            <TableRow key={submission.question.id}>
                                <TableCell className="text-muted-foreground">{submission.question.id}</TableCell>
                                <TableCell className="font-medium">{submission.question.title}</TableCell>
                                <TableCell className={difficultyColor(submission.question.difficulty)}>
                                    {submission.question.difficulty}
                                </TableCell>
                                <TableCell>
                                    <Link
                                        to={`/questions/${submission.question.id}`}
                                        className="text-sm text-primary underline-offset-4 hover:underline"
                                    >
                                        View
                                    </Link>
                                </TableCell>
                                <TableCell>
                                    <a
                                        href={generateLeetcodeURL(submission.question.id)}
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

export default ListQuestionSubmissionsPage;

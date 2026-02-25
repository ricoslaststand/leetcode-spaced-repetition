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
import { getProblemSubmissionsV2 } from "../api";
import { Button } from "../components/ui/button";
import ProblemDifficultyTag from "../components/ProblemDifficultyTag";
import ConfidenceLevelTag from "../components/ConfidenceLevelTag";
import { generateLinkForLeetcode } from "../lib/leetcodeUtils";

const formatDate = (dateStr: string): string => {
    if (!dateStr) return "—"
    const datePart = dateStr.split("T")[0]
    const [year, month, day] = datePart.split("-").map(Number)
    if (isNaN(year) || isNaN(month) || isNaN(day)) return datePart
    const d = new Date(Date.UTC(year, month - 1, day))
    return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric", timeZone: "UTC" })
}

const formatTimeTaken = (seconds: number): string => {
    if (!seconds) return "—"
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    if (m === 0) return `${s}s`
    if (s === 0) return `${m}m`
    return `${m}m ${s}s`
}


const ListProblemSubmissionsPage = () => {
    const [problems, setProblems] = useState<any[]>([])

    useEffect(() => {
        (async () => {
            const data = await getProblemSubmissionsV2([])
            setProblems(data?.data || []);
        })()
    }, [])

    return (
        <div>
            <div className="flex items-center justify-between mb-6">
                <h1 className="text-2xl font-semibold tracking-tight">All Submissions</h1>
                <Button asChild>
                    <Link to="/problems/submissions/new">+ New Submission</Link>
                </Button>
            </div>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-24">#</TableHead>
                        <TableHead>Title</TableHead>
                        <TableHead>Difficulty</TableHead>
                        <TableHead>Confidence</TableHead>
                        <TableHead>Date</TableHead>
                        <TableHead>Time Taken</TableHead>
                        <TableHead>Submissions</TableHead>
                        <TableHead>LeetCode</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {problems.length === 0 ? (
                        <TableRow>
                            <TableCell colSpan={8} className="text-center text-muted-foreground py-8">
                                No submissions yet.
                            </TableCell>
                        </TableRow>
                    ) : (
                        problems.map(submission => (
                            <TableRow key={submission.problem.id}>
                                <TableCell className="text-muted-foreground">{submission.problem.id}</TableCell>
                                <TableCell className="font-medium">{submission.problem.title}</TableCell>
                                <TableCell>
                                    <ProblemDifficultyTag difficulty={submission.problem.difficulty} />
                                </TableCell>
                                <TableCell>
                                    <ConfidenceLevelTag confidenceLevel={submission.confidenceLevel} />
                                </TableCell>
                                <TableCell className="text-sm text-muted-foreground">
                                    {formatDate(submission.submittedAt)}
                                </TableCell>
                                <TableCell className="text-sm text-muted-foreground">
                                    {formatTimeTaken(submission.timeTaken)}
                                </TableCell>
                                <TableCell>
                                    <Link
                                        to={`/problems/${submission.problem.id}/submissions`}
                                        className="text-sm text-primary underline-offset-4 hover:underline"
                                    >
                                        View
                                    </Link>
                                </TableCell>
                                <TableCell>
                                    <a
                                        href={generateLinkForLeetcode(submission.problem.slug)}
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

export default ListProblemSubmissionsPage;

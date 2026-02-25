import { Link } from "@tanstack/react-router"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import ProblemDifficultyTag from "../components/ProblemDifficultyTag"
import { useDashboard } from "../hooks/api"
import { generateLinkForLeetcode } from "../lib/leetcodeUtils"
import type { ProblemReviewItem } from "../models/Problem"

const formatDueDate = (isoStr: string): string => {
    if (!isoStr) return "—"
    const datePart = isoStr.split("T")[0]
    const [year, month, day] = datePart.split("-").map(Number)
    if (isNaN(year) || isNaN(month) || isNaN(day)) return datePart
    const d = new Date(Date.UTC(year, month - 1, day))
    return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric", timeZone: "UTC" })
}

const ProblemTableRow = ({ item, extraCell }: { item: ProblemReviewItem; extraCell: React.ReactNode }) => (
    <TableRow>
        <TableCell className="text-muted-foreground">{item.problemId}</TableCell>
        <TableCell className="font-medium">{item.title}</TableCell>
        <TableCell><ProblemDifficultyTag difficulty={item.difficulty} /></TableCell>
        <TableCell className="text-sm text-muted-foreground">{extraCell}</TableCell>
        <TableCell>
            <Link
                to={`/problems/${item.problemId}/submissions`}
                className="text-sm text-primary underline-offset-4 hover:underline"
            >
                View
            </Link>
        </TableCell>
        <TableCell>
            <a
                href={generateLinkForLeetcode(item.slug)}
                target="_blank"
                rel="noreferrer"
                className="text-sm text-primary underline-offset-4 hover:underline"
            >
                Open
            </a>
        </TableCell>
    </TableRow>
)

const DASHBOARD_TABLE_LIMIT = 10

const DashboardPage = () => {
    const { data, isLoading } = useDashboard(DASHBOARD_TABLE_LIMIT)

    return (
        <div className="space-y-10">
            <h1 className="text-2xl font-semibold tracking-tight">Dashboard</h1>

            <section>
                <div className="flex items-center gap-2 mb-3">
                    <h2 className="text-lg font-medium">Due for Review</h2>
                    {!isLoading && (data?.overdueCount ?? 0) > 0 && (
                        <span className="text-sm text-muted-foreground">
                            ({data!.overdueCount} total)
                        </span>
                    )}
                </div>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-24">#</TableHead>
                            <TableHead>Title</TableHead>
                            <TableHead>Difficulty</TableHead>
                            <TableHead>Due</TableHead>
                            <TableHead>Submissions</TableHead>
                            <TableHead>LeetCode</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            <TableRow>
                                <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                                    Loading...
                                </TableCell>
                            </TableRow>
                        ) : (data?.due ?? []).length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                                    No problems due for review.
                                </TableCell>
                            </TableRow>
                        ) : (
                            (data?.due ?? []).map(item => (
                                <ProblemTableRow
                                    key={item.problemId}
                                    item={item}
                                    extraCell={formatDueDate(item.due)}
                                />
                            ))
                        )}
                    </TableBody>
                </Table>
            </section>

            <section>
                <h2 className="text-lg font-medium mb-3">Most at Risk of Forgetting</h2>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-24">#</TableHead>
                            <TableHead>Title</TableHead>
                            <TableHead>Difficulty</TableHead>
                            <TableHead>Stability</TableHead>
                            <TableHead>Submissions</TableHead>
                            <TableHead>LeetCode</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            <TableRow>
                                <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                                    Loading...
                                </TableCell>
                            </TableRow>
                        ) : (data?.lowStability ?? []).length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                                    No reviewed problems yet.
                                </TableCell>
                            </TableRow>
                        ) : (
                            (data?.lowStability ?? []).map(item => (
                                <ProblemTableRow
                                    key={item.problemId}
                                    item={item}
                                    extraCell={item.stability.toFixed(2)}
                                />
                            ))
                        )}
                    </TableBody>
                </Table>
            </section>
        </div>
    )
}

export default DashboardPage

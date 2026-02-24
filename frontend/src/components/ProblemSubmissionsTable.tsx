import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"

type ProblemSubmissionsTableProps = {
    submissions: any[]
};

const confidenceBadgeVariant = (level: string): "default" | "secondary" | "destructive" | "outline" => {
    switch (level?.toLowerCase()) {
        case "easy": return "default"
        case "good": return "secondary"
        case "hard": return "outline"
        case "again": return "destructive"
        default: return "outline"
    }
}

const formatTimeTaken = (seconds: number): string => {
    if (!seconds) return "—"
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    if (m === 0) return `${s}s`
    if (s === 0) return `${m}m`
    return `${m}m ${s}s`
}

const formatDate = (dateStr: string): string => {
    if (!dateStr) return "—"
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr.split("T")[0]
    return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" })
}

export default function ProblemSubmissionsTable({ submissions }: ProblemSubmissionsTableProps) {
    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>Date</TableHead>
                    <TableHead>Confidence</TableHead>
                    <TableHead className="text-right">Time Taken</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {submissions.length === 0 ? (
                    <TableRow>
                        <TableCell colSpan={3} className="text-center text-muted-foreground py-8">
                            No submissions recorded.
                        </TableCell>
                    </TableRow>
                ) : (
                    submissions.map((submission, i) => (
                        <TableRow key={i}>
                            <TableCell className="text-muted-foreground text-sm">
                                {formatDate(submission.date)}
                            </TableCell>
                            <TableCell>
                                <Badge variant={confidenceBadgeVariant(submission.confidenceLevel)}>
                                    {submission.confidenceLevel}
                                </Badge>
                            </TableCell>
                            <TableCell className="text-right text-sm text-muted-foreground">
                                {formatTimeTaken(submission.timeTaken)}
                            </TableCell>
                        </TableRow>
                    ))
                )}
            </TableBody>
        </Table>
    )
}

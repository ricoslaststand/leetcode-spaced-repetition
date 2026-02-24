import { Badge } from "@/components/ui/badge"

interface Props {
    difficulty: string
}

const difficultyBadgeClass = (d: string) => {
    const normalized = d.toLowerCase()
    if (normalized === "easy") return "bg-green-100 text-green-700 border-green-200"
    if (normalized === "medium") return "bg-yellow-100 text-yellow-700 border-yellow-200"
    return "bg-red-100 text-red-700 border-red-200"
}

export default function ProblemDifficultyTag({ difficulty }: Props) {
    return (
        <Badge className={difficultyBadgeClass(difficulty)}>
            {difficulty.charAt(0).toUpperCase() + difficulty.slice(1).toLowerCase()}
        </Badge>
    )
}

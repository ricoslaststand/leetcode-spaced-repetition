import { Badge } from "@/components/ui/badge"

interface Props {
    confidenceLevel: string
}

const confidenceBadgeClass = (level: string) => {
    const normalized = level.toLowerCase()
    if (normalized === "easy") return "bg-green-100 text-green-700 border-green-200"
    if (normalized === "good") return "bg-yellow-100 text-yellow-700 border-yellow-200"
    if (normalized === "hard") return "bg-orange-100 text-orange-700 border-orange-200"
    return "bg-red-100 text-red-700 border-red-200"
}

export default function ConfidenceLevelTag({ confidenceLevel }: Props) {
    return (
        <Badge className={confidenceBadgeClass(confidenceLevel)}>
            {confidenceLevel.charAt(0).toUpperCase() + confidenceLevel.slice(1).toLowerCase()}
        </Badge>
    )
}

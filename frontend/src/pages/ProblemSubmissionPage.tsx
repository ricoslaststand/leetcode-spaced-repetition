import React, { useEffect, useState } from 'react';
import { Link } from '@tanstack/react-router';

import { z } from 'zod';
import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from "@hookform/resolvers/zod"
import { toast } from "sonner"

import type { ConfidenceLevel as ConfidenceLevelType } from '../models/Problem';
import { ConfidenceLevel, confidenceLevelToString } from '../models/Problem';
import { Input } from '../components/ui/input';
import { Slider } from '../components/ui/slider';
import { Field, FieldGroup, FieldLabel } from '../components/ui/field';
import { Button } from '../components/ui/button';
import DurationInput from '../components/DurationInput';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { createProblemSubmission, getProblemSubmissionsV2, importSubmissions } from '../api';
import { generateLeetcodeURL } from '../lib/leetcodeUtils';

const ConfidenceLevelMemes: { level: ConfidenceLevelType; text: string }[] = [
    {
        level: ConfidenceLevel.Again,
        text: "I have no clue what's going on"
    },
    {
        level: ConfidenceLevel.Hard,
        text: "I see how they did it, but I did not see that coming"
    },
    {
        level: ConfidenceLevel.Good,
        text: "Things are starting to click..."
    },
    {
        level: ConfidenceLevel.Easy,
        text: "You did it, buddy"
    }
]

const confidenceLevelLabels = ["Again", "Hard", "Good", "Easy"]

const formatDate = (dateStr: string): string => {
    if (!dateStr) return "—"
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr.split("T")[0]
    return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" })
}

const formatTimeTaken = (seconds: number): string => {
    if (!seconds) return "—"
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    if (m === 0) return `${s}s`
    if (s === 0) return `${m}m`
    return `${m}m ${s}s`
}

const difficultyColor = (d: string) =>
    d === "Easy" ? "text-emerald-600 font-medium" :
    d === "Medium" ? "text-amber-500 font-medium" :
    "text-red-500 font-medium"

const formSchema = z.object({
    problemId: z.string(),
    confidenceLevel: z.number().min(ConfidenceLevel.Again).max(ConfidenceLevel.Easy),
    timeTaken: z.number().min(0).optional()
})

const ProblemSubmissionPage: React.FC = () => {
    const [importFile, setImportFile] = useState<File | null>(null)
    const [isImporting, setIsImporting] = useState(false)
    const [problems, setProblems] = useState<any[]>([])

    const fetchSubmissions = async () => {
        const data = await getProblemSubmissionsV2([])
        setProblems(data?.data || [])
    }

    useEffect(() => { fetchSubmissions() }, [])

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            problemId: "",
            confidenceLevel: ConfidenceLevel.Good,
        }
    })

    const onImport = async () => {
        if (!importFile) return
        setIsImporting(true)
        try {
            const result = await importSubmissions(importFile)
            if (result.errors.length === 0) {
                toast.success(`Imported ${result.imported} submission${result.imported !== 1 ? 's' : ''}`)
            } else {
                toast.warning(
                    `Imported ${result.imported}, skipped ${result.errors.length} row${result.errors.length !== 1 ? 's' : ''} with errors`
                )
            }
            setImportFile(null)
            await fetchSubmissions()
        } catch {
            toast.error("Failed to import file. Please check the file format and try again.")
        } finally {
            setIsImporting(false)
        }
    }

    const onSubmit = async (data: z.infer<typeof formSchema>) => {
        const timeTaken = data.timeTaken && data.timeTaken > 0 ? data.timeTaken : undefined
        await createProblemSubmission(
            parseInt(data.problemId),
            confidenceLevelToString(data.confidenceLevel as ConfidenceLevelType),
            timeTaken
        )
        toast.success("Submission logged!")
        form.reset()
        await fetchSubmissions()
    }

    return (
        <div>
            <div>
                <h1 className="text-2xl font-semibold tracking-tight mb-6">
                    Log a submission
                </h1>
                <form onSubmit={form.handleSubmit(onSubmit)}>
                    <Controller
                        name="problemId"
                        control={form.control}
                        render={({ field, fieldState }) => (
                            <Field data-invalid={fieldState.invalid}>
                                <FieldLabel htmlFor={field.name}>Problem #</FieldLabel>
                                <Input
                                    {...field}
                                    placeholder="e.g. 42"
                                    aria-invalid={fieldState.invalid}
                                />
                            </Field>
                        )}
                    />

                    <FieldGroup className="grid grid-cols-2 gap-6 my-6">
                        <Controller
                            name="confidenceLevel"
                            control={form.control}
                            render={({ field }) => (
                                <Field>
                                    <FieldLabel>Confidence Level</FieldLabel>
                                    <div className="pt-2">
                                        <Slider
                                            value={[field.value]}
                                            step={1}
                                            min={ConfidenceLevel.Again}
                                            max={ConfidenceLevel.Easy}
                                            onValueChange={val => field.onChange(val[0])}
                                        />
                                    </div>
                                    <div className="mt-2">
                                        <span className="text-sm font-medium">
                                            {confidenceLevelLabels[field.value - 1]}
                                        </span>
                                        <p className="text-xs text-muted-foreground mt-0.5">
                                            {ConfidenceLevelMemes[field.value - 1].text}
                                        </p>
                                    </div>
                                </Field>
                            )}
                        />

                        <Controller
                            name="timeTaken"
                            control={form.control}
                            render={({ field }) => (
                                <DurationInput
                                    label="Time Taken"
                                    valueSeconds={field.value ?? 0}
                                    onChangeSeconds={field.onChange}
                                />
                            )}
                        />
                    </FieldGroup>

                    <Button disabled={!form.formState.isValid} type="submit">
                        Log Submission
                    </Button>
                </form>

                <div className="flex items-center gap-3 my-8">
                    <div className="h-px flex-1 bg-border" />
                    <span className="text-xs text-muted-foreground">or</span>
                    <div className="h-px flex-1 bg-border" />
                </div>

                <div>
                    <h2 className="text-sm font-medium mb-3">Import from Excel</h2>
                    <p className="text-xs text-muted-foreground mb-4">
                        Upload an .xlsx file with columns: Problem #, Date (YYYY-MM-DD), Time Taken (optional), Confidence (1–4).
                    </p>
                    <div className="flex items-center gap-3">
                        <Input
                            type="file"
                            accept=".xlsx"
                            onChange={e => setImportFile(e.target.files?.[0] ?? null)}
                        />
                        <Button
                            type="button"
                            disabled={!importFile || isImporting}
                            onClick={onImport}
                        >
                            {isImporting ? "Importing..." : "Import"}
                        </Button>
                    </div>
                </div>
            </div>

            <div className="mt-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-6">All Submissions</h2>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-24">#</TableHead>
                            <TableHead>Title</TableHead>
                            <TableHead>Difficulty</TableHead>
                            <TableHead>Date</TableHead>
                            <TableHead>Duration</TableHead>
                            <TableHead>Submissions</TableHead>
                            <TableHead>LeetCode</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {problems.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={7} className="text-center text-muted-foreground py-8">
                                    No submissions yet.
                                </TableCell>
                            </TableRow>
                        ) : (
                            problems.map(submission => (
                                <TableRow key={submission.problem.id}>
                                    <TableCell className="text-muted-foreground">{submission.problem.id}</TableCell>
                                    <TableCell className="font-medium">{submission.problem.title}</TableCell>
                                    <TableCell className={difficultyColor(submission.problem.difficulty)}>
                                        {submission.problem.difficulty}
                                    </TableCell>
                                    <TableCell className="text-sm text-muted-foreground">
                                        {formatDate(submission.submittedAt)}
                                    </TableCell>
                                    <TableCell className="text-sm text-muted-foreground">
                                        {formatTimeTaken(submission.timeTaken)}
                                    </TableCell>
                                    <TableCell>
                                        <Link
                                            to={`/problems/${submission.problem.id}`}
                                            className="text-sm text-primary underline-offset-4 hover:underline"
                                        >
                                            View
                                        </Link>
                                    </TableCell>
                                    <TableCell>
                                        <a
                                            href={generateLeetcodeURL(submission.problem.id)}
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
        </div>
    )
}

export default ProblemSubmissionPage;

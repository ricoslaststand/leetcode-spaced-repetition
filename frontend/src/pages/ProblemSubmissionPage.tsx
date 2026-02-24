import React, { useRef, useState } from 'react';

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
import { createProblemSubmission, importSubmissions } from '../api';

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


const formSchema = z.object({
    problemId: z.string(),
    confidenceLevel: z.number().min(ConfidenceLevel.Again).max(ConfidenceLevel.Easy),
    timeTaken: z.number().min(0).optional()
})

const ProblemSubmissionPage: React.FC = () => {
    const [importFile, setImportFile] = useState<File | null>(null)
    const [isImporting, setIsImporting] = useState(false)
    const fileInputRef = useRef<HTMLInputElement>(null)

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
            if (fileInputRef.current) fileInputRef.current.value = ''
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
                            ref={fileInputRef}
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

        </div>
    )
}

export default ProblemSubmissionPage;

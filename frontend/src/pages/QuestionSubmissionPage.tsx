import React from 'react';

import { z } from 'zod';
import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from "@hookform/resolvers/zod"
import parse from 'parse-duration'
import { toast } from "sonner"

import type { ConfidenceLevel as ConfidenceLevelType } from '../models/Question';
import { ConfidenceLevel, confidenceLevelToString } from '../models/Question';
import { Input } from '../components/ui/input';
import { Slider } from '../components/ui/slider';
import { Field, FieldGroup, FieldLabel } from '../components/ui/field';
import { Button } from '../components/ui/button';
import { createQuestionSubmission } from '../api';

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
    questionId: z.string(),
    confidenceLevel: z.number().min(ConfidenceLevel.Again).max(ConfidenceLevel.Easy),
    timeTaken: z.number().min(0)
})

const QuestionSubmissionPage: React.FC = () => {
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            questionId: "",
            confidenceLevel: ConfidenceLevel.Good,
        }
    })

    const onSubmit = async (data: z.infer<typeof formSchema>) => {
        await createQuestionSubmission(
            parseInt(data.questionId),
            confidenceLevelToString(data.confidenceLevel as ConfidenceLevelType),
            data.timeTaken
        )
        toast.success("Submission logged!")
        form.reset()
    }

    return (
        <div className="max-w-lg">
            <h1 className="text-2xl font-semibold tracking-tight mb-6">
                Log a submission
            </h1>
            <form onSubmit={form.handleSubmit(onSubmit)}>
                <Controller
                    name="questionId"
                    control={form.control}
                    render={({ field, fieldState }) => (
                        <Field data-invalid={fieldState.invalid}>
                            <FieldLabel htmlFor={field.name}>Question #</FieldLabel>
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
                            <Field>
                                <FieldLabel htmlFor="timeTaken">Time Taken</FieldLabel>
                                <Input
                                    id="timeTaken"
                                    placeholder="e.g. 15m30s"
                                    onBlur={e => {
                                        const value = parse(e.currentTarget.value)
                                        if (value !== null) {
                                            field.onChange(Math.floor(value / 1_000))
                                        } else {
                                            field.onChange(undefined)
                                        }
                                    }}
                                />
                            </Field>
                        )}
                    />
                </FieldGroup>

                <Button disabled={!form.formState.isValid} type="submit">
                    Log Submission
                </Button>
            </form>
        </div>
    )
}

export default QuestionSubmissionPage;

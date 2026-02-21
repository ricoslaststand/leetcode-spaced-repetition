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

const ConfidenceLevelMemes: { level: ConfidenceLevelType; meme: string; text: string }[] = [
    {
        level: ConfidenceLevel.Again,
        meme: "simpsons_repeat_stuff.gif",
        text: "I have no clue what's going on"
    },
    {
        level: ConfidenceLevel.Hard,
        meme: "drake_explaining.gif",
        text: "I see how they did it, but I did not see that coming"
    },
    {
        level: ConfidenceLevel.Good,
        meme: "exploding_brain.gif",
        text: "Things are starting to click..."
    },
    {
        level: ConfidenceLevel.Easy,
        meme: "great_gatsy_nod.gif",
        text: "You did it, buddy"
    }
]

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
        toast.success("Question submission created!")
        form.reset()
    }

    return (
        <div>
            <h1>Which question did you complete today?</h1>
            <div className="my-4">
                <form onSubmit={form.handleSubmit(onSubmit)} >
                    <Controller
                        name="questionId"
                        control={form.control}
                        render={({ field, fieldState }) => (
                            <Field data-invalid={fieldState.invalid}>
                                <FieldLabel htmlFor={field.name}>Question #</FieldLabel>
                                <Input
                                    {...field}
                                    aria-invalid={fieldState.invalid}
                                />
                            </Field>
                        )}
                    />
                    <FieldGroup className="grid grid-cols-2 my-4">
                        <Controller
                            name="confidenceLevel"
                            control={form.control}
                            render={({ field }) => (
                                <Field>
                                    <FieldLabel>ConfidenceLevel</FieldLabel>
                                    <Slider
                                        value={[field.value]}
                                        step={1}
                                        min={ConfidenceLevel.Again}
                                        max={ConfidenceLevel.Easy}
                                        onValueChange={val=> field.onChange(val[0])}
                                    />
                                    <p>{ConfidenceLevelMemes[field.value - 1].text}</p>
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
                    <Button disabled={!form.formState.isValid} type="submit" variant="outline">Create Submission</Button>
                </form>
            </div>
        </div>
    )
}

export default QuestionSubmissionPage;

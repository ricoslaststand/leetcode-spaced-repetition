import React, { useEffect, useRef, useState } from "react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

/**
 * DurationInput
 *
 * Advanced UX Pattern:
 * - Two numeric inputs (minutes + seconds)
 * - Auto-carry logic (e.g., 90 sec -> 1 min 30 sec)
 * - Keyboard arrows adjust values (native for type=number; plus extra shortcuts)
 * - Prevent seconds > 59 (we normalize on change/blur)
 * - Backend-friendly: emits totalSeconds
 * - Renders a compact normalized value string
 *
 * Usage:
 *   <DurationInput valueSeconds={duration} onChangeSeconds={setDuration} />
 */
type DurationInputProps = {
  label?: string;
  valueSeconds: number;
  onChangeSeconds: (seconds: number) => void;
  minSeconds?: number;
  maxSeconds?: number;
  id?: string;
};

function DurationInput({
  label = "Duration",
  valueSeconds,
  onChangeSeconds,
  minSeconds = 0,
  maxSeconds = Number.POSITIVE_INFINITY,
  id = "duration",
}: DurationInputProps) {
  const minutesRef = useRef<HTMLInputElement | null>(null);
  const secondsRef = useRef<HTMLInputElement | null>(null);

  // Keep local fields so typing is smooth (especially when user enters "" or partial).
  const [minutesStr, setMinutesStr] = useState<string>("0");
  const [secondsStr, setSecondsStr] = useState<string>("0");

  // Initialize/refresh local strings from canonical valueSeconds.
  useEffect(() => {
    const clamped = clamp(valueSeconds, minSeconds, maxSeconds);
    const { minutes, seconds } = splitSeconds(clamped);
    setMinutesStr(String(minutes));
    setSecondsStr(String(seconds));
  }, [valueSeconds, minSeconds, maxSeconds]);

  function commitFromFields(nextMinutesStr: string, nextSecondsStr: string) {
    // Allow empty while typing; treat as 0 on commit.
    const m = parseIntSafe(nextMinutesStr);
    const s = parseIntSafe(nextSecondsStr);

    // Normalize seconds overflow/underflow (carry/borrow).
    let total = m * 60 + s;

    // Clamp to bounds.
    total = clamp(total, minSeconds, maxSeconds);

    // Normalize again into fields (ensures seconds 0..59).
    const { minutes, seconds } = splitSeconds(total);
    setMinutesStr(String(minutes));
    setSecondsStr(String(seconds));

    onChangeSeconds(total);
  }

  function handleMinutesChange(v: string) {
    // Only allow digits and optional leading minus (we'll normalize/borrow).
    if (!/^[-]?\d*$/.test(v)) return;
    setMinutesStr(v);
    commitFromFields(v, secondsStr);
  }

  function handleSecondsChange(v: string) {
    if (!/^[-]?\d*$/.test(v)) return;
    setSecondsStr(v);
    commitFromFields(minutesStr, v);
  }

  function handleBlur() {
    // Normalize on blur in case user leaves empty.
    commitFromFields(minutesStr, secondsStr);
  }

  // Keyboard helpers:
  // - ArrowUp/Down on seconds borrows/carries by 1 sec; Shift adjusts by 10 sec.
  // - Ctrl/Cmd+ArrowUp/Down adjusts minutes by 1 (Shift by 5).
  function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>, field: "minutes" | "seconds") {
    const isMod = e.ctrlKey || e.metaKey;
    const shift = e.shiftKey;

    if (e.key === "ArrowUp" || e.key === "ArrowDown") {
      const dir = e.key === "ArrowUp" ? 1 : -1;

      // If modifier, adjust minutes.
      if (isMod) {
        e.preventDefault();
        const stepMinutes = shift ? 5 : 1;
        const delta = dir * stepMinutes * 60;
        onChangeSeconds(clamp(valueSeconds + delta, minSeconds, maxSeconds));
        if (field === "minutes") {
          minutesRef.current?.focus();
          minutesRef.current?.select();
        }
        return;
      }

      // Otherwise adjust seconds (fine-grained).
      if (field === "seconds") {
        e.preventDefault();
        const stepSeconds = shift ? 10 : 1;
        const delta = dir * stepSeconds;
        onChangeSeconds(clamp(valueSeconds + delta, minSeconds, maxSeconds));
        secondsRef.current?.focus();
        secondsRef.current?.select();
        return;
      }
    }

    // Convenience: if user types ':' in minutes field, jump to seconds.
    if (field === "minutes" && e.key === ":") {
      e.preventDefault();
      secondsRef.current?.focus();
      secondsRef.current?.select();
    }
  }

  return (
    <div className="space-y-2">
      <Label htmlFor={id}>
        {label}
      </Label>

      <div className="flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-2 rounded-2xl border px-3 py-2 shadow-sm">
          <Input
            ref={minutesRef}
            id={id}
            inputMode="numeric"
            type="number"
            className="w-20 bg-transparent text-sm"
            value={minutesStr}
            onChange={(e) => handleMinutesChange(e.target.value)}
            onBlur={handleBlur}
            onKeyDown={(e) => handleKeyDown(e, "minutes")}
            aria-label="Minutes"
          />
          <span className="text-sm text-muted-foreground">min</span>

          <div className="h-5 w-px bg-border mx-1" />

          <Input
            ref={secondsRef}
            inputMode="numeric"
            type="number"
            className="w-20 bg-transparent text-sm"
            value={secondsStr}
            onChange={(e) => handleSecondsChange(e.target.value)}
            onBlur={handleBlur}
            onKeyDown={(e) => handleKeyDown(e, "seconds")}
            aria-label="Seconds"
          />
          <span className="text-sm text-muted-foreground">sec</span>
        </div>

      </div>
    </div>
  );
}

function clamp(n: number, min: number, max: number) {
  if (Number.isNaN(n)) return min;
  return Math.min(Math.max(n, min), max);
}

function splitSeconds(totalSeconds: number) {
  const t = Math.max(0, Math.floor(totalSeconds));
  const minutes = Math.floor(t / 60);
  const seconds = t % 60;
  return { minutes, seconds };
}

function parseIntSafe(s: string) {
  if (s == null || s.trim() === "" || s === "-") return 0;
  const n = parseInt(s, 10);
  return Number.isFinite(n) ? n : 0;
}

export default DurationInput;

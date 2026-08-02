# LeetCode Buddy

A single-user web app for practising LeetCode problems on a spaced repetition schedule.

Grinding through LeetCode teaches you patterns that you quietly forget a few weeks later — you
solve a problem, mark it done, and by the time it shows up in an interview the trick is gone.
LeetCode Buddy treats every problem you've attempted as a flashcard instead of a checkbox. Each
time you work a problem you log a submission with a confidence rating — *Again*, *Hard*, *Good*,
or *Easy* — plus optional notes on how long it took and which language you used. From there the
app decides when you should see that problem again.

Scheduling is handled by [FSRS](https://github.com/open-spaced-repetition/go-fsrs), the same
algorithm family used by modern Anki. Every problem carries its own card state: a due date, a
*stability* score representing how durably you know it, an intrinsic difficulty, and counts of
repetitions and lapses. Your confidence rating is fed in as the FSRS grade, so a problem you
flailed at comes back within days while one you solved cleanly gets pushed weeks out. The
dashboard reads directly off that state, with one tab listing everything due for review and
another ranking the problems most at risk of being forgotten by lowest stability.

Day to day, you browse a seeded catalogue of LeetCode problems filtered by topic tag and
difficulty, jump straight out to the problem on leetcode.com, and log a submission when you're
done. Every problem has its own history page showing past attempts, time taken, and current
stability, and there's a combined view of all submissions across problems. If you've been
tracking practice in a spreadsheet already, an Excel import backfills that history in bulk —
each row is replayed through the scheduler so your due dates reflect everything you've done, not
just what you've logged since installing the app.

The backend is a Go API built on Gin and backed by PostgreSQL, with Goose migrations, structured
logging, Prometheus metrics, and generated Swagger docs. The frontend is a React single-page app
using TanStack Router, Tailwind, and shadcn/ui. The app is deliberately single-user and holds no
login code of its own: authentication is delegated to an authenticating reverse proxy (Authelia
via Traefik's ForwardAuth), and the API only accepts requests that the proxy has already
identified as the deployment's owner.

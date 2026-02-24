import './App.css'

import { createRootRoute, createRoute, createRouter, redirect, RouterProvider } from '@tanstack/react-router'
import ProblemSubmissionPage from './pages/ProblemSubmissionPage'
import ProblemsPage from './pages/ProblemsPage'
import ProblemMetadataPage from './pages/ProblemMetadataPage'
import ListProblemSubmissionsPage from './pages/ListProblemSubmissionsPage'
import { Toaster } from 'sonner'
import { AuthenticatedLayout } from './layouts/AuthenticatedLayout'

const rootRoute = createRootRoute()
const layoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'layout',
  component: () => <AuthenticatedLayout />,
})

const homeRedirectRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: '/',
  beforeLoad: () => { throw redirect({ to: '/problems' }) }
})

const problemsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: 'problems',
  component: () => <AuthenticatedLayout />
})
const problemsListRoute = createRoute({
  getParentRoute: () => problemsRoute,
  path: '/',
  component: () => <ProblemsPage />,
})

const allSubmissionsRoute = createRoute({
  getParentRoute: () => problemsRoute,
  path: 'submissions',
  component: () => <ListProblemSubmissionsPage />
})

const newSubmissionRoute = createRoute({
  getParentRoute: () => problemsRoute,
  path: 'submissions/new',
  component: () => <ProblemSubmissionPage />
})

const problemSubmissionsRoute = createRoute({
  getParentRoute: () => problemsRoute,
  path: '$problemId/submissions',
  component: () => <ProblemMetadataPage />
})

const routeTree = rootRoute.addChildren([
  layoutRoute.addChildren([
    homeRedirectRoute,
  ]),
  problemsRoute.addChildren([
    problemsListRoute,
    allSubmissionsRoute,
    newSubmissionRoute,
    problemSubmissionsRoute,
  ]),
])
const router = createRouter({ routeTree })

export default function App() {
  return (
    <>
      <Toaster position="top-right" richColors />
      <RouterProvider router={router} />
    </>
  )
}

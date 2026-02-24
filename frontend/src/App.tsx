import './App.css'

import { createRootRoute, createRoute, createRouter, RouterProvider } from '@tanstack/react-router'
import ProblemSubmissionPage from './pages/ProblemSubmissionPage'
import ProblemsPage from './pages/ProblemsPage'
import ProblemMetadataPage from './pages/ProblemMetadataPage'
import ListProblemSubmissionsPage from './pages/ListProblemSubmissionsPage'
import { Toaster } from 'sonner'
import { AuthenticatedLayout } from './layouts/AuthenticatedLayout'

const rootRoute = createRootRoute()
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: () => <AuthenticatedLayout />,
})

const problemSubmissionsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: 'submissions',
  component: () => <ProblemSubmissionPage />
})

const problemsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: 'problems',
  component: () => <AuthenticatedLayout />
})
const submissionsRoute = createRoute({
  getParentRoute: () => problemsRoute,
  path: 'submissions',
  component: () => <ListProblemSubmissionsPage />
})

const problemsListRoute = createRoute({
  getParentRoute: () => problemsRoute,
  path: '/',
  component: () => <ProblemsPage />,
})

const problemMetadataRoute = createRoute({
  getParentRoute: () => problemsRoute,
  path: '$problemId',
  component: () => <ProblemMetadataPage />
})

const routeTree = rootRoute.addChildren([
  indexRoute.addChildren([
    problemSubmissionsRoute,
  ]),
  problemsRoute.addChildren([
    problemsListRoute,
    problemMetadataRoute
  ]),
  submissionsRoute,
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

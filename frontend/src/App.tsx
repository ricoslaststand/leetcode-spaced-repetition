import './App.css'

import { createRootRoute, createRoute, createRouter, redirect, RouterProvider } from '@tanstack/react-router'
import ProblemSubmissionPage from './pages/ProblemSubmissionPage'
import ProblemsPage from './pages/ProblemsPage'
import ProblemMetadataPage from './pages/ProblemMetadataPage'
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

const problemSubmissionsRoute = createRoute({
  getParentRoute: () => layoutRoute,
  path: 'submissions',
  component: () => <ProblemSubmissionPage />
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

const problemMetadataRoute = createRoute({
  getParentRoute: () => problemsRoute,
  path: '$problemId',
  component: () => <ProblemMetadataPage />
})

const routeTree = rootRoute.addChildren([
  layoutRoute.addChildren([
    homeRedirectRoute,
    problemSubmissionsRoute,
  ]),
  problemsRoute.addChildren([
    problemsListRoute,
    problemMetadataRoute
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

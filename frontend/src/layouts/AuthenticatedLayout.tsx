import { Outlet } from "@tanstack/react-router"

import { MainNav } from "@/components/NavigationBar"

export function AuthenticatedLayout() {
  return (
    <div className="min-w-full">
      <header className="border-b">
        <div className="mx-auto max-w-7xl px-6 py-4">
          <MainNav />
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-6 py-8">
        <Outlet />
      </main>
    </div>
  )
}

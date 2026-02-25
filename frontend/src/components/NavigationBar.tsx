import { useRouterState } from "@tanstack/react-router"
import { cn } from "@/lib/utils"
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu"

export function MainNav() {
  const { location } = useRouterState()
  const pathname = location.pathname

  const activeClass = "bg-accent text-accent-foreground"

  return (
    <div className="flex items-center gap-8">
      <a href="/" className="text-sm font-semibold tracking-tight">
        LeetCode Buddy
      </a>
      <NavigationMenu>
        <NavigationMenuList>
          <NavigationMenuItem>
            <NavigationMenuLink
              href="/"
              active={pathname === '/'}
              className={cn(navigationMenuTriggerStyle(), pathname === '/' && activeClass)}
            >
              Dashboard
            </NavigationMenuLink>
          </NavigationMenuItem>

          <NavigationMenuItem>
            <NavigationMenuLink
              href="/problems"
              active={pathname === '/problems'}
              className={cn(navigationMenuTriggerStyle(), pathname === '/problems' && activeClass)}
            >
              Problems
            </NavigationMenuLink>
          </NavigationMenuItem>

          <NavigationMenuItem>
            <NavigationMenuLink
              href="/problems/submissions"
              active={pathname.startsWith('/problems/submissions')}
              className={cn(navigationMenuTriggerStyle(), pathname.startsWith('/problems/submissions') && activeClass)}
            >
              Submissions
            </NavigationMenuLink>
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>
    </div>
  )
}

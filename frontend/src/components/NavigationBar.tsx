import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu"

export function MainNav() {
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
              className={navigationMenuTriggerStyle()}
            >
              Dashboard
            </NavigationMenuLink>
          </NavigationMenuItem>

          <NavigationMenuItem>
            <NavigationMenuLink
              href="/problems"
              className={navigationMenuTriggerStyle()}
            >
              Problems
            </NavigationMenuLink>
          </NavigationMenuItem>

          <NavigationMenuItem>
            <NavigationMenuLink
              href="/problems/submissions"
              className={navigationMenuTriggerStyle()}
            >
              Submissions
            </NavigationMenuLink>
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>
    </div>
  )
}

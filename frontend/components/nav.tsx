"use client"

import * as React from "react"
import { Moon, Sun, Cake, LogOut } from "lucide-react"
import { useTheme } from "next-themes"
import { usePathname } from "next/navigation"
import { useAuth } from "@/src/context/AuthContext"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

// ModeToggle component to handle theme switching
export function ModeToggle() {
  const { setTheme } = useTheme()

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="icon">
          <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
          <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
          <span className="sr-only">Toggle theme</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setTheme("light")}>Light</DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("dark")}>Dark</DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("system")}>System</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

// Navigation component
const Nav: React.FC = () => {
  const pathname = usePathname()
  const { logout } = useAuth()
  const isDashboard = pathname === "/dashboard"

  return (
    <nav className="flex justify-between items-center p-4">
      <div className="flex items-center">
        <a href="/" className="text-xl font-bold">
          <Cake />
        </a>
      </div>
      <div className="flex items-center space-x-4">
        {isDashboard && (
          <Button variant="outline" className="shadow-none" onClick={logout}>
            <LogOut />
            <span className="sr-only">Logout</span>
          </Button>
        )}
        <ModeToggle  />
      </div>
    </nav>
  )
}

export default Nav

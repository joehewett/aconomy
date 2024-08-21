"use client"

import * as React from "react"

import { cn } from "../lib/utils"
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  navigationMenuTriggerStyle,
} from "./ui/navigation-menu"
import { Button } from "./ui/button"

const components: { title: string; href: string; description: string }[] = [
  {
    title: "Agents 👩‍🌾",
    href: "",
    description:
      "Agents are the entities that interact with the simulation. They are powered by Language Models.",
  },
  {
    title: "Wheat 🌾",
    href: "",
    description:
      "Wheat is a primary resource in the simulation. It is used to feed workers.",
  },
  {
    title: "Gold 🪙",
    href: "",
    description:
      "Gold is a primary resource in the simulation. It is used to purchase buildings and workers, and a certain quantity is required to win the game.",
  },
  {
    title: "Farms 🚜",
    href: "",
    description: "Farms are a building, which can be manned or unmanned. They produce wheat when manned, which requires a worker.",
  },
  {
    title: "Mine ⛏️",
    href: "",
    description:
      "A mine is a building, which can be manned or unmanned. They produce gold when manned, which requires a worker.",
  },
  {
    title: "Workers 👷",
    href: "",
    description:
      "Workers can be bought for gold and assigned to buildings. Manned buildings produce resources. Workers consume more wheat when they are working.",
  },
]

interface ActionBarProps {
  StartGame: () => void
  loading: boolean
  started: boolean
  stopGame: () => void
}

export function ActionBar({ StartGame, loading, stopGame, started }: ActionBarProps) {
  return (
    <NavigationMenu>
      <NavigationMenuList>
        <NavigationMenuItem>
          <NavigationMenuTrigger>Info</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid gap-3 p-6 md:w-[400px] lg:w-[500px] lg:grid-cols-[.75fr_1fr]">
              <li className="row-span-3">
                <NavigationMenuLink asChild>
                  <a
                    className="flex h-full w-full select-none flex-col justify-end rounded-md bg-gradient-to-b from-muted/50 to-muted p-6 no-underline outline-none focus:shadow-md"
                    href="/"
                  >
                    {/* <Icons.logo className="h-6 w-6" /> */}
                    <div className="mb-2 mt-4 text-lg font-medium">
                      aconomy🌾
                    </div>
                    <p className="text-sm leading-tight text-muted-foreground">
                      An agent-based simulation of a medieval economy, designed to test the abilities of AI agents to manage resources, strategise, plan, co-operate and compete.
                    </p>
                  </a>
                </NavigationMenuLink>
              </li>
              <SpanListItem title="">
                To get started, paste your OpenAI API key in the input field and click the "Start Simulation" button.
              </SpanListItem>
              <ClickableListItem href="https://github.com/joehewett/aconomy" title="GitHub 🔗">
                Check out the source code and submit issues here.
              </ClickableListItem>
              <ClickableListItem href="https://simulated.land" title="Explanation 🔗">
                Click here for a blog post explaining the simulation.
              </ClickableListItem>
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem>
          <NavigationMenuTrigger>Rules & Resources</NavigationMenuTrigger>
          <NavigationMenuContent>
            <ul className="grid w-[400px] gap-3 p-4 md:w-[500px] md:grid-cols-2 lg:w-[600px] ">
              {components.map((component) => (
                <SpanListItem
                  key={component.title}
                  title={component.title}
                >
                  {component.description}
                </SpanListItem>
              ))}
            </ul>
          </NavigationMenuContent>
        </NavigationMenuItem>
        <NavigationMenuItem>
          {!started && (
            <Button onClick={() => StartGame()} className="px-6" variant="outline">
              Start{loading ? 'ing' : ''} Simulation
              {loading && <span className="animate-spin ml-2">@</span>}
            </Button>

          ) || (
              <Button onClick={() => stopGame()} className="px-6 ml-4" variant="outline">
                Stop Simulation
              </Button>
            )}
        </NavigationMenuItem>
      </NavigationMenuList>
    </NavigationMenu>
  )
}

const ClickableListItem = React.forwardRef<
  React.ElementRef<"a">,
  React.ComponentPropsWithoutRef<"a">
>(({ className, title, children, ...props }, ref) => {
  return (
    <li>
      <NavigationMenuLink asChild>
        <a
          ref={ref}
          className={cn(
            "block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground",
            className
          )}
          {...props}
        >
          <div className="text-sm font-medium leading-none">{title}</div>
          <p className="text-sm leading-snug text-muted-foreground text-left">
            {children}
          </p>
        </a>
      </NavigationMenuLink>
    </li>
  )
})
ClickableListItem.displayName = "ClickableListItem"

const SpanListItem = React.forwardRef<
  React.ElementRef<"li">,
  React.ComponentPropsWithoutRef<"li">
>(({ className, title, children, ...props }, ref) => {
  return (
    <span
      ref={ref}
      className={cn(
        "block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none",
        className
      )}
      {...props}
    >
      <div className="text-sm font-medium leading-none">{title}</div>
      <p className="text-sm leading-snug text-muted-foreground text-left">
        {children}
      </p>
    </span>
  )
})
SpanListItem.displayName = "SpanListItem"

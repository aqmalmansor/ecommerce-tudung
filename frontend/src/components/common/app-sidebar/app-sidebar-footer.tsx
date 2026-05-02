"use client";

import { ChevronsUpDown, CircleUser, LogOut, User } from "lucide-react";
import { useSession, signOut } from "next-auth/react";

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui";
import { Role } from "@/types";

export function AppSidebarFooter() {
  const { data: session } = useSession();

  const name = session?.user?.name ?? "User";
  const email = session?.user?.email ?? "";
  const role = session?.user?.role;

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger
            render={
              <SidebarMenuButton
                size="lg"
                className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              />
            }
          >
            <div className="flex aspect-square size-8 items-center justify-center rounded-full bg-sidebar-primary text-sidebar-primary-foreground">
              <CircleUser className="size-4" />
            </div>
            <div className="flex flex-col gap-0.5 leading-none">
              <span className="font-medium truncate">{name}</span>
              <span className="text-xs text-muted-foreground truncate">
                {email}
              </span>
            </div>
            <ChevronsUpDown className="ml-auto" />
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 font-sans"
            align="start"
            side="top"
            sideOffset={4}
          >
            <DropdownMenuGroup>
              <DropdownMenuLabel className="p-0 font-normal">
                <div className="flex items-center gap-3 px-2 py-2">
                  <div className="flex size-10 items-center justify-center rounded-full bg-sidebar-primary text-sidebar-primary-foreground shrink-0">
                    <CircleUser className="size-5" />
                  </div>
                  <div className="flex flex-col min-w-0">
                    <span className="font-medium text-sm truncate">{name}</span>
                    <span className="text-xs text-muted-foreground truncate">
                      {email}
                    </span>
                  </div>
                </div>
              </DropdownMenuLabel>
            </DropdownMenuGroup>
            {role === Role.CLIENT && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                  <DropdownMenuItem render={<a href="/customer/profile" />}>
                    <User className="size-4" />
                    Profile
                  </DropdownMenuItem>
                </DropdownMenuGroup>
              </>
            )}
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem
                onClick={() => signOut({ callbackUrl: "/sign-in" })}
                className="text-destructive focus:text-destructive"
              >
                <LogOut className="size-4 text-destructive!" stroke="red" />
                Logout
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}

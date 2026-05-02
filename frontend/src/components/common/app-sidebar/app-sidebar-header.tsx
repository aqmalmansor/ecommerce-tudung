"use client";

import { Store } from "lucide-react";

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui";
import { useSession } from "next-auth/react";
import { Role } from "@/types";
import { METADATA } from "@/lib/constants";

export function AppSidebarHeader() {
  const { data: session } = useSession();

  const role = session?.user.role;
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        {role && (
          <SidebarMenuButton
            size="lg"
            render={
              <a
                href={`/${role === Role.ADMIN ? "admin" : "customer"}/dashboard`}
              />
            }
          >
            <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
              <Store className="size-4" />
            </div>
            <div className="flex flex-col gap-0.5 leading-none">
              <span className="font-semibold">{METADATA.name}</span>
              <span className="text-xs text-muted-foreground">
                {METADATA.description}
              </span>
            </div>
          </SidebarMenuButton>
        )}
      </SidebarMenuItem>
    </SidebarMenu>
  );
}

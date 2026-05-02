"use client";

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarProvider,
  SidebarFooter,
} from "@/components/ui";
import { usePathname } from "next/navigation";
import { AppSidebarFooter } from "./app-sidebar-footer";
import { AppSidebarHeader } from "./app-sidebar-header";
import { CustomSidebarProps } from "./config";

export function CustomSidebar({ navItems, ...props }: CustomSidebarProps) {
  const pathname = usePathname();

  const checkPathIsActive = (url: string): boolean => {
    if (url === "/dashboard") return pathname === url;
    return pathname === url || pathname.startsWith(url + "/");
  };

  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <AppSidebarHeader />
      </SidebarHeader>
      <SidebarContent>
        {navItems.map((item) => (
          <SidebarGroup key={item.title}>
            <SidebarGroupLabel className="font-sans">
              {item.title}
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {item.children.map((item) => (
                  <SidebarMenuItem className="font-sans" key={item.title}>
                    <SidebarMenuButton
                      render={<a href={item.url} />}
                      isActive={checkPathIsActive(item.url)}
                    >
                      {item.title}
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
      <SidebarRail />
      <SidebarFooter>
        <AppSidebarFooter />
      </SidebarFooter>
    </Sidebar>
  );
}

export const AppSidebar = ({
  children,
  navItems,
}: {
  children: React.ReactNode;
  navItems: CustomSidebarProps["navItems"];
}) => {
  return (
    <SidebarProvider>
      <CustomSidebar navItems={navItems} />
      {children}
    </SidebarProvider>
  );
};

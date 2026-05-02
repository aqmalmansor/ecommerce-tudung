import { ComponentProps } from "react";
import { Sidebar } from "@/components/ui";
import { NavItem } from "@/types";

export type CustomSidebarProps = ComponentProps<typeof Sidebar> & {
  navItems: NavItem[];
};

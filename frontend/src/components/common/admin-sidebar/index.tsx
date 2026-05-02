import { ReactNode } from "react";
import { AppSidebar } from "../app-sidebar";
import { navItems } from "./config";

export const AdminSidebar = ({ children }: { children: ReactNode }) => (
  <AppSidebar navItems={navItems}>{children}</AppSidebar>
);

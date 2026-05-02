import { NavItem } from "@/types";

export const navItems: NavItem[] = [
  {
    title: "Dashboard",
    url: "/dashboard",
    children: [
      {
        title: "Home",
        url: "/customer/dashboard",
      },
      {
        title: "Settings",
        url: "/customer/settings",
      },
    ],
  },
];

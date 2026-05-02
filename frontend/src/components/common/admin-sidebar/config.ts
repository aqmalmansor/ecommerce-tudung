import { NavItem } from "@/types";

export const navItems: NavItem[] = [
  {
    title: "Dashboard",
    url: "/dashboard",
    children: [
      {
        title: "Home",
        url: "/admin/dashboard",
      },
      {
        title: "Products",
        url: "/admin/products",
      },
      {
        title: "Customers",
        url: "/admin/customers",
      },
      {
        title: "Orders",
        url: "/admin/orders",
      },
      {
        title: "Report",
        url: "/admin/report",
      },
    ],
  },
];

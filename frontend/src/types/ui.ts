export interface NavItem {
  title: string;
  url: string;
  children: Omit<NavItem, "children">[];
}

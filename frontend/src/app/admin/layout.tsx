import { ReactNode } from "react";

import { AuthGuard } from "@/components/core";

const StaffAdminLayout = ({ children }: { children: ReactNode }) => {
  return <AuthGuard>{children}</AuthGuard>;
};

export default StaffAdminLayout;

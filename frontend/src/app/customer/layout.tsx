import { ReactNode } from "react";

import { AuthGuard } from "@/components/core";
import { SidebarInset } from "@/components/ui";
import { ClientSidebar } from "@/components/common";

const ClientLayout = ({ children }: { children: ReactNode }) => {
  return (
    <AuthGuard>
      <ClientSidebar>
        <SidebarInset>{children}</SidebarInset>
      </ClientSidebar>
    </AuthGuard>
  );
};

export default ClientLayout;

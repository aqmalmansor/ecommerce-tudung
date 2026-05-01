"use client";
import { FC, PropsWithChildren, useEffect } from "react";
import { useSession, signOut } from "next-auth/react";
import { useRouter } from "next/navigation";

import { ERROR_NAME } from "@/lib/constants";

export const AuthGuard: FC<PropsWithChildren> = ({ children }) => {
  const { data: session, status } = useSession();
  const router = useRouter();

  useEffect(() => {
    if (session?.error === ERROR_NAME.REFRESH_ACCESS_TOKEN) {
      signOut({ callbackUrl: "/sign-in" });
    }
  }, [session]);

  if (status === "loading") {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (status === "unauthenticated") {
    router.push("/sign-in");
    return null;
  }

  if (session?.error === ERROR_NAME.REFRESH_ACCESS_TOKEN) {
    return (
      <div className="text-center mt-8">
        Session expired. Redirecting to login...
      </div>
    );
  }

  return children;
};

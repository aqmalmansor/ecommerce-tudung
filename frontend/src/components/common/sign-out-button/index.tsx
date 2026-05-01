"use client";

import { signOut } from "next-auth/react";

export const SignOutButton = () => (
  <button
    onClick={() => signOut({ callbackUrl: "/sign-in" })}
    className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
  >
    Sign Out
  </button>
);

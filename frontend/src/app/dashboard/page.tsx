import { SignOutButton } from "@/components/common";
import { AuthGuard } from "@/components/core";
import { auth } from "@/lib/auth";
import Link from "next/link";

export default async function DashboardPage() {
  const session = await auth();

  return (
    <AuthGuard>
      <div className="container mx-auto px-4 py-8">
        <div className="bg-white rounded-lg shadow-md p-6 relative">
          <div className="absolute top-6 right-6">
            <SignOutButton />
          </div>
          <h1 className="text-3xl font-bold mb-4 text-gray-900">Dashboard</h1>
          <div className="mb-6">
            <p className="text-gray-900">
              Welcome back, <strong>{session?.user?.name}</strong>!
            </p>
            <p className="text-sm text-gray-900">
              Role: <span className="font-bold">{session?.user?.role}</span>
            </p>
          </div>
          <h3 className="text-md text-gray-900">Quick Links</h3>
          <div className="flex flex-row gap-3 flex-wrap items-center">
            <Link href="/profile" className="text-gray-900">
              Profile
            </Link>
          </div>
        </div>
      </div>
    </AuthGuard>
  );
}

import { LoginForm } from "@/components/pages/sign-in";
import Link from "next/link";

export const metadata = {
  title: "Login - My App",
  description: "Sign in to your account",
};

export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full">
        <h2 className="text-3xl font-bold text-center mb-8 text-gray-900">
          Sign In
        </h2>
        <LoginForm />
        <div className="mt-3">
          <Link href="/sign-up" className="text-gray-900">
            Sign Up
          </Link>
        </div>
      </div>
    </div>
  );
}

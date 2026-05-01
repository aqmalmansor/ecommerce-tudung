"use client";

import { useState } from "react";
import { AuthApi } from "@/services/auth-api";
import { useRouter } from "next/navigation";

export default function SignUpPage() {
  const [email, setEmail] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState("");
  const router = useRouter();

  const generateRandomEmail = () => {
    const randomNumber = Math.floor(Math.random() * 1000000)
      .toString()
      .padStart(6, "0");
    return `dev-${randomNumber}@xyz.com`;
  };

  const handleGenerateEmail = () => {
    setEmail(generateRandomEmail());
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setMessage("");

    const testEmail = email || generateRandomEmail();
    const fixedPassword = "TestPassword123!";

    try {
      await AuthApi.signUp({
        email: testEmail,
        password: fixedPassword,
        name: `Aqmal Mansor ${Math.random()}`,
      });

      router.push(`/sign-in?email=${encodeURIComponent(testEmail)}`);
    } catch (e) {
      console.log(e);
      setMessage("An error occurred during sign up");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Test Sign Up
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            For testing auth system
          </p>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div className="mb-4">
              <label
                htmlFor="email"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Email
              </label>
              <div className="flex space-x-2">
                <input
                  id="email"
                  name="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="flex-1 appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                  placeholder="Email will be auto-generated if empty"
                />
                <button
                  type="button"
                  onClick={handleGenerateEmail}
                  className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  Generate
                </button>
              </div>
            </div>
            <div className="mb-4">
              <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-700 mb-2"
              >
                Password (Fixed)
              </label>
              <input
                id="password"
                name="password"
                type="text"
                value="TestPassword123!"
                disabled
                className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-500 bg-gray-100 sm:text-sm"
              />
            </div>
          </div>
          {message && (
            <div
              className={`p-3 rounded-md text-sm ${
                message.startsWith("Error")
                  ? "bg-red-100 text-red-700"
                  : "bg-green-100 text-green-700"
              }`}
            >
              {message}
            </div>
          )}
          <div>
            <button
              type="submit"
              disabled={isLoading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? "Processing..." : "Test Sign Up"}
            </button>
          </div>
          <div className="text-center">
            <div className="text-xs text-gray-500 space-y-1">
              <div>Fixed Password: TestPassword123!</div>
              <div>
                Random Email Format: aqmal.mansor+dev-XXXXXX@ayp-group.com
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}

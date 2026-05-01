import NextAuth, { Session, User } from "next-auth";
import { JWT } from "next-auth/jwt";
import CredentialsProvider from "next-auth/providers/credentials";

import { AuthApi } from "@/services/auth-api";

import { A_DAY, ENV, ERROR_NAME } from "../constants";
import { DEBUG } from "./config";

async function refreshAccessToken(token: JWT) {
  try {
    if (!token.refreshToken) throw new Error("No refresh token available");

    const response = await AuthApi.refreshToken(token.refreshToken);

    return {
      ...token,
      accessToken: response.accessToken,
      refreshToken: response.refreshToken || token.refreshToken,
      accessTokenExpires: Date.now() + response.expiresIn * 1000,
    };
  } catch (error) {
    DEBUG.FAILED_REFRESH_TOKEN_GENERATION(error);
    return {
      ...token,
      error: ERROR_NAME.REFRESH_ACCESS_TOKEN,
    };
  }
}

async function authenticateUser(email: string, password: string) {
  try {
    const response = await AuthApi.signIn({
      email,
      password,
    });

    return response;
  } catch (error) {
    DEBUG.USER_AUTHENTICATION_ERROR(error);
    return null;
  }
}

export const authOptions = {
  providers: [
    CredentialsProvider({
      name: "credentials",
      credentials: {
        email: {},
        password: {},
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          return null;
        }

        const authResult = await authenticateUser(
          credentials.email.toString(),
          credentials.password.toString()
        );

        if (!authResult) return null;

        const user = {
          id: authResult.user?.id,
          name: authResult.user?.name,
          email: authResult.user?.email,
          role: authResult.user?.role,
          accessToken: authResult.accessToken,
          refreshToken: authResult.refreshToken,
          accessTokenExpires: Date.now() + authResult.expiresIn * 1000,
        };

        DEBUG.AUTHORIZE_FINAL_USER_OBJECT(user);
        return user;
      },
    }),
  ],

  callbacks: {
    async jwt({ token, user }: { token: JWT; user: User }) {
      if (user) {
        DEBUG.JWT_INITIAL_TOKEN(token, user);
        const newToken = {
          ...token,
          role: user.role,
          email: user.email,
          name: user.name,
          accessToken: user.accessToken,
          refreshToken: user.refreshToken,
          accessTokenExpires: user.accessTokenExpires,
        };
        DEBUG.JWT_MASSAGED_TOKEN(newToken);
        return newToken;
      }

      DEBUG.JWT_CHECK_TOKEN(token);

      if (Date.now() < token.accessTokenExpires) {
        return token;
      }

      return await refreshAccessToken(token);
    },

    async session({ session, token }: { session: Session; token: JWT }) {
      session.user.id = token.id;
      session.user.role = token.role;
      session.accessToken = token.accessToken;
      session.error = token.error;

      return session;
    },
  },

  secret: ENV.JWT_SECRET,
  session: {
    strategy: "jwt" as const,
    maxAge: 30 * A_DAY,
  },
};

export const { handlers, signIn, signOut, auth } = NextAuth(authOptions);

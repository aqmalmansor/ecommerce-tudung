import { User } from "next-auth";
import { JWT } from "next-auth/jwt";

import { DEBUG_ENABLED } from "../constants";

export const DEBUG = {
  JWT_INITIAL_TOKEN: (token: JWT, user: User) => {
    if (DEBUG_ENABLED) {
      console.log("JWT callback - Initial token", token);
      console.log("JWT callback - Initial user", user);
    }
  },
  JWT_MASSAGED_TOKEN: (token: JWT) => {
    if (DEBUG_ENABLED) {
      console.log("JWT callback - Updated token", token);
    }
  },
  JWT_CHECK_TOKEN: (token: JWT) => {
    if (DEBUG_ENABLED) {
      console.log("JWT callback - Check Token", {
        hasRefreshToken: !!token.refreshToken,
        accessTokenExpires: token.accessTokenExpires,
        currentTime: Date.now(),
        isExpired: Date.now() >= token.accessTokenExpires,
      });
    }
  },
  AUTHORIZE_FINAL_USER_OBJECT: (data: Record<string, unknown>) => {
    if (DEBUG_ENABLED) {
      console.log("Authorize returning user object:", data);
    }
  },
  USER_AUTHENTICATION_ERROR: (error: unknown) => {
    if (DEBUG_ENABLED) {
      console.error("Authentication error:", error);
    }
  },
  FAILED_REFRESH_TOKEN_GENERATION: (error: unknown) => {
    if (DEBUG_ENABLED) {
      console.log("Failed to regenerate access token: ", error);
    }
  },
};

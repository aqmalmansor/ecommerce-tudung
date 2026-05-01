export const A_SECOND = 1;
export const A_MINUTE = 60 * A_SECOND;
export const AN_HOUR = 60 * A_MINUTE;
export const A_DAY = 24 * AN_HOUR;

export const ENV = {
  BASE_URL: process.env.NEXT_PUBLIC_BASE_URL!,
  API_BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL!,
  JWT_SECRET: process.env.JWT_SECRET!,
};

export const DEBUG_ENABLED = process.env.NODE_ENV !== "production";

export const KY_CONFIG = {
  TIMEOUT: 30_000,
  RETRIES: 2,
};

export const CATCH_ERROR = {
  COMMON: "Something went wrong!",
};

export const ERROR_NAME = {
  REFRESH_ACCESS_TOKEN: "RefreshAccessTokenError",
};

export const URL_QUERY_KEY = {
  FROM: "from",
};

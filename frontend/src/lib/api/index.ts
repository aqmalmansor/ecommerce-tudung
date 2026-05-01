import ky from "ky";
import { Session } from "next-auth";

import { ENV, KY_CONFIG } from "../constants";

export interface BaseApiAttributes {
  path: string;
  body?: unknown;
  headers?: Record<string, string>;
  token?: string;
  session?: Session;
}

export type BaseApiMethods =
  | "get"
  | "patch"
  | "post"
  | "delete"
  | "all"
  | "head"
  | "put"
  | "options";

export const ApiClient = ky.create({
  prefixUrl: ENV.API_BASE_URL,
  credentials: "include",
  timeout: KY_CONFIG.TIMEOUT,
  retry: {
    limit: KY_CONFIG.RETRIES,
    methods: ["get", "head", "options"],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
  },
  hooks: {
    beforeRequest: [
      async (request) => {
        if (
          !request.headers.has("Content-Type") &&
          request.method !== "GET" &&
          request.method !== "HEAD"
        ) {
          request.headers.set("Content-Type", "application/json");
        }
      },
    ],
    afterResponse: [
      async (_request, _options, response) => {
        if (response.status === 401) {
          if (window) window.location.href = "/sign-in";
        }
        return response;
      },
    ],
  },
});

const BaseApiClient = (): {
  [K in BaseApiMethods]: <T>(
    baseApiAttributes: BaseApiAttributes
  ) => Promise<T>;
} => {
  const buildOptions = (attrs: BaseApiAttributes) => {
    const { body, headers, session } = attrs;

    const finalHeaders = { ...headers };

    if (session) {
      finalHeaders.Authorization = `Bearer ${session.accessToken}`;
    }

    const options: Record<string, unknown> = {
      headers: finalHeaders,
      throwHttpErrors: false,
    };

    if (body) {
      options.body = JSON.stringify(body);
    }

    return options;
  };

  return {
    get: async <T>(attrs: BaseApiAttributes): Promise<T> => {
      const response = await ApiClient.get<T>(attrs.path, buildOptions(attrs));

      const data = (await response.json()) as T;

      if (!response.ok) throw data;
      return data;
    },
    post: async <T>(attrs: BaseApiAttributes): Promise<T> => {
      const response = await ApiClient.post(attrs.path, buildOptions(attrs));
      const data = (await response.json()) as T;

      if (!response.ok) throw data;

      return data;
    },
    patch: async <T>(attrs: BaseApiAttributes): Promise<T> => {
      const response = await ApiClient.patch<T>(
        attrs.path,
        buildOptions(attrs)
      );
      const data = (await response.json()) as T;

      if (!response.ok) throw data;

      return data;
    },
    delete: async <T>(attrs: BaseApiAttributes): Promise<T> => {
      const response = await ApiClient.delete<T>(
        attrs.path,
        buildOptions(attrs)
      );
      const data = (await response.json()) as T;

      if (!response.ok) throw data;

      return data;
    },
    all: async <T>() => Promise.resolve() as Promise<T>,
    head: async <T>() => Promise.resolve() as Promise<T>,
    put: async <T>() => Promise.resolve() as Promise<T>,
    options: async <T>() => Promise.resolve() as Promise<T>,
  };
};

export const BaseApi = BaseApiClient();

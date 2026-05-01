import { BaseApi } from "@/lib/api";
import { IUser } from "@/types";

interface RegisterBodyRequest {
  name: string;
  email: string;
  password: string;
}

interface LoginBodyRequest {
  email: string;
  password: string;
}

interface RefreshTokenResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface LoginResponse {
  user: IUser;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export const AuthApi = {
  signUp: (body: RegisterBodyRequest) =>
    BaseApi.post({
      path: "api/auth/register",
      body,
    }),
  signIn: (body: LoginBodyRequest) =>
    BaseApi.post<LoginResponse>({ path: "api/auth/login", body }),
  refreshToken: (token: string) =>
    BaseApi.post<RefreshTokenResponse>({
      path: "api/auth/refresh",
      body: { token },
    }),
};

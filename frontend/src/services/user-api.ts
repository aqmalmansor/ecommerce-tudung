import { Session } from "next-auth";

import { IUser } from "@/types";
import { BaseApi } from "@/lib/api";

export interface GetUserProfileResponse {
  user: IUser;
}

export const UserApi = {
  getProfile: (session: Session) =>
    BaseApi.get<GetUserProfileResponse>({
      path: "api/user/profile",
      session,
    }),
};

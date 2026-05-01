export enum Role {
  CLIENT = "CLIENT",
  ADMIN = "ADMIN",
}

export interface IUser {
  id: number;
  email: string;
  name: string;
  role: Role;
}

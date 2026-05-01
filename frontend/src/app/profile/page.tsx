import { redirect } from "next/navigation";
import { UserApi } from "@/services/user-api";
import { auth } from "@/lib/auth";

export default async function Page() {
  const session = await auth();

  if (!session) redirect("/sign-in");

  const user = await UserApi.getProfile(session);

  return (
    <div>
      <pre>{JSON.stringify(user, null, 2)}</pre>
      <pre>{JSON.stringify(session, null, 2)}</pre>
    </div>
  );
}

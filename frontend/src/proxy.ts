import { NextAuthRequest } from "next-auth";
import { NextURL } from "next/dist/server/web/next-url";

import { auth } from "@/lib/auth";
import { ENV, URL_QUERY_KEY } from "@/lib/constants";

type URL = {
  path: string;
  redirect?: string | null;
};

const constructUrl = ({ path, redirect }: URL): NextURL => {
  const url = new NextURL(path, ENV.BASE_URL);

  if (redirect) {
    url.searchParams.set(URL_QUERY_KEY.FROM, redirect);
  }

  return url;
};

export default auth((req: NextAuthRequest) => {
  const url = req.nextUrl.clone();

  const publicRoutes = ["/sign-in", "/sign-up"];

  if (!req.auth && !publicRoutes.includes(req.nextUrl.pathname)) {
    const newUrl = new URL("/sign-in", req.nextUrl.origin);
    return Response.redirect(newUrl);
  }

  if (req.auth) {
    const permittedRoutes = ["/dashboard", "/profile"];

    if (!permittedRoutes.includes(url.pathname)) {
      return Response.redirect(constructUrl({ path: permittedRoutes[0] }));
    }
  }
});

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico|svg).*)"],
};

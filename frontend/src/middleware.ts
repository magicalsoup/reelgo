import { User } from '@/lib/types'
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { cookies } from 'next/headers';

export const config = {
    matcher: [
      /*
       * Match all request paths except for the ones starting with:
       * - api (API routes)
       * - _next/static (static files)
       * - _next/image (image optimization files)
       * - favicon.ico, sitemap.xml, robots.txt (metadata files)
       */
      {
        source:
          '/((?!api|_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)',
        missing: [
          { type: 'header', key: 'next-router-prefetch' },
          { type: 'header', key: 'purpose', value: 'prefetch' },
        ],
      },
   
      {
        source:
          '/((?!api|_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)',
        has: [
          { type: 'header', key: 'next-router-prefetch' },
          { type: 'header', key: 'purpose', value: 'prefetch' },
        ],
      },
   
      {
        source:
          '/((?!api|_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)',
        has: [{ type: 'header', key: 'x-present' }],
        missing: [{ type: 'header', key: 'x-missing', value: 'prefetch' }],
      },
    ],
}

function isHomeRequest(request: NextRequest) {
  return request.nextUrl.pathname === "/"
}

function isAuthRequest(request: NextRequest) {
  return request.nextUrl.pathname.startsWith("/login") || request.nextUrl.pathname.startsWith("/signup")
}

export async function middleware(request: NextRequest) {
  
    const cookieStore = (await cookies())

    const res = await fetch(`${process.env.NEXT_PUBLIC_REEL_GO_SERVER_API_ENDPOINT}/user`, {
        method: "GET",
        credentials: "include",
        headers: {
            Cookie: cookieStore.toString()
        }
    })

    const user: User = await res?.json().catch((err) => {
      return null
    })

    const loggedIn = res.ok && user

    if (!loggedIn && !isAuthRequest(request) && !isHomeRequest(request)) {
      return NextResponse.redirect(new URL("/login", request.url))
    }


    if (loggedIn && (isAuthRequest(request) || isHomeRequest(request))) {
      return NextResponse.redirect(new URL("/dashboard", request.url)) 
    }

    if (loggedIn && request.nextUrl.pathname.startsWith("/dashboard") && !user.verified) {
      return NextResponse.redirect(new URL(`/link?uid=${user.uid}`, request.url))
    }
    return NextResponse.next()
}
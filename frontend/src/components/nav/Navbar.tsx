"use client"

import { usePathname } from "next/navigation"
import { Button } from "@/components/ui/button"
import LogOutButton from "../dashboard/LogOutButton"
import { useRouter } from "next/navigation"

export default function Navbar() {
    const pathName = usePathname()
    const router = useRouter()
    const show = pathName == "/" && !pathName.startsWith("/login") && !pathName.startsWith("/signup")

    return (
        <div className="sticky absolute top-0 outline outline-1 outline-gray-400 bg-white flex justify-center">
            <div className="flex justify-between w-full max-w-screen-xl h-20 items-center px-8 py-4">
                <button className="logo font-bold text-4xl" onClick={() => router.push("/login")}>reelgo</button>
                {show && 
                    <div className="flex gap-x-2">
                        <Button variant="outline" onClick={() => router.push("/signup")}>Sign Up</Button>
                        <Button onClick={() => router.push("/login")}>Log in</Button>
                    </div>
                }
                {!show && <LogOutButton/>}
            </div>
        </div>
    )
}
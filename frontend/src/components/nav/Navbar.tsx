"use client"

import { usePathname } from "next/navigation"
import { Button } from "@/components/ui/button"

export default function Navbar() {
    const pathName = usePathname()
    const show = pathName == "/" && !pathName.startsWith("/login") && !pathName.startsWith("/signup")

    return (
        <div className="sticky absolute top-0 outline outline-1 outline-gray-400 flex justify-center">
            <div className="flex justify-between w-full max-w-[1280px] h-20 items-center px-8 py-4">
                <strong className="logo font-bold text-4xl">reelgo</strong>
                {show && 
                    <div className="flex gap-x-2">
                        <Button variant="outline">Sign Up</Button>
                        <Button>Log in</Button>
                    </div>
                }
            </div>
        </div>
    )
}
"use client"

import Navbar from "@/components/nav/Navbar";
import OTP from "./otp";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";

export default function Home() {

    const searchParams = useSearchParams()
    const uidStr = searchParams.get("uid") ?? ""
    const uid = parseInt(uidStr)

    const cmd = `!verify:${uid}`
    
    return (
        <div className="flex flex-col h-screen gap-y-12 font-sans">
            <Navbar/>
            <div className="flex flex-col h-full justify-center items-center gap-y-6">
                <OTP uid={uid}/>
                <div className="text-xl flex flex-col items-center gap-y-2">
                    <p>You're almost done! copy the following command and send it as a message to <a className="underline text-sky-500" href="https://www.instagram.com/reelgo.app/">reelgo.app</a> on instagram and type back the 6 digit authentication code</p>
                    <Button variant="outline" className="text-xl w-fit" onClick={async () => await navigator.clipboard.writeText(cmd)}>
                        Copy Command
                    </Button>
                    <p> </p>
                </div>
            </div>
        </div>
    )
}
"use client"
import OTP from "./otp";
import { useSearchParams } from "next/navigation";

export default function Home() {

    const searchParams = useSearchParams()
    const uidStr = searchParams.get("uid") ?? ""
    const uid = parseInt(uidStr)
    
    return (
        <div className="flex flex-col h-screen justify-center items-center">
            <h1>You're almost done! copy and send the message "!verify:{uid}" to reelapp.go on instagram and type back the 6 digit authentication code {uid} </h1>
            <OTP uid={uid}/>
        </div>
    )
}
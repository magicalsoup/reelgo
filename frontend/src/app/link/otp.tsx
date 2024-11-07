"use client"

import { InputOTP, InputOTPGroup, InputOTPSlot, InputOTPSeparator } from "@/components/ui/input-otp";
import { redirect } from "next/navigation";
import { useState, useEffect } from "react";

export default function OTP({uid} : {uid: number | undefined}) {
    const [verificationCode, setVerificationCode] = useState("")
    const [message, setMessage] = useState("");

    useEffect(() => {
        async function checkCode(code: string) {
            const api_url = `${process.env.NEXT_PUBLIC_REEL_GO_SERVER_API_ENDPOINT}/auth`
            const req = await fetch(api_url, {
                method: "POST",
                credentials: "include",
                body: JSON.stringify({
                    uid: uid ?? "",
                    code: code,                  
                })
            })

            if (req.ok) {
                setMessage("You're all set! redirecting you to dashboard...");
                redirect("/dashboard");
            } else {
                setMessage("It seems you entered the wrong code, try again");
            }
        }

        if (verificationCode.length == 6) {
            // send verification code to server to link
            checkCode(verificationCode)
        }
    }, [verificationCode])
    return (
        <>
            <p className="h-8">{message}</p>
            <InputOTP maxLength={6}
                value={verificationCode}
                onChange={(value) => setVerificationCode(value)}>
                <InputOTPGroup>
                    <InputOTPSlot index={0} />
                    <InputOTPSlot index={1} />
                    <InputOTPSlot index={2} />
                </InputOTPGroup>
                <InputOTPSeparator />
                <InputOTPGroup>
                    <InputOTPSlot index={3} />
                    <InputOTPSlot index={4} />
                    <InputOTPSlot index={5} />
                </InputOTPGroup>
            </InputOTP>
        </>
    )
}
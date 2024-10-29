"use client"

import { InputOTP, InputOTPGroup, InputOTPSlot, InputOTPSeparator } from "@/components/ui/input-otp";
import { redirect } from "next/navigation";
import { useState, useEffect } from "react";

export default function OTP({uid} : {uid: string | undefined}) {
    const [verificationCode, setVerificationCode] = useState("")
    const [message, setMessage] = useState("");

    useEffect(() => {
        async function checkCode(code: string) {
            const api_url = `${process.env.REEL_GO_SERVER_API_ENDPOINT}/auth`
            const req = await fetch(api_url, {
                method: "GET",
                headers: {
                    "Server-Auth-Token": process.env.SERVER_AUTH_TOKEN ?? "",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    uid: uid ?? "",
                    code: code,                  
                })
            })

            if (req.status == 200) {
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
            {message}
        </>
    )
}
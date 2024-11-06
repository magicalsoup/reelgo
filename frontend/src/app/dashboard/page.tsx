"use client"
import LogOutButton from "@/components/dashboard/LogOutButton";
import { useEffect } from "react";

export default function Home() {

    useEffect(() => {

        async function fetchTrips() {
            const api_url = `${process.env.NEXT_PUBLIC_REEL_GO_SERVER_API_ENDPOINT}/trips`

            const res = await fetch(api_url, {
                credentials: "include",
            })

            console.log(res)

            if (!res.ok) {
                console.error("something went wrong in fetch trips ", res.status, res.statusText, await res.text())
                return
            }
        
            const data = await res.json()
            console.log(data)
        }

        fetchTrips()
    }, [])

    return (
        <div>
            <h1>Your dashboard</h1>
            <LogOutButton/>
        </div>
    )
}
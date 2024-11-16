"use client"
import { useEffect } from "react";

export default function Home() {

    useEffect(() => {
        async function fetchTrips() {
            const api_url = `${process.env.NEXT_PUBLIC_REEL_GO_SERVER_API_ENDPOINT}/trips`

            const res = await fetch(api_url, {
                credentials: "include",
            })

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
        <div className="h-screen">
        <div className="flex flex-col h-full font-sans p-32">
            <h1 className="text-5xl">Your Dashboard</h1>
        </div>
      </div>
    )
}
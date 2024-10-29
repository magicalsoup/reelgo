import * as React from "react"
import { cookies } from "next/headers"
import Login from "./login"
 

export default async function Home() {

    const cookieStore = await cookies();

    return (
        <div className="flex justify-center items-center h-screen">
            <Login cookieStore={cookieStore}/>
        </div>
  )
}

import * as React from "react"
import { cookies } from "next/headers"
import Login from "./signup"
import SignUp from "./signup";
 

export default async function Home() {

    const cookieStore = await cookies();

    return (
        <div className="flex justify-center items-center h-screen">
            <SignUp cookieStore={cookieStore}/>
        </div>
  )
}

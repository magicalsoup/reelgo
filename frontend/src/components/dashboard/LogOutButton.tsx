
import { Button } from "../ui/button";
import { redirect } from "next/navigation";


export default function LogOutButton() {
    async function logoutUser() {
        // const cookieStore = await cookies()
        const api_url = `${process.env.NEXT_PUBLIC_REEL_GO_SERVER_API_ENDPOINT}/logout`
        const res = await fetch(api_url, {
            credentials: "include"
        })
        if (!res.ok) {
            // do some error handling here
            console.error("could not log user out", await res.text())
            return 
        }
        redirect("/")
    }

    return (
        <Button onClick={() => logoutUser()}>Log Out</Button>
    )
}
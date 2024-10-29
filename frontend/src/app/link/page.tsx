
import { USER_ID_COOKIE_NAME } from "@/lib/cookies";
import { cookies } from "next/headers";
import OTP from "./otp";

export default async function Home() {
    
    const cookieStore = await cookies();
    const uid = cookieStore.get(USER_ID_COOKIE_NAME)?.value;
    
    return (
        <div className="flex flex-col h-screen justify-center items-center">
            <h1>You're almost done! send a message "!verify" to reelapp.go on instagram and type back the 6 digit authentication code</h1>
            <OTP uid={uid}/>
        </div>
    )
}
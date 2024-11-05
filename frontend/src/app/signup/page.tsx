"use client"

import * as React from "react"
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input"
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { redirect } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
    email: z.string(),
    password: z.string()
})
  
const crypto = require("crypto")

export default function Home() {
    const [errorMessage, setErrorMessage] = React.useState<String>("")

    async function loginUser (values: z.infer<typeof formSchema>) {
        const hashedPassword = crypto.createHash("sha256").update(values.password).digest("hex")
        const api_url = `${process.env.NEXT_PUBLIC_REEL_GO_SERVER_API_ENDPOINT}/signup`

        console.log("api_url: ", api_url)
    
        const res = await fetch(api_url, {
            method: "POST",
            credentials: 'include',
            body: JSON.stringify({
                email: values.email,
                hashedPassword: hashedPassword
            })
        })
    
        console.log("[res set cookie headers]", res.headers.getSetCookie())
    
        if (!res.ok) {
            console.error("cors stupid ", res)
            return false
        }
        
        if (res.status == 201) { // status created
            setErrorMessage("")
        } else {
            setErrorMessage("error signing you up")
        }
    }

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            email: "",
            password: "",
        },
    })


    return (
        <div className="flex justify-center items-center h-screen">
            <Card className="w-[350px]">
                <CardHeader>
                    <CardTitle>Start planning with ReelGo!</CardTitle>
                </CardHeader>
                <CardContent className="flex flex-col gap-y-2">
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(loginUser)} className="space-y-4">
                        <FormField
                        control={form.control}
                        name="email"
                        render={({ field }) => (
                            <FormItem>
                            <FormLabel>email</FormLabel>
                            <FormControl>
                                <Input placeholder="abc@example.com" {...field}/>
                            </FormControl>
                            <FormMessage />
                            </FormItem>
                        )}
                        />
                        <FormField
                        control={form.control}
                        name="password"
                        render={({ field }) => (
                            <FormItem>
                            <FormLabel>password</FormLabel>
                            <FormControl>
                                <Input placeholder="password" type="password" {...field}/>
                            </FormControl>
                            <FormMessage />
                            </FormItem>
                        )}
                        />
                        <FormDescription>
                            {errorMessage}
                        </FormDescription>
                        <div className="flex justify-between">
                            <Button type="submit" variant="outline">Sign up</Button>
                        </div>
                    </form>
                </Form>
                </CardContent>
            </Card>
        </div>
  )
}

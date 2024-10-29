"use client"

import * as React from "react"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { Form, FormField, FormItem, FormLabel, FormControl, FormDescription, FormMessage } from "@/components/ui/form"
import { useState } from "react"
import { redirect } from "next/navigation"
import { User } from "@/lib/types"
import { USER_BEARER_TOKEN_COOKIE_NAME, USER_ID_COOKIE_NAME } from "@/lib/cookies"
import { ReadonlyRequestCookies } from "next/dist/server/web/spec-extension/adapters/request-cookies"
 
const formSchema = z.object({
  email: z.string(),
  password: z.string()
})

const crypto = require("crypto")

export default  function Login({ cookieStore } : { cookieStore: ReadonlyRequestCookies }) {
    const [errorMessage, setErrorMessage] = useState<String>("")

    async function loginUser (values: z.infer<typeof formSchema>) {
        const hashedPassword = crypto.createHash("sha256").update(values.password).digest("hex")
        const api_url = `${process.env.REEL_GO_SERVER_API_ENDPOINT}/login`
        const res = await fetch(api_url, {
            method: "POST",
            headers: {
                "Server-Auth-Token": process.env.SERVER_AUTH_TOKEN ?? "No Token",
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: values.email,
                password: hashedPassword
            })
        })

        const userData: User = await res.json() // contains the user data

        // get a sessionToken and store it in the cookie store
        if (res.status == 200) {
            
            cookieStore.set(USER_BEARER_TOKEN_COOKIE_NAME, userData.bearerToken, {expires: userData.expiry_time});
            cookieStore.set(USER_ID_COOKIE_NAME, userData.id.toString(), { expires: userData.expiry_time});

            setErrorMessage("")
            redirect("/dashboard")
        } else {
            setErrorMessage("error, cannot sign you in")
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
        <Card className="w-[350px]">
            <CardHeader>
                <CardTitle>Welcome Back!</CardTitle>
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
                        <Button type="submit">Sign in!</Button>
                    </div>
                </form>
            </Form>
            </CardContent>
        </Card>
  )
}

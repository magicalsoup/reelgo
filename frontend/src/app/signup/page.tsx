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
import sha256Hash from "@/lib/hash";
import { useRouter } from "next/navigation";

const formSchema = z.object({
    name: z.string(),
    email: z.string(),
    password: z.string()
})

export default function Home() {
    const [errorMessage, setErrorMessage] = React.useState<String>("")
    const router = useRouter()

    async function signupUser (values: z.infer<typeof formSchema>) {
        const hashedPassword = sha256Hash(values.password)
        const api_url = `${process.env.NEXT_PUBLIC_REEL_GO_SERVER_API_ENDPOINT}/signup`
    
        const res = await fetch(api_url, {
            method: "POST",
            credentials: 'include',
            body: JSON.stringify({
                name: values.name,
                email: values.email,
                hashedPassword: hashedPassword
            })
        })


        if (!res.ok) {
            setErrorMessage("error signing you up")
            return
        }

        const user = await res.json()
        
        setErrorMessage("redirecting you to setup...")
        router.push(`/link?uid=${user.uid}`)
    }

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            name: "",
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
                    <form onSubmit={form.handleSubmit(signupUser)} className="space-y-4">
                        <FormField
                        control={form.control}
                        name="name"
                        render={({ field }) => (
                            <FormItem>
                            <FormLabel>name</FormLabel>
                            <FormControl>
                                <Input placeholder="name" {...field}/>
                            </FormControl>
                            <FormMessage />
                            </FormItem>
                        )}
                        />
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

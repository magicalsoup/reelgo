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
 
const formSchema = z.object({
  username: z.string(),
  password: z.string()
})

export default function Home() {

    const [errorMessage, setErrorMessage] = useState<String>("")

    function loginUser (values: z.infer<typeof formSchema>) {
        console.log(values.username + " " + values.password)
        setErrorMessage("Logging you in...")
    }


    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            username: "",
            password: "",
        },
    })

    return (
    <div className="flex justify-center items-center h-screen">
        <Card className="w-[350px]">
            <CardHeader>
                <CardTitle>Welcome Back!</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-y-2">
            <Form {...form}>
                <form onSubmit={form.handleSubmit(loginUser)} className="space-y-4">
                    <FormField
                    control={form.control}
                    name="username"
                    render={({ field }) => (
                        <FormItem>
                        <FormLabel>username</FormLabel>
                        <FormControl>
                            <Input placeholder="username" {...field}/>
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
    </div>
  )
}

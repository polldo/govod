import Layout from '@/components/layout'
import Logout from '@/components/logout'
import Head from 'next/head'
import { useSession } from '@/session/context'
import { useEffect } from 'react'
import { useState } from 'react'
import { useRouter } from 'next/router'
import { useFetch } from '@/services/fetch'
import Image from 'next/image'
import { toast } from 'react-hot-toast'

type Course = {
    name: string
    description: string
    image: string
}

function Card(props: Course) {
    return (
        <a
            href="#"
            className="flex w-full flex-col items-center rounded-lg border border-gray-200 bg-white shadow hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-700 md:max-w-xl md:flex-row"
        >
            <Image
                className="w-full rounded-t-lg border border-red-800 object-contain md:w-20"
                alt=""
                src={props.image}
                width={80}
                height={32}
            />

            <div className="flex flex-col justify-between p-4 leading-normal">
                <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{props.name}</h5>
                <p className="mb-3 font-normal text-gray-700 dark:text-gray-400">{props.description}</p>
            </div>
        </a>
    )
}

export default function Dashboard() {
    const [courses, setCourses] = useState<Course[]>([])
    const router = useRouter()
    const { isLoggedIn, isLoading } = useSession()
    const fetch = useFetch()

    useEffect(() => {
        fetch('http://mylocal.com:8000/courses/owned')
            .then((res) => {
                if (!res.ok) {
                    throw new Error()
                }
                return res.json()
            })
            .then((data) => setCourses(data))
            .catch(() => {
                toast.error('Something went wrong')
            })
    }, [fetch])

    if (isLoading) {
        return null
    }

    if (!isLoggedIn) {
        router.push('login')
        return null
    }

    return (
        <>
            <Head>
                <title>Dashboard</title>
            </Head>
            <Layout>
                <div className="flex w-1/2 flex-col">
                    <div>
                        <p>Hello, this is your dashboard!</p>
                        <Logout></Logout>
                    </div>

                    <div className="flex flex-col items-center space-y-5 pt-6 pb-6">
                        {courses && courses.map((course) => <Card {...course} key={course.name} />)}
                        {!courses && <p>No courses yet.</p>}
                    </div>
                </div>
            </Layout>
        </>
    )
}

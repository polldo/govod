import Layout from '@/components/layout'
import Head from 'next/head'
import { useSession } from '@/session/context'
import { useRouter } from 'next/router'
import Image from 'next/image'
import Link from 'next/link'
import useSWR from 'swr'
import { Course } from '@/services/types'

function Card(props: Course) {
    return (
        <Link
            href={`/dashboard/course/${encodeURIComponent(props.id)}`}
            className="mx-auto flex w-2/3 max-w-3xl flex-col items-center rounded-lg border border-gray-200 bg-white pt-5 shadow hover:bg-gray-100 md:w-1/2 lg:w-2/3"
        >
            <Image
                className="max-w-1/6 w-2/3 rounded-t-lg object-contain md:m-8 "
                alt=""
                src={props.imageUrl}
                width={80}
                height={32}
            />

            <div className="flex w-full flex-col p-4 leading-normal">
                <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{props.name}</h5>
                <p className="mb-6 mt-2 font-normal text-gray-700 dark:text-gray-400">{props.description}</p>
            </div>
        </Link>
    )
}

export default function Dashboard() {
    const router = useRouter()
    const { isLoggedIn, isLoading } = useSession()

    const { data: courses } = useSWR<Course[]>(isLoggedIn ? '/courses/owned' : null)

    if (isLoading) {
        return null
    }

    if (!isLoggedIn) {
        router.push('/login')
        return null
    }

    return (
        <>
            <Head>
                <title>Dashboard</title>
            </Head>
            <Layout>
                <div className="flex w-full flex-col">
                    <div className="flex w-full p-10">
                        <p className="mx-auto text-xl">Hello, this is your dashboard!</p>
                    </div>

                    <div className="grid items-stretch space-y-5 p-4 pt-6 md:grid-cols-2 md:gap-y-10 md:space-y-0 lg:grid-cols-3">
                        {courses?.map((course) => (
                            <Card {...course} key={course.name} />
                        ))}
                    </div>

                    {courses?.length == 0 && (
                        <div className="mx-auto flex w-1/2 rounded-sm border border-blue-900 p-2">
                            <p>No courses yet...</p>
                        </div>
                    )}
                </div>
            </Layout>
        </>
    )
}

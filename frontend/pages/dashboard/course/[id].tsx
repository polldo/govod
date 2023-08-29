import Layout from '@/components/layout'
import Head from 'next/head'
import Image from 'next/image'
import { useSession } from '@/session/context'
import { useEffect } from 'react'
import { useState } from 'react'
import { useRouter } from 'next/router'
import { useFetch } from '@/services/fetch'
import { toast } from 'react-hot-toast'

type Course = {
    name: string
    description: string
    image: string
}

type Video = {
    name: string
    description: string
}

export default function DashboardCourse() {
    const [course, setCourse] = useState<Course>()
    const [videos, setVideos] = useState<Video[]>()
    const { isLoggedIn, isLoading } = useSession()
    const fetch = useFetch()
    const router = useRouter()
    const { id } = router.query

    useEffect(() => {
        if (!router.isReady) {
            return
        }

        fetch('http://mylocal.com:8000/courses/' + id)
            .then((res) => {
                if (!res.ok) {
                    throw new Error()
                }
                return res.json()
            })
            .then((data) => setCourse(data))
            .catch(() => {
                toast.error('Something went wrong')
            })
    }, [id, fetch, router.isReady])

    useEffect(() => {
        if (!router.isReady) {
            return
        }

        fetch('http://mylocal.com:8000/courses/' + id + '/videos')
            .then((res) => {
                if (!res.ok) {
                    throw new Error()
                }
                return res.json()
            })
            .then((data) => setVideos(data))
            .catch(() => {
                toast.error('Something went wrong')
            })
    }, [id, fetch, router.isReady])

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
                <title>Course - {course?.name}</title>
            </Head>
            <Layout>
                <div className="flex w-1/2 flex-col">
                    <p>Enjoy the course</p>
                    <p>{course?.name}</p>
                    <br></br>
                    <p>{course?.description}</p>

                    <div className="flex flex-col items-center space-y-5 pt-6 pb-6">
                        {videos && videos.map((video) => <Card {...video} key={video.name} />)}
                    </div>
                </div>
            </Layout>
        </>
    )
}

function Card(props: Video) {
    return (
        <a
            href="#"
            className="flex w-1/2 flex-col items-center rounded-lg border border-gray-200 bg-white shadow hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-700 md:max-w-xl md:flex-row"
        >
            <Image
                className="w-full rounded-t-lg border border-red-800 object-contain md:w-20"
                alt=""
                src={''}
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
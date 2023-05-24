import Layout from '@/components/layout'
import Logout from '@/components/logout'
import Head from 'next/head'
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

export default function Dashboard() {
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
                <title>Dashboard</title>
            </Head>
            <Layout>
                <div className="flex w-1/2 flex-col">
                    <div>
                        <p>Hello, this is your dashboard!</p>
                        <Logout></Logout>
                    </div>

                    <p>{course?.name}</p>
                    <br></br>
                    <p>{course?.description}</p>

                    {videos && videos.map((video) => <p key={video.name}> {video.name} </p>)}
                </div>
            </Layout>
        </>
    )
}

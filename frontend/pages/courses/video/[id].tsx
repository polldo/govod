import Layout from '@/components/layout'
import VideoJS from '@/components/videoplayer'
import Head from 'next/head'
import Link from 'next/link'
import { useEffect } from 'react'
import { useState } from 'react'
import { useRouter } from 'next/router'
import { fetcher } from '@/services/fetch'
import { toast } from 'react-hot-toast'
import React from 'react'

type Course = {
    name: string
}

type Video = {
    id: string
    index: number
    course_id: string
    name: string
    description: string
}

export default function CourseDetails() {
    const [video, setVideo] = useState<Video>()
    const [url, setUrl] = useState<string>()
    const [course, setCourse] = useState<Course>()

    const router = useRouter()
    const { id } = router.query

    useEffect(() => {
        if (!router.isReady) {
            return
        }
        fetcher
            .fetch('http://mylocal.com:8000/videos/' + id + '/free')
            .then((res) => {
                return res.json()
            })
            .then((data) => {
                setVideo(data.video)
                setCourse(data.course)
                setUrl(data.url)
            })
            .catch(() => {
                toast.error('Something went wrong')
            })
    }, [id, router.isReady])

    const videoJsOptions = {
        controls: true,
        responsive: true,
        fluid: true,
        fill: true,
        sources: [
            {
                type: 'video/youtube',
                src: url,
            },
        ],
    }

    if (!video) {
        return null
    }

    return (
        <>
            <Head>
                <title>Video - {video?.name}</title>
            </Head>
            <Layout>
                <div className="mt-10 flex w-2/4 flex-col">
                    <div className="w-full">{url && <VideoJS options={videoJsOptions} onReady={() => {}} />}</div>

                    <div className="mx-16 mt-5 flex flex-col p-4 sm:mx-20">
                        <h2 className="text-base font-bold sm:text-xl">{video.name}</h2>
                        <p className="mt-2 text-base italic sm:text-xl">{video.description}</p>
                        <Link
                            href={`/courses/${encodeURIComponent(video.course_id)}`}
                            className="mt-2 w-20 cursor-pointer text-sm text-blue-500 underline"
                        >
                            {course?.name}
                        </Link>
                    </div>
                </div>
            </Layout>
        </>
    )
}

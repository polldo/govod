import Layout from '@/components/layout'
import VideoJS from '@/components/videoplayer'
import Head from 'next/head'
import Link from 'next/link'
import { useEffect } from 'react'
import { useState } from 'react'
import { useRouter } from 'next/router'
import { useFetch } from '@/services/fetch'
import { toast } from 'react-hot-toast'
import React from 'react'
import { useSession } from '@/session/context'

type Course = {
    name: string
}

type Video = {
    id: string
    index: number
    course_id: string
    name: string
    description: string
    free: boolean
}

export default function CourseDetails() {
    const [video, setVideo] = useState<Video>()
    const [url, setUrl] = useState<string>()
    const [course, setCourse] = useState<Course>()
    const [videos, setVideos] = useState<Video[]>()
    const { isLoggedIn, isLoading } = useSession()
    const fetch = useFetch()
    const router = useRouter()
    const { id } = router.query
    const playerRef = React.useRef(null)

    useEffect(() => {
        if (!router.isReady) {
            return
        }
        fetch('http://mylocal.com:8000/videos/' + id + '/full')
            .then((res) => {
                if (!res.ok) {
                    throw new Error()
                }
                return res.json()
            })
            .then((data) => {
                setVideo(data.video)
                setCourse(data.course)
                setVideos(data.all_videos)
                setUrl(data.url)
            })
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

    const handlePlayerReady = (player: any) => {
        playerRef.current = player
        player.on('waiting', () => {})
        player.on('dispose', () => {})
    }

    return (
        <>
            <Head>
                <title>Video - {video?.name}</title>
            </Head>
            <Layout>
                <div className="h-full w-full">
                    <div className="flex w-full ">
                        <div className="ml-16 mt-10 flex w-1/2 flex-row ">
                            <div className="w-full">
                                {url && <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />}
                            </div>
                        </div>
                        <div className="mx-0 mt-10 flex w-1/4 flex-col border border-black">
                            {videos &&
                                videos
                                    .slice()
                                    .sort((a, b) => a.index - b.index)
                                    .map((vid) => (
                                        <Link key={vid.name} href={`/dashboard/video/${encodeURIComponent(vid.id)}`}>
                                            {vid.name} {vid.index == video!.index ? '<' : ''}
                                        </Link>
                                    ))}
                            {!videos && <p>No videos here.</p>}
                        </div>
                    </div>
                    <div className="mx-20 mt-5 flex flex-col p-4">
                        <h2 className="text-xl font-bold">{video?.name}</h2>
                        <p className="mt-2 text-base italic">{video?.description}</p>
                        {video && (
                            <Link
                                href={`/dashboard/course/${encodeURIComponent(video?.course_id)}`}
                                className="mt-2 block cursor-pointer text-sm text-blue-500 underline"
                            >
                                {course?.name}
                            </Link>
                        )}
                    </div>
                </div>
            </Layout>
        </>
    )
}

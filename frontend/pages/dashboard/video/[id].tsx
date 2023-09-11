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

    var next = ''
    var prev = ''
    if (video && videos) {
        const sorted = videos.slice().sort((a, b) => a.index - b.index)

        // Videos index starts from 1.
        const idx = video.index - 1
        const isLast = idx === sorted.length - 1
        const isFirst = idx === 0

        next = isLast ? '' : sorted[idx + 1].id
        prev = isFirst ? '' : sorted[idx - 1].id
    }

    return (
        <>
            <Head>
                <title>Video - {video?.name}</title>
            </Head>
            <Layout>
                <div className="h-full w-full">
                    <div className="flex w-full">
                        <div className="mx-16 mt-10 flex w-full lg:mr-0 lg:ml-auto lg:w-[700px] xl:w-[850px]">
                            <div className="w-full">
                                {url && <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />}
                            </div>
                        </div>

                        <div className="mx-0 mt-10 hidden w-[300px] flex-col overflow-y-scroll border border-black lg:mr-auto lg:flex lg:max-h-[394px] xl:max-h-[478px]">
                            {videos &&
                                videos
                                    .slice()
                                    .sort((a, b) => a.index - b.index)
                                    .map((vid) => (
                                        <Link
                                            key={vid.name}
                                            href={`/dashboard/video/${encodeURIComponent(vid.id)}`}
                                            className="mx-2 my-2 text-sm"
                                        >
                                            {vid.name} {vid.index == video!.index ? '<' : ''}
                                        </Link>
                                    ))}
                            {!videos && <p>No videos here.</p>}
                        </div>
                    </div>

                    <div className="mx-16 flex justify-between lg:hidden">
                        {prev != '' ? (
                            <Link href={`/dashboard/video/${encodeURIComponent(prev)}`}>prev</Link>
                        ) : (
                            <button disabled={true} className="text-gray-400">
                                prev
                            </button>
                        )}
                        {next != '' ? (
                            <Link href={`/dashboard/video/${encodeURIComponent(next)}`}>next</Link>
                        ) : (
                            <button disabled={true} className="text-gray-400">
                                next
                            </button>
                        )}
                    </div>

                    <div className="mx-16 mt-5 flex flex-col p-4 sm:mx-20">
                        <h2 className="text-base font-bold sm:text-xl">{video?.name}</h2>
                        <p className="mt-2 text-base italic sm:text-xl">{video?.description}</p>
                        {video && (
                            <Link
                                href={`/dashboard/course/${encodeURIComponent(video?.course_id)}`}
                                className="mt-2 w-20 cursor-pointer text-sm text-blue-500 underline"
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

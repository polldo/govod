import Layout from '@/components/layout'
import VideoJS from '@/components/videoplayer'
import Head from 'next/head'
import Link from 'next/link'
import { useEffect } from 'react'
import { useMemo } from 'react'
import { useRef } from 'react'
import { useRouter } from 'next/router'
import { fetcher } from '@/services/fetch'
import React from 'react'
import { useSession } from '@/session/context'
import useSWR from 'swr'
import Image from 'next/image'
import { ProgressBar } from '@/components/progressbar'
import { Course, Video, Progress } from '@/services/types'

type ProgressMap = {
    [videoId: string]: number
}

export default function CourseDetails() {
    const { isLoggedIn, isLoading } = useSession()

    const router = useRouter()
    const { id } = router.query

    const { data } = useSWR(id ? `/videos/${id}/full` : null)
    const video: Video = data?.video
    const videos: Video[] = data?.allVideos
    const course: Course = data?.course
    const url: string = data?.url

    const progress = useMemo(() => {
        let map: ProgressMap = {}
        if (data?.allProgress) {
            data.allProgress.forEach((progressItem: Progress) => {
                map[progressItem.videoId] = progressItem.progress
            })
        }
        return map
    }, [data])

    // Refs are needed to synchronise react with the video player.
    const progressRef = useRef<number>(0)
    const lastProgressRef = useRef<number>(0)
    const startRef = useRef<number>(0)

    useEffect(() => {
        // The starting time of the video will depend on the user progress.
        startRef.current = progress[video?.id] || 0
    }, [progress, video])

    // Send any NEW progress every 20 seconds.
    useEffect(() => {
        if (!video) {
            return
        }
        const interval = setInterval(() => {
            if (progressRef.current === lastProgressRef.current) {
                return
            }
            fetcher
                .fetch('/videos/' + video.id + '/progress', {
                    method: 'PUT',
                    body: JSON.stringify({ progress: progressRef.current }),
                })
                .then(() => {
                    lastProgressRef.current = progressRef.current
                })
                .catch()
        }, 20000)
        return () => {
            clearInterval(interval)
        }
    }, [video])

    const handlePlayerReady = (player: any) => {
        player.on('loadstart', () => {
            // Setting the poster to empty string seems to be the
            // only solution to keep it updated when the URL is changed.
            player.poster('')
            const tot: number = player.duration()
            // Start the video from the beginning if it was completed, otherwise videojs remains blocked.
            const adjust = startRef.current === 100 ? 0 : startRef.current
            const start = (tot * adjust) / 100
            player.currentTime(start)
            player.play()
        })
        player.on('timeupdate', () => {
            // Floor is better than ceiling when dealing with long videos (to avoid going too far).
            progressRef.current = Math.floor((player.currentTime() * 100) / player.duration())
        })
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

    if (isLoading) {
        return null
    }

    if (!isLoggedIn) {
        router.push('/login')
        return null
    }

    if (!video || !videos) {
        return null
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
                            {videos
                                .slice()
                                .sort((a, b) => a.index - b.index)
                                .map((vid) => (
                                    <Link
                                        key={vid.name}
                                        href={`/dashboard/video/${encodeURIComponent(vid.id)}`}
                                        className={`flex flex-row items-center p-2 py-2 text-sm ${
                                            vid.index % 2 ? 'bg-gray-200' : 'bg-gray-300'
                                        }`}
                                    >
                                        {vid.index === video.index ? '> ' : ''}
                                        <Image
                                            className="m-2 mr-4 w-10"
                                            alt=""
                                            src={vid.imageUrl}
                                            width={80}
                                            height={32}
                                        />
                                        <div className="flex w-full flex-col">
                                            {vid.name}
                                            {vid.index !== video.index ? (
                                                <div className="w-2/3">
                                                    <ProgressBar percent={progress[vid.id] || 0} />
                                                </div>
                                            ) : null}
                                        </div>
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

                    {video && (
                        <div className="mx-16 mt-5 flex flex-col p-4 sm:mx-20">
                            <h2 className="text-base font-bold sm:text-xl">{video.name}</h2>
                            <p className="mt-2 text-base italic sm:text-xl">{video.description}</p>
                            <Link
                                href={`/dashboard/course/${encodeURIComponent(video.courseId)}`}
                                className="mt-2 cursor-pointer text-sm text-blue-500 underline"
                            >
                                {course?.name}
                            </Link>
                        </div>
                    )}
                </div>
            </Layout>
        </>
    )
}

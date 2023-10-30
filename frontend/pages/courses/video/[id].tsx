import Layout from '@/components/layout'
import VideoJS from '@/components/videoplayer'
import Head from 'next/head'
import Link from 'next/link'
import { useRouter } from 'next/router'
import React from 'react'
import useSWR from 'swr'
import { Video, Course } from '@/services/types'

export default function CourseDetails() {
    const router = useRouter()
    const { id } = router.query

    const { data } = useSWR(id ? `/videos/${id}/free` : null)
    const video: Video = data?.video
    const course: Course = data?.course
    const url: string = data?.url

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

    if (!video || !url) {
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

                    <div className="mx-auto mt-5 flex w-full flex-col p-4">
                        <h2 className="text-base font-bold sm:text-xl">{video.name}</h2>
                        <p className="mt-2 text-base italic sm:text-xl">{video.description}</p>
                        <Link
                            href={`/courses/${encodeURIComponent(video.courseId)}`}
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

import Layout from '@/components/layout'
import Head from 'next/head'
import Image from 'next/image'
import Link from 'next/link'
import { useRouter } from 'next/router'
import useSWR from 'swr'
import { CourseCard } from '@/components/coursecard'
import { Course, Video } from '@/services/types'

export default function CourseDetails() {
    const router = useRouter()
    const { id } = router.query

    const { data: course } = useSWR<Course>(id ? `/courses/${id}` : null)
    const { data: videos } = useSWR<Video[]>(id ? `/courses/${id}/videos` : null)

    if (!course || !videos) {
        return null
    }

    return (
        <>
            <Head>
                <title>Course - {course?.name}</title>
            </Head>
            <Layout>
                <div className="flex w-full flex-col">
                    <CourseCard course={course}></CourseCard>

                    <div className="flex w-full flex-col">
                        <div className="flex flex-col items-center space-y-5 pt-6 pb-6">
                            {videos && videos.map((video) => <Card {...video} key={video.name} />)}
                        </div>
                    </div>
                </div>
            </Layout>
        </>
    )
}

function Card(props: Video) {
    return (
        <div className="flex w-full flex-col items-center">
            {props.free && (
                <div className="w-2/3 rounded-t-lg border bg-green-600 text-xs text-white md:max-w-3xl">
                    <p className="mx-5 font-semibold">Free</p>
                </div>
            )}

            <div className="flex w-2/3 flex-col items-center rounded-lg border border-gray-200 bg-white shadow md:max-w-3xl md:flex-row">
                <Image className="m-2 w-20 object-contain" alt="" src={props.imageUrl} width={80} height={32} />

                <div className="flex w-full flex-col items-center justify-between md:flex-row">
                    <div className="flex flex-col p-4 leading-normal">
                        <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900">{props.name}</h5>
                        <p className="mb-3 font-normal text-gray-700">{props.description}</p>
                    </div>

                    {props.free && (
                        <Link
                            href={`/courses/video/${encodeURIComponent(props.id)}`}
                            className="m-2 h-1/2 rounded bg-blue-700 p-4 font-semibold text-white hover:bg-blue-900"
                        >
                            Play
                        </Link>
                    )}
                </div>
            </div>
        </div>
    )
}

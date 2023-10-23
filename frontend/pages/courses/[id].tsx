import Layout from '@/components/layout'
import Head from 'next/head'
import Image from 'next/image'
import Link from 'next/link'
import { useRouter } from 'next/router'
import useSWR from 'swr'

type Course = {
    name: string
    description: string
    image_url: string
}

type Video = {
    id: string
    name: string
    description: string
    free: boolean
    image_url: string
}

export default function CourseDetails() {
    const router = useRouter()
    const { id } = router.query

    const { data: course } = useSWR<Course>(id ? `/courses/${id}` : null)
    const { data: videos } = useSWR<Video[]>(id ? `/courses/${id}/videos` : null)

    const free: Video[] = videos?.filter((video) => video.free) || []

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
                    <div className="my-10 mx-auto flex w-2/3 flex-col ">
                        <Image
                            className="my-5 mx-auto object-contain "
                            alt=""
                            src={course.image_url}
                            width={80}
                            height={32}
                        />
                        <div className="mx-auto flex w-2/3 flex-col">
                            <h5 className="mx-auto text-2xl font-bold text-gray-900">{course.name}</h5>
                            <p className="mx-auto">{course.description}</p>
                        </div>
                    </div>

                    {free && free.length > 0 && (
                        <div className="flex w-full flex-col">
                            <p className="mx-auto">Free Samples</p>
                            <div className="flex w-full flex-col items-center space-y-5 pt-6 pb-6">
                                {free.map((video) => (
                                    <Card {...video} key={video.name} />
                                ))}
                            </div>
                        </div>
                    )}

                    <div className="flex w-full flex-col">
                        <p className="mx-auto">All videos contained in this course</p>
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
        <div className="flex w-1/2 flex-col items-center rounded-lg border border-gray-200 bg-white shadow hover:bg-gray-100 md:max-w-xl md:flex-row">
            <Image
                className="w-full rounded-t-lg border border-red-800 object-contain md:w-20"
                alt=""
                src={props.image_url}
                width={80}
                height={32}
            />

            <div className="flex flex-col justify-between p-4 leading-normal">
                <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{props.name}</h5>
                <p className="mb-3 font-normal text-gray-700 dark:text-gray-400">{props.description}</p>
            </div>

            {props.free && (
                <Link
                    href={`/courses/video/${encodeURIComponent(props.id)}`}
                    className="rounded bg-blue-700 p-4 font-semibold text-white hover:bg-blue-900"
                >
                    Play
                </Link>
            )}
        </div>
    )
}

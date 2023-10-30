import Layout from '@/components/layout'
import Head from 'next/head'
import Image from 'next/image'
import { useSession } from '@/session/context'
import { CourseCard } from '@/components/coursecard'
import { useRouter } from 'next/router'
import Link from 'next/link'
import useSWR from 'swr'
import { ProgressBar } from '@/components/progressbar'
import { Course, Video, Progress } from '@/services/types'

type ProgressMap = {
    [videoId: string]: number
}

export default function DashboardCourse() {
    const { isLoggedIn, isLoading } = useSession()
    const router = useRouter()
    const { id } = router.query

    const { data: course } = useSWR<Course>(id ? `/courses/${id}` : null)
    const { data: videos } = useSWR<Video[]>(id ? `/courses/${id}/videos` : null)

    const { data: progressData } = useSWR<Progress[]>(id ? `/courses/${id}/progress` : null)
    let progress: ProgressMap = {}
    progressData?.forEach((p: Progress) => {
        progress[p.videoId] = p.progress
    })

    if (isLoading || !course || !videos) {
        return null
    }

    if (!isLoggedIn) {
        router.push('/login')
        return null
    }

    return (
        <>
            <Head>
                <title>Course - {course.name}</title>
            </Head>
            <Layout>
                <div className="flex w-full flex-col">
                    <CourseCard course={course}></CourseCard>
                    <div className="flex w-full flex-col items-center space-y-5 pt-6 pb-6">
                        {videos &&
                            videos.map((video) => (
                                <Card {...video} progress={progress[video.id] || 0} key={video.name} />
                            ))}
                    </div>
                </div>
            </Layout>
        </>
    )
}

type CardProps = Video & {
    progress: number
}

function Card(props: CardProps) {
    return (
        <Link
            href={`/dashboard/video/${props.id}`}
            className="flex w-2/3 flex-col items-center rounded-lg border border-gray-200 bg-white shadow hover:bg-gray-100 md:max-w-xl md:flex-row"
        >
            <div className="flex w-full flex-col items-center">
                <Image
                    className="m-2 w-20 rounded-t-lg object-contain"
                    alt=""
                    src={props.imageUrl}
                    width={80}
                    height={32}
                />

                <div className="flex flex-col justify-between p-4 leading-normal">
                    <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
                        {props.name}
                    </h5>
                    <p className="mb-3 font-normal text-gray-700 dark:text-gray-400">{props.description}</p>
                </div>

                <p>Progress: {props.progress}%</p>
                <ProgressBar percent={props.progress || 0} />
            </div>
        </Link>
    )
}

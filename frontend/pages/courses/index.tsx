import Layout from '@/components/layout'
import Head from 'next/head'
import Image from 'next/image'
import Link from 'next/link'
import { useEffect, useState } from 'react'

type Course = {
    id: string
    name: string
    description: string
    image: string
}

function Card(props: Course) {
    return (
        <Link
            href={`/courses/${encodeURIComponent(props.id)}`}
            className="flex w-1/2 flex-col items-center rounded-lg border border-gray-200 bg-white shadow hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-700 md:max-w-xl md:flex-row"
        >
            <Image
                className="w-full rounded-t-lg border border-red-800 object-contain md:w-20"
                alt=""
                src={props.image}
                width={80}
                height={32}
            />

            <div className="flex flex-col justify-between p-4 leading-normal">
                <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{props.name}</h5>
                <p className="mb-3 font-normal text-gray-700 dark:text-gray-400">{props.description}</p>
            </div>
        </Link>
    )
}

export default function Courses() {
    const [courses, setCourses] = useState<Course[]>()

    useEffect(() => {
        fetch('http://mylocal.com:8000/courses')
            .then((res) => res.json())
            .then((data) => setCourses(data))
    }, [])

    if (!courses) {
        return null
    }

    return (
        <>
            <Head>
                <title>Courses</title>
            </Head>
            <Layout>
                <div className="flex flex-col items-center space-y-5 pt-6 pb-6">
                    {courses.map((course) => (
                        <Card {...course} key={course.name} />
                    ))}
                </div>
            </Layout>
        </>
    )
}

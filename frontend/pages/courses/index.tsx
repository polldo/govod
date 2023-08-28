import Layout from '@/components/layout'
import Head from 'next/head'
import Image from 'next/image'
import Link from 'next/link'
import { useEffect, useState } from 'react'
import { useSession } from '@/session/context'
import { toast } from 'react-hot-toast'
import { useFetch } from '@/services/fetch'

type Course = {
    id: string
    name: string
    description: string
    image: string
}

type Cart = {
    items: CartItem[]
}

type CartItem = {
    course_id: string
}

type CardProps = Course & {
    isOwned: boolean
    isInCart: boolean
    isLoggedIn: boolean
}

function Card(props: CardProps) {
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
                {props.isOwned ? (
                    <button onClick={() => (window.location.href = `/courses/${props.id}`)}>Go to Course</button>
                ) : props.isInCart ? (
                    <button disabled>Already in Cart</button>
                ) : (
                    <button onClick={() => {}}>{props.isLoggedIn ? 'Add to Cart' : 'Login to Purchase'}</button>
                )}
            </div>
        </Link>
    )
}

export default function Courses() {
    const [courses, setCourses] = useState<Course[]>()
    const [cartCourses, setCartCourses] = useState<string[]>([])
    const [isLoadingCart, setIsLoadingCart] = useState<boolean>(true)
    const [ownedCourses, setOwnedCourses] = useState<string[]>([])
    const [isLoadingOwned, setIsLoadingOwned] = useState<boolean>(true)
    const { isLoggedIn, isLoading } = useSession()
    const fetch = useFetch()

    useEffect(() => {
        fetch('http://mylocal.com:8000/courses')
            .then((res) => res.json())
            .then((data) => setCourses(data))
    }, [fetch])

    useEffect(() => {
        if (!isLoggedIn) {
            return
        }

        fetch('http://mylocal.com:8000/cart')
            .then((res) => {
                if (!res.ok) {
                    throw new Error()
                }
                return res.json()
            })
            .then((data: Cart) => {
                const incart = data.items.map((item: CartItem) => item.course_id)
                setCartCourses(incart)
                setIsLoadingCart(false)
            })
            .catch(() => {
                toast.error('Something went wrong')
            })

        fetch('http://mylocal.com:8000/courses/owned')
            .then((res) => {
                if (!res.ok) {
                    throw new Error()
                }
                return res.json()
            })
            .then((data: Course[]) => {
                const owned = data.map((item: Course) => item.id)
                setOwnedCourses(owned)
                setIsLoadingOwned(false)
            })
            .catch(() => {
                toast.error('Something went wrong')
            })
    }, [fetch, isLoggedIn])

    if (isLoading) {
        return null
    }

    if (isLoggedIn && (isLoadingCart || isLoadingOwned)) {
        return null
    }

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
                        <Card
                            {...course}
                            isOwned={ownedCourses.includes(course.id)}
                            isInCart={cartCourses.includes(course.id)}
                            isLoggedIn={isLoggedIn}
                            key={course.name}
                        />
                    ))}
                </div>
            </Layout>
        </>
    )
}

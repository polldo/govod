import Layout from '@/components/layout'
import Head from 'next/head'
import Image from 'next/image'
import Link from 'next/link'
import { useSession } from '@/session/context'
import { toast } from 'react-hot-toast'
import { fetcher } from '@/services/fetch'
import useSWR from 'swr'
import { useRouter } from 'next/router'

type Course = {
    id: string
    name: string
    description: string
    image_url: string
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
    onAddToCart: () => void
}

function Card(props: CardProps) {
    const router = useRouter()

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        if (props.isOwned) {
            router.push(`/dashboard/course/${props.id}`)
            return
        }
        if (props.isInCart) {
            // Do nothing.
            return
        }
        if (!props.isLoggedIn) {
            router.push(`/login`)
            return
        }

        props.onAddToCart()
    }

    const linkURL: string = props.isOwned
        ? `/dashboard/course/${encodeURIComponent(props.id)}`
        : `/courses/${encodeURIComponent(props.id)}`

    return (
        <Link
            href={linkURL}
            className="flex w-1/2 flex-col items-center rounded-lg border border-gray-200 bg-white shadow hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-700 md:max-w-xl md:flex-row"
        >
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
                {props.isOwned ? (
                    <button onClick={handleSubmit} className="w-full rounded bg-blue-500 p-2 font-semibold text-white">
                        Go to Course
                    </button>
                ) : props.isInCart ? (
                    <button disabled className="w-full rounded bg-gray-500 p-2 font-semibold text-white">
                        In Cart
                    </button>
                ) : (
                    <button onClick={handleSubmit} className="w-full rounded bg-green-500 p-2 font-semibold text-white">
                        {props.isLoggedIn ? 'Add to Cart' : 'Login to buy'}
                    </button>
                )}
            </div>
        </Link>
    )
}

export default function Courses() {
    const router = useRouter()
    const { isLoggedIn, isLoading } = useSession()

    const { data: courses } = useSWR<Course[]>('/courses')

    const { data: cartData, mutate: cartMutate } = useSWR<Cart>(isLoggedIn ? '/cart' : null)
    const cartCourses = cartData ? cartData.items.map((item: CartItem) => item.course_id) : []

    const { data: ownedData } = useSWR<Course[]>(isLoggedIn ? '/courses/owned' : null)
    const ownedCourses = ownedData ? ownedData.map((item: Course) => item.id) : []

    const handleAddToCart = (courseID: string) => {
        fetcher
            .fetch('/cart/items', {
                method: 'PUT',
                body: JSON.stringify({ course_id: courseID }),
            })
            .then(() => {
                cartMutate()
                router.push('/cart')
            })
            .catch(() => {
                toast.error('Something went wrong')
            })
    }

    if (!courses || isLoading) {
        return null
    }

    const coursePriority = (course: Course) => {
        if (ownedCourses.includes(course.id)) {
            return 2
        }
        if (cartCourses.includes(course.id)) {
            return 1
        }
        return 0
    }

    const sortedCourses = [...courses].sort((a, b) => {
        return coursePriority(a) - coursePriority(b)
    })

    return (
        <>
            <Head>
                <title>Courses</title>
            </Head>
            <Layout>
                <div className="flex flex-col items-center space-y-5 pt-6 pb-6">
                    {sortedCourses.map((course) => (
                        <Card
                            {...course}
                            isOwned={ownedCourses.includes(course.id)}
                            isInCart={cartCourses.includes(course.id)}
                            isLoggedIn={isLoggedIn}
                            onAddToCart={() => {
                                handleAddToCart(course.id)
                            }}
                            key={course.name}
                        />
                    ))}
                </div>
            </Layout>
        </>
    )
}

import Layout from '@/components/layout'
import Head from 'next/head'
import { useState } from 'react'
import { useEffect } from 'react'
import { useRouter } from 'next/router'
import { useSession } from '@/session/context'
import { useFetch } from '@/services/fetch'
import Image from 'next/image'
import { toast } from 'react-hot-toast'

type CartType = {
    items: Item[]
}

type Item = {
    course_id: string
}

type Course = {
    id: string
    name: string
    price: number
}

type CoursesMap = {
    [courseId: string]: Course
}

type Product = {
    id: string
    name: string
    image: string
    price: number
}

type ProductCardProps = {
    product: Product
    onDelete: (x: string) => void
}

function ProductCard(props: ProductCardProps) {
    return (
        <div className="mb-2 flex flex-col items-center justify-between rounded bg-white p-4 shadow sm:flex-row">
            <button
                onClick={() => props.onDelete(props.product.id)}
                className="rounded bg-red-500 px-4 py-2 text-white"
            >
                x
            </button>
            <Image src={props.product.image} alt={props.product.name} className="h-20 w-20" />
            <div className="ml-4">{props.product.name}</div>
            <div className="font-bold">${props.product.price}</div>
        </div>
    )
}

export default function Cart() {
    const [cart, setCart] = useState<CartType>()
    const [courses, setCourses] = useState<CoursesMap>({})
    const { isLoggedIn, isLoading } = useSession()
    const router = useRouter()
    const fetch = useFetch()

    useEffect(() => {
        if (!isLoggedIn) {
            return
        }

        fetch('http://mylocal.com:8000/cart')
            .then((res) => res.json())
            .then((data: CartType) => {
                setCart(data)

                const courseFetches = data.items.map((item) => {
                    return fetch(`http://mylocal.com:8000/courses/${item.course_id}`)
                        .then((res) => res.json())
                        .then((course: Course) => {
                            setCourses((prevCourses) => ({ ...prevCourses, [course.id]: course }))
                        })
                })

                Promise.all(courseFetches).catch(() => {
                    toast.error('Something went wrong')
                })
            })
            .catch(() => {
                toast.error('Something went wrong')
            })
    }, [fetch, isLoggedIn])

    if (isLoading) {
        return null
    }

    if (!isLoggedIn) {
        router.push('/login')
        return null
    }

    return (
        <>
            <Head>
                <title>Cart</title>
            </Head>
            <Layout>
                <div className="flex w-full flex-col items-center justify-center p-4">
                    <h1 className="mb-4 text-3xl font-bold">Shopping Cart</h1>
                    <div className="w-full bg-gray-100 p-4 sm:w-1/2">
                        {cart &&
                            cart.items.map((item) => {
                                const course = courses[item.course_id]
                                if (!course) return null

                                return (
                                    <ProductCard
                                        key={item.course_id}
                                        product={{
                                            id: course.id,
                                            name: course.name,
                                            image: '',
                                            price: course.price,
                                        }}
                                        onDelete={(id) => {
                                            // Handle deletion logic here
                                            console.log('Deleting item with Course ID:', id)
                                        }}
                                    />
                                )
                            })}
                        <button className="mt-4 w-full rounded bg-green-500 p-4 text-white">Checkout</button>
                    </div>
                </div>
            </Layout>
        </>
    )
}

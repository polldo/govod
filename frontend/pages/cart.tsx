import Layout from '@/components/layout'
import Head from 'next/head'
import { useState } from 'react'
import { useEffect } from 'react'
import { useRouter } from 'next/router'
import { useSession } from '@/session/context'
import { useFetch } from '@/services/fetch'
import Image from 'next/image'
import { toast } from 'react-hot-toast'
import { PayPalButtons, usePayPalScriptReducer } from '@paypal/react-paypal-js'

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
    image_url: string
}

type CoursesMap = {
    [courseId: string]: Course
}

type Product = {
    id: string
    name: string
    image_url: string
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
            <Image
                className="h-20 w-20"
                alt={props.product.name}
                src={props.product.image_url}
                width={80}
                height={32}
            />
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
    const [{ isPending, isResolved }] = usePayPalScriptReducer()

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
            .then((data: CartType) => {
                setCart(data)

                const courseFetches = data.items.map((item) => {
                    return fetch(`http://mylocal.com:8000/courses/${item.course_id}`)
                        .then((res) => {
                            if (!res.ok) {
                                throw new Error()
                            }
                            return res.json()
                        })
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

    const handleDeleteItem = async (id: string) => {
        try {
            const res = await fetch(`http://mylocal.com:8000/cart/items/${id}`, {
                method: 'DELETE',
            })

            if (!res.ok) {
                throw new Error()
            }

            if (cart) {
                setCart({
                    ...cart,
                    items: cart.items.filter((item) => item.course_id !== id),
                })
            }
        } catch (err) {
            toast.error('Something went wrong')
        }
    }

    const handleStripeCheckout = async (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        e.preventDefault()
        try {
            const res = await fetch(`http://mylocal.com:8000/orders/stripe`, { method: 'POST' })
            if (!res.ok) {
                throw new Error()
            }
            const data = await res.json()
            window.location.href = data
        } catch (err) {
            toast.error('Something went wrong')
        }
    }

    const handlePaypalCheckout = async () => {
        try {
            const res = await fetch(`http://mylocal.com:8000/orders/paypal`, { method: 'POST' })
            if (!res.ok) {
                throw new Error()
            }
            const data = await res.json()
            return data.id
        } catch (err) {
            toast.error('Something went wrong')
        }
    }

    const handlePaypalCapture = async (capture: { orderID: string }) => {
        try {
            const res = await fetch(`http://mylocal.com:8000/orders/paypal/${capture.orderID}/capture`, {
                method: 'POST',
            })
            if (!res.ok) {
                throw new Error()
            }
            window.location.href = `/dashboard`
        } catch (err) {
            toast.error('Something went wrong')
        }
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
                                            image_url: course.image_url,
                                            price: course.price,
                                        }}
                                        onDelete={handleDeleteItem}
                                    />
                                )
                            })}

                        <div className="flex flex-col items-center gap-4 md:flex-row md:justify-center">
                            <button
                                onClick={handleStripeCheckout}
                                className={`rounded p-2 text-white md:w-1/3 ${
                                    cart && cart.items.length > 0
                                        ? 'bg-green-500 hover:bg-green-600'
                                        : 'bg-green-200 hover:bg-green-300'
                                }`}
                                disabled={cart?.items.length === 0}
                            >
                                Checkout
                            </button>

                            {isPending || !isResolved ? (
                                <button disabled={true} className="rounded bg-gray-500 p-2 text-white md:w-1/3">
                                    PayPal ...loading...
                                </button>
                            ) : (
                                <PayPalButtons
                                    disabled={cart?.items.length == 0}
                                    className="rounded p-2 md:w-1/3"
                                    createOrder={handlePaypalCheckout}
                                    onApprove={handlePaypalCapture}
                                    style={{ layout: 'vertical' }}
                                    fundingSource="paypal"
                                />
                            )}
                        </div>
                    </div>
                </div>
            </Layout>
        </>
    )
}

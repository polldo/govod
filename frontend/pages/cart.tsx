import Layout from '@/components/layout'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useSession } from '@/session/context'
import { fetcher } from '@/services/fetch'
import Image from 'next/image'
import { toast } from 'react-hot-toast'
import { PayPalButtons, usePayPalScriptReducer } from '@paypal/react-paypal-js'
import useSWR from 'swr'

type Cart = {
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

type CartCourseProps = {
    course: string
    onDelete: (x: string) => void
}

function CartCourse(props: CartCourseProps) {
    const { data: course } = useSWR<Course>(`/courses/${props.course}`)

    if (!course) {
        return null
    }

    return (
        <div className="mb-2 flex flex-col items-center justify-between rounded bg-white p-4 shadow sm:flex-row">
            <button onClick={() => props.onDelete(course.id)} className="rounded bg-red-500 px-4 py-2 text-white">
                x
            </button>
            <Image className="h-20 w-20" alt={course.name} src={course.image_url} width={80} height={32} />
            <div className="ml-4">{course.name}</div>
            <div className="font-bold">${course.price}</div>
        </div>
    )
}

export default function CartPage() {
    const { isLoggedIn, isLoading } = useSession()
    const router = useRouter()
    const [{ isPending, isResolved }] = usePayPalScriptReducer()

    const { data: cart, mutate } = useSWR<Cart>(isLoggedIn ? '/cart' : null)

    if (isLoading) {
        return null
    }

    if (!isLoggedIn) {
        router.push('/login')
        return null
    }

    const handleDeleteItem = async (id: string) => {
        try {
            await fetcher.fetch(`/cart/items/${id}`, {
                method: 'DELETE',
            })

            mutate()
        } catch (err) {
            toast.error('Something went wrong')
        }
    }

    const handleStripeCheckout = async (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        e.preventDefault()
        try {
            const res = await fetcher.fetch(`/orders/stripe`, { method: 'POST' })
            const data = await res.json()
            window.location.href = data
        } catch (err) {
            toast.error('Something went wrong')
        }
    }

    const handlePaypalCheckout = async () => {
        try {
            const res = await fetcher.fetch(`/orders/paypal`, { method: 'POST' })
            const data = await res.json()
            return data.id
        } catch (err) {
            toast.error('Something went wrong')
        }
    }

    const handlePaypalCapture = async (capture: { orderID: string }) => {
        try {
            await fetcher.fetch(`/orders/paypal/${capture.orderID}/capture`, {
                method: 'POST',
            })
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
                                return (
                                    <CartCourse
                                        key={item.course_id}
                                        course={item.course_id}
                                        onDelete={handleDeleteItem}
                                    />
                                )
                            })}

                        {cart?.items.length === 0 && <p className="p-5 text-center">Nothing in the cart..</p>}

                        <div className="flex flex-col items-center gap-4 md:flex-row md:justify-center">
                            <button
                                onClick={handleStripeCheckout}
                                className={`rounded p-2 text-white md:w-1/3 ${
                                    cart && cart.items.length > 0
                                        ? 'bg-green-700 hover:bg-green-900'
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

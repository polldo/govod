import Layout from '@/components/layout'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useSession } from '@/session/context'
import { fetcher } from '@/services/fetch'
import Image from 'next/image'
import { toast } from 'react-hot-toast'
import { PayPalButtons, usePayPalScriptReducer } from '@paypal/react-paypal-js'
import useSWR from 'swr'
import { Cart, Course } from '@/services/types'

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
        <div className="mb-2 flex flex-col rounded bg-white p-4 shadow sm:flex-row sm:items-center sm:justify-between">
            <button className="w-1/5" onClick={() => props.onDelete(course.id)}>
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="30" height="30" viewBox="0 0 48 48">
                    <path d="M 24 4 C 20.491685 4 17.570396 6.6214322 17.080078 10 L 6.5 10 A 1.50015 1.50015 0 1 0 6.5 13 L 8.6367188 13 L 11.15625 39.029297 C 11.43025 41.862297 13.785813 44 16.632812 44 L 31.367188 44 C 34.214187 44 36.56875 41.862297 36.84375 39.029297 L 39.363281 13 L 41.5 13 A 1.50015 1.50015 0 1 0 41.5 10 L 30.919922 10 C 30.429604 6.6214322 27.508315 4 24 4 z M 24 7 C 25.879156 7 27.420767 8.2681608 27.861328 10 L 20.138672 10 C 20.579233 8.2681608 22.120844 7 24 7 z M 19.5 18 C 20.328 18 21 18.671 21 19.5 L 21 34.5 C 21 35.329 20.328 36 19.5 36 C 18.672 36 18 35.329 18 34.5 L 18 19.5 C 18 18.671 18.672 18 19.5 18 z M 28.5 18 C 29.328 18 30 18.671 30 19.5 L 30 34.5 C 30 35.329 29.328 36 28.5 36 C 27.672 36 27 35.329 27 34.5 L 27 19.5 C 27 18.671 27.672 18 28.5 18 z"></path>
                </svg>
            </button>
            <Image className="mx-auto h-20 w-20" alt={course.name} src={course.imageUrl} width={80} height={32} />
            <div className="mx-auto ">{course.name}</div>
            <div className="mx-auto font-bold">${course.price}</div>
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
                                        key={item.courseId}
                                        course={item.courseId}
                                        onDelete={handleDeleteItem}
                                    />
                                )
                            })}

                        {cart?.items.length === 0 && <p className="p-5 text-center">Nothing in the cart..</p>}

                        <div className="flex w-full flex-col items-center gap-1 p-2">
                            <button
                                onClick={handleStripeCheckout}
                                className={`flex w-full flex-row justify-center rounded p-2 text-white md:w-1/2 ${
                                    cart && cart.items.length > 0
                                        ? 'bg-green-700 hover:bg-green-900'
                                        : 'bg-green-200 hover:bg-green-300'
                                }`}
                                disabled={cart?.items.length === 0}
                            >
                                <svg
                                    width="24"
                                    height="24"
                                    xmlns="http://www.w3.org/2000/svg"
                                    fill-rule="evenodd"
                                    clip-rule="evenodd"
                                >
                                    <path d="M22 3c.53 0 1.039.211 1.414.586s.586.884.586 1.414v14c0 .53-.211 1.039-.586 1.414s-.884.586-1.414.586h-20c-.53 0-1.039-.211-1.414-.586s-.586-.884-.586-1.414v-14c0-.53.211-1.039.586-1.414s.884-.586 1.414-.586h20zm1 8h-22v8c0 .552.448 1 1 1h20c.552 0 1-.448 1-1v-8zm-15 5v1h-5v-1h5zm13-2v1h-3v-1h3zm-10 0v1h-8v-1h8zm-10-6v2h22v-2h-22zm22-1v-2c0-.552-.448-1-1-1h-20c-.552 0-1 .448-1 1v2h22z" />
                                </svg>
                                <p className="mx-3 font-bold">Card</p>
                            </button>

                            {isPending || !isResolved ? (
                                <button disabled={true} className="w-full rounded bg-gray-500 text-white md:w-1/2">
                                    PayPal ...loading...
                                </button>
                            ) : (
                                <PayPalButtons
                                    disabled={cart?.items.length == 0}
                                    className="w-full rounded md:w-1/2 "
                                    createOrder={handlePaypalCheckout}
                                    onApprove={handlePaypalCapture}
                                    style={{ layout: 'horizontal' }}
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

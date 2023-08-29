import type { PropsWithChildren } from 'react'
import Link from 'next/link'
import { useSession } from '@/session/context'
import Logout from '@/components/logout'
import Cart from '@/components/cart'

export default function Layout(props: PropsWithChildren) {
    return (
        <>
            <div className="h-screen">
                <Navbar />
                <main className="overflow-none flex justify-center">{props.children}</main>
            </div>
        </>
    )
}

function Navbar() {
    const { isLoggedIn, isLoading } = useSession()
    return (
        <nav className="bg-gray-900 py-4">
            <div className="container mx-auto">
                <div className="flex items-center justify-between">
                    <Link className="text-xl font-bold text-white" href="/">
                        Govod
                    </Link>
                    <div className="flex flex-row items-center">
                        <Link className="rounded px-3 py-2 text-gray-400 hover:text-white" href="/courses">
                            Courses
                        </Link>

                        {isLoggedIn && (
                            <Link className="rounded px-3 py-2 text-gray-400 hover:text-white" href="/dashboard">
                                Dashboard
                            </Link>
                        )}

                        {isLoggedIn && (
                            <div className="px-3 py-2 ">
                                <Cart></Cart>
                            </div>
                        )}

                        {isLoggedIn && (
                            <div className="px-3 py-2 ">
                                <Logout></Logout>
                            </div>
                        )}

                        {!isLoggedIn && (
                            <Link className="w-full rounded bg-blue-800 p-2 font-semibold text-white" href="/login">
                                Login
                            </Link>
                        )}
                    </div>
                </div>
            </div>
        </nav>
    )
}

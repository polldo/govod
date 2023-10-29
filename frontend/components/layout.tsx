import type { PropsWithChildren } from 'react'
import Link from 'next/link'
import { useState } from 'react'
import { useEffect } from 'react'
import { useRef } from 'react'
import { useSession } from '@/session/context'
import Logout from '@/components/logout'
import Cart from '@/components/cart'
import { Transition } from '@headlessui/react'

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

function Links() {
    const { isLoggedIn } = useSession()
    return (
        <>
            <Link className="rounded text-gray-300 hover:text-white" href="/courses">
                Courses
            </Link>

            {isLoggedIn && (
                <Link className="rounded text-gray-300 hover:text-white" href="/dashboard">
                    Dashboard
                </Link>
            )}

            {isLoggedIn && <Cart></Cart>}

            {isLoggedIn && (
                <div className="w-1/2">
                    <Logout></Logout>
                </div>
            )}

            {!isLoggedIn && (
                <Link
                    className="w-1/2 rounded bg-green-700 p-2 text-center font-semibold text-white hover:bg-green-900"
                    href="/login"
                >
                    Login
                </Link>
            )}
        </>
    )
}

function Navbar() {
    const [isOpen, setIsOpen] = useState(false)
    const navbarRef = useRef<HTMLDivElement | null>(null)

    useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (navbarRef.current && !navbarRef.current.contains(event.target as Node)) {
                setIsOpen(false)
            }
        }

        document.addEventListener('mousedown', handleClickOutside)
        return () => {
            document.removeEventListener('mousedown', handleClickOutside)
        }
    }, [])

    return (
        <nav className="bg-blue-900 p-4" ref={navbarRef}>
            <div className="container mx-auto">
                <div className="flex items-center justify-between">
                    <Link href="/" className="text-xl font-bold text-white">
                        Govod
                    </Link>
                    <div className="hidden items-center space-x-4 md:flex">
                        <Links />
                    </div>
                    <div className="md:hidden">
                        <button onClick={() => setIsOpen(!isOpen)}>
                            <svg
                                className="h-6 w-6 text-white"
                                fill="none"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path d="M4 6h16M4 12h16m-7 6h7"></path>
                            </svg>
                        </button>
                    </div>
                </div>
                <Transition
                    show={isOpen}
                    enter="transition ease-out duration-120 transform"
                    enterFrom="opacity-0 scale-95"
                    enterTo="opacity-100 scale-100"
                    leave="transition ease-in duration-120 transform"
                    leaveFrom="opacity-100 scale-100"
                    leaveTo="opacity-0 scale-95"
                >
                    <div className="flex flex-col items-center space-y-3 p-4 md:hidden">
                        <Links></Links>
                    </div>
                </Transition>
            </div>
        </nav>
    )
}

import type { PropsWithChildren } from 'react'

export default function Layout(props: PropsWithChildren) {
    return (
        <>
            <Navbar />
            <main className="overflow-none flex h-screen justify-center">
                <div className="flex h-full w-full flex-col md:max-w-2xl">{props.children}</div>
            </main>
        </>
    )
}

function Navbar() {
    return (
        <nav className="bg-gray-900 py-4">
            <div className="container mx-auto">
                <div className="flex justify-between">
                    <div className="text-xl font-bold text-white">My Website</div>
                    <div>
                        <a className="rounded px-3 py-2 text-gray-400 hover:text-white" href="#">
                            Home
                        </a>
                        <a className="rounded px-3 py-2 text-gray-400 hover:text-white" href="#">
                            About
                        </a>
                        <a className="rounded px-3 py-2 text-gray-400 hover:text-white" href="#">
                            Contact
                        </a>
                    </div>
                </div>
            </div>
        </nav>
    )
}

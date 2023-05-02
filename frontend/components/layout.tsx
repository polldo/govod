import type { PropsWithChildren } from 'react'

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
    return (
        <nav className="bg-gray-900 py-4">
            <div className="container mx-auto">
                <div className="flex justify-between">
                    <div className="text-xl font-bold text-white">Govod</div>
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

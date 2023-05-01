import type { PropsWithChildren } from 'react'

export default function Layout(props: PropsWithChildren) {
    return (
        <>
            <main className="overflow-none flex h-screen justify-center">
                <div className="flex h-full w-full flex-col border-x border-slate-400 md:max-w-2xl">
                    {props.children}
                </div>
            </main>
        </>
    )
}

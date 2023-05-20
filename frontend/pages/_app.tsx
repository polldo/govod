import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { SessionProvider } from '@/session/context'

export default function App({ Component, pageProps }: AppProps) {
    return (
        <SessionProvider {...pageProps}>
            <Component {...pageProps} />
        </SessionProvider>
    )
}

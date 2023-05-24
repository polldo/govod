import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { SessionProvider } from '@/session/context'
import { Toaster } from 'react-hot-toast'

export default function App({ Component, pageProps }: AppProps) {
    return (
        <SessionProvider {...pageProps}>
            <Toaster position="bottom-center" />
            <Component {...pageProps} />
        </SessionProvider>
    )
}

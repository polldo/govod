import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { SessionProvider } from '@/session/context'
import { Toaster } from 'react-hot-toast'
import { PayPalScriptProvider } from '@paypal/react-paypal-js'
import { fetcher } from '@/services/fetch'
import { useEffect } from 'react'
import { useSession } from '@/session/context'

// Set the unauthenticated callback of 'fetch.Fetcher'.
function FetchInterceptor() {
    const { updateSession } = useSession()

    useEffect(() => {
        fetcher.setUnauthenticated(() => {
            updateSession()
        })
    }, [updateSession])

    return null
}

export default function App({ Component, pageProps }: AppProps) {
    return (
        <>
            <SessionProvider {...pageProps}>
                <FetchInterceptor></FetchInterceptor>
                <Toaster position="bottom-center" />
                <PayPalScriptProvider
                    options={{
                        clientId: process.env.NEXT_PUBLIC_PAYPAL_CLIENT_ID || '',
                        currency: 'USD',
                    }}
                >
                    <Component {...pageProps} />
                </PayPalScriptProvider>
            </SessionProvider>
        </>
    )
}

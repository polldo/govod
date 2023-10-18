import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { SessionProvider } from '@/session/context'
import { Toaster } from 'react-hot-toast'
import { PayPalScriptProvider } from '@paypal/react-paypal-js'
import { fetcher } from '@/services/fetch'
import { useEffect } from 'react'
import { useSession } from '@/session/context'
import { SWRConfig } from 'swr'
import { toast } from 'react-hot-toast'

// Configure the custom fetcher.
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
                <Toaster position="bottom-center" />
                <FetchInterceptor></FetchInterceptor>
                <SWRConfig
                    value={{
                        fetcher: (url: string) => fetcher.fetch(url).then((res) => res.json()),
                        onError: () => {
                            toast.error('Something went wrong')
                        },
                    }}
                >
                    <PayPalScriptProvider
                        options={{
                            clientId: process.env.NEXT_PUBLIC_PAYPAL_CLIENT_ID || '',
                            currency: 'USD',
                        }}
                    >
                        <Component {...pageProps} />
                    </PayPalScriptProvider>
                </SWRConfig>
            </SessionProvider>
        </>
    )
}

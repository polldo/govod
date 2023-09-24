import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import { SessionProvider } from '@/session/context'
import { Toaster } from 'react-hot-toast'
import { PayPalScriptProvider } from '@paypal/react-paypal-js'

export default function App({ Component, pageProps }: AppProps) {
    return (
        <SessionProvider {...pageProps}>
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
    )
}

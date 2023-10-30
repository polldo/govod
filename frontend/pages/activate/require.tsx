import Layout from '@/components/layout'
import { Spinner } from '@/components/spinner'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import { useCallback } from 'react'
import { useState } from 'react'
import { fetcher } from '@/services/fetch'
import { TokenRequest } from '@/services/types'

export default function Require() {
    const router = useRouter()
    const { email } = router.query
    const [reloading, setReloading] = useState<boolean>(false)

    const handleEmail = useCallback(() => {
        if (!router.isReady) {
            return
        }
        const { email } = router.query
        const body: TokenRequest = {
            email: typeof email === 'string' ? email : '',
            scope: 'activation',
        }

        fetcher
            .fetch('/tokens', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body),
            })
            .catch((err) => {
                console.log(err)
            })
        setReloading(true)
    }, [router.query, router.isReady])

    useEffect(() => {
        if (!reloading) return
        const timer = setTimeout(() => setReloading(false), 1500)
        return () => clearTimeout(timer)
    }, [reloading])

    useEffect(() => {
        handleEmail()
    }, [handleEmail])

    return (
        <>
            <Head>
                <title>Activation required</title>
            </Head>
            <Layout>
                <div className="flex items-center justify-center py-32">
                    <div className="flex flex-col rounded-lg border border-gray-300 bg-gray-100 p-6 text-center">
                        <h1 className="mb-4 text-2xl font-bold">Account Activation Required</h1>
                        <p className="p-2 text-lg">
                            An account has been created with the email
                            <strong className="ml-1 text-blue-600">{email}</strong>. <br></br>Please check your email to
                            activate your account.
                        </p>

                        {reloading ? (
                            <div className="mx-auto flex ">
                                <Spinner />
                            </div>
                        ) : (
                            <button
                                onClick={handleEmail}
                                className="mx-auto w-1/2 rounded bg-blue-700 p-2 font-semibold text-white hover:bg-blue-900"
                            >
                                Send email again
                            </button>
                        )}
                    </div>
                </div>
            </Layout>
        </>
    )
}

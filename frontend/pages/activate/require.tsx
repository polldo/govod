import Layout from '@/components/layout'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import { useCallback } from 'react'
import { fetcher } from '@/services/fetch'

type ActivateBody = {
    Email: string
    Scope: string
}

export default function Require() {
    const router = useRouter()
    const { email } = router.query

    const handleEmail = useCallback(() => {
        if (!router.isReady) {
            return
        }
        const { email } = router.query
        const body: ActivateBody = {
            Email: typeof email === 'string' ? email : '',
            Scope: 'activation',
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
    }, [router.query, router.isReady])

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
                    <div className="rounded-lg border border-gray-300 bg-gray-100 p-6 text-center">
                        <h1 className="mb-4 text-2xl font-bold">Account Activation Required</h1>
                        <p className="text-lg">
                            An account has been created with the email{' '}
                            <strong className="text-blue-600">{email}</strong>. <br></br>Please check your email to
                            activate your account.
                        </p>
                        <button
                            onClick={handleEmail}
                            className="w-full rounded bg-blue-500 p-2 font-semibold text-white"
                        >
                            Send email again
                        </button>
                    </div>
                </div>
            </Layout>
        </>
    )
}

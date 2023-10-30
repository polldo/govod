import Layout from '@/components/layout'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import { useState } from 'react'
import { useCallback } from 'react'
import { useSession } from '@/session/context'
import { fetcher, ResponseError } from '@/services/fetch'
import { ActivationToken } from '@/services/types'

export default function Confirm() {
    const [activated, setActivated] = useState<boolean>(false)
    const [error, setError] = useState<string>('')
    const { updateSession } = useSession()
    const router = useRouter()

    const handleSubmit = useCallback(async () => {
        if (!router.isReady) {
            return
        }
        const { token } = router.query
        const body: ActivationToken = {
            token: typeof token === 'string' ? token : '',
        }

        try {
            await fetcher.fetch('/tokens/activate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body),
            })
            updateSession()
            setActivated(true)
            setTimeout(() => {
                router.push('/login')
            }, 1500)
        } catch (err) {
            setError('Something went wrong')
            if (err instanceof ResponseError) {
                if (err.status === 422) {
                    setError('Invalid token')
                }
            }
        }
    }, [router, updateSession])

    useEffect(() => {
        handleSubmit()
    }, [handleSubmit])

    return (
        <>
            <Head>
                <title>Activation confirm</title>
            </Head>
            <Layout>
                <div className="flex items-center justify-center py-32">
                    <div className="rounded-lg border border-gray-300 bg-gray-100 p-6 text-center">
                        <h1 className="mb-4 text-2xl font-bold">Account Activation</h1>
                        <p className="text-lg">
                            Your account is being activated. <br></br>Please wait...
                        </p>
                        <br></br>
                        {error && <p className="mb-4 text-sm text-red-500">{error}</p>}
                        {activated && <p className="mb-4 text-sm text-blue-500">activated... Redirecting..</p>}
                    </div>
                </div>
            </Layout>
        </>
    )
}

import Layout from '@/components/layout'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import { useState } from 'react'
import { useCallback } from 'react'

type Token = {
    Token: string
}

export default function Confirm() {
    const [activated, setActivated] = useState<boolean>(false)
    const [error, setError] = useState<string>('')
    const router = useRouter()

    const handleSubmit = useCallback(async () => {
        if (!router.isReady) {
            return
        }
        const { token } = router.query
        const body: Token = {
            Token: typeof token === 'string' ? token : '',
        }

        try {
            const res = await fetch('http://mylocal.com:8000/tokens/activate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body),
            })
            if (res.status === 422) {
                throw new Error('Invalid token')
            }
            if (!res.ok) {
                throw new Error('Something went wrong')
            }

            setActivated(true)
            setTimeout(() => {
                router.push('/login')
            }, 1500)
        } catch (err) {
            if (err instanceof Error) {
                setError(err.message)
            } else {
                setError('Something went wrong')
            }
        }
    }, [router])

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

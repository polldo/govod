import Layout from '@/components/layout'
import Head from 'next/head'
import { useState } from 'react'
import { fetcher, ResponseError } from '@/services/fetch'
import { TokenRequest } from '@/services/types'

export default function Require() {
    const [email, setEmail] = useState<string>('')
    const [error, setError] = useState<string>('')
    const [isSent, setIsSent] = useState<boolean>(false)

    const handleEmailChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmail(event.target.value)
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setError('')

        const body: TokenRequest = {
            email: email,
            scope: 'recovery',
        }

        try {
            await fetcher.fetch('/tokens', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body),
            })
            setIsSent(true)
        } catch (err) {
            setError('Something went wrong')
            if (err instanceof ResponseError) {
                if (err.status === 422) {
                    setError(err.message)
                }
                if (err.status === 429) {
                    setError('Retry in few seconds')
                }
            }
        }
    }

    return (
        <>
            <Head>
                <title>Reset password</title>
            </Head>

            <Layout>
                <div className="my-12 flex items-center justify-center bg-gray-100">
                    <form onSubmit={handleSubmit} className="w-full rounded bg-white p-6 shadow-md sm:w-96">
                        <h1 className="mb-4 text-2xl font-semibold">Reset password</h1>
                        {error && <p className="mb-4 text-sm text-red-500">{error}</p>}

                        <input
                            type="email"
                            value={email}
                            onChange={handleEmailChange}
                            className="mb-4 block w-full rounded bg-gray-100 p-2"
                            placeholder="Email"
                            required
                            disabled={isSent}
                        />

                        {isSent ? (
                            <p className="rounded-lg bg-blue-300 p-4">
                                An email has been sent to your email address. Please follow the instructions.
                            </p>
                        ) : (
                            <button
                                type="submit"
                                className="w-full rounded bg-blue-700 p-2 font-semibold text-white hover:bg-blue-900"
                            >
                                Send
                            </button>
                        )}
                    </form>
                </div>
            </Layout>
        </>
    )
}

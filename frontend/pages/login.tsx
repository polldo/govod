import Layout from '@/components/layout'
import Head from 'next/head'
import { useState } from 'react'
import { Buffer } from 'buffer'
import { useRouter } from 'next/router'
import { useSession } from '@/session/context'
import Link from 'next/link'
import { fetcher, ResponseError } from '@/services/fetch'

export default function Login() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const { isLoggedIn, isLoading, updateSession } = useSession()
    const router = useRouter()

    if (isLoading) {
        return null
    }

    if (isLoggedIn) {
        router.push('/dashboard')
        return null
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setError('')

        try {
            await fetcher.fetch('/auth/login', {
                method: 'POST',
                headers: {
                    Authorization: `Basic ${Buffer.from(`${email}:${password}`).toString('base64')}`,
                },
            })

            updateSession()
        } catch (err) {
            setError('Something went wrong')
            if (err instanceof ResponseError) {
                if (err.status === 401) {
                    setError('Invalid credentials')
                }
                if (err.status === 423) {
                    router.push({ pathname: '/activate/require', query: { email } })
                    return
                }
            }
        }
    }

    const handleGoogleLogin = async (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        e.preventDefault()

        try {
            const res = await fetcher.fetch('/auth/oauth-login/google')
            const data = await res.json()

            // No need to call 'login', because after the oauth login the user will be
            // redirected and the whole app will be reloaded.
            window.location.href = data
        } catch (err) {
            setError('Something went wrong')
        }
    }

    const handleEmailChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmail(event.target.value)
    }

    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPassword(event.target.value)
    }

    return (
        <>
            <Head>
                <title>Login</title>
            </Head>
            <Layout>
                <div className="my-12 flex items-center justify-center bg-gray-100">
                    <form onSubmit={handleSubmit} className="w-full rounded bg-white p-6 shadow-md sm:w-96">
                        <h1 className="mb-4 text-2xl font-semibold">Login</h1>
                        {error && <p className="mb-4 text-sm text-red-500">{error}</p>}

                        <input
                            type="email"
                            value={email}
                            onChange={handleEmailChange}
                            className="mb-4 block w-full rounded bg-gray-100 p-2 hover:bg-gray-200"
                            placeholder="Email"
                            required
                        />

                        <input
                            type="password"
                            value={password}
                            onChange={handlePasswordChange}
                            className="mb-4 block w-full rounded bg-gray-100 p-2 hover:bg-gray-200"
                            placeholder="Password"
                            required
                        />

                        <button
                            type="submit"
                            className="w-full rounded bg-blue-700 p-2 font-semibold text-white hover:bg-blue-900"
                        >
                            Login
                        </button>

                        <button
                            onClick={handleGoogleLogin}
                            className="w-full rounded bg-red-500 p-2 font-semibold text-white hover:bg-red-700"
                        >
                            Login with Google
                        </button>

                        <Link href="/password/reset" className="mb-2 text-sm text-blue-500 hover:underline">
                            Forgot password?
                        </Link>

                        <div className="mt-2 flex flex-col">
                            <p className="mx-auto"> -- or --</p>
                            <Link
                                href={`/signup`}
                                className="mt-4 w-full rounded bg-gray-500 p-2 text-center font-semibold text-white hover:bg-gray-700"
                            >
                                <p>Signup</p>
                            </Link>
                        </div>
                    </form>
                </div>
            </Layout>
        </>
    )
}

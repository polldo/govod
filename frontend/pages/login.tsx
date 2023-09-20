import Layout from '@/components/layout'
import Head from 'next/head'
import { useState } from 'react'
import { Buffer } from 'buffer'
import { useRouter } from 'next/router'
import { useSession } from '@/session/context'
import { useFetch } from '@/services/fetch'
import Link from 'next/link'

export default function Login() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const { isLoggedIn, isLoading, updateSession } = useSession()
    const router = useRouter()
    const fetch = useFetch()

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
            const res = await fetch('http://mylocal.com:8000/auth/login', {
                method: 'POST',
                headers: {
                    Authorization: `Basic ${Buffer.from(`${email}:${password}`).toString('base64')}`,
                },
            })

            if (res.status === 401) {
                throw new Error('Invalid credentials')
            }
            if (res.status === 423) {
                router.push({ pathname: '/activate/require', query: { email } })
                return
            }
            if (!res.ok) {
                throw new Error('Something went wrong')
            }
            updateSession()
        } catch (err) {
            if (err instanceof Error) {
                setError(err.message)
            } else {
                setError('Something went wrong')
            }
        }
    }

    const handleGoogleLogin = async (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        e.preventDefault()

        try {
            const res = await fetch('http://mylocal.com:8000/auth/oauth-login/google', { method: 'GET' })
            const data = await res.json()

            // No need to call 'login', because after the oauth login the user will be
            // redirected and the whole app will be reloaded.
            window.location.href = data
        } catch (err) {
            if (err instanceof Error) {
                setError(err.message)
            } else {
                setError('Something went wrong')
            }
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
                            className="mb-4 block w-full rounded bg-gray-100 p-2"
                            placeholder="Email"
                            required
                        />

                        <input
                            type="password"
                            value={password}
                            onChange={handlePasswordChange}
                            className="mb-4 block w-full rounded bg-gray-100 p-2"
                            placeholder="Password"
                            required
                        />

                        <button type="submit" className="w-full rounded bg-blue-500 p-2 font-semibold text-white">
                            Login
                        </button>

                        <button
                            onClick={handleGoogleLogin}
                            className="w-full rounded bg-red-500 p-2 font-semibold text-white"
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
                                className="mt-4 w-full rounded bg-gray-500 p-2 text-center font-semibold text-white"
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

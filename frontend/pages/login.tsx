import Layout from '@/components/layout'
import Head from 'next/head'
import { useState } from 'react'
import { Buffer } from 'buffer'
import { useRouter } from 'next/router'
import { useSession } from '@/session/context'
import { useFetch } from '@/services/fetch'

export default function Login() {
    const router = useRouter()
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')

    const { isLoggedIn, isLoading, login, logout } = useSession()
    const fetch = useFetch()

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
                throw new Error('Activate your account to login')
            }
            if (!res.ok) {
                throw new Error('Something went wrong')
            }

            const data = await res.json()
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
        // console.log(email, password)

        try {
            const res = await fetch('http://mylocal.com:8000/auth/oauth-login/google', {
                method: 'GET',
                // TODO: understand if this credentials include is normal for frontend that use session for authentication.
                credentials: 'include',
            })

            // if (res.status === 401) {
            //     throw new Error('Invalid credentials')
            // }
            // if (res.status === 423) {
            //     throw new Error('Activate your account to login')
            // }
            // if (!res.ok) {
            //     throw new Error('Something went wrong')
            // }

            const data = await res.json()
            console.log(data)
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
                    </form>
                </div>
            </Layout>
        </>
    )
}

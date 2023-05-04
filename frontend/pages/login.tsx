import Layout from '@/components/layout'
import Head from 'next/head'
import { useState } from 'react'
import { Buffer } from 'buffer'

export default function Login() {
    const [email, setEmail] = useState('ss@ss')
    const [password, setPassword] = useState('ss')
    const [error, setError] = useState('')

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        setError('')

        fetch('http://127.0.0.1:8080/auth/login', {
            method: 'POST',
            headers: {
                Authorization: `Basic ${Buffer.from(`${email}:${password}`).toString('base64')}`,
            },
        })
            .then((res) => res.json())
            .then((data) => console.log(data))
            .catch((err) => {
                setError(err.message)
            })
    }

    const handleEmailChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setEmail(event.target.value)
    }

    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPassword(event.target.value)
    }

    const handleGoogleLogin = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        e.preventDefault()
        console.log(email, password)
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

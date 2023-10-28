import Layout from '@/components/layout'
import Head from 'next/head'
import { useState } from 'react'
import { useRouter } from 'next/router'
import { useSession } from '@/session/context'
import { fetcher, ResponseError } from '@/services/fetch'

export default function Signup() {
    const [name, setName] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState<string>('')
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
            const res = await fetcher.fetch('/auth/signup', {
                method: 'POST',
                body: JSON.stringify({
                    email: email,
                    name: name,
                    password: password,
                }),
            })

            const data = await res.json()
            if (data.active) {
                updateSession()
                return
            }
            router.push({
                pathname: '/activate/require',
                query: { email },
            })
        } catch (err) {
            setError('Something went wrong')
            if (err instanceof ResponseError) {
                if (err.status === 409) {
                    setError('Email already exists')
                }
                if (err.status === 422) {
                    setError(err.message)
                }
            }
        }
    }

    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName(event.target.value)
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
                <title>Signup</title>
            </Head>
            <Layout>
                <div className="my-12 flex items-center justify-center bg-gray-100">
                    <form onSubmit={handleSubmit} className="w-full rounded bg-white p-6 shadow-md sm:w-96">
                        <h1 className="mb-4 text-2xl font-semibold">Signup</h1>
                        {error && <p className="mb-4 text-sm text-red-500">{error}</p>}
                        <input
                            type="text"
                            value={name}
                            onChange={handleNameChange}
                            className="mb-4 block w-full rounded bg-gray-100 p-2"
                            placeholder="Name"
                            required
                        />
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
                        <button
                            type="submit"
                            className="w-full rounded bg-blue-700 p-2 font-semibold text-white hover:bg-blue-900"
                        >
                            Signup
                        </button>
                    </form>
                </div>
            </Layout>
        </>
    )
}

import Layout from '@/components/layout'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useState } from 'react'

type ConfirmBody = {
    Token: string
    Password: string
    Password_Confirm: string
}

export default function Confirm() {
    const [password, setPassword] = useState<string>('')
    const [passwordConfirm, setPasswordConfirm] = useState<string>('')
    const [error, setError] = useState<string>('')
    const [success, setSuccess] = useState<boolean>(false)
    const router = useRouter()

    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPassword(event.target.value)
    }

    const handlePasswordConfirmChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPasswordConfirm(event.target.value)
    }

    const handleClick = (e: React.FormEvent) => {
        e.preventDefault()
        router.push('/login')
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setError('')

        const { token } = router.query
        const body: ConfirmBody = {
            Token: typeof token === 'string' ? token : '',
            Password: password,
            Password_Confirm: passwordConfirm,
        }

        try {
            const res = await fetch('http://mylocal.com:8000/tokens/recover', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body),
            })
            if (res.status === 422) {
                const data = await res.json()
                throw new Error(data.error)
            }
            if (!res.ok) {
                throw new Error('Something went wrong')
            }

            setSuccess(true)
        } catch (err) {
            if (err instanceof Error) {
                setError(err.message)
            } else {
                setError('Something went wrong')
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
                            type="password"
                            value={password}
                            onChange={handlePasswordChange}
                            className="mb-4 block w-full rounded bg-gray-100 p-2"
                            placeholder="Password"
                            required
                        />

                        <input
                            type="password"
                            value={passwordConfirm}
                            onChange={handlePasswordConfirmChange}
                            className="mb-4 block w-full rounded bg-gray-100 p-2"
                            placeholder="Password confirmation"
                            required
                        />

                        {success ? (
                            <>
                                <p className="rounded-lg bg-blue-300 p-4">Your password has been correctly reset.</p>
                                <button
                                    onClick={handleClick}
                                    className="w-full rounded bg-blue-500 p-2 font-semibold text-white"
                                >
                                    Go to Login
                                </button>
                            </>
                        ) : (
                            <button type="submit" className="w-full rounded bg-blue-500 p-2 font-semibold text-white">
                                Reset password
                            </button>
                        )}
                    </form>
                </div>
            </Layout>
        </>
    )
}

import Layout from '@/components/layout'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useState } from 'react'
import { fetcher, ResponseError } from '@/services/fetch'
import { PasswordToken } from '@/services/types'

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
        const body: PasswordToken = {
            token: typeof token === 'string' ? token : '',
            password: password,
            passwordConfirm: passwordConfirm,
        }

        try {
            await fetcher.fetch('/tokens/recover', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body),
            })
            setSuccess(true)
        } catch (err) {
            setError('Something went wrong')
            if (err instanceof ResponseError) {
                if (err.status === 422) {
                    setError(err.message)
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
                                    className="w-full rounded bg-blue-700 p-2 font-semibold text-white hover:bg-blue-900"
                                >
                                    Go to Login
                                </button>
                            </>
                        ) : (
                            <button
                                type="submit"
                                className="w-full rounded bg-blue-700 p-2 font-semibold text-white hover:bg-blue-900"
                            >
                                Reset password
                            </button>
                        )}
                    </form>
                </div>
            </Layout>
        </>
    )
}

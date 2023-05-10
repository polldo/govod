import Layout from '@/components/layout'
import Head from 'next/head'
import { useRouter } from 'next/router'

export default function Activate() {
    const router = useRouter()
    const { email } = router.query

    return (
        <>
            <Head>
                <title>Activate</title>
            </Head>
            <Layout>
                <div className="flex items-center justify-center py-32">
                    <div className="rounded-lg border border-gray-300 bg-gray-100 p-6 text-center">
                        <h1 className="mb-4 text-2xl font-bold">Account Activation Required</h1>
                        <p className="text-lg">
                            An account has been created with the email{' '}
                            <strong className="text-blue-600">{email}</strong>. <br></br>Please check your email to
                            activate your account.
                        </p>
                    </div>
                </div>
            </Layout>
        </>
    )
}

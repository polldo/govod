import Layout from '@/components/layout'
import Logout from '@/components/logout'
import Head from 'next/head'
import { useSession } from '@/session/context'
import { useRouter } from 'next/router'

export default function Dashboard() {
    const router = useRouter()
    const { isLoggedIn, isLoading } = useSession()

    if (isLoading) {
        return null
    }

    if (!isLoggedIn) {
        router.push('login')
        return null
    }

    return (
        <>
            <Head>
                <title>Dashboard</title>
            </Head>
            <Layout>
                <div>
                    <p>Hello, this is your dashboard!</p>
                    <Logout></Logout>
                </div>
            </Layout>
        </>
    )
}

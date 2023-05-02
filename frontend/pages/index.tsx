import Layout from '@/components/layout'
import Head from 'next/head'

export default function Home() {
    return (
        <>
            <Layout>
                <Head>
                    <title>Govod</title>
                </Head>
                <Hero />
            </Layout>
        </>
    )
}

function Hero() {
    return (
        <div className="fullHeight scroll-smooth">
            <div className="bg-cover bg-center py-20" style={{ backgroundImage: 'url("/hero-bg.jpg")' }}>
                <div className="container mx-auto text-center">
                    <h1 className="mb-4 text-4xl font-bold text-white">Welcome to My Website</h1>
                    <p className="mb-8 text-white">
                        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus non turpis sed nisi convallis
                        feugiat.
                    </p>
                    <button className="rounded bg-blue-500 py-2 px-4 font-bold text-white hover:bg-blue-700">
                        Get Started
                    </button>
                </div>
            </div>
        </div>
    )
}

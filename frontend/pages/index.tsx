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
        <>
            <div className="bg-gray-100 py-16">
                <div className="mx-auto flex max-w-7xl flex-col px-4 sm:px-6 lg:flex-row lg:px-8">
                    <div className="flex items-center text-center lg:w-1/2">
                        <div className="mx-auto w-1/2">
                            <h1 className="mb-6 text-6xl font-extrabold tracking-tight text-blue-800 ">
                                &iexcl; govod !
                            </h1>
                            <p className="mb-6 text-lg font-extrabold tracking-tight text-gray-600 ">
                                WebApp to sell video on demand
                            </p>
                            <p className="mb-8 font-serif text-lg text-gray-600">
                                Start selling courses now. Customize the website with your own style. Make content and
                                enrich other people&#39;s life.
                            </p>
                        </div>
                    </div>
                    <div className="mx-auto lg:w-1/3">
                        <iframe
                            className="h-full w-full rounded-lg object-cover shadow-lg"
                            src="https://www.youtube.com/embed/446E-r0rXHI"
                            title="YouTube video player"
                            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                            allowFullScreen
                        ></iframe>
                    </div>
                </div>
            </div>
        </>
    )
}

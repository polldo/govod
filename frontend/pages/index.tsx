import Head from 'next/head'
import Image from 'next/image'
import styles from '@/styles/Home.module.css'
import Footer from '@/components/Footer/Footer'

export default function Home() {
    return (
        <>
            <Head>
                <title>govod</title>
                <meta name="description" content="govod - webapp for video on demand" />
                <meta name="viewport" content="width=device-width, initial-scale=1" />
                <link rel="icon" href="/favicon.ico" />
            </Head>

            <main className={styles.main}>
                <div className={styles.description}></div>

                <div className={styles.center}>
                    <Image
                        className={styles.logo}
                        src="/next.svg"
                        alt="Next.js Logo"
                        width={180}
                        height={37}
                        priority
                    />
                </div>

                <div className={styles.grid}></div>
            </main>

            <Footer />
        </>
    )
}

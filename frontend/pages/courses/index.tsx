import Image from 'next/image'

type Course = {
    title: string
    description: string
    image: string
}

const data = [
    { title: 'ok', description: 'full description', image: '/francesca1.png' },
    { title: 'second', description: 'second description', image: '/francesca1.png' },
    { title: 'third', description: 'third description', image: '/francesca1.png' },
]

function Card(props: Course) {
    return (
        <a
            href="#"
            className="flex w-1/2 flex-col items-center rounded-lg border border-gray-200 bg-white shadow hover:bg-gray-100 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-700 md:max-w-xl md:flex-row"
        >
            <Image
                className="rounded-t-lg border border-red-800"
                // className="h-96 w-full rounded-t-lg object-cover md:h-auto md:w-48 md:rounded-none md:rounded-l-lg"
                alt=""
                src={props.image}
                width={128}
                height={128}
            />
            <div className="flex flex-col justify-between p-4 leading-normal">
                <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{props.title}</h5>
                <p className="mb-3 font-normal text-gray-700 dark:text-gray-400">{props.description}</p>
            </div>
        </a>
    )
}

export default function Courses() {
    return (
        <div className="flex flex-col items-center space-y-5 pt-6 pb-6">
            {data.map((fullPost) => (
                // <PostView {...fullPost} key={fullPost.post.id} />
                <Card {...fullPost} key={fullPost.title} />
            ))}
        </div>
    )
}

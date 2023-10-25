import Image from 'next/image'

type Course = {
    name: string
    description: string
    image_url: string
}

type CourseCardProps = {
    course: Course
}

export function CourseCard(props: CourseCardProps) {
    return (
        <div className="my-10 mx-auto flex w-2/3 flex-col ">
            <Image
                className="my-5 mx-auto object-contain "
                alt=""
                src={props.course.image_url}
                width={80}
                height={32}
            />
            <div className="mx-auto flex w-2/3 flex-col">
                <h5 className="mx-auto text-xl font-bold text-gray-900 md:text-2xl">{props.course.name}</h5>
                <p className="mx-auto">{props.course.description}</p>
            </div>
        </div>
    )
}

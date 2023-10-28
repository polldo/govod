export function ProgressBar(props: { percent: number }) {
    return (
        <>
            <div className="mt-1 h-2.5 w-full rounded-full bg-gray-100">
                <div className="h-2.5 rounded-full bg-blue-600" style={{ width: props.percent.toString() + '%' }}></div>
            </div>
        </>
    )
}

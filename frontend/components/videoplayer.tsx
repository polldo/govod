import React, { useEffect, useRef } from 'react'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'
import Player from 'video.js/dist/types/player'
import 'videojs-youtube'

type VideoJSProps = {
    options: any
    onReady?: (player: Player) => void
}

export default function VideoJS(props: VideoJSProps) {
    const videoRef = useRef<HTMLDivElement | null>(null)
    const playerRef = useRef<Player | null>(null)

    useEffect(() => {
        if (!playerRef.current && videoRef.current) {
            const videoElement = document.createElement('video-js')
            videoRef.current.appendChild(videoElement)

            const player = videojs(videoElement, props.options, () => {
                videojs.log('player is ready')
                props.onReady && props.onReady(player)
            })

            playerRef.current = player
        } else if (playerRef.current) {
            playerRef.current.autoplay(props.options.autoplay)
            playerRef.current.src(props.options.sources)
        }
    }, [videoRef, props])

    useEffect(() => {
        return () => {
            const player = playerRef.current
            if (player && !player.isDisposed()) {
                player.dispose()
                playerRef.current = null
            }
        }
    }, [playerRef])

    return (
        <div data-vjs-player>
            <div ref={videoRef} />
        </div>
    )
}

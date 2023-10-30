export type Course = {
    id: string
    name: string
    description: string
    image_url: string
    price: number
}

export type Video = {
    id: string
    course_id: string
    index: number
    name: string
    description: string
    free: boolean
    image_url: string
}

export type Cart = {
    items: CartItem[]
}

export type CartItem = {
    course_id: string
}

export type Progress = {
    video_id: string
    progress: number
}

export type ActivationToken = {
    Token: string
}

export type TokenRequest = {
    Email: string
    Scope: string
}

export type PasswordToken = {
    Token: string
    Password: string
    Password_Confirm: string
}

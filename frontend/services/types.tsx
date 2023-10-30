export type Course = {
    id: string
    name: string
    description: string
    imageUrl: string
    price: number
}

export type Video = {
    id: string
    courseId: string
    index: number
    name: string
    description: string
    free: boolean
    imageUrl: string
}

export type Cart = {
    items: CartItem[]
}

export type CartItem = {
    courseId: string
}

export type Progress = {
    videoId: string
    progress: number
}

export type ActivationToken = {
    token: string
}

export type TokenRequest = {
    email: string
    scope: string
}

export type PasswordToken = {
    token: string
    password: string
    passwordConfirm: string
}

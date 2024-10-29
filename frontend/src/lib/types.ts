export type User = {
    id: number;
    bearerToken: string; // auth token
    expiryTime: number; // ns when user auth expires
}

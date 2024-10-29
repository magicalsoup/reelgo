export type User = {
    id: number;
    bearerToken: string; // auth token
    expiry_time: number; // ns when user auth expires
}

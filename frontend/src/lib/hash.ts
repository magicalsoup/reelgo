

export default function sha256Hash(str: string): string {
    const crypto = require('crypto')
    return crypto.createHash('sha256').update(str).digest("hex")
}


export function padString(n: number) {
    return n < 10 ? '0' + n.toFixed(0) : n.toFixed(0)
}
export function normalizeSiteUrl(value: unknown): string {
  try {
    const url = new URL(String(value))
    return ['https:', 'http:'].includes(url.protocol) && !url.username && !url.password ? url.origin : ''
  }
  catch {
    return ''
  }
}

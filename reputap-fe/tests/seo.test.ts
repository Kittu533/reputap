import { describe, expect, test } from 'bun:test'
import { normalizeSiteUrl } from '../shared/utils/site-url'

describe('SEO site URL', () => {
  test('normalizes the public origin for canonical and sitemap URLs', () => {
    expect(normalizeSiteUrl('https://reputap.example/path?q=1#product')).toBe('https://reputap.example')
    expect(normalizeSiteUrl('https://reputap.example/')).toBe('https://reputap.example')
  })
  test('omits invalid origins and credentials', () => {
    for (const value of ['', undefined, 'not-a-url', 'javascript:alert(1)', 'https://user:password@example.com']) {
      expect(normalizeSiteUrl(value)).toBe('')
    }
  })
})

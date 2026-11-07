import { normalizeSiteUrl } from '#shared/utils/site-url'

export default defineEventHandler((event) => {
  const siteUrl = normalizeSiteUrl(useRuntimeConfig(event).public.siteUrl)
  if (!siteUrl) {
    throw createError({ statusCode: 503, statusMessage: 'Site URL is not configured' })
  }
  setHeader(event, 'Content-Type', 'application/xml; charset=utf-8')
  return `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>${siteUrl}/</loc></url>
</urlset>`
})

import { spawn } from 'node:child_process'
import assert from 'node:assert/strict'

const origin = 'https://reputap-alpha.vercel.app'
const preview = 'http://127.0.0.1:3197'
const server = spawn(process.execPath, ['.output/server/index.mjs'], {
  env: { ...process.env, PORT: '3197', HOST: '127.0.0.1', NUXT_PUBLIC_SITE_URL: origin },
  stdio: ['ignore', 'pipe', 'pipe'],
})
let log = ''
server.stderr.on('data', chunk => { log += chunk })
try {
  let response
  for (let attempt = 0; attempt < 40; attempt++) {
    try {
      response = await fetch(preview)
      break
    }
    catch {
      await new Promise(resolve => setTimeout(resolve, 150))
    }
  }
  assert.ok(response, log || 'Preview server unavailable')
  assert.equal(response.status, 200)
  const html = await response.text()
  assert.ok(html.includes('lang="id"'))
  assert.ok(html.includes(`rel="canonical" href="${origin}/"`))
  assert.ok(html.includes('index, follow, max-image-preview:large'))
  assert.ok(html.includes(`${origin}/images/reputap-display.png`))
  const schema = JSON.parse(html.match(/<script[^>]*type="application\/ld\+json"[^>]*>([\s\S]*?)<\/script>/)[1])
  const items = schema['@graph'].find(item => item['@type'] === 'CollectionPage').mainEntity.itemListElement
  assert.deepEqual(items.map(item => item.item.offers.price), [79900, 50900])
  assert.ok(items.every(item => item.item.offers.availability === 'https://schema.org/PreOrder'))
  const login = await (await fetch(`${preview}/login`)).text()
  assert.ok(login.includes('noindex, nofollow'))
  const sitemap = await (await fetch(`${preview}/sitemap.xml`)).text()
  assert.ok(sitemap.includes(`<loc>${origin}/</loc>`))
  assert.ok(!sitemap.includes('/login'))
  const robots = await (await fetch(`${preview}/robots.txt`)).text()
  assert.ok(robots.includes(`Sitemap: ${origin}/sitemap.xml`))
  console.log('PASS: SSR metadata, canonical, social image, product prices/preorder, admin noindex, sitemap, robots.txt')
}
finally {
  server.kill()
}

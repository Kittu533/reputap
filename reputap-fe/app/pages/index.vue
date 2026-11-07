<script setup lang="ts">
import { normalizeSiteUrl } from '#shared/utils/site-url'
definePageMeta({ layout: false })
const config = useRuntimeConfig()
const number = String(config.public.whatsappNumber).trim()
const canOrder = /^\d{8,15}$/.test(number)
const products = [
  { id: 'display-google-review', name: 'Display Google Review NFC + QR', price: 79900, material: 'Akrilik transparan 2 mm · 10 × 15 cm · NFC + QR', image: '/images/reputap-display.png', tag: 'Display meja · L-shape' },
  { id: 'kartu-google-review', name: 'Kartu Google Review NFC + QR', price: 50900, material: 'Akrilik 2–3 mm · 10 × 12 cm · NFC + QR', image: '/images/reputap-card.png', tag: 'Rounded · Desain minimalis' },
]
const siteUrl = normalizeSiteUrl(config.public.siteUrl)
const title = 'Display & Kartu Google Review NFC + QR | REPUTAP'
const description = 'Jual display akrilik Google Review Rp79.900 dan kartu NFC + QR Rp50.900. Pilihan desain custom untuk bisnis Anda. Pre-order langsung via WhatsApp REPUTAP.'
const imageUrl = siteUrl ? `${siteUrl}/images/reputap-display.png` : undefined
useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description,
  ogType: 'website',
  ogLocale: 'id_ID',
  ogSiteName: 'REPUTAP',
  ogUrl: siteUrl ? `${siteUrl}/` : undefined,
  ogImage: imageUrl,
  ogImageWidth: 1280,
  ogImageHeight: 853,
  ogImageAlt: 'Display akrilik Google Review REPUTAP dengan NFC dan QR',
  twitterCard: 'summary_large_image',
  twitterTitle: title,
  twitterDescription: description,
  twitterImage: imageUrl,
})
useHead({
  link: siteUrl ? [{ rel: 'canonical', href: `${siteUrl}/` }] : [],
  script: siteUrl ? [{
    type: 'application/ld+json',
    innerHTML: JSON.stringify({
      '@context': 'https://schema.org',
      '@graph': [
        { '@type': 'Organization', '@id': `${siteUrl}/#organization`, name: 'REPUTAP', url: `${siteUrl}/`, contactPoint: { '@type': 'ContactPoint', telephone: `+${number}`, contactType: 'sales', availableLanguage: 'Indonesian' } },
        { '@type': 'WebSite', '@id': `${siteUrl}/#website`, name: 'REPUTAP', url: `${siteUrl}/`, inLanguage: 'id-ID', publisher: { '@id': `${siteUrl}/#organization` } },
        { '@type': 'CollectionPage', '@id': `${siteUrl}/#webpage`, name: title, description, url: `${siteUrl}/`, inLanguage: 'id-ID', isPartOf: { '@id': `${siteUrl}/#website` }, mainEntity: {
          '@type': 'ItemList',
          itemListElement: products.map((product, index) => ({
            '@type': 'ListItem', position: index + 1,
            item: { '@type': 'Product', '@id': `${siteUrl}/#${product.id}`, name: product.name, description: product.material, image: `${siteUrl}${product.image}`, url: `${siteUrl}/#${product.id}`, brand: { '@type': 'Brand', name: 'REPUTAP' }, offers: { '@type': 'Offer', price: product.price, priceCurrency: 'IDR', availability: 'https://schema.org/PreOrder', url: `${siteUrl}/#${product.id}`, seller: { '@id': `${siteUrl}/#organization` } } },
          })),
        } },
      ],
    }).replace(/</g, '\\u003c'),
  }] : [],
})
const rupiah = (price: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
const steps = [
  { title: 'Pilih yang cocok.', text: 'Temukan bentuk dan material yang pas untuk bisnis Anda.' },
  { title: 'Ngobrol via WhatsApp.', text: 'Konfirmasi jumlah, desain, harga, dan estimasi pengerjaan bersama admin.' },
  { title: 'Kami siapkan pesanan.', text: 'Produksi dimulai setelah detail dan pembayaran disepakati. Lalu, siap dikirim.' },
]
function orderLink(name: string) {
  return `https://wa.me/${number}?text=${encodeURIComponent(`Halo REPUTAP! Saya ingin pre-order ${name}. Boleh info desain, estimasi produksi, dan total harganya?\nJumlah: \nNama bisnis: `)}`
}
</script>

<template>
  <div class="shop">
    <a class="skip" href="#main">Langsung ke konten</a>
    <div class="announcement">DIBUAT UNTUK BISNIS ANDA <span>✳</span> PRE-ORDER COLLECTION</div>
    <header class="site-header">
      <div class="nav wrap">
      <a href="#" class="brand" aria-label="REPUTAP beranda"><span class="brand-mark">R</span>REPUTAP</a>
      <nav aria-label="Navigasi utama"><a href="#produk">Produk</a><a href="#cara-pesan">Cara pesan</a><NuxtLink to="/login" class="admin">Login admin</NuxtLink></nav>
      </div>
    </header>
    <main id="main">
      <section class="hero wrap">
        <div class="hero-copy">
          <p class="eyebrow">SMART REVIEW FOR YOUR BUSINESS</p>
          <h1>Display & kartu<br>Google Review.<br><em>Tap NFC. Scan QR.</em></h1>
          <p class="intro">Permudah pelanggan membagikan pengalaman di Google Maps. Cukup scan QR atau tap NFC, dengan display yang menyatu dengan identitas bisnis Anda.</p>
          <div class="hero-actions"><a href="#produk" class="button">Temukan produk Anda</a><a href="#cara-pesan" class="text-link">Cara pesan <UIcon name="i-lucide-arrow-down" /></a></div>
          <div class="hero-notes"><span><UIcon name="i-lucide-nfc" /> Tap dengan NFC</span><span><UIcon name="i-lucide-scan-line" /> Scan dengan QR</span></div>
        </div>
        <div class="hero-art">
          <img src="/images/reputap-display.png" alt="REPUTAP Display Google Review: stand akrilik NFC dan QR, pilihan desain, ukuran dan material" width="1280" height="853" fetchpriority="high">
          <span class="visual-label">01 / DISPLAY GOOGLE REVIEW</span>
          <span class="hero-badge"><UIcon name="i-lucide-nfc" /> NFC + Dynamic QR</span>
        </div>
      </section>
      <div class="strip"><div class="wrap"><span>Untuk bisnis lokal</span><span>✳</span><span>Desain sesuai brand</span><span>✳</span><span>Pesan via WhatsApp</span></div></div>
      <section id="produk" class="catalog wrap">
        <div class="heading"><div><p class="eyebrow">THE COLLECTION / 01</p><h2>Pilih display atau kartu<br>Google Review Anda.</h2></div><p>Display untuk meja kasir atau kartu ringkas untuk berinteraksi dengan pelanggan. Dua bentuk, sama mudahnya.</p></div>
        <p class="note">Pilih desain Standar REPUTAP, Custom Brand, atau Custom Premium. Detail desain dan estimasi pre-order dikonfirmasi melalui WhatsApp.</p>
        <div class="products">
          <article v-for="product in products" :key="product.name" :id="product.id" class="product-card">
            <a :href="product.image" target="_blank" rel="noopener noreferrer" class="product-art" :aria-label="`Lihat detail ${product.name} (tab baru)`">
              <img :src="product.image" :alt="`${product.name}: spesifikasi, varian desain, dan tampak produk`" width="1280" height="853" loading="lazy">
            </a>
            <div class="product-info"><p class="product-category">{{ product.tag }}</p><h3>{{ product.name }}</h3><p class="material">{{ product.material }}</p><p class="price">{{ rupiah(product.price) }} <span>/ unit</span></p><p class="product-status"><span class="status-dot" /> Open pre-order</p>
              <a v-if="canOrder" :href="orderLink(product.name)" target="_blank" rel="noopener noreferrer" class="order">Pre-order via WhatsApp</a>
              <button v-else class="order" disabled>WhatsApp segera tersedia <UIcon name="i-lucide-message-circle" /></button>
            </div>
          </article>
        </div>
        <p v-if="!canOrder" class="availability" role="status">Pemesanan belum dibuka. Kontak WhatsApp toko sedang disiapkan.</p>
      </section>
      <section id="cara-pesan" class="how"><div class="wrap"><div class="heading"><div><p class="eyebrow">NO RIBET / 02</p><h2>Dari pilihan Anda,<br>ke meja bisnis Anda.</h2></div><p>Pilih produk, tentukan desain, lalu diskusikan pesanan langsung bersama kami.</p></div><div class="steps"><article v-for="(step, index) in steps" :key="step.title"><span>0{{ index + 1 }}</span><h3>{{ step.title }}</h3><p>{{ step.text }}</p></article></div></div></section>
      <section class="faq wrap"><div class="faq-heading"><p class="eyebrow">SEBELUM PESAN / 03</p><h2>Yang sering<br>ditanyakan.</h2><p>Kenali produknya sebelum menentukan pilihan Anda.</p></div><div class="faq-items"><details><summary>Bagaimana cara menggunakan produknya?</summary><p>Pelanggan menempelkan ponsel yang mendukung NFC ke produk. Untuk produk dengan QR, pelanggan juga bisa memindai kode melalui kamera ponsel.</p></details><details><summary>Bisa pakai desain brand sendiri?</summary><p>Diskusikan logo, warna, dan kebutuhan desain dengan admin melalui WhatsApp. Pilihan desain dan biaya akhir dikonfirmasi sebelum produksi.</p></details><details><summary>Berapa lama proses pre-order?</summary><p>Estimasi bergantung pada produk, jumlah, dan kesiapan desain. Admin akan mengonfirmasi jadwal produksi dan pengiriman sebelum pesanan disepakati.</p></details></div></section>
      <section class="closing wrap"><p class="eyebrow">LET’S MAKE A CONNECTION.</p><h2>Bisnis Anda.<br>Sentuhan berikutnya.</h2><a href="#produk" class="button">Pilih produk & mulai PO</a><span aria-hidden="true">✳</span></section>
    </main>
    <footer class="wrap"><a href="#" class="brand"><span class="brand-mark">R</span>REPUTAP</a><p>Smart Review for Your Business.</p><span>© {{ new Date().getFullYear() }} REPUTAP</span></footer>
  </div>
</template>

<style scoped>
.shop {
  --ink: #101c32;
  --blue: #0866eb;
  --muted: #52627b;
  --line: #dce6f3;
  --soft: #f2f7ff;
  background: #fff;
  color: var(--ink);
  font-family: Arial, Helvetica, sans-serif;
  overflow: clip;
}
.wrap { width: min(100%, 1280px); margin-inline: auto; padding-inline: 48px; }
.announcement { padding: 11px 16px; background: #0d2347; color: #fff; text-align: center; font-size: 10px; letter-spacing: 1.8px; }
.announcement span { margin-inline: 20px; color: #8cbdff; }
.site-header { position: sticky; top: 0; z-index: 30; background: #ffffffed; backdrop-filter: blur(16px); border-bottom: 1px solid var(--line); }
.nav { display: flex; align-items: center; justify-content: space-between; min-height: 90px; gap: 24px; }
.brand { display: inline-flex; align-items: center; gap: 9px; font-size: 25px; font-weight: 900; letter-spacing: -1px; flex-shrink: 0; }
.brand-mark { color: var(--blue); font-size: 43px; font-style: italic; line-height: 1; }
.nav nav { display: flex; align-items: center; gap: 30px; font-size: 13px; font-weight: 600; }
.nav nav a { transition: color .2s, background .2s; }
.nav nav a:hover { color: var(--blue); }
.admin { padding: 11px 18px; border: 1px solid var(--line); border-radius: 8px; }
.admin:hover { background: var(--soft); }
.hero { display: grid; grid-template-columns: 1fr 1.08fr; align-items: center; gap: 60px; padding-block: 76px 90px; position: relative; }
.hero::before { content: ''; position: absolute; width: 600px; height: 600px; right: -230px; top: -90px; border-radius: 50%; background: radial-gradient(circle, #e6f0ff, transparent 67%); pointer-events: none; }
.hero-copy, .hero-art { min-width: 0; position: relative; }
.eyebrow { color: var(--blue); font-size: 10px; font-weight: 700; line-height: 1.6; letter-spacing: 1.8px; margin-bottom: 22px; }
h1 { font-size: clamp(40px, 4.65vw, 60px); font-weight: 750; line-height: 1.09; letter-spacing: -2.8px; }
h1 em { color: var(--blue); font-style: normal; }
.intro { max-width: 425px; margin-block: 24px 30px; color: var(--muted); font-size: 15px; line-height: 1.85; }
.hero-actions { display: flex; flex-wrap: wrap; gap: 18px 24px; align-items: center; }
.button { display: inline-flex; align-items: center; justify-content: center; gap: 20px; padding: 17px 23px; background: var(--blue); color: #fff; border: 1px solid transparent; border-radius: 9px; font-size: 13px; font-weight: 700; transition: transform .25s, box-shadow .25s, background .25s; }
.button:hover { transform: translateY(-3px); background: #0054cc; box-shadow: 0 10px 24px #0866eb26; }
.button:active { transform: translateY(0); }
.button .iconify, .text-link .iconify, .order .iconify { transition: transform .25s; }
.button:hover .iconify, .order:hover .iconify { transform: translate(3px, -3px); }
.text-link { display: inline-flex; align-items: center; gap: 8px; font-size: 12px; font-weight: 600; padding-block: 12px; }
.text-link:hover { color: var(--blue); }
.text-link:hover .iconify { transform: translateY(3px); }
.hero-notes { display: flex; gap: 24px; margin-top: 28px; font-size: 11px; color: var(--muted); }
.hero-notes span { display: flex; align-items: center; gap: 7px; }
.hero-notes .iconify { color: var(--blue); font-size: 16px; }
.hero-art { padding: 13px; padding-bottom: 44px; background: #fff; border: 1px solid var(--line); border-radius: 20px; box-shadow: 0 24px 65px #143c7914; transform: rotate(2deg); }
.hero-art img { width: 100%; height: auto; display: block; border-radius: 11px; }
.visual-label { display: block; position: absolute; bottom: 15px; left: 20px; font-size: 8px; letter-spacing: 1.5px; color: var(--muted); }
.hero-badge { position: absolute; bottom: -20px; right: 20px; display: flex; align-items: center; gap: 8px; background: #fff; border: 1px solid var(--line); color: var(--blue); box-shadow: 0 8px 20px #0754c419; border-radius: 10px; padding: 13px 18px; font-size: 12px; font-weight: 700; }
.hero-badge .iconify { font-size: 21px; }
.strip { border-block: 1px solid var(--line); position: relative; background: var(--soft); }
.strip::before { content: ''; position: absolute; top: -1px; inset-inline: 0; height: 3px; background: linear-gradient(to right, #4285f4 25%, #ea4335 25% 50%, #fbbc05 50% 75%, #34a853 75%); }
.strip .wrap { display: flex; align-items: center; justify-content: space-between; padding-block: 23px; gap: 12px; font-size: 12px; font-weight: 600; }
.strip span:nth-child(even) { color: var(--blue); }
.catalog { padding-block: 96px; }
section[id] { scroll-margin-top: 115px; }
.heading { display: flex; align-items: end; justify-content: space-between; gap: 40px; margin-bottom: 28px; }
h2 { font-size: clamp(32px, 3.4vw, 44px); font-weight: 650; letter-spacing: -1.8px; line-height: 1.15; }
.heading > p { max-width: 340px; color: var(--muted); font-size: 14px; line-height: 1.85; }
.note { max-width: 760px; margin-bottom: 30px; color: var(--muted); font-size: 12px; line-height: 1.8; }
.products { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 28px; }
.product-card { scroll-margin-top: 145px; border: 1px solid var(--line); border-radius: 18px; padding: 12px; background: #fff; transition: transform .3s, box-shadow .3s, border-color .3s; }
.product-art { display: block; overflow: hidden; border-radius: 10px; background: var(--soft); }
.product-art img { display: block; width: 100%; height: auto; transition: transform .5s ease; }
.product-info { padding: 25px 14px 12px; }
.product-category { margin-bottom: 10px; color: var(--blue); font-size: 10px; font-weight: 700; letter-spacing: 1.3px; text-transform: uppercase; }
.product-info h3 { font-size: 25px; font-weight: 650; letter-spacing: -.8px; margin-bottom: 12px; }
.material { color: var(--muted); font-size: 12px; line-height: 1.7; }
.price { display: flex; align-items: center; flex-wrap: wrap; gap: 7px; font-size: 26px; font-weight: 700; margin-block: 20px 10px; }
.price > span:last-child { color: var(--muted); font-size: 12px; font-weight: 400; }
.product-status { display: flex; align-items: center; gap: 7px; color: var(--muted); font-size: 11px; margin-bottom: 22px; }
.status-dot { width: 6px; height: 6px; border-radius: 50%; background: #168447; }
.order { display: flex; width: 100%; align-items: center; justify-content: space-between; gap: 12px; padding: 15px 17px; border: 1px solid #15803d; border-radius: 8px; color: #fff; background: #15803d; font-size: 13px; font-weight: 700; transition: background .2s; }
.order:not(:disabled):hover { background: #166534; border-color: #166534; }
.order:disabled { color: #596b83; background: #f6f8fb; border-color: var(--line); cursor: not-allowed; }
.availability { margin-top: 24px; padding: 16px 20px; border-left: 3px solid var(--blue); background: var(--soft); color: var(--muted); font-size: 12px; line-height: 1.7; border-radius: 0 8px 8px 0; }
.how { background: #084fbd; background-image: radial-gradient(ellipse at top right, #2479ef, transparent 65%); color: #fff; padding-block: 80px; }
.how .eyebrow { color: #c5dcff; }
.how .heading > p { color: #d7e6ff; }
.steps { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 24px; margin-top: 48px; }
.steps article { padding: 26px; border: 1px solid #ffffff38; border-radius: 14px; background: #ffffff08; transition: background .3s, transform .3s; }
.steps article > span { display: block; font-size: 13px; color: #bad8ff; margin-bottom: 32px; }
.steps h3 { font-size: 21px; font-weight: 600; letter-spacing: -.5px; margin-bottom: 14px; }
.steps p { color: #d7e6ff; font-size: 13px; line-height: 1.85; }
.faq { display: grid; grid-template-columns: .8fr 1.2fr; gap: 72px; padding-block: 96px; }
.faq-heading > p:last-child { color: var(--muted); font-size: 13px; line-height: 1.8; margin-top: 24px; max-width: 270px; }
.faq details { border-bottom: 1px solid var(--line); padding-block: 25px; }
.faq details:first-child { border-top: 1px solid var(--line); }
.faq summary { display: flex; align-items: center; justify-content: space-between; gap: 20px; cursor: pointer; font-size: 14px; font-weight: 600; list-style: none; }
.faq summary::-webkit-details-marker { display: none; }
.faq summary::after { content: '+'; display: grid; place-items: center; width: 26px; height: 26px; flex-shrink: 0; border-radius: 50%; background: var(--soft); color: var(--blue); font-size: 21px; font-weight: 400; transition: transform .25s, background .25s; }
.faq details[open] summary::after { transform: rotate(45deg); background: #dceaff; }
.faq details p { font-size: 13px; line-height: 1.85; color: var(--muted); padding-right: 35px; margin-top: 16px; }
.closing { width: min(calc(100% - 96px), 1184px); padding: 48px; border: 1px solid #d5e5ff; border-radius: 22px; background: linear-gradient(115deg, #edf5ff, #e0edff); position: relative; isolation: isolate; overflow: hidden; }
.closing h2 { margin-bottom: 28px; }
.closing > span { position: absolute; z-index: -1; right: 70px; top: 6px; color: var(--blue); opacity: .12; font-size: 270px; line-height: 1.2; }
footer { display: flex; justify-content: space-between; align-items: center; gap: 20px; padding-block: 48px 38px; }
footer p, footer > span { color: var(--muted); font-size: 11px; }
.skip { position: absolute; top: -100px; }
.skip:focus { top: 10px; left: 10px; background: white; padding: 12px; z-index: 10; }
a:focus-visible, button:focus-visible, summary:focus-visible { outline: 3px solid var(--blue); outline-offset: 5px; }
@media (hover: hover) {
  .product-card:hover { transform: translateY(-6px); border-color: #a9cafa; box-shadow: 0 18px 45px #143c7912; }
  .product-art:hover img { transform: scale(1.035); }
  .steps article:hover { background: #ffffff13; transform: translateY(-4px); }
}
@keyframes enter { from { opacity: 0; transform: translateY(22px); } to { opacity: 1; transform: translateY(0); } }
@keyframes showcase { from { opacity: 0; transform: translateY(28px) rotate(0deg); } to { opacity: 1; transform: translateY(0) rotate(2deg); } }
@keyframes float { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-6px); } }
@media (prefers-reduced-motion: no-preference) {
  .hero-copy > * { animation: enter .7s both; }
  .hero-copy h1 { animation-delay: .08s; }
  .hero-copy .intro { animation-delay: .16s; }
  .hero-actions { animation-delay: .24s; }
  .hero-notes { animation-delay: .32s; }
  .hero-art { animation: showcase .9s .15s both; }
  .hero-badge { animation: float 4s 1.1s ease-in-out 3; }
  .faq details[open] p { animation: enter .25s ease-out; }
  @supports (animation-timeline: view()) {
    .heading, .product-card, .steps article, .faq-heading, .faq-items, .closing { animation: enter both linear; animation-timeline: view(); animation-range: entry 0% entry 20%; }
  }
}
@media (max-width: 1000px) {
  .wrap { padding-inline: 28px; }
  .hero { gap: 32px; padding-block: 60px 75px; }
  h1 { font-size: clamp(34px, 4.8vw, 48px); letter-spacing: -1.8px; }
  .hero-actions { gap: 12px; }
  .hero-notes { gap: 15px; }
  .steps article { padding: 20px; }
  .steps h3 { font-size: 18px; }
  .faq { gap: 40px; }
  .closing { width: calc(100% - 56px); padding: 36px; }
}
@media (max-width: 760px) {
  .hero { grid-template-columns: 1fr; gap: 45px; padding-block: 45px 65px; }
  h1 { font-size: clamp(36px, 7.5vw, 56px); letter-spacing: -2px; }
  .intro { max-width: 530px; }
  .hero-art { max-width: 580px; width: 100%; justify-self: center; }
  .heading { align-items: start; flex-direction: column; gap: 20px; }
  .heading > p { max-width: 520px; }
  .products { gap: 20px; }
  .product-info { padding: 20px 7px 7px; }
  .product-info h3 { font-size: 21px; }
  .order { font-size: 12px; padding: 13px; }
  .steps { grid-template-columns: 1fr; gap: 16px; margin-top: 30px; }
  .steps article { padding: 24px; }
  .steps article > span { margin-bottom: 18px; }
  .faq { grid-template-columns: 1fr; gap: 32px; }
  .faq-heading h2 br { display: none; }
  .faq-heading > p:last-child { max-width: none; margin-top: 16px; }
  footer { flex-wrap: wrap; }
  footer p { width: 100%; order: 3; }
}
@media (max-width: 540px) {
  .wrap { padding-inline: 20px; }
  .announcement { font-size: 8px; letter-spacing: 1px; }
  .announcement span { margin-inline: 7px; }
  section[id] { scroll-margin-top: 145px; }
  .nav { min-height: 80px; flex-wrap: wrap; justify-content: center; gap: 12px; padding-block: 16px; }
  .brand { font-size: 21px; }
  .brand-mark { font-size: 34px; }
  .nav nav { justify-content: center; gap: 22px; font-size: 11px; }
  .admin { padding: 8px 12px; }
  .hero { padding-top: 36px; }
  .eyebrow { font-size: 9px; letter-spacing: 1.3px; margin-bottom: 18px; }
  .intro { font-size: 14px; }
  .hero-notes { flex-wrap: wrap; gap: 12px 20px; }
  .hero-art { padding: 8px 8px 38px; border-radius: 14px; }
  .hero-art img { border-radius: 8px; }
  .visual-label { left: 14px; font-size: 7px; }
  .hero-badge { right: 12px; bottom: -22px; font-size: 10px; padding: 10px 12px; }
  .strip .wrap { flex-wrap: wrap; justify-content: center; gap: 10px; padding-block: 20px; font-size: 10px; }
  .catalog, .faq { padding-block: 60px; }
  .products { grid-template-columns: 1fr; }
  .product-info { padding: 22px 10px 10px; }
  .material { font-size: 12px; }
  .how { padding-block: 55px; }
  .closing { width: calc(100% - 40px); padding: 30px 24px; border-radius: 16px; }
  .closing > span { right: -25px; top: 40px; font-size: 190px; }
  .closing h2 { font-size: 34px; }
  footer { gap: 16px; padding-block: 34px; }
  footer > span { font-size: 10px; }
}
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after { animation: none !important; transition: none !important; }
  .hero-art { transform: none; }
}
</style>

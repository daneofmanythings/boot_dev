const { argv } = require('node:process')
const { crawlPage } = require('./crawl')
const { printReports } = require('./report')

async function main() {
  if (argv.length != 3) {
    console.log(`need one html argument. got=${argv}`)
    return
  }

  const baseURL = process.argv[2]

  console.log(`Begin crawling on url: ${argv[2]}`)

  pages = await crawlPage(new URL(baseURL), new URL(baseURL), {})

  printReports(pages)
  return
}

main()

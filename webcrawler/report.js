

function printReports(pages) {
  console.log("Print Reports is starting...")

  const sortedPagesArray = sortPagesByCount(pages)
  for (const page of sortedPagesArray) {
    console.log(`Found ${page[1]} internal links to ${page[0]}`)
  }
}

function sortPagesByCount(pages) {
  // this is going to be slow. idc
  const pagesArray = Object.entries(pages)
  pagesArray.sort((a, b) => {
    return b[1] - a[1]
  })
  return pagesArray
}

module.exports = {
  printReports
}

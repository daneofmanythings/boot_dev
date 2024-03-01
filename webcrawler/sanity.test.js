const { recordOutputs } = require('./sanity.js')
const { test, expect } = require('@jest/globals')

test('squaring numbers', () => {
  const inputs = [1, 2, 3]
  const expected = { "1": 1, "2": 4, "3": 9 }
  const actual = recordOutputs((input) => {
    return input * input
  }, inputs)
  expect(actual).toStrictEqual(expected)
})

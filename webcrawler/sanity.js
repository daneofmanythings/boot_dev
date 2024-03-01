function recordOutputs(func, inputs) {
  outputs = {}
  for (const input of inputs) {
    outputs[input] = func(input)
  }
  return outputs
}

module.exports = {
  recordOutputs,
}

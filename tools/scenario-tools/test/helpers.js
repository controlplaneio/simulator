const { resolve, join } = require('path')
const { spy } = require('sinon')

function fixture (name) {
  return resolve(join(__dirname, 'fixtures', name))
}

function testoutput (name) {
  return resolve(join(__dirname, 'output', name))
}

function createSpyingLogger () {
  return { warn: spy(), info: spy(), error: spy() }
}

module.exports = {
  fixture,
  testoutput,
  createSpyingLogger
}

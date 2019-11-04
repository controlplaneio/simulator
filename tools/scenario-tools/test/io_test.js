const {resolve, join} = require('path')
const {openSync, readFileSync} = require('fs')
const test = require('ava')
const {loadHintsFile, writeHintsFile, findScenarioHintsFiles} = require('../lib/io')

function fixture(name) {
  return resolve(join(__dirname, 'fixtures', name))
}

function testoutput(name) {
  return resolve(join(__dirname, 'output', name))

}

test('loadHintsFile parses a yaml file', t => {
  const actual = loadHintsFile(fixture('test.yaml'))
  t.deepEqual(actual, {test: 'test'})
})

test('writeHintsFile serializes a js object to yaml', t => {
  const input = {test: 'test'}
  const outputFile = testoutput('writeHintsTest.yaml')

  // create an empty test output file to overwrite
  openSync(outputFile, 'w');

  writeHintsFile(input, outputFile)
  t.deepEqual(readFileSync(outputFile, 'utf-8'), 'test: test\n')
})

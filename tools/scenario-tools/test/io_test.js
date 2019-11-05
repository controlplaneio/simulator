const { resolve, join } = require('path')
const { openSync, readFileSync } = require('fs')
const test = require('ava')
const { loadYamlFile, writeYamlFile, getProgress, saveProgress } = require('../lib/io')

function fixture (name) {
  return resolve(join(__dirname, 'fixtures', name))
}

function testoutput (name) {
  return resolve(join(__dirname, 'output', name))
}

test('loadYamlFile parses a yaml file', t => {
  const actual = loadYamlFile(fixture('test.yaml'))
  t.deepEqual({ test: 'test' }, actual)
})

test('writeYamlFile serializes a js object to yaml', t => {
  const input = { test: 'test' }
  const outputFile = testoutput('writeYamlTest.yaml')

  // create an empty test output file to overwrite
  openSync(outputFile, 'w')

  writeYamlFile(input, outputFile)
  t.deepEqual('test: test\n', readFileSync(outputFile, 'utf-8'))
})

test('getProgress parses a JSON file', t => {
  const actual = getProgress(fixture('test-progress.json'))
  t.deepEqual({ 'Task 1': 1 }, actual)
})

test('saveProgress writes a JSON file', t => {
  const input = { test: 'test' }
  const outputFile = testoutput('saveProgres.json')

  saveProgress(input, outputFile)
  t.deepEqual('{"test":"test"}', readFileSync(outputFile, 'utf-8'))
})

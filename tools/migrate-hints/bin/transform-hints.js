#!/usr/bin/env node

// Usage: run this in the root of the simulator dir
//
// `npm install js-yaml`
// `node transform-hints.js`
//
const yaml = require('js-yaml')
const {truncateSync, writeFileSync, readFileSync, readdirSync} = require('fs')
const {resolve,join} = require('path')
const {transformV0ToV1} = require('../lib/transforms')

// Mixing sync and async but it doesnt matter as this is a quick n dirty script
function loadHintsFile(p) {
  console.log(`Loading ${p}`)
  const doc = yaml.safeLoad(readFileSync(p, 'utf8'));
  console.log(doc)
  return doc
}

function writeHintsFile(hints, p) {
  console.log(`Writing transformed file ${p}`)
  const contents = yaml.safeDump(hints)
  truncateSync(p)
  writeFileSync(p, contents)
}

function findScenarioHintsFiles(scenariosDir) {
  const absPath = resolve(scenariosDir)

  return readdirSync(absPath, { withFileTypes: true })
    .filter(dirent => dirent.isDirectory())
    .map(dirent => join(absPath, dirent.name, 'hints.yaml'))
}

hintsFiles = findScenarioHintsFiles('./simulation-scripts/scenario')
hintsFiles.forEach(hintsFile => {
  const original = loadHintsFile(hintsFile)
  const transformed = transformV0ToV1(original)
  writeHintsFile(transformed, hintFile)
})


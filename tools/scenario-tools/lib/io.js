const yaml = require('js-yaml')
const { truncateSync, writeFileSync, readFileSync, readdirSync } = require('fs')
const { resolve, join } = require('path')
const { createLogger } = require('../lib/logger')

const logger = createLogger({})

// Loads and parses a `hints.yaml` from the supplied absolute path.
// Returns an object representing the yaml file
function loadYamlFile (p) {
  logger.info(`Loading ${p}`)
  const doc = yaml.safeLoad(readFileSync(p, 'utf8'))
  return doc
}

// Serializes the supplied `hints` to YAML and overwrites an existing
// `hints.yaml` file to the supplied path `p` with the YAMl
function writeYamlFile (hints, p) {
  logger.info(`Writing transformed file ${p}`)
  const contents = yaml.safeDump(hints)
  truncateSync(p)
  writeFileSync(p, contents)
}

// Given a relative path to a scenario directory, scans for scenarios
// Returns a list of absolute paths to all `hints.yaml` files
function findScenarioHintsFiles (scenariosDir) {
  const absPath = resolve(scenariosDir)

  return readdirSync(absPath, { withFileTypes: true })
    .filter(dirent => dirent.isDirectory())
    .map(dirent => join(absPath, dirent.name, 'hints.yaml'))
}

module.exports = {
  loadYamlFile,
  writeYamlFile,
  findScenarioHintsFiles
}

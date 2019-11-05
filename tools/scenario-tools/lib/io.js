const yaml = require('js-yaml')
const { truncateSync, writeFileSync, readFileSync, readdirSync, existsSync } = require('fs')
const { resolve, join } = require('path')
const { createLogger } = require('../lib/logger')

const PROGRESS_FILE_PATH = '/progress.json'

// TODO(rem): retrieve from S3 Bucket
function getProgress (p = PROGRESS_FILE_PATH) {
  const absPath = resolve(p)
  if (!existsSync(absPath)) {
    writeFileSync(absPath, '{}')
  }

  const contents = readFileSync(absPath, 'utf-8')
  return JSON.parse(contents)
}

// TODO(rem): write to S3 Bucket
function saveProgress (progress, p = PROGRESS_FILE_PATH) {
  const absPath = resolve(p)
  const contents = JSON.stringify(progress)
  if (existsSync(absPath)) {
    truncateSync(absPath)
  }

  writeFileSync(absPath, contents)
}

// Loads and parses a `hints.yaml` from the supplied absolute path.
// Returns an object representing the yaml file
function loadYamlFile (p) {
  const doc = yaml.safeLoad(readFileSync(p, 'utf8'))
  return doc
}

// Serializes the supplied `hints` to YAML and overwrites an existing
// `hints.yaml` file to the supplied path `p` with the YAMl
function writeYamlFile (hints, p) {
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
  getProgress,
  saveProgress,
  loadYamlFile,
  writeYamlFile,
  findScenarioHintsFiles
}

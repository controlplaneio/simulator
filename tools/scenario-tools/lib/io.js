const yaml = require('js-yaml')
const {
  truncateSync,
  writeFileSync,
  readFileSync,
  readdirSync
} = require('fs')
const { get, post } = require('axios')
const http = require('http-status-codes')
const { resolve, join } = require('path')
const { createLogger } = require('../lib/logger')

const logger = createLogger({})

const PROGRESS_FILE_PATH = '/progress.json'
const PROGRESS_SERVER_URL = 'http://localhost:51234'

async function getProgress (name, p = PROGRESS_FILE_PATH) {
  let response
  try {
    response = await get(PROGRESS_SERVER_URL + '?scenario=' + name)
    return response.data
  } catch (e) {
    if (e.response === undefined) {
      // Unknown error
      throw e
    }

    if (e.response.status === http.NOT_FOUND) {
      const progress = {
        name: name,
        currentTask: null,
        tasks: []
      }

      try {
        await saveProgress(progress)
        return progress
      } catch (e) {
        logger.error('Error POSTing initial progress', { response: e.response })
        return
      }
    }

    if (e.response.status !== http.OK) {
      logger.error('Unexpected HTTP status getting progress', { response })
    }
  }
}

async function saveProgress (progress, p = PROGRESS_FILE_PATH) {
  let response
  try {
    response = await post(PROGRESS_SERVER_URL, progress)
    return
  } catch (e) {
    if (e.response === undefined) {
      // Unknown error
      throw e
    }

    if (response.status !== http.OK) {
      const message = 'Unexpected HTTP status POSTing progress'
      logger.error(message, { response })
      throw new Error(message)
    }

    logger.error('Error POSTing progress', { error: e, progress })
    throw e
  }
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
function findScenarioFiles (scenariosDir, filename) {
  const absPath = resolve(scenariosDir)

  return readdirSync(absPath, { withFileTypes: true })
    .filter(dirent => dirent.isDirectory())
    .map(dirent => join(absPath, dirent.name, filename))
}

module.exports = {
  getProgress,
  saveProgress,
  loadYamlFile,
  writeYamlFile,
  findScenarioFiles,
  PROGRESS_SERVER_URL
}

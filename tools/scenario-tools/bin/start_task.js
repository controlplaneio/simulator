#!/usr/bin/env node

const { createLogger } = require('../lib/logger')
const { cloneArray } = require('../lib/helpers')
const { startTask } = require('../.lib/tasks.js')

// const args = process.argv.slice(2)
const logger = createLogger({})

// Global fallback to handle and pretty print unhandled errors
process.on('uncaughtException', (err) => {
  logger.error(`Uncaught exception: ${err}`)

  if (err && err.stack) {
    logger.error(err.stack)
  }

  process.exit(1)
})

// Global fallback to handle and pretty print unhandled promise rejection errors
process.on('unhandledRejection', (err) => {
  logger.error(`Promise rejection: ${err}`)

  if (err && err.stack) {
    logger.error(err.stack)
  }

  process.exit(1)
})

const args = cloneArray(process.argv)

args.shift() // remove `node` from argv
args.shift() // remove `scenario.js` from argv

if (args.length !== 1) {
  logger.error('Please provide the task you wish to start. These are listed in the instructions')
}

const success = startTask(args[0])
if (success) {
  process.exit(0)
} else {
  process.exit(1)
}

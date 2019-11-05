#!/usr/bin/env node

const os = require('os')

const { createLogger } = require('../lib/logger')
const { parse, showUsage } = require('../lib/args')
const { migrate } = require('../lib/migrate')
const { showHints, nextHint } = require('../lib/hints')

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

// Only show dev options when not in the attack container
function showHelp () {
  let dev = false
  if (os.hostname() !== 'attack') {
    dev = true
  }

  showUsage(dev)
}

function cloneArray (arr) {
  return arr.slice(0)
}

const args = cloneArray(process.argv)

args.shift() // remove `node` from argv
args.shift() // remove `scenario.js` from argv

try {
  const { command, options } = parse(args)

  // TODO(rem): This needs tidying up and pulling into its own module
  // let's check we are happy with the UX first
  if (command === 'migrate') {
    migrate(options)
    process.exit(0)
  } else if (command === 'show-hints') {
    showHints(options.task)
    process.exit(0)
  } else if (command === 'next-hint') {
    nextHint(options.task)
    process.exit(0)
  } else {
    showHelp()
  }
} catch (e) {
  logger.error('Unrecognised cli arguments - try running the \'help\' command')
  logger.error(e.message)
  showHelp()
}

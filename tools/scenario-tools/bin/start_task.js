#!/usr/bin/env node

const { createLogger } = require('../lib/logger')
const { cloneArray } = require('../lib/helpers')
const { startTask } = require('../lib/tasks.js')

const logger = createLogger({})

require('../lib/error-handler')

const args = cloneArray(process.argv)

args.shift() // remove `node` from argv
args.shift() // remove `scenario.js` from argv

if (process.argv0 === 'start_task' && args.length !== 1) {
  logger.error('Please provide the task you wish to start. These are listed in the instructions')
  process.exit(1)
}

// ignore the result - it will be true/false depending on whether the user
// actually switched task
startTask(args[0]).then(_ => {
  process.exit(0)
}, reason => {
  logger.error(reason.message)
  process.exit(1)
})

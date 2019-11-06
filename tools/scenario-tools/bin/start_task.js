#!/usr/bin/env node

const { createLogger } = require('../lib/logger')
const { cloneArray } = require('../lib/helpers')
const { startTask } = require('../lib/tasks.js')

const logger = createLogger({})

require('../lib/error-handler')

const args = cloneArray(process.argv)

args.shift() // remove `node` from argv
args.shift() // remove `scenario.js` from argv

if (args.length !== 1) {
  logger.error('Please provide the task you wish to start. These are listed in the instructions')
  process.exit(1)
}

const success = startTask(args[0])
if (success) {
  process.exit(0)
} else {
  process.exit(1)
}

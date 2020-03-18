#!/usr/bin/env node

const { createLogger } = require('../lib/logger')
const { getCurrentTask } = require('../lib/tasks.js')
const { nextHint } = require('../lib/hints.js')

const logger = createLogger({})
require('../lib/error-handler')

getCurrentTask().then(task => {
  if (task === undefined) {
    logger.error(
      'You have not started a task.  Please run `start_task` to select your task')
    process.exit(1)
  }

  nextHint(Number(task)).then(_ => {
    process.exit(0)
  }, reason => {
    logger.errror(reason.message)
    process.exit(1)
  }, reason => {
    logger.errror(reason.message)
    process.exit(1)
  })
})

const { createLogger } = require('../lib/logger')

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

const { inspect } = require('util')
const { MESSAGE, SPLAT } = require('triple-beam')
const winston = require('winston')
const { format } = require('logform')

const formatLogMessage = format((info, opts) => {
  const depth = opts.depth || null
  if (info[SPLAT]) {
    for (const splat of info[SPLAT]) {
      info.message += '\n' + inspect(splat, false, depth, opts.colorize)
    }
  }

  info[MESSAGE] = `${info.level}:${info.message}`

  return info
})

// Creates a new logger with the supplied options.
// Options are:
// - colorize - whethers to colourise the output
// - level - the logging level
// - depth - how deep to inspect js objects
function createLogger (options) {
  const formats = []

  // default to coloured output
  options = Object.assign({ colorize: true }, options)

  if (options.colorize) {
    formats.push(winston.format.colorize())
  }

  formats.push(winston.format.align())
  formats.push(formatLogMessage({ colorize: options.colorize }))

  const transport = new winston.transports.Console({
    level: options.level || 'info',
    format: winston.format.combine(...formats)
  })

  return winston.createLogger({ transports: [transport] })
}

module.exports.createLogger = createLogger

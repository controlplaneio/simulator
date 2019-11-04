const commandLineArgs = require('command-line-args')
const commandLineCommands = require('command-line-commands')
const {commandLineUsage} = require('command-line-usage')

function parse(argv) {
  const commands = [ 'migrate', 'show-hints', 'next-hint' ]

  const { command, argv: remaining } = commandLineCommands(commands, argv)

  let options = {}

  if (command.name === 'migrate') {
    const migrateArguments = [ { name: 'all', alias: 'a', type: Boolean } ]

    options = commandLineArgs(migrateArguments, {
      stopAtFirstUnknown: true,
      argv: remaining
    })
  }

  if (command === 'show-hints' || command === 'next-hint') {
    const hintArguments = [ { name: 'task', alias: 't', type: String } ]

    options = commandLineArgs(hintArguments, {
      stopAtFirstUnknown: true,
      argv: remaining
    })
  }

  return { command, options }
}

module.exports.parse = parse

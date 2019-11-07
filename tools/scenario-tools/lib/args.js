const commandLineArgs = require('command-line-args')
const commandLineCommands = require('command-line-commands')
const commandLineUsage = require('command-line-usage')

function createArgumentError (msg) {
  const err = new Error(msg)
  err.type = 'IncorrectArguments'
  return err
}

function parse (argv) {
  const commands = ['migrate', 'help']

  const { command, argv: remaining } = commandLineCommands(commands, argv)

  let options = {}

  if (command === 'migrate') {
    const migrateArguments = [{
      name: 'name',
      alias: 'n',
      type: String
    }, {
      name: 'all',
      alias: 'a',
      type: Boolean,
      defaultValue: false
    }
    ]

    options = commandLineArgs(migrateArguments, {
      stopAtFirstUnknown: true,
      argv: remaining
    })

    if (!options.all && !options.name) {
      throw createArgumentError('You must supply one of --all or --name arguments')
    } else if (options.all && options.name) {
      throw createArgumentError('You cannot supply both --all or --name arguments')
    }
  }

  return { command, options }
}

function showUsage () {
  const sections = [
    {
      header: 'Scenario Tool',
      content: 'Helper to interact with the current scenario'
    }, {
      header: 'help',
      content: 'Print this usage guide.'
    }, {
      header: 'migrate',
      content: 'Helper for (mass) migration of scenario tasks.yaml files',
      optionList: [{
        name: 'all',
        typeLabel: '{underline file}',
        description: 'The task to show the hints for'
      }, {
        name: 'name',
        typeLabel: '{underline file}',
        description: 'The name of the scenario to migrate'
      }
      ]
    }
  ]

  const usage = commandLineUsage(sections)
  console.log(usage)
}

module.exports = {
  parse,
  showUsage
}

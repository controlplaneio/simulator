const commandLineArgs = require('command-line-args')
const commandLineCommands = require('command-line-commands')
const commandLineUsage = require('command-line-usage')

function createArgumentError (msg) {
  const err = new Error(msg)
  err.type = 'IncorrectArguments'
  return err
}

function parse (argv) {
  const commands = ['migrate', 'show-hints', 'next-hint', 'help']

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

  if (command === 'show-hints' || command === 'next-hint') {
    const hintArguments = [{ name: 'task', alias: 't', type: String }]

    options = commandLineArgs(hintArguments, {
      stopAtFirstUnknown: true,
      argv: remaining
    })

    if (!options.task) {
      throw createArgumentError('Task is a required argument')
    }
  }

  return { command, options }
}

function showUsage (dev) {
  const sections = [
    {
      header: 'Scenario Tool',
      content: 'Helper to interact with the current scenario'
    }, {
      header: 'show-hints',
      content: 'Shows any and all hints already seen for the supplied task',
      optionList: [
        {
          name: 'task',
          typeLabel: '{underline file}',
          description: 'The task to show the hints you have already seen'
        }
      ]
    }, {
      header: 'next-hint',
      content: 'Shows the next hint for the supplied task',
      optionList: [
        {
          name: 'task',
          typeLabel: '{underline file}',
          description: 'The task to show the hints for'
        }
      ]
    }, {
      header: 'help',
      content: 'Print this usage guide.'
    }
  ]

  // hide the devtools if we are in the attack container
  if (dev) {
    sections.push({
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
    })
  }
  const usage = commandLineUsage(sections)
  console.log(usage)
}

module.exports = {
  parse,
  showUsage
}

const commandLineArgs = require('command-line-args')
const commandLineCommands = require('command-line-commands')
const commandLineUsage = require('command-line-usage')

function parse(argv) {
  const commands = [ 'migrate', 'show-hints', 'next-hint', 'help' ]

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

function showUsage(dev) {
  const sections = [
    {
      header: 'Scenario Tool',
      content: 'Helper to interact with the current scenario'
    }, {
      header: 'show-hints',
      content: 'Shows any and all hints already seen for the given task',
      optionList: [
        {
          name: 'task',
          typeLabel: '{underline file}',
          description: 'The task to show the hints for'
        }
      ]
    }, {
      header: 'show-hints',
      content: 'Shows any and all hints already seen for the given task',
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

  if (dev) {
    sections.push({
      header: 'migrate',
      content: 'Helper for (mass) migration of scenario tasks.yaml files',
      optionList: [ {
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

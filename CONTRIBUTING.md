# Contributing to Simulator

👍🎉 We're thrilled that you're interested in contributing! 🎉👍

`Simulator` is under Apache 2.0 license and welcomes contributions through GitHub pull requests.

This document provides guidelines for contributing to `Simulator`. While we adhere to certain standards to maintain quality, we welcome your PRs and are happy to collaborate to make them fit our standards. Suggestions for improving this guide are also welcome via pull requests.

## Table Of Contents

- [Code of Conduct](#code-of-conduct)
- [I Don't Want To Read This Whole Thing I Just Have a Question!!!](#i-dont-want-to-read-this-whole-thing-i-just-have-a-question)
- [What Should I Know Before I Get Started?](#what-should-i-know-before-i-get-started)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
    - [Before Submitting a Bug Report](#before-submitting-a-bug-report)
    - [How Do I Submit a (Good) Bug Report?](#how-do-i-submit-a-good-bug-report)
  - [Suggesting Enhancements](#suggesting-enhancements)
    - [Before Submitting an Enhancement Suggestion](#before-submitting-an-enhancement-suggestion)
    - [How Do I Submit A (Good) Enhancement Suggestion?](#how-do-i-submit-a-good-enhancement-suggestion)
  - [Your First Code Contribution](#your-first-code-contribution)
    - [Development](#development)
  - [Pull Requests](#pull-requests)
- [Style Guides](#style-guides)
  - [Git Commit Messages](#git-commit-messages)
  - [General Style Guide](#general-style-guide)
  - [GoLang Style Guide](#golang-style-guide)
  - [Documentation Style Guide](#documentation-style-guide)

---

## Code of Conduct

All participants in this project are governed by our [Code of Conduct](CODE_OF_CONDUCT.md). We expect everyone to abide by this code. Please report any unacceptable behavior to [andy@control-plane.io](mailto:andy@control-plane.io.).

## I Don't Want To Read This Whole Thing I Just Have a Question!!!

We have an official message board with a detailed FAQ and where the community chimes in with helpful advice if you have questions.

We also have an issue template for questions [here](https://github.com/controlplaneio/simulator/issues/new).

## What Should I Know Before I Get Started?

The Simulator project is centered around a single central repo - [Simulator](https://github.com/controlplaneio/simulator). This is the core component of the project, encompassing all functionalities of the Simulator. Feel free to explore the repository to understand more about how the Simulator works and where you might be able to contribute.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for `Simulator`. Following these guidelines helps maintainers and the
community understand your report, reproduce the behaviour, and find related reports.

Before creating bug reports, please check [this list](#before-submitting-a-bug-report) as you might find out that you
don't need to create one. When you are creating a bug report, please [include as many details as possible](#how-do-i-submit-a-good-bug-report).
Fill out the issue template for bugs, the information it asks for helps us resolve issues faster.

> **Note:** If you find a **Closed** issue that seems like it is the same thing that you're experiencing, open a new issue
> and include a link to the original issue in the body of your new one.

#### Before Submitting a Bug Report

- **Perform a [cursory search](https://github.com/search?q=+is:issue+user:controlplaneio)** to see if the problem has already been reported. If it has **and the issue is still open**, add a comment to the existing issue instead of opening a new one

#### How Do I Submit a (Good) Bug Report?

Bugs are tracked as [GitHub issues](https://guides.github.com/features/issues/). Create an issue on that repository and provide the following information by filling in the issue template [here](https://github.com/controlplaneio/simulator/issues/new).

Explain the problem and include additional details to help maintainers reproduce the problem:

- **Use a clear and descriptive title** for the issue to identify the problem
- **Describe the exact steps which reproduce the problem** in as many details as possible. For example, start by explaining how you started `simulator` and what you did until you noticed an error or the unexpected behaviour. Additional commmands, configuration files and logs, can help us better understand the problem.
- **Provide specific examples to demonstrate the steps**. Include links to files or GitHub projects, or copy/pasteable snippets, which you use in those examples. If you're providing snippets in the issue, use [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines).
- **Describe the behaviour you observed after following the steps** and point out what exactly is the problem with that behaviour
- **Explain which behaviour you expected to see instead and why.**

Provide more context by answering these questions:

- **Did the problem start happening recently** (e.g. after updating to a new version of Simulator) or was this always a problem ?
- If the problem started happening recently, **can you reproduce the problem in an older version of Simulator?** What is the most recent version in which the problem does not happen? You can download older versions of Simulator from  [the releases page](https://github.com/controlplaneio/simulator/releases)
- **Can you reliably reproduce the issue?** If not, please provide details about how often the problem happens and under which conditions it normally happens

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Simulator, including completely new features and minor improvements to existing functionality. Following these guidelines helps maintainers and the community understand your suggestion and find related suggestions.

Before creating enhancement suggestions, please check [this list](#before-submitting-an-enhancement-suggestion) as you might find out that you don't need to create one. When you are creating an enhancement suggestion, please [include as many details as possible](#how-do-i-submit-a-good-enhancement-suggestion). Fill in the template feature request template, including the steps that you imagine you would take if the feature you're requesting existed.

#### Before Submitting an Enhancement Suggestion

- **Perform a [cursory search](https://github.com/search?q=+is:issue+user:controlplaneio)** to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one.

#### How Do I Submit A (Good) Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://guides.github.com/features/issues/). Create an issue on this
repository and provide the following information:

- **Use a clear and descriptive title** for the issue to identify the suggestion
- **Provide a step-by-step description of the suggested enhancement** in as many details as possible
- **Provide specific examples to demonstrate the steps**. Include copy/pasteable snippets which you use in those examples, as [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines)
- **Describe the current behaviour** and **explain which behaviour you expected to see instead** and why ?
- **Explain why this enhancement would be useful** to most Simulator users and isn't something that can or should be implemented
  as a separate community project
- **List some other tools where this enhancement exists.**
- **Specify which version of Simulator you're using.** You can get the exact version by running `simulator version` in your terminal
- **Specify the name and version of the OS you're using.**

### Your First Code Contribution

Unsure where to begin contributing to `Simulator`? You can start by looking through these `Good First Issue` and `Help Wanted`
issues:

- [Good First Issue issues][good_first_issue] - issues which should only require a few lines of code, and a test or two
- [Help wanted issues][help_wanted] - issues which should be a bit more involved than `Good First Issue` issues

Both issue lists are sorted by total number of comments. While not perfect, number of comments is a reasonable proxy for impact a given change will have.

#### Development

- Simulator is written in [Go programming language](https://golang.org/doc/).
- Please follow the instructions [here](https://go.dev/doc/install) to install Go on your machine
- Please install `golangcli-lint` by following the instructions [here](https://golangci-lint.run/usage/install/#local-installation)

### Pull Requests

The process described here has several goals:

- Maintain the quality of `Simulator`
- Fix problems that are important to users
- Engage the community in working toward the best possible Simulator
- Enable a sustainable system for Simulator's maintainers to review contributions

Please follow these steps to have your contribution considered by the maintainers:

<!-- markdownlint-disable no-inline-html -->

1. Follow all instructions in the template
2. Follow the [style guides](#style-guides)
3. After you submit your pull request, verify that all [status checks](https://help.github.com/articles/about-status-checks/)
   are passing
   <details>
    <summary>What if the status checks are failing?</summary>
    If a status check is failing, and you believe that the failure is unrelated to your change, please leave a comment on
    the pull request explaining why you believe the failure is unrelated. A maintainer will re-run the status check for
    you. If we conclude that the failure was a false positive, then we will open an issue to track that problem with our
    status check suite.
   </details>

<!-- markdownlint-enable no-inline-html -->

While the prerequisites above must be satisfied prior to having your pull request reviewed, the reviewer(s) may ask you to
complete additional tests, or other changes before your pull request can be ultimately accepted.

## Style Guides

### Git Commit Messages

- It's strongly preferred you [GPG Verify][commit_signing] your commits if you can
- Follow [Conventional Commits](https://www.conventionalcommits.org)
- Use the present tense ("add feature" not "added feature")
- Use the imperative mood ("move cursor to..." not "moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

### General Style Guide

Look at installing an `.editorconfig` plugin or configure your editor to match the `.editorconfig` file in the root of the
repository.

### GoLang Style Guide

All Go code is linted with [golangci-lint](https://golangci-lint.run/).

For formatting rely on `gofmt` to handle styling.

### Documentation Style Guide

All markdown code is linted with [markdownlint-cli](https://github.com/igorshubovych/markdownlint-cli).

[good_first_issue]:https://github.com/controlplaneio/simulator/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22+sort%3Acomments-desc
[help_wanted]: https://github.com/controlplaneio/simulator/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22help+wanted%22

[commit_signing]: https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/managing-commit-signature-verification
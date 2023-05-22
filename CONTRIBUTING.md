If you are new to contributing to open-source projects, we want to encourage you to take a look at this great guide made by GitHub on how to contribute to open-source projects [https://opensource.guide/how-to-contribute/](https://opensource.guide/how-to-contribute/).

Feel free to reach out to any of the maintainers or other community members if you have any questions.

## Contributing to Debricked CLI

### Projects setup and Architecture

Requirements

    Go version 1.20 or higher

### If you have questions

- First check and search the current issues and see if they can help answer your questions.
- If you then still feel the need to ask a question and need clarification, we recommend the following:
    - Open an Issue.
    - Provide as much context as you can about what you’re running into.
    - Provide project and platform versions, depending on what seems relevant.

### Feature Request

You're also welcome to submit ideas for enhancements to Debricked CLI. When doing so, please search the issue list to
see if the enhancement has already been filed. If it has, vote for it (add a reaction to it) and optionally add a
comment with your perspective on the idea.

## Ready? Let’s go!

To get started, [fork this project](https://github.com/debricked/cli/fork) to your own git.


 Make sure to keep your fork up to date as well. You can do so by: 
 
`$ git remote add debricked-cli https://github.com/debricked/cli.git`

`$ git fetch debricked-cli`

`$ git checkout main`

`$ git rebase debricked-cli/main`

`$ git push --force-with-lease`

## Doing your work

When you start implementing your changes make sure that you create a new branch before. It’s the proper way to do things and also helps keep everything neatly organised from the master branch. This way it’s less headache inducing managing multiple PRs for every task completed.

Create Branch:

`$ git checkout -b my-cool-branch`

## Submitting a Pull Request

### Tidy it up

Before submitting your PRs, we would appreciate it if you took some time to clean up your branch. This makes it a lot easier to test, accept and merge your changes into the master branch.

1. Commit the change(s) and push to your fork

```
$ git add .
$ git commit -S -m "This is a cool commit"
$ git push -u origin my-cool-branch
```

2. [Submit a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request).

# Bug Reports

**Before Submitting a Bug Report**

Make sure to include as much details as possible by using our general guidelines below before submitting a bug issue. 

- Make sure that your fork is up to date.
- Determine if your bug is really a bug and not an error on your side.
- Make sure not to create duplicate issues, search the current issues before creating a new one.
- Collect information about the bug:
    - Possibly your input and the output.
    - Can you reliably reproduce the issue? And can you also reproduce it with older versions?

> Never report security related issues, vulnerabilities or bugs to the issue tracker, or elsewhere in public. Instead sensitive bugs must be sent by email to **[security@debricked.com](mailto:security@debricked.com)**

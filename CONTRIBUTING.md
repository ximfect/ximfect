# ximfect Contributing Guidelines

Below are some basic guidelines for different ways you can contribute to the development of ximfect.

These are not perfect, but they should give you a general idea of what to expect if you want to contribute.

Remember to make yourself familiar with the [Code of Conduct](CODE_OF_CONDUCT.md).

## Discussing issues/PRs?

* Discussing issues and PRs is contributing too!
* While discussing:
  1. Respect others
  2. Have some common sense
  3. Remember that not everyone is the best, and they can make some mistakes too

## Making an issue?

* Check if an open issue for the problem exists.
* If not, create one. Remember to use the templates we provide you, and try to give all the needed information. 

## Making a pull request?

* Check if an open PR for the bug exists.
* If not, make sure to:
  1. Run `go vet` and fix any errors it shows 
  2. Run `gofmt` to format your changes
  3. While commiting, write a list of all changes you've made in the descriptions of your commits.
     
     Example:
     
     ```sh
     $ git commit -m "Fix bugs + optimise packages
     >
     > * Bug A, B and C (issues #123, #234, #567) no longer occur.
     > * Optimised package decoding.
     > * Fixed the `about-package` action not showing last file in the package."
     ```

## Thank you for contributing

For the most part, ximfect was developed by a single developer, and any and all contributions are nothing but welcome.

Even if you didn't contribute, thank you for taking the time to read this.

~ qeaml

# act

Run your [GitHub Actions](https://developer.github.com/actions/) locally! Why would you want to do this? Two reasons:

- **Fast Feedback** - Rather than having to commit/push every time you want to test out the changes you are making to your `.github/workflows/` files (or for any changes to embedded GitHub actions), you can use `act` to run the actions locally. The [environment variables](https://help.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables#default-environment-variables) and [filesystem](https://help.github.com/en/actions/reference/virtual-environments-for-github-hosted-runners#filesystems-on-github-hosted-runners) are all configured to match what GitHub provides.
- **Local Task Runner** - I love [make](https://en.wikipedia.org/wiki/Make_(software)). However, I also hate repeating myself. With `act`, you can use the GitHub Actions defined in your `.github/workflows/` to replace your `Makefile`!

## install act

+ [https://github.com/nektos/act](https://github.com/nektos/act)

```bash
curl -s https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
···
# How do I contribute code to OpenIM

<p align="center">
  <a href="./CONTRIBUTING.md">Englist</a> · 
  <a href="./CONTRIBUTING-zh_CN.md">中文</a> · 
  <a href="docs/contributing/CONTRIBUTING-UA.md">Українська</a> · 
  <a href="docs/contributing/CONTRIBUTING-CS.md">Česky</a> · 
  <a href="docs/contributing/CONTRIBUTING-HU.md">Magyar</a> · 
  <a href="docs/contributing/CONTRIBUTING-ES.md">Español</a> · 
  <a href="docs/contributing/CONTRIBUTING-FA.md">فارسی</a> · 
  <a href="docs/contributing/CONTRIBUTING-FR.md">Français</a> · 
  <a href="docs/contributing/CONTRIBUTING-DE.md">Deutsch</a> · 
  <a href="docs/contributing/CONTRIBUTING-PL.md">Polski</a> · 
  <a href="docs/contributing/CONTRIBUTING-ID.md">Indonesian</a> · 
  <a href="docs/contributing/CONTRIBUTING-FI.md">Suomi</a> · 
  <a href="docs/contributing/CONTRIBUTING-ML.md">മലയാളം</a> · 
  <a href="docs/contributing/CONTRIBUTING-JP.md">日本語</a> · 
  <a href="docs/contributing/CONTRIBUTING-NL.md">Nederlands</a> · 
  <a href="docs/contributing/CONTRIBUTING-IT.md">Italiano</a> · 
  <a href="docs/contributing/CONTRIBUTING-RU.md">Русский</a> · 
  <a href="docs/contributing/CONTRIBUTING-PTBR.md">Português (Brasil)</a> · 
  <a href="docs/contributing/CONTRIBUTING-EO.md">Esperanto</a> · 
  <a href="docs/contributing/CONTRIBUTING-KR.md">한국어</a> · 
  <a href="docs/contributing/CONTRIBUTING-AR.md">العربي</a> · 
  <a href="docs/contributing/CONTRIBUTING-VN.md">Tiếng Việt</a> · 
  <a href="docs/contributing/CONTRIBUTING-DA.md">Dansk</a> · 
  <a href="docs/contributing/CONTRIBUTING-GR.md">Ελληνικά</a> · 
  <a href="docs/contributing/CONTRIBUTING-TR.md">Türkçe</a>
</p>

</div>

</p>

So, you want to hack on open-im-server? Yay!

First of all, thank you for considering contributing to our project! We appreciate your time and effort, and we value any contribution, whether it's reporting a bug, suggesting a new feature, or submitting a pull request.

![Hello OpenIM Image](assets/demo/hello-openim.png)

> Use `make demo` start contributing fast.

This document provides guidelines and best practices to help you contribute effectively.

## 📇Topics

- [How do I contribute code to OpenIM](#how-do-i-contribute-code-to-openim)
  - [📇Topics](#topics)
  - [What we expect of you](#what-we-expect-of-you)
  - [Code of ConductCode of Conduct](#code-of-conductcode-of-conduct)
      - [Code and doc contribution](#code-and-doc-contribution)
      - [Where should I start?](#where-should-i-start)
      - [Design documents](#design-documents)
  - [Getting Started](#getting-started)
  - [Style and Specification](#style-and-specification)
      - [Reporting security issues](#reporting-security-issues)
      - [Reporting general issues](#reporting-general-issues)
      - [Commit Rules](#commit-rules)
      - [PR Description](#pr-description)
      - [Docs Contribution](#docs-contribution)
  - [Engage to help anything](#engage-to-help-anything)
  - [Release version](#release-version)
  - [Contact Us](#contact-us)

## What we expect of you

We hope that anyone can join open-im-server , even if you are a student, writer, translator

Please meet the minimum version of the Go language published in [go.mod](./go.mod). If you want to manage the Go language version, we provide tools tHow do I contribute code to OpenIMo install [gvm](https://github.com/moovweb/gvm) in our [Makefile](./Makefile)

You'd better use Linux OR WSL as the development environment, Linux with [Makefile](./Makefile) can help you quickly build and test open-im-server project.

If you are familiar with [Makefile](./Makefile) , you can easily see the clever design of the open-im-server Makefile. Storing the necessary tools such as golangci in the `/tools` directory can avoid some tool version issues.

The [Makefile](./Makefile) is for every developer, even if you don't know how to use the Makefile tool, don't worry, we provide two great commands to get you up to speed with the Makefile architecture, `make help` and `make help-all`, it can reduce problems of the developing environment.

In accordance with the naming conventions adopted by OpenIM and drawing reference from the Google Naming Conventions as per the guidelines available at https://google.github.io/styleguide/go/, the following expectations for naming practices within the project are set forth: 

+ https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/code-conventions.md
+ https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/go-code.md


## Code of ConductCode of Conduct

#### Code and doc contribution

Every action to make project open-im-server better is encouraged. On GitHub, every improvement for open-im-server could be via a [PR](https://github.com/openimsdk/open-im-server/pulls) (short for pull request).

+ If you find a typo, try to fix it!
+ If you find a bug, try to fix it!
+ If you find some redundant codes, try to remove them!
+ If you find some test cases missing, try to add them!
+ If you could enhance a feature, please **DO NOT** hesitate!
+ If you find code implicit, try to add comments to make it clear!
+ If you find code ugly, try to refactor that!
+ If you can help to improve documents, it could not be better!
+ If you find document incorrect, just do it and fix that!
+ ...

#### Where should I start?

+ If you are new to the project, don't know how to contribute open-im-server, please check out the [good first issue](https://github.com/openimsdk/open-im-server/issues?q=is%3Aopen+label%3A"good+first+issue"+sort%3Aupdated-desc) label.
+ You should be good at filtering the open-im-server issue tags and finding the ones you like, such as [RFC](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+RFC+label%3ARFC) for big initiatives, features for [feature](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+label%3Afeature) proposals, and [bug](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+label%3Abug+) fixes.
+ If you are looking for something to work on, check out our [open issues](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc).
+ If you have an idea for a new feature, please [open an issue](https://github.com/openimsdk/open-im-server/issues/new/choose), and we can discuss it.

#### Design documents

For any substantial design, there should be a well-crafted design document. This document is not just a simple record, but also a detailed description and manifestation, which can help team members better understand the design thinking and grasp the design direction. In the process of writing the design document, we can choose to use tools such as `Google Docs` or `Notion`, and even mark RFC in [issues](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+RFC+label%3ARFC) or [discussions](https://github.com/openimsdk/open-im-server/discussions) for better collaboration. Of course, after completing the design document, we should also add it to our [Shared Drive](https://drive.google.com/drive/) and notify the appropriate working group to let everyone know of its existence. Only by doing so can we maximize the effectiveness of the design document and provide strong support for the smooth progress of the project.

Anybody can access the shared Drive for reading. To get access to comment. Once you've done that, head to the [shared Drive](https://drive.google.com/) and behold all the docs.

In addition to that, we'd love to invite you to [join our Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q) where you can play with your imagination, tell us what you're working on, and get a quick response.

When documenting a new design, we recommend a 2-step approach:

1. Use the short-form RFC template to outline your ideas and get early feedback.
2. Once you have received sufficient feedback and consensus, you may use the longer-form design doc template to specify and discuss your design in more details.

In order to contribute a feature to open-im-server you'll need to go through the following steps:

+ Discuss your idea with the appropriate [working groups](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q) on the working group's Slack channel.
+ Once there is general agreement that the feature is useful, create a GitHub issue to track the discussion. The issue should include information about the requirements and use cases that it is trying to address.
+ Include a discussion of the proposed design and technical details of the implementation in the issue.

But keep in mind that there is no guarantee of it being accepted and so it is usually best to get agreement on the idea/design before time is spent coding it. However, sometimes seeing the exact code change can help focus discussions, so the choice is up to you.

## Getting Started

To propose PR for the open-im-server item, we assume you have registered a GitHub ID. Then you could finish the preparation in the following steps:

1. Fork the repository(open-im-server)

2. **CLONE** your own repository to main locally. Use `git clone https://github.com/<your-username>/open-im-server.git` to clone repository to your local machine. Then you can create new branches to finish the change you wish to make.

3. **Initialize Git Hooks with `make init-githooks`**

   After cloning the repository, it's recommended to set up Git hooks to streamline your workflow and ensure your contributions adhere to OpenIM's community standards. Git hooks are scripts that run automatically every time a particular event occurs in a Git repository, such as before a commit or push. To initialize Git hooks for the OpenIM server repository, use the `make init-githooks` command.

   - **Enabling Git Hooks Mode**: By running `make init-githooks` and entering `1` when prompted, you enable Git hooks mode. This action will generate a series of hooks within the `.git/hooks/` directory. These hooks impose certain checks on your commits and pushes, ensuring they meet the quality and standards expected by the OpenIM community. For instance, commit hooks might enforce a specific commit message format, while push hooks could check for code style or linting issues.

   - **Benefits for First-Time Contributors**: This setup is especially beneficial for new contributors. It guides you to make professional, community-standard-compliant Pull Requests (PRs) and PR descriptions right from your first contribution. By automating checks and balances, it reduces the chances of common mistakes and speeds up the review process.

   - **Disabling Git Hooks**: If for any reason you wish to remove the Git hooks, simply run `make init-githooks` again and enter `2` when prompted. This will delete the existing Git hooks, removing the automatic checks and constraints from your Git operations. However, keep in mind that manually ensuring your contributions adhere to community standards without the aid of Git hooks requires diligence.
 
   > [!NOTE] Utilizing Git hooks through the `make init-githooks` command is a straightforward yet powerful way to ensure your contributions are consistent and high-quality. It's a step towards fostering a professional and efficient development environment in the OpenIM project.



4. **Set Remote** upstream to be `https://github.com/openimsdk/open-im-server.git` using the following two commands:

   ```bash
   ❯ git remote add upstream https://github.com/openimsdk/open-im-server.git
   ❯ git remote set-url --push upstream no-pushing
   ```

   With this remote setting, you can check your git remote configuration like this:

   ```bash
   ❯ git remote -v
   origin     https://github.com/<your-username>/open-im-server.git (fetch)
   origin     https://github.com/<your-username>/open-im-server.git (push)
   upstream   https://github.com/openimsdk/open-im-server.git (fetch)
   upstream   no-pushing (push)
   ```

   Adding this, we can easily synchronize local branches with upstream branches.

5. Create a new branch for your changes (use a descriptive name, such as `fix-bug-123` or `add-new-feature`).

   ```bash
   ❯ cd open-im-server
   ❯ git fetch upstream
   ❯ git checkout upstream/main
   ```

   Create a new branch: 

   ```bash
   ❯ git checkout -b <new-branch>
   ```

   Make any change on the `new-branch` then use [Makefile](./Makefile) build and test your codes.


6. **Commit your changes** to your local branch, lint before committing and commit with sign-off

   ```bash
   ❯ git rebase upstream/main
   ❯ make lint	  # golangci-lint run -c .golangci.yml
   ❯ git add -A  # add changes to staging
   ❯ git commit -a -s -m "message for your changes" # -s adds a Signed-off-by trailer
   ```

7. **Push your branch**  to your forked repository, it is recommended to have only one commit for a PR.

   ```bash
   # sync up with upstream
   ❯ git fetch upstream main
   ❯ git rebase upstream/main
   ❯ 
   ❯ git rebase -i	<commit-id> # rebase with interactive mode to squash your commits into a single one
   ❯ git push # push to the remote repository, if it's a first time push, run git push --set-upstream origin <new-branch># sync up with upstream
   ❯ git fetch upstream main
   git rebase upstream/main
   
   ❯ git rebase -i	<commit-id> # rebase with interactive mode to squash your commits into a single one
   ❯ git push # push to the remote repository, if it's a first time push, run git push --set-upstream origin <new-branch>
   ```

   You can also use `git commit -s --amend && git push -f` to update modifications on the previous commit.

   If you have developed multiple features in the same branch, you should create PR separately by rebasing to the main branch between each push:

   ```bash
   # create new branch, for example git checkout -b feature/infra
   ❯ git checkout -b <new branch>
   # update some code, feature1
   ❯ git add -A
   ❯ git commit -m -s "feat: feature one"
   ❯ git push # if it's first time push, run git push --set-upstream origin <new-branch>
   # then create pull request, and merge
   # update some new feature, feature2, rebase main branch first.
   ❯ git rebase upstream/main # rebase the current branch to upstream/main branch
   ❯ git add -A
   ❯ git commit -m -s "feat: feature two"
   ```

    **Verifying Your Pull Request with `make all` Command**

   Before verifying, you may need to complete the basic deployment of OpenIM to get familiar with the deployment status of OpenIM. Please read [this deployment document](https://docs.openim.io/zh-Hans/guides/gettingStarted/imSourceCodeDeployment), which will tell you how to deploy OpenIM middleware and OpenIM services in detail

   Before submitting your Pull Request (PR), it's crucial to ensure that it passes all the necessary checks and verifications to maintain the quality and integrity of the OpenIM server project. To facilitate this process, we have encapsulated a series of validation steps into the `make all` command.

   - **Purpose of `make all` Command**: The `make all` command serves as a comprehensive pre-PR verification tool. It sequentially executes a variety of tasks designed to scrutinize your changes from multiple angles, ensuring they are ready for submission.

   - **Included Commands**:
     - `tidy`: Cleans up the module by removing unused dependencies.
     - `gen`: Generates necessary files from templates or specifications, ensuring that your codebase is up-to-date with its dependencies.
     - `add-copyright`: Checks for and adds copyright notices to files, ensuring compliance with legal requirements.
     - `verify`: Verifies the integrity and consistency of the code, dependencies, and various checks.
     - `test-api`: Runs API tests to ensure that your changes do not break any existing functionality and adhere to the expected behaviors.
     - `lint`: Analyzes the code for potential stylistic or programming errors, enforcing the project's coding standards.
     - `cover`: Measures the code coverage of tests, helping you understand how much of the code is being tested.
     - `restart`: (Optionally) restarts services or applications to ensure that changes are correctly applied and functioning in a live environment.

   - **Executing the Command**: To run the `make all` command, simply navigate to the root directory of your cloned repository in your terminal and execute:
     ```bash
     make all
     ```
     This command will sequentially perform all the listed actions, outputting any warnings or errors encountered during the process. It's a vital step to catch any issues early and ensure your contribution meets the quality standards set by the OpenIM community.

   - **Benefits**: By using `make all` for pre-PR verification, you significantly increase the likelihood of your PR being accepted on the first review. It not only demonstrates your commitment to quality but also streamlines the review process by minimizing back-and-forth due to common issues that can be caught automatically.


    **Troubleshooting Git Push Failures**

    When working with Git, encountering errors during push operations is not uncommon. Two primary reasons you might face push failures are due to firewall restrictions or authentication issues. Here’s how you can troubleshoot and resolve these problems.

    **Firewall Errors**

    If you're behind a corporate firewall or your network restricts certain types of traffic, you might encounter issues when trying to push your changes via HTTPS. This is because firewalls can block the ports used by the HTTPS protocol. To resolve this issue, you can configure Git to use a proxy.

    If you have a local proxy server set up, you can direct Git to use it by setting the `https_proxy` and `http_proxy` environment variables. Open your terminal or command prompt and run the following commands:

    ```bash
    export https_proxy="http://127.0.0.1:7890"
    export http_proxy="http://127.0.0.1:7890"
    ```

    Replace `127.0.0.1:7890` with the address and port of your proxy server. These commands set the proxy for the current session. If you want to make these changes permanent, add them to your `.bashrc`, `.bash_profile`, or equivalent shell configuration file.

    **Using SSH Instead of HTTPS**

    An alternative to using HTTPS is to set up an SSH connection for Git operations. SSH connections are often not blocked by firewalls and do not require proxy settings. Additionally, SSH provides a secure channel and can simplify the authentication process since it relies on SSH keys rather than username and password credentials.

    To use SSH with Git, you first need to generate an SSH key pair and add the public key to your GitHub account (or another Git hosting service).

    1. **Generate SSH Key Pair**: Open your terminal and run `ssh-keygen -t rsa -b 4096 -C "your_email@example.com"`, replacing `your_email@example.com` with your email. Press enter to accept the default file location and passphrase prompts.

    2. **Add SSH Key to SSH-Agent**: Ensure the ssh-agent is running with `eval "$(ssh-agent -s)"` and then add your SSH private key to the ssh-agent using `ssh-add ~/.ssh/id_rsa`.

    3. **Add SSH Key to GitHub**: Copy your SSH public key to your clipboard with `cat ~/.ssh/id_rsa.pub | clip` (Windows) or `pbcopy < ~/.ssh/id_rsa.pub` (Mac). Go to GitHub, navigate to Settings > SSH and GPG keys, and add a new SSH key, pasting your key into the field provided.

    4. **Switch to SSH in Your Repository**: Change your repository's remote URL from HTTPS to SSH. You can find the SSH URL in your repository settings on GitHub and use `git remote set-url origin git@github.com:username/repository.git` to switch.

    **Authentication Errors**

    If you're experiencing authentication errors, it might be due to missing or incorrect credentials. Ensure you have added your SSH key to your Git hosting service. You can test your SSH connection with `ssh -T git@github.com` (replace `github.com` with your Git hosting service's domain). If successful, you'll receive a welcome message.

    For HTTPS users, check that your username and password (or personal access token for services like GitHub that no longer accept password authentication for Git operations) are correct.

8. **Open a pull request** to `openimsdk/open-im-server:main`

   It is recommended to review your changes before filing a pull request. Check if your code doesn't conflict with the main branch and no redundant code is included.
   
   > [!TIP] There is a [good blog post documenting](https://nsddd.top/posts/participating-in-this-project/) the entire push contribution process.


## Style and Specification

We divide the problem into security and general problems:

#### Reporting security issues

Security issues are always treated seriously. As our usual principle, we discourage anyone to spread security issues. If you find a security issue of open-im-server, please do not discuss it in public and even do not open a public issue.

Instead we encourage you to send us a private email to info@openim.io to report this.

#### Reporting general issues

To be honest, we regard every user of open-im-serveras a very kind contributor. After experiencing open-im-server, you may have some feedback for the project. Then feel free to open an issue via [NEW ISSUE](https://github.com/openimsdk/open-im-server/issues/new/choose).

Since we collaborate project open-im-server in a distributed way, we appreciate **WELL-WRITTEN**, **DETAILED**, **EXPLICIT** issue reports. To make the communication more efficient, we wish everyone could search if your issue is an existing one in the searching list. If you find it existing, please add your details in comments under the existing issue instead of opening a brand new one.

To make the issue details as standard as possible, we setup an [ISSUE TEMPLATE](https://github.com/OpenIMSDK/.github/tree/main/.github/ISSUE_TEMPLATE) for issue reporters. You can find three kinds of issue templates there: question, bug report and feature request. Please **BE SURE** to follow the instructions to fill fields in template.

**There are a lot of cases when you could open an issue:**

+ bug report
+ feature request
+ open-im-server performance issues 
+ feature proposal
+ feature design
+ help wanted
+ doc incomplete
+ test improvement
+ any questions on open-im-server project 
+ and so on 

Also, we must be reminded when submitting a new question about open-im-server, please remember to remove the sensitive data from your post. Sensitive data could be password, secret key, network locations, private business data and so on.

#### Commit Rules

Actually in open-im-server, we take two rules serious when committing:

**🥇 Commit Message:**

Commit message could help reviewers better understand what the purpose of submitted PR is. It could help accelerate the code review procedure as well. We encourage contributors to use **EXPLICIT** commit message rather than ambiguous message. In general, we advocate the following commit message type:

We use [Semantic Commits](https://www.conventionalcommits.org/en/v1.0.0/) to make it easier to understand what a commit does and to build pretty changelogs. Please use the following prefixes for your commits:

+ `docs: xxxx`. For example, "docs: add docs about storage installation".
+ `feature: xxxx`.For example, "feature: make result show in sorted order".
+ `bugfix: xxxx`. For example, "bugfix: fix panic when input nil parameter".
+ `style: xxxx`. For example, "style: format the code style of Constants.java".
+ `refactor: xxxx.` For example, "refactor: simplify to make codes more readable".
+ `test: xxx`. For example, "test: add unit test case for func InsertIntoArray".
+ `chore: xxx.` For example, "chore: integrate travis-ci". It's the type of mantainance change.
+ other readable and explicit expression ways.

On the other side, we discourage contributors from committing message like the following ways:

+ ~~fix bug~~
+ ~~update~~
+ ~~add doc~~

**🥈 Commit Content:**

Commit content represents all content changes included in one commit. We had better include things in one single commit which could support reviewer's complete review without any other commits' help.

In another word, contents in one single commit can pass the CI to avoid code mess. In brief, there are two minor rules for us to keep in mind:

1. avoid very large change in a commit.
2. complete and reviewable for each commit.
3. words are written in lowercase English, not uppercase English or other languages such as Chinese.

No matter what the commit message, or commit content is, we do take more emphasis on code review.

An example for this could be:

```bash
❯ git commit -a -s -m "docs: add a new section to the README"
```

#### PR Description

PR is the only way to make change to open-im-server project files. To help reviewers better get your purpose, PR description could not be too detailed. We encourage contributors to follow the [PR template](https://github.com/OpenIMSDK/.github/tree/main/.github/PULL_REQUEST_TEMPLATE.md) to finish the pull request.

You can find some very formal PR in [RFC](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+RFC+label%3ARFC) issues and learn about them.

**📖 Opening PRs:**

+ As long as you are working on your PR, please mark it as a draft
+ Please make sure that your PR is up-to-date with the latest changes in `main`
+ Mention the issue that your PR is addressing (Fix: #{ID_1}, #{ID_2})
+ Make sure that your PR passes all checks

**🈴 Reviewing PRs:**

+ Be respectful and constructive 
+ Assign yourself to the PR (comment `/assign`)
+ Check if all checks are passing
+ Suggest changes instead of simply commenting on found issues
+ If you are unsure about something, ask the author
+ If you are not sure if the changes work, try them out
+ Reach out to other reviewers if you are unsure about something
+ If you are happy with the changes, approve the PR
+ Merge the PR once it has all approvals and the checks are passing

**⚠️ DCO check:**

We have a DCO check that runs on every pull request to ensure code quality and maintainability. This check verifies that the commit has been signed off, indicating that you have read and agreed to the provisions of the Developer Certificate of Origin. If you have not yet signed off on the commit, you can use the following command to sign off on the last commit you made:

```bash
❯ git commit --amend --signoff
```

Please note that signing off on a commit is a commitment that you have read and agreed to the provisions of the Developer Certificate of Origin. If you have not yet read this document, we strongly recommend that you take some time to read it carefully. If you have any questions about the content of this document, or if you need further assistance, please contact an administrator or relevant personnel.

You can also automate signing off your commits by adding the following to your `.zshrc` or `.bashrc`:

```go
git() {
  if [ $# -gt 0 ] && [[ "$1" == "commit" ]] ; then
     shift
     command git commit --signoff "$@"
  else
     command git "$@"
  fi
}
```


#### Docs Contribution

The documentation for open-im-server includes:

+ [README.md](https://github.com/openimsdk/open-im-server/blob/main/README.md): This file includes the basic information and instructions for getting started with open-im-server.
+ [CONTRIBUTING.md](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md): This file contains guidelines for contributing to open-im-server's codebase, such as how to submit issues, pull requests, and code reviews.
+ [Official Documentation](https://doc.rentsoft.cn/): This is the official documentation for open-im-server, which includes comprehensive information on all of its features, configuration options, and troubleshooting tips.

Please obey the following rules to better format the docs, which would greatly improve the reading experience.

1. Please do not use Chinese punctuations in English docs, and vice versa.
2. Please use upper case letters where applicable, like the first letter of sentences / headings, etc.
3. Please specify a language for each Markdown code blocks, unless there's no associated languages.
4. Please insert a whitespace between Chinese and English words.
5. Please use the correct case for technical terms, such as using `HTTP` instead of http, `MySQL` rather than mysql, `Kubernetes` instead of kubernetes, etc.
6. Please check if there's any typos in the docs before submitting PRs.

## Engage to help anything

We choose GitHub as the primary place for open-im-server to collaborate. So the latest updates of open-im-server are always here. Although contributions via PR is an explicit way to help, we still call for any other ways.

+ reply to other's [issues](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc) if you could;
+ help solve other user's problems;
+ help review other's [PR](https://github.com/openimsdk/open-im-server/pulls?q=is%3Apr+is%3Aopen+sort%3Aupdated-desc) design; 
+ discuss about open-im-server to make things clearer;
+ advocate [open-im-server](https://google.com/search?q=open-im-server) technology beyond GitHub;
+ write blogs on open-im-server and so on.

In a word, **ANY HELP IS CONTRIBUTION.**

## Release version

Releases of open-im-server are done using [Release Please](https://github.com/googleapis/release-please) and [GoReleaser](https://goreleaser.com/). The workflow looks like this:

🎯 A PR is merged to the `main` branch:

+ Release please is triggered, creates or updates a new release PR
+ This is done with every merge to main, the current release PR is updated every time

🎯 Merging the 'release please' PR to `main`:

+ Release please is triggered, creates a new release and updates the changelog based on the commit messages
+ GoReleaser is triggered, builds the binaries and attaches them to the release
+ Containers are created and pushed to the container registry

With the next relevant merge, a new release PR will be created and the process starts again

**👀 Manually setting the version:**

If you want to manually set the version, you can create a PR with an empty commit message that contains the version number in the commit message. For example:

Such a commit can get produced as follows: 

````bash
❯ git commit --allow-empty -m "chore: release 0.0.3" -m "Release-As: 0.0.3
````

For the complex release process, in fact, and encapsulation as CICD, you only need to tag locally, and then publish the tag to OpenIM github to complete the entire OpenIM release process.

In addition to CICD, we also do a complex release command locally, which can help you complete the full platform compilation, testing, and release to Minio, just by using the `make release` command.
Please [read the detailed documents](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/release.md)


## Contact Us

We value close connections with our users, developers, and contributors here at open-im-server. With a large community and maintainer team, we're always here to help and support you. Whether you're looking to join our community or have any questions or suggestions, we welcome you to get in touch with us.

Our most recommended way to get in touch is through [Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q). Even if you're in China, Slack is usually not blocked by firewalls, making it an easy way to connect with us. Our Slack community is the ideal place to discuss and share ideas and suggestions with other users and developers of open-im-server. You can ask technical questions, seek help, or share your experiences with other users of open-im-server.

In addition to Slack, we also offer the following ways to get in touch:

+ <a href="https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q" target="_blank"><img src="https://img.shields.io/badge/slack-%40OpenIMSDKCore-informational?logo=slack&style=flat-square"></a>:  We also have Slack channels for you to communicate and discuss. To join, visit https://slack.com/ and join our [👀 open-im-server slack](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q) team channel.
+ <a href="https://mail.google.com/mail/u/0/?fs=1&tf=cm&to=4closetool3@gmail.com" target="_blank"><img src="https://img.shields.io/badge/gmail-%40OOpenIMSDKCore?style=social&logo=gmail"></a>: Get in touch with us on [Gmail](info@openim.io). If you have any questions or issues that need resolving, or any suggestions and feedback for our open source projects, please feel free to contact us via email.
+ <a href="https://doc.rentsoft.cn/" target="_blank"><img src="https://img.shields.io/badge/%E5%8D%9A%E5%AE%A2-%40OpenIMSDKCore-blue?style=social&logo=Octopus%20Deploy"></a>: Read our [blog](https://doc.rentsoft.cn/). Our blog is a great place to stay up-to-date with open-im-server projects and trends. On the blog, we share our latest developments, tech trends, and other interesting information.
+ <a href="https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg" target="_blank"><img src="https://img.shields.io/badge/%E5%BE%AE%E4%BF%A1-OpenIMSDKCore-brightgreen?logo=wechat&style=flat-square"></a>: Add [Wechat](https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg) and indicate that you are a user or developer of open-im-server. We will process your request as soon as possible.

Whether you're looking to join our community or have any questions or suggestions, we welcome you to get in touch with us.

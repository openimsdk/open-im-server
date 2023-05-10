#  参与Open-IM-Server的建设

怎么样,你想要为Open-IM-Server做出贡献吗?好耶!
首先,感谢你考虑为我们的项目做出贡献!我们十分感激你的时间和精力,我们重视任何形式的贡献,无论是报告错误、对新功能有所建议,或提交PR。
本文档提供指导和最佳的实践教程,帮助您有效地做出贡献。

## 📇标题

- [我们期望你做什么](#我们期望你做什么)
- [行为准则](#行为准则)
- [上手IM](#上手IM)
- [风格和规范](#风格和规范)
- [参与任何形式的帮助](#参与任何形式的帮助)
- [发布版本](#发布版本)
- [联系我们](#联系我们)

## 我们期望你做什么

我们希望任何人都可以加入Open-IM-Server,即使您是一名学生、作家或翻译。 

请满足 [go.mod](./go.mod) 中发布的Go语言的最低版本要求。如果您想管理Go语言版本,我们在 [Makefile](./Makefile) 中提供安装 [gvm](https://github.com/moovweb/gvm) 的工具。

最好使用Linux OR WSL作为开发环境，Linux配合 [Makefile](./Makefile) 可以帮助您快速构建和测试Open-IM-Server项目。  

如果您熟悉 [Makefile](./Makefile) , 您可以轻松看到Open-IM-Server Makefile的巧妙设计。将必要的工具(如golangci)存储在`/tools`目录中可以避免一些工具的版本问题。 

[Makefile](./Makefile) 适用于每个开发者，即使您不知道如何使用Makefile工具,不要担心,我们提供了两个很棒的命令来让您熟悉Makefile架构:`make help`和`make help-all`, 它可以减少开发环境的问题。 

## 行为准则

#### 代码与文档贡献

我们鼓励每一件让Open-IM-Server更好的行动。在Github上，每一个对Open-IM-Server的改进都可以通过 [PR](https://github.com/Open-IM-Server/pulls) 来改进。

+ 如果发现一个错误,尝试修复它!
+ 如果发现一个bug,尝试修复它!
+ 如果发现一些冗余的代码,尝试删除它们!
+ 如果发现缺少一些测试用例,尝试添加它们! 
+ 如果可以增强一个功能,请不要犹豫!
+ 如果发现代码隐晦,尝试添加注释使其清晰!
+ 如果发现丑陋的代码,尝试重构它!
+ 如果可以帮助改进文档,那就更好不过了! 
+ 如果发现文档不正确,请修复它!
+ ...

#### 我应该从哪里开始

+ 如果您刚开始接触该项目,不知道如何为Open-IM-Server做出贡献,请查看 [good first issue](https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aopen+label%3A"good+first+issue"+sort%3Aupdated-desc) 标签.
+ 您应该对过滤Open-IM-Server问题标签很熟练,找到您喜欢的标签,例如 [RFC](https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aissue+is%3Aopen+RFC+label%3ARFC) 用于大型计划与功能的 [feature](https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aissue+label%3Afeature) 提案, 以及 [bug](https://github.com/{github/issues?q=is%3Aissue+label%3Abug+) 修复。
+ 如果您正在寻找可以工作的内容,请查看我们的 [open issues](https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc).
+ 如果您有一个新的功能想法,请 [开启一个 issue](https://github.com/OpenIMSDK/Open-IM-Server/issues/new/choose), 我们可以讨论它

#### 设计文档

对于任何重大设计,都应该有一份精心制作的设计文档。这份文档不仅仅是一份简单的记录,也是一份详细的描述和阐释,它可以帮助团队成员更好地理解设计思路和掌握设计方向。在编写设计文档的过程中,我们可以选择使用诸如 `Google Docs` 或者 `Notion`, 甚至可以在 [issues](https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aissue+is%3Aopen+RFC+label%3ARFC) 或 [discussions](https://github.com/OpenIMSDK/Open-IM-Server/discussions) 标记RFC以进行更好的协作。当然,完成设计文档后,我们还应将其添加到我们的 [Shared Drive](https://drive.google.com/drive/) 中并通知相应的工作组,让所有人都知道它的存在。只有这样,我们才能最大限度地发挥设计文档的效果,为项目的顺利进展提供强有力的支持。

任何人都可以访问共享驱动器进行读取。要获得评论权限,请完成注册。完成注册后,前往[shared Drive](https://drive.google.com/) 查看所有文档。

除此之外,我们也热切邀请您 [加入我们的 Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg) 在那里您可以发挥想象力,告诉我们您正在做什么,并快速得到回复。

在记录新的设计时,我们建议采取两步方法:
1. 使用简短的RFC模板概述您的想法，并尽早获得反馈。 
2. 一旦您获得足够的反馈和共识,您可以使用更长的设计文档模板更详细地指定和讨论您的设计。

为了向Open-IM-Server贡献功能,您需要遵循以下步骤:

+  与适当的 [工作组](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg) 在Slack频道上讨论您的想法。
+ 一旦普遍认为该功能是有用的,请在GitHub上创建议题来跟踪讨论。该议题应包括该功能试图解决的需求和使用案例的信息。
+ 在议题中包括拟议设计和实现的技术细节的讨论。

但请记住,提交PR并不能保证它被接受,所以通常最好在花时间编码之前就与我们达成对该想法/设计的认同。然而,有时看到确切的代码变更对我们的集体讨论有帮助,所以选择权在您。 

## 上手IM

为了向Open-IM-Server提出Pull Request,我们假设您已注册GitHub ID。然后您可以按以下步骤完成准备工作:
1. Fork该仓库(Open-IM-Server) 
2. **CLONE**您自己的仓库到本地主分支。使用 `git clone https://github.com/<your-username>/Open-IM-Server.git`将仓库克隆到本地机器。然后您可以创建新分支来完成您希望进行的更改。

3. **Set Remote** 上游为 `https://github.com/OpenIMSDK/Open-IM-Server.git` 使用以下两个命令:

   ```bash
   ❯ git remote add upstream https://github.com/OpenIMSDK/Open-IM-Server.git
   ❯ git remote set-url --push upstream no-pushing
   ```

   通过这样的远程设置，您可以这样检查您的git远程配置

   ```bash
   ❯ git remote -v
   origin     https://github.com/<your-username>/Open-IM-Server.git (fetch)
   origin     https://github.com/<your-username>/Open-IM-Server.git (push)
   upstream   https://github.com/OpenIMSDK/Open-IM-Server.git (fetch)
   upstream   no-pushing (push)
   ```

   添加此内容后,我们可以轻松地将本地分支与上游分支同步。

4. 为您的更改创建新分支(使用描述性名称,如`fix-bug-123`或`add-new-feature`)。

   ```bash
   ❯ cd Open-IM-Server
   ❯ git fetch upstream
   ❯ git checkout upstream/main
   ```

   创建新分支:

   ```bash
   ❯ git checkout -b <new-branch>
   ```

   在`new-branch`上进行任何更改,然后使用[Makefile](./Makefile)构建和测试您的代码。
   
5. **Commit your changes** 到本地分支,提交之前进行Lint,并签名提交

   ```bash
   ❯ git rebase upstream/main
   ❯ make link	  # golangci-lint run -c .golangci.yml
   ❯ git add -A  # add changes to staging
   ❯ git commit -a -s -m "message for your changes" # -s adds a Signed-off-by trailer
   ```

6. **Push your branch**  推送到您的Fork仓库,只提交一个PR。

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

  您也可以使用 `git commit -s --amend && git push -f` 更新前一次提交中的修改。

  如果您在同一分支中开发了多个功能，则应在每次推送之间对主分支进行rebase，并单独创建PR

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
   # then create pull request, and merge
   ```

7. **Open a pull request** 到 `OpenIMSDK/Open-IM-Server:main`

  建议在提交拉取请求之前检查您的更改。检查您的代码是否与主分支冲突,并且未包含重复的代码。

## 风格和规范

我们将问题划分为安全问题和一般问题:

#### 报告安全问题

安全问题一向被我们严肃对待。按照我们的惯例原则,我们不赞成任何人散布安全问题。如果您发现Open-IM-Server的安全问题,请不要在公开场合讨论,甚至不要公开问题。
如果遇到安全问题,我们鼓励您发送私人电子邮件到winxu81@gmail.com报告此问题。

#### 报告一般问题

老实说,我们认为Open-IM-Server的每一个用户都是非常善良的贡献者。体验Open-IM-Server之后,您可能会对该项目有一些反馈。那么请通过 [NEW ISSUE](https://github.com/OpenIMSDK/Open-IM-Server/issues/new/choose)开启问题。

由于我们以分布式的方式协作Open-IM-Server项目,我们非常感激**写得好**,**详细**,**明确**的问题报告。为了使交流更加高效,我们希望每个人都能在搜索列表中搜索您的问题是否已存在。如果您找到现有问题,请在现有问题的评论下添加您的详细信息,而不是打开一个全新的问题。

为了使问题详情尽可能标准化,我们为问题报告者设置了 [ISSUE模板](https://github.com/OpenIMSDK/Open-IM-Server/tree/main/.github/ISSUE_TEMPLATE)。您可以在那里找到三种问题模板:问题、bug报告和功能请求。请**确保**遵循说明填写模板中的字段。

**您可以在许多情况下打开问题:**
+ bug报告 
+ 功能请求 
+ Open-IM-Server性能问题 
+ 功能建议 
+ 功能设计 
+ 需要帮助
+ 文档不完整 
+ 测试改进 
+ 关于Open-IM-Server项目的任何问题
+ 等等 

另外,在提交关于Open-IM-Server的新问题时,请记住从您的帖子中删除敏感数据。敏感数据可以是密码、密钥、网络位置、私人商业数据等。

#### 提交规则

实际上在Open-IM-Server中,我们在提交时遵循两条规则:

**🥇 提交消息:**

提交消息可以帮助评论者更好地理解提交的PR的目的。它还可以帮助加快代码审查过程。我们鼓励贡献者使用**明确**的提交消息而不是模糊的消息。一般来说,我们提倡以下提交消息类型:

我们使用 [语义提交](https://www.conventionalcommits.org/en/v1.0.0/) 使得更易于理解提交做了什么以及生成漂亮的更新日志。请为您的提交使用以下前缀:

+ `docs: xxxx`. 例如,"docs:添加存储安装文档"。
+ `feature: xxxx`.例如,"feature:使结果以排序顺序显示"。 
+ `bugfix: xxxx`. 例如,"bugfix:修复输入空参数时的panic"。
+ `style: xxxx`. 例如,"style:格式化Constants.java的代码样式"。
+ `refactor: xxxx.` 例如,"refactor:简化以使代码更易读"。 
+ `test: xxx`. 例如,"test:为func InsertIntoArray添加单元测试用例"。
+ `chore: xxx.` 例如,"chore:集成travis-ci"。这是维护更改的类型。
+ 其他可读且明确的表达方式。

另一方面,我们不鼓励贡献者以以下方式提交消息:

+ ~~fix bug~~
+ ~~update~~
+ ~~add doc~~

**🥈 提交内容:**

提交内容表示一个提交中包含的所有内容更改。我们最好在一个单独的提交中包含能够支持评论者完整审查，而无需任何其他提交帮助的内容。

换句话说,一个单独提交中的内容可以通过CI以避免代码混乱。简而言之,我们要记住的两个小规则是:

1. 避免在一次提交中进行非常大的更改。 
2. 每次提交都要完整且可审查。 
3. 字词使用小写英语,而不是大写英语或其他语言(如中文)。
无论提交消息还是提交内容如何,我们都更加重视代码审查。

举个例子:

```bash
❯ git commit -a -s -m "docs: add a new section to the README"
```

#### PR描述

PR是对Open-IM-Server项目文件进行更改的唯一方式。为了帮助评论者更好地了解您的目的,PR描述不应过于详细。我们鼓励贡献者遵循 [PR模板](https://github.com/OpenIMSDK/Open-IM-Server/tree/main/.github/PULL_REQUEST_TEMPLATE.md) 来拉取请求。

您可以在 [RFC](https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aissue+is%3Aopen+RFC+label%3ARFC) 问题中找到一些非常正式的PR,并了解它们。

**📖 打开PR:**

+ 只要您正在处理PR,请将其标记为草稿。 
+ 请确保您的PR与`main`的最新更改同步。 
+  mention您的PR要解决的问题(Fix:#{ID_1},#{ID_2})。
+ 确保您的PR通过所有检查。

**🈴 审阅PR:**

+ 报错尊重、建设性的原则 
+ 要尊重和建设性的
+ 将自己分配给 PR
+ 检查所有检查是否通过
+ 提出更改建议而不仅仅评论发现的问题。
+ 如果您对某些事情不确定,请向作者提问
+ 如果您不确定更改是否有效,请尝试更改
+ 如果您对某些事情不确定,请联系其他审阅人员
+ 如果您对更改感到满意,请批准 PR
+ 合并PR一次都有批准并通过检查

**⚠️ DCO检测:**

我们在每个拉取请求上运行DCO检查以确保代码质量和可维护性。此检查验证提交已签名,表示您已阅读并同意开发者证书原则的规定。如果您尚未对提交进行签名,可以使用以下命令对您最后进行的提交进行签名:

```bash
❯ git commit --amend --signoff
```

请注意,在提交上签名是一个承诺,即您已经阅读并同意开发者证书原则的规定。如果您还没有阅读这份文件,我们强烈建议您花些时间仔细阅读。如果您对本文件的内容或需要进一步帮助有任何疑问,请联系管理员或相关人员。

您也可以通过将以下内容添加到 `.zshrc` 或 `.bashrc` 中来自动签名您的提交:

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


#### 文档贡献
Open-IM-Server的文档包括:
+ [README.md](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/README.md):此文件包含有关入门Open-IM-Server的基本信息和说明。
+ [CONTRIBUTING.md](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/CONTRIBUTING.md):此文件包含有关为Open-IM-Server代码库做出贡献的指南,例如如何提交问题、拉取请求和代码审查。
+ [DEVELOPGUIDE.md](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/DEVELOPGUIDE.md):此文件提供了有关开发Open-IM-Server的更深入指南,包括有关项目架构、编码约定和测试实践的信息。
+ [官方文档](https://doc.rentsoft.cn/):这是Open-IM-Server的官方文档,其中包括对所有功能、配置选项和故障排除提示的全面信息。


请遵守以下规则以更好地格式化文档,这将极大地提高阅读体验。

1. 请不要在英文文档中使用中文标点,反之亦然。
2. 请在适当的地方使用大写字母,例如句子/标题的第一个字母等。
3. 请为每个Markdown代码块指定语言,除非没有关联的语言。
4. 请在中文和英语单词之间插入空格。 
5.请为技术术语使用正确的大小写,例如使用 `HTTP` 而不是http， `MySQL` 而不是mysql, `Kubernetes` i而不是kubernetes等
6. 请在提交PR之前检查文档中是否有任何拼写错误。

## 参与任何形式的帮助

我们选择GitHub作为Open-IM-Server主要的协作场所。因此,Open-IM-Server的最新更新总是在这里。虽然通过PR是明确的帮助方式,但我们仍然呼吁任何其他方式的帮助。

+ 尽可能回复他人的 [issues](https://github.com/OpenIMSDK/Open-IM-Server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc)。 
+ 帮助解决其他用户的问题；
+ 帮助审核他人的 [PR](https://github.com/OpenIMSDK/Open-IM-Server/pulls?q=is%3Apr+is%3Aopen+sort%3Aupdated-desc) 设计; 
+ 讨论Open-IM-Server，讨论通常会使问题变得清晰; 
+ 在Github之外宣传 [Open-IM-Server](google.com/search?q=Open-IM-Server) 技术 ;
+ 在Open-IM-Server上写博客等等。

简而言之,**任何帮助都是贡献。**   

## 发布版本

Open-IM-Server的版本发布使用 [Release Please](https://github.com/googleapis/release-please) 和 [GoReleaser](https://goreleaser.com/)完成。工作流程如下:

🎯 PR合并到`main`分支:

+ Release please被触发,创建或更新新的版本PR 
+ 每次合并到主分支时,都会执行此操作,当前的发行版PR每次都会更新

🎯 合并'velease please' PR到`main`:

+ Release please被触发,根据提交消息创建新版本和更新更改日志 
+ 触发GoReleaser,构建二进制文件和将其附加到版本 
+ 创建容器并推送到容器注册表

下一次相关合并,将创建新的发行版PR,然后过程重新开始。

**👀手动设置版本号：**

如果您想手动设置版本号,可以创建一个包含版本号的空提交消息的PR。举个例子: 


这样的PR可以按照如下方式产生

````bash
❯ git commit --allow-empty -m "chore: release 0.0.3" -m "Release-As: 0.0.3
````

## 联系我们

我们非常重视与用户,开发者和贡献者的紧密联系。拥有庞大的社区和维护团队,我们随时准备提供帮助和支持。无论您是想加入我们的社区,还是有任何疑问或建议,我们欢迎您与我们取得联系。

我们最推荐的方式是通过 [Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg)。即使您在中国,Slack通常也不会被防火墙封锁,这使其成为与我们联系的简便方式。我们的Slack社区是与其他Open-IM-Server用户和开发者讨论和共享想法和建议的理想场所。您可以提出技术问题,寻求帮助,或与其他Open-IM-Server用户分享您的体验。

除了Slack,我们还提供以下联系方式:

+ <a href="https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg" target="_blank"><img src="https://img.shields.io/badge/slack-%40OpenIMSDKCore-informational?logo=slack&style=flat-square"></a>:   我们还有Slack频道供您沟通和讨论。要加入,请访问https://slack.com/并加入我们的 [👀 Open-IM-Server slack](https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg) 团队频道。
+ <a href="https://mail.google.com/mail/u/0/?fs=1&tf=cm&to=4closetool3@gmail.com" target="_blank"><img src="https://img.shields.io/badge/gmail-%40OOpenIMSDKCore?style=social&logo=gmail"></a>: 通过 [Gmail](winxu81@gmail.com)与我们联系。如果您有任何需要解决的问题或问题,或者对我们的开源项目有任何建议和反馈,请随时通过电子邮件与我们联系。 
+ <a href="https://doc.rentsoft.cn/" target="_blank"><img src="https://img.shields.io/badge/%E5%8D%9A%E5%AE%A2-%40OpenIMSDKCore-blue?style=social&logo=Octopus%20Deploy"></a>: 阅读我们的 [blog](https://doc.rentsoft.cn/)。 我们的博客是跟上Open-IM-Server项目和趋势的理想场所。在博客上,我们分享最新进展,技术趋势和其他有趣的信息。
+ <a href="https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg" target="_blank"><img src="https://img.shields.io/badge/%E5%BE%AE%E4%BF%A1-OpenIMSDKCore-brightgreen?logo=wechat&style=flat-square"></a>: 添加 [微信](https://github.com/OpenIMSDK/OpenIM-Docs/blob/main/docs/images/WechatIMG20.jpeg) 并表明您是Open-IM-Server的用户或开发者。我们将尽快处理您的请求。

无论您是想加入我们的社区,还是有任何疑问或建议, 我们都欢迎您与我们的接触交流。

#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Usage Instructions: https://github.com/OpenIMSDK/Open-IM-Server/tree/main/docs/contrib/git_cherry-pick.md

# Checkout a PR from GitHub. (Yes, this is sitting in a Git tree. How
# meta.) Assumes you care about pulls from remote "upstream" and
# checks them out to a branch named:
#  automated-cherry-pick-of-<pr>-<target branch>-<timestamp>

set -o errexit
set -o nounset
set -o pipefail

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/lib/init.sh"

REPO_ROOT="$(git rev-parse --show-toplevel)"
declare -r REPO_ROOT
cd "${REPO_ROOT}"

STARTINGBRANCH=$(git symbolic-ref --short HEAD)
declare -r STARTINGBRANCH
declare -r REBASEMAGIC="${REPO_ROOT}/.git/rebase-apply"
DRY_RUN=${DRY_RUN:-""}
REGENERATE_DOCS=${REGENERATE_DOCS:-""}
UPSTREAM_REMOTE=${UPSTREAM_REMOTE:-upstream}
FORK_REMOTE=${FORK_REMOTE:-origin}
MAIN_REPO_ORG=${MAIN_REPO_ORG:-$(git remote get-url "$UPSTREAM_REMOTE" | awk '{gsub(/http[s]:\/\/|git@/,"")}1' | awk -F'[@:./]' 'NR==1{print $3}')}
MAIN_REPO_NAME=${MAIN_REPO_NAME:-$(git remote get-url "$UPSTREAM_REMOTE" | awk '{gsub(/http[s]:\/\/|git@/,"")}1' | awk -F'[@:./]' 'NR==1{print $4}')}

if [[ -z ${GITHUB_USER:-} ]]; then
  openim::log::error_exit "Please export GITHUB_USER=<your-user> (or GH organization, if that's where your fork lives)"
fi

if ! command -v gh > /dev/null; then
  openim::log::error_exit "Can't find 'gh' tool in PATH, please install from https://github.com/cli/cli"
fi

if [[ "$#" -lt 2 ]]; then
  echo "${0} <remote branch> <pr-number>...: cherry pick one or more <pr> onto <remote branch> and leave instructions for proposing pull request"
  echo
  echo "  Checks out <remote branch> and handles the cherry-pick of <pr> (possibly multiple) for you."
  echo "  Examples:"
  echo "    $0 upstream/release-v3.1 12345        # Cherry-picks PR 12345 onto upstream/release-v3.1 and proposes that as a PR."
  echo "    $0 upstream/release-v3.1 12345 56789  # Cherry-picks PR 12345, then 56789 and proposes the combination as a single PR."
  echo
  echo "  Set the DRY_RUN environment var to skip git push and creating PR."
  echo "  This is useful for creating patches to a release branch without making a PR."
  echo "  When DRY_RUN is set the script will leave you in a branch containing the commits you cherry-picked."
  echo
  echo "  Set the REGENERATE_DOCS environment var to regenerate documentation for the target branch after picking the specified commits."
  echo "  This is useful when picking commits containing changes to API documentation."
  echo
  echo "  Set UPSTREAM_REMOTE (default: upstream) and FORK_REMOTE (default: origin)"
  echo "  to override the default remote names to what you have locally."
  echo
  echo "  For merge process info, see https://github.com/OpenIMSDK/Open-IM-Server/tree/main/docs/contrib/git_cherry-pick.md"
  exit 2
fi

# Checks if you are logged in. Will error/bail if you are not.
gh auth status

if git_status=$(git status --porcelain --untracked=no 2>/dev/null) && [[ -n "${git_status}" ]]; then
  openim::log::error_exit "!!! Dirty tree. Clean up and try again."
fi

if [[ -e "${REBASEMAGIC}" ]]; then
  openim::log::error_exit "!!! 'git rebase' or 'git am' in progress. Clean up and try again."
fi

declare -r BRANCH="$1"
shift 1
declare -r PULLS=( "$@" )

function join { local IFS="$1"; shift; echo "$*"; }
PULLDASH=$(join - "${PULLS[@]/#/#}") # Generates something like "#12345-#56789"
declare -r PULLDASH
PULLSUBJ=$(join " " "${PULLS[@]/#/#}") # Generates something like "#12345 #56789"
declare -r PULLSUBJ

openim::log::status "Updating remotes..."
git remote update "${UPSTREAM_REMOTE}" "${FORK_REMOTE}"

if ! git log -n1 --format=%H "${BRANCH}" >/dev/null 2>&1; then
  openim::log::error " '${BRANCH}' not found. The second argument should be something like ${UPSTREAM_REMOTE}/release-0.21."
  openim::log::error "    (In particular, it needs to be a valid, existing remote branch that I can 'git checkout'.)"
  exit 1
fi

NEWBRANCHREQ="automated-cherry-pick-of-${PULLDASH}" # "Required" portion for tools.
declare -r NEWBRANCHREQ
NEWBRANCH="$(echo "${NEWBRANCHREQ}-${BRANCH}" | sed 's/\//-/g')"
declare -r NEWBRANCH
NEWBRANCHUNIQ="${NEWBRANCH}-$(date +%s)"
declare -r NEWBRANCHUNIQ
openim::log::info "+++ Creating local branch ${NEWBRANCHUNIQ}"

cleanbranch=""
gitamcleanup=false
function return_to_kansas {
  if [[ "${gitamcleanup}" == "true" ]]; then
    echo
    openim::log::status "Aborting in-progress git am."
    git am --abort >/dev/null 2>&1 || true
  fi

  # return to the starting branch and delete the PR text file
  if [[ -z "${DRY_RUN}" ]]; then
    echo
    openim::log::status "Returning you to the ${STARTINGBRANCH} branch and cleaning up."
    git checkout -f "${STARTINGBRANCH}" >/dev/null 2>&1 || true
    if [[ -n "${cleanbranch}" ]]; then
      git branch -D "${cleanbranch}" >/dev/null 2>&1 || true
    fi
  fi
}
trap return_to_kansas EXIT

SUBJECTS=()
function make-a-pr() {
  local rel
  rel="$(basename "${BRANCH}")"
  echo
  openim::log::status "Creating a pull request on GitHub at ${GITHUB_USER}:${NEWBRANCH}"

  local numandtitle
  numandtitle=$(printf '%s\n' "${SUBJECTS[@]}")
  prtext=$(cat <<EOF
Cherry pick of ${PULLSUBJ} on ${rel}.

${numandtitle}

For details on the cherry pick process, see the [cherry pick requests](https://github.com/OpenIMSDK/Open-IM-Server/tree/main/docs/contrib/git_cherry-pick.md) page.

\`\`\`release-note

\`\`\`
EOF
)

  gh pr create --title="Automated cherry pick of ${numandtitle}" --body="${prtext}" --head "${GITHUB_USER}:${NEWBRANCH}" --base "${rel}" --repo="${MAIN_REPO_ORG}/${MAIN_REPO_NAME}"
}

git checkout -b "${NEWBRANCHUNIQ}" "${BRANCH}"
cleanbranch="${NEWBRANCHUNIQ}"

gitamcleanup=true
for pull in "${PULLS[@]}"; do
  openim::log::status "Downloading patch to /tmp/${pull}.patch (in case you need to do this again)"

  curl -o "/tmp/${pull}.patch" -sSL "https://github.com/${MAIN_REPO_ORG}/${MAIN_REPO_NAME}/pull/${pull}.patch"
  echo
  openim::log::status "About to attempt cherry pick of PR. To reattempt:"
  echo "  $ git am -3 /tmp/${pull}.patch"
  echo
  git am -3 "/tmp/${pull}.patch" || {
    conflicts=false
    while unmerged=$(git status --porcelain | grep ^U) && [[ -n ${unmerged} ]] \
      || [[ -e "${REBASEMAGIC}" ]]; do
      conflicts=true # <-- We should have detected conflicts once
      echo
      openim::log::status "Conflicts detected:"
      echo
      (git status --porcelain | grep ^U) || echo "!!! None. Did you git am --continue?"
      echo
      openim::log::status "Please resolve the conflicts in another window (and remember to 'git add / git am --continue')"
      read -p "+++ Proceed (anything other than 'y' aborts the cherry-pick)? [y/n] " -r
      echo
      if ! [[ "${REPLY}" =~ ^[yY]$ ]]; then
        echo "Aborting." >&2
        exit 1
      fi
    done

    if [[ "${conflicts}" != "true" ]]; then
      echo "!!! git am failed, likely because of an in-progress 'git am' or 'git rebase'"
      exit 1
    fi
  }

  # set the subject
  subject=$(grep -m 1 "^Subject" "/tmp/${pull}.patch" | sed -e 's/Subject: \[PATCH//g' | sed 's/.*] //')
  SUBJECTS+=("#${pull}: ${subject}")

  # remove the patch file from /tmp
  rm -f "/tmp/${pull}.patch"
done
gitamcleanup=false

# Re-generate docs (if needed)
if [[ -n "${REGENERATE_DOCS}" ]]; then
  echo
  echo "Regenerating docs..."
  if ! scripts/generate-docs.sh; then
    echo
    echo "scripts/gendoc.sh FAILED to complete."
    exit 1
  fi
fi

if [[ -n "${DRY_RUN}" ]]; then
  openim::log::error "!!! Skipping git push and PR creation because you set DRY_RUN."
  echo "To return to the branch you were in when you invoked this script:"
  echo
  echo "  git checkout ${STARTINGBRANCH}"
  echo
  echo "To delete this branch:"
  echo
  echo "  git branch -D ${NEWBRANCHUNIQ}"
  exit 0
fi

if git remote -v | grep ^"${FORK_REMOTE}" | grep "${MAIN_REPO_ORG}/${MAIN_REPO_NAME}.git"; then
  echo "!!! You have ${FORK_REMOTE} configured as your ${MAIN_REPO_ORG}/${MAIN_REPO_NAME}.git"
  echo "This isn't normal. Leaving you with push instructions:"
  echo
  openim::log::status "First manually push the branch this script created:"
  echo
  echo "  git push REMOTE ${NEWBRANCHUNIQ}:${NEWBRANCH}"
  echo
  echo "where REMOTE is your personal fork (maybe ${UPSTREAM_REMOTE}? Consider swapping those.)."
  echo "OR consider setting UPSTREAM_REMOTE and FORK_REMOTE to different values."
  echo
  make-a-pr
  cleanbranch=""
  exit 0
fi

echo
openim::log::status "I'm about to do the following to push to GitHub (and I'm assuming ${FORK_REMOTE} is your personal fork):"
echo
echo "  git push ${FORK_REMOTE} ${NEWBRANCHUNIQ}:${NEWBRANCH}"
echo
read -p "+++ Proceed (anything other than 'y' aborts the cherry-pick)? [y/n] " -r
if ! [[ "${REPLY}" =~ ^[yY]$ ]]; then
  echo "Aborting." >&2
  exit 1
fi

git push "${FORK_REMOTE}" -f "${NEWBRANCHUNIQ}:${NEWBRANCH}"
make-a-pr
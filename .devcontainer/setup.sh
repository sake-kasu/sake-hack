#!/usr/bin/env zsh

set -euo pipefail

echo "🧪 Goのバージョンを確認:"
go version

# TestContainers統合テスト用: Dockerソケットの権限確認
if [ -S /var/run/docker.sock ]; then
  echo "🐳 Dockerソケット検出: TestContainers統合テスト利用可能"
  DOCKER_GID=$(stat -c '%g' /var/run/docker.sock)
  echo "   Dockerグループ ID: $DOCKER_GID"
  groups | grep -q docker && echo "   ✅ nodeユーザーはdockerグループに所属しています" || echo "   ⚠️  nodeユーザーはdockerグループに所属していません"
else
  echo "⚠️  Dockerソケットが見つかりません（統合テストは実行できません）"
fi

ZSHRC="/home/node/.zshrc"

cat << 'EOF' >> "$ZSHRC"

# Git補完・プロンプト
source ~/.zsh/git-prompt.sh
fpath=(~/.zsh $fpath)
autoload -U compinit
compinit -u

# カラー補完
autoload -U colors
colors
zstyle ':completion:*' list-colors "${LS_COLORS}"

# autosuggestions
source ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh

# 補完設定
setopt complete_in_word
zstyle ':completion:*:default' menu select=1
zstyle ':completion::complete:*' use-cache true
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Z}'
setopt list_packed

# コマンド修正提案
setopt correct
SPROMPT="correct: %R -> %r ? [Yes/No/Abort/Edit] => "

# Git PS1 プロンプト設定
GIT_PS1_SHOWDIRTYSTATE=true
GIT_PS1_SHOWUNTRACKEDFILES=true
GIT_PS1_SHOWSTASHSTATE=true
GIT_PS1_SHOWUPSTREAM=auto

# プロンプト表示
setopt PROMPT_SUBST
PS1='%F{green}%n@%m%f: %F{cyan}%~%f %F{red}$(__git_ps1 "(%s)")%f'$'\n''\$ '
EOF

echo "✅ DevContainerの設定完了"

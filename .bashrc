export TMPDIR=$PREFIX/tmp
alias claude="proot --bind=$PREFIX/tmp:/tmp $(which claude)"
. "$HOME/.cargo/env"
export PATH="/data/data/com.termux/files/home/.npm-global/bin:$PATH"

# Claude Code 완벽한 분신 시스템 별칭
alias check-logic='source ~/.claude/aliases.sh && check_logic'
alias save-proof='source ~/.claude/aliases.sh && save_proof'
alias verify-vanilla='source ~/.claude/aliases.sh && verify_vanilla'
alias red-team='source ~/.claude/aliases.sh && red_team'
alias refactor-first='source ~/.claude/aliases.sh && refactor_first'
alias prove-it='source ~/.claude/aliases.sh && prove_it'
alias map-index='bash ~/.claude/hooks/auto-map-index.sh'
alias termux-optimize='bash ~/.claude/hooks/termux-optimize.sh'

# GitHub Token (Blogger automation 및 저장소 관리)
export GITHUB_TOKEN="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"


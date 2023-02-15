motd() {
    # shellcheck disable=SC2046
    echo "$(tput setaf 3)
    db   dD db    db d8888b. d88888b .d8888. d888888b .88b  d88.
    88 ,8P' 88    88 88  \`8D 88'     88'  YP   \`88'   88'YbdP\`88
    88,8P   88    88 88oooY' 88ooooo \`8bo.      88    88  88  88
    88\`8b   88    88 88~~~b. 88~~~~~   \`Y8b.    88    88  88  88
    88 \`88. 88b  d88 88   8D 88.     db   8D   .88.   88  88  88
    YP   YD ~Y8888P' Y8888P' Y88888P \`8888Y' Y888888P YP  YP  YP
$(tput setaf 3)$(figlet $(hostname))
    $(tput sgr0)"
}

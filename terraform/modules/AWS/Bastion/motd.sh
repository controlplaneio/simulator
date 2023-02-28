motd() {
    if [[ "${KUBESIM:-}" != "" ]]; then
      return
    fi

    # shellcheck disable=SC2046
    echo "$(tput setaf 3)
                                                        |>>>
                                                        |
                                                    _  _|_  _
                                                  |;|_|;|_|;|
                                                  \\\\.    .  /
                                                  \\\\:  .  /
                                                    ||:   |
                                                    ||:.  |
                                                    ||:  .|
                                                    ||:   |       \\,/
                                                    ||: , |            /\`\\
                                                    ||:   |
                                                    ||: . |
                      __                            _||_   |
        $(tput setaf 4)     ____----    \\----__            __ ----     -----\\              ___
        -~--~                   ----____---/                  ------_____---    -----~~
$(tput setaf 3)$(figlet $(hostname))
        $(tput sgr0)"
}

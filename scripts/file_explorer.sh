fe() {
    local out="/tmp/fe-output-$$"
    rm -f "$out"  # Clean up old file just in case

    /home/will/Projects/go/file-explorer/go-file-explorer "$out"
	clear

    if [[ -f "$out" ]]; then
        local cmd
        cmd=$(<"$out")
        rm -f "$out"

        if [[ -n "$cmd" ]]; then
            eval "$cmd"
        else
            echo "No command received from file explorer." >&2
        fi
    else
        echo "Output file not found." >&2
    fi
}

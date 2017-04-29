needroot :
	@[ "$(shell id -u)" -eq "0" ] || exit 1

notroot :
	@[ "$(shell id -u)" != "0" ] || exit 1

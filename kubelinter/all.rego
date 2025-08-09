package kubelinter

# METADATA
# description: Run all rules
# entrypoint: true
deny if {
	data.kubelinter.rules[_].deny
}

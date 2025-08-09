package kubelinter.template.pdbmaxunavailable

import data.kubelinter.objectkinds.is_poddisruptionbudget

deny contains msg if {
	is_poddisruptionbudget
	input.spec.maxUnavailable
	maxUnavailable := input.spec.maxUnavailable
	maxUnavailable == 0
	msg := "MaxUnavailable is set to 0"
}

deny contains msg if {
	is_poddisruptionbudget
	input.spec.maxUnavailable
	maxUnavailable := input.spec.maxUnavailable
	maxUnavailable == "0%"
	msg := "MaxUnavailable is set to 0"
}

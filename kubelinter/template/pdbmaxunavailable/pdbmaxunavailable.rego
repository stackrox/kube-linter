package kubelinter.template.pdbmaxunavailable

import data.kubelinter.objectkinds.is_poddisruptionbudget

deny contains msg if {
	is_poddisruptionbudget
	max_unavailable := input.spec.maxUnavailable
	max_unavailable == 0
	msg := "MaxUnavailable is set to 0"
}

deny contains msg if {
	is_poddisruptionbudget
	max_unavailable := input.spec.maxUnavailable
	max_unavailable == "0%"
	msg := "MaxUnavailable is set to 0"
}

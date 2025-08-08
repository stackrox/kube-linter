package kubelinter.template.pdbunhealthypodevictionpolicy

import kubelinter.objectkinds.is_poddisruptionbudget

deny contains msg if {
	is_poddisruptionbudget
	not input.spec.unhealthyPodEvictionPolicy
	msg := "unhealthyPodEvictionPolicy is not explicitly set"
}
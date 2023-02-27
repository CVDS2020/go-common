package flag

import "gitee.com/sy_183/common/generic"

func TestFlag[F generic.Integer](flag F, mask F) bool {
	return flag&mask != 0
}

func MaskFlag[F generic.Integer](flag F, mask F) F {
	return flag | mask
}

func MaskFlagP[F generic.Integer](flagP *F, mask F) {
	*flagP |= mask
}

func UnmaskFlag[F generic.Integer](flag F, mask F) F {
	return flag & ^mask
}

func UnmaskFlagP[F generic.Integer](flagP *F, mask F) {
	*flagP &= ^mask
}

func SwapFlagMask[F generic.Integer](flag F, unmask, mask F) F {
	return (flag & ^unmask) | mask
}

func SwapFlagPMask[F generic.Integer](flagP *F, unmask, mask F) {
	*flagP = (*flagP & ^unmask) | mask
}

package sgr

import "strings"

var (
	// BoldFlag                 = 8
	// FaintFlag                = 9
	// ItalicFlag               = 10
	// FrakturFlag              = 11
	// ReversedFlag             = 12
	// NormalIntensityFlag      = 13
	// NotItalicBlackLetterFlag = 14
	// NotReversedFlag          = 15
	cacheTable1 [256]string

	// UnderlineFlag        = 16
	// DoublyUnderlinedFlag = 17
	// SlowBlinkFlag        = 18
	// RapidBlinkFlag       = 19
	// OverlinedFlag        = 20
	// NotUnderlinedFlag    = 21
	// NotBlinkingFlag      = 22
	// NotOverlinedFlag     = 23
	cacheTable2 [256]string

	// ConcealFlag                    = 24
	// CrossedOutFlag                 = 25
	// DisableProportionalSpacingFlag = 26
	// RevealFlag                     = 27
	// NotCrossedOutFlag              = 28
	// ProportionalSpacingFlag        = 29
	cacheTable3 [256]string

	// FramedFlag                  = 32
	// EncircledFlag               = 33
	// SuperscriptFlag             = 34
	// SubscriptFlag               = 35
	// NotFramedEncircledFlag      = 36
	// NotSuperscriptSubscriptFlag = 37
	cacheTable4 [256]string

	// IdeogramUnderlineFlag       = 40
	// IdeogramDoubleUnderlineFlag = 41
	// IdeogramOverlineFlag        = 42
	// IdeogramDoubleOverlineFlag  = 43
	// IdeogramStressMarkingFlag   = 44
	// NoIdeogramAttributesFlag    = 45
	cacheTable5 [256]string
)

func initCacheTable1() {
	for i := 0; i < 256; i++ {
		var flags1S []string
		if i&0b00100000 != 0 {
			flags1S = append(flags1S, NormalIntensityCode)
		} else if i&0b00000001 != 0 {
			flags1S = append(flags1S, BoldCode)
		} else if i&0b00000010 != 0 {
			flags1S = append(flags1S, FaintCode)
		}
		if i&0b01000000 != 0 {
			flags1S = append(flags1S, NotItalicBlackLetterCode)
		} else if i&0b00000100 != 0 {
			flags1S = append(flags1S, ItalicCode)
		} else if i&0b00001000 != 0 {
			flags1S = append(flags1S, FrakturCode)
		}
		if i&0b10000000 != 0 {
			flags1S = append(flags1S, NotReversedCode)
		} else if i&0b00010000 != 0 {
			flags1S = append(flags1S, ReversedCode)
		}
		cacheTable1[i] = strings.Join(flags1S, ";")
	}
}

func initCacheTable2() {
	for i := 0; i < 256; i++ {
		var flags2S []string
		if i&0b00100000 != 0 {
			flags2S = append(flags2S, NotUnderlinedCode)
		} else if i&0b00000001 != 0 {
			flags2S = append(flags2S, UnderlineCode)
		} else if i&0b00000010 != 0 {
			flags2S = append(flags2S, DoublyUnderlinedCode)
		}
		if i&0b01000000 != 0 {
			flags2S = append(flags2S, NotBlinkingCode)
		} else if i&0b00000100 != 0 {
			flags2S = append(flags2S, SlowBlinkCode)
		} else if i&0b00001000 != 0 {
			flags2S = append(flags2S, RapidBlinkCode)
		}
		if i&0b10000000 != 0 {
			flags2S = append(flags2S, NotOverlinedCode)
		} else if i&0b00010000 != 0 {
			flags2S = append(flags2S, OverlinedCode)
		}
		cacheTable2[i] = strings.Join(flags2S, ";")
	}
}

func initCacheTable3() {
	for i := 0; i < 256; i++ {
		var flags3S []string
		if i&0b00001000 != 0 {
			flags3S = append(flags3S, RevealCode)
		} else if i&0b00000001 != 0 {
			flags3S = append(flags3S, ConcealCode)
		}
		if i&0b00010000 != 0 {
			flags3S = append(flags3S, NotCrossedOutCode)
		} else if i&0b00000010 != 0 {
			flags3S = append(flags3S, CrossedOutCode)
		}
		if i&0b00100000 != 0 {
			flags3S = append(flags3S, DisableProportionalSpacingCode)
		} else if i&0b00000100 != 0 {
			flags3S = append(flags3S, ProportionalSpacingCode)
		}
		cacheTable3[i] = strings.Join(flags3S, ";")
	}
}

func initCacheTable4() {
	for i := 0; i < 256; i++ {
		var flags4S []string
		if i&0b00010000 != 0 {
			flags4S = append(flags4S, NotFramedEncircledCode)
		} else if i&0b00000001 != 0 {
			flags4S = append(flags4S, FramedCode)
		} else if i&0b00000010 != 0 {
			flags4S = append(flags4S, EncircledCode)
		}
		if i&0b00100000 != 0 {
			flags4S = append(flags4S, NotSuperscriptSubscriptCode)
		} else if i&0b00000100 != 0 {
			flags4S = append(flags4S, SuperscriptCode)
		} else if i&0b00001000 != 0 {
			flags4S = append(flags4S, SubscriptCode)
		}
		cacheTable4[i] = strings.Join(flags4S, ";")
	}
}

func initCacheTable5() {
	for i := 0; i < 256; i++ {
		var flags5S []string
		if i&0b00100000 != 0 {
			flags5S = append(flags5S, NoIdeogramAttributesCode)
		} else if i&0b00000001 != 0 {
			flags5S = append(flags5S, IdeogramUnderlineCode)
		} else if i&0b00000010 != 0 {
			flags5S = append(flags5S, IdeogramDoubleUnderlineCode)
		} else if i&0b00000100 != 0 {
			flags5S = append(flags5S, IdeogramOverlineCode)
		} else if i&0b00001000 != 0 {
			flags5S = append(flags5S, IdeogramDoubleOverlineCode)
		} else if i&0b00010000 != 0 {
			flags5S = append(flags5S, IdeogramStressMarkingCode)
		}
		cacheTable5[i] = strings.Join(flags5S, ";")
	}
}

func init() {
	initCacheTable1()
	initCacheTable2()
	initCacheTable3()
	initCacheTable4()
	initCacheTable5()
}

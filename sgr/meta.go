package sgr

type SGRMeta struct {
	Flag Flag
	Name string
	Code string
}

var (
	ResetMeta = SGRMeta{Flag: Reset, Name: ResetName, Code: ResetCode}

	BoldMeta                 = SGRMeta{Flag: Bold, Name: BoldName, Code: BoldCode}
	FaintMeta                = SGRMeta{Flag: Faint, Name: FaintName, Code: FaintCode}
	ItalicMeta               = SGRMeta{Flag: Italic, Name: ItalicName, Code: ItalicCode}
	FrakturMeta              = SGRMeta{Flag: Fraktur, Name: FrakturName, Code: FrakturCode}
	ReversedMeta             = SGRMeta{Flag: Reversed, Name: ReversedName, Code: ReversedCode}
	NormalIntensityMeta      = SGRMeta{Flag: NormalIntensity, Name: NormalIntensityName, Code: NormalIntensityCode}
	NotItalicBlackLetterMeta = SGRMeta{Flag: NotItalicBlackLetter, Name: NotItalicBlackLetterName, Code: NotItalicBlackLetterCode}
	NotReversedMeta          = SGRMeta{Flag: NotReversed, Name: NotReversedName, Code: NotReversedCode}

	UnderlineMeta        = SGRMeta{Flag: Underline, Name: UnderlineName, Code: UnderlineCode}
	DoublyUnderlinedMeta = SGRMeta{Flag: DoublyUnderlined, Name: DoublyUnderlinedName, Code: DoublyUnderlinedCode}
	SlowBlinkMeta        = SGRMeta{Flag: SlowBlink, Name: SlowBlinkName, Code: SlowBlinkCode}
	RapidBlinkMeta       = SGRMeta{Flag: RapidBlink, Name: RapidBlinkName, Code: RapidBlinkCode}
	OverlinedMeta        = SGRMeta{Flag: Overlined, Name: OverlinedName, Code: OverlinedCode}
	NotUnderlinedMeta    = SGRMeta{Flag: NotUnderlined, Name: NotUnderlinedName, Code: NotUnderlinedCode}
	NotBlinkingMeta      = SGRMeta{Flag: NotBlinking, Name: NotBlinkingName, Code: NotBlinkingCode}
	NotOverlinedMeta     = SGRMeta{Flag: NotOverlined, Name: NotOverlinedName, Code: NotOverlinedCode}

	ConcealMeta                    = SGRMeta{Flag: Conceal, Name: ConcealName, Code: ConcealCode}
	CrossedOutMeta                 = SGRMeta{Flag: CrossedOut, Name: CrossedOutName, Code: CrossedOutCode}
	DisableProportionalSpacingMeta = SGRMeta{Flag: DisableProportionalSpacing, Name: DisableProportionalSpacingName, Code: DisableProportionalSpacingCode}
	RevealMeta                     = SGRMeta{Flag: Reveal, Name: RevealName, Code: RevealCode}
	NotCrossedOutMeta              = SGRMeta{Flag: NotCrossedOut, Name: NotCrossedOutName, Code: NotCrossedOutCode}
	ProportionalSpacingMeta        = SGRMeta{Flag: ProportionalSpacing, Name: ProportionalSpacingName, Code: ProportionalSpacingCode}

	FramedMeta                  = SGRMeta{Flag: Framed, Name: FramedName, Code: FramedCode}
	EncircledMeta               = SGRMeta{Flag: Encircled, Name: EncircledName, Code: EncircledCode}
	SuperscriptMeta             = SGRMeta{Flag: Superscript, Name: SuperscriptName, Code: SuperscriptCode}
	SubscriptMeta               = SGRMeta{Flag: Subscript, Name: SubscriptName, Code: SubscriptCode}
	NotFramedEncircledMeta      = SGRMeta{Flag: NotFramedEncircled, Name: NotFramedEncircledName, Code: NotFramedEncircledCode}
	NotSuperscriptSubscriptMeta = SGRMeta{Flag: NotSuperscriptSubscript, Name: NotSuperscriptSubscriptName, Code: NotSuperscriptSubscriptCode}

	IdeogramUnderlineMeta       = SGRMeta{Flag: IdeogramUnderline, Name: IdeogramUnderlineName, Code: IdeogramUnderlineCode}
	IdeogramDoubleUnderlineMeta = SGRMeta{Flag: IdeogramDoubleUnderline, Name: IdeogramDoubleUnderlineName, Code: IdeogramDoubleUnderlineCode}
	IdeogramOverlineMeta        = SGRMeta{Flag: IdeogramOverline, Name: IdeogramOverlineName, Code: IdeogramOverlineCode}
	IdeogramDoubleOverlineMeta  = SGRMeta{Flag: IdeogramDoubleOverline, Name: IdeogramDoubleOverlineName, Code: IdeogramDoubleOverlineCode}
	IdeogramStressMarkingMeta   = SGRMeta{Flag: IdeogramStressMarking, Name: IdeogramStressMarkingName, Code: IdeogramStressMarkingCode}
	NoIdeogramAttributesMeta    = SGRMeta{Flag: NoIdeogramAttributes, Name: NoIdeogramAttributesName, Code: NoIdeogramAttributesCode}

	FgBlackMeta   = SGRMeta{Flag: FgBlack, Name: FgBlackName, Code: FgBlackCode}
	FgRedMeta     = SGRMeta{Flag: FgRed, Name: FgRedName, Code: FgRedCode}
	FgGreenMeta   = SGRMeta{Flag: FgGreen, Name: FgGreenName, Code: FgGreenCode}
	FgYellowMeta  = SGRMeta{Flag: FgYellow, Name: FgYellowName, Code: FgYellowCode}
	FgBlueMeta    = SGRMeta{Flag: FgBlue, Name: FgBlueName, Code: FgBlueCode}
	FgMagentaMeta = SGRMeta{Flag: FgMagenta, Name: FgMagentaName, Code: FgMagentaCode}
	FgCyanMeta    = SGRMeta{Flag: FgCyan, Name: FgCyanName, Code: FgCyanCode}
	FgWhiteMeta   = SGRMeta{Flag: FgWhite, Name: FgWhiteName, Code: FgWhiteCode}

	BgBlackMeta   = SGRMeta{Flag: BgBlack, Name: BgBlackName, Code: BgBlackCode}
	BgRedMeta     = SGRMeta{Flag: BgRed, Name: BgRedName, Code: BgRedCode}
	BgGreenMeta   = SGRMeta{Flag: BgGreen, Name: BgGreenName, Code: BgGreenCode}
	BgYellowMeta  = SGRMeta{Flag: BgYellow, Name: BgYellowName, Code: BgYellowCode}
	BgBlueMeta    = SGRMeta{Flag: BgBlue, Name: BgBlueName, Code: BgBlueCode}
	BgMagentaMeta = SGRMeta{Flag: BgMagenta, Name: BgMagentaName, Code: BgMagentaCode}
	BgCyanMeta    = SGRMeta{Flag: BgCyan, Name: BgCyanName, Code: BgCyanCode}
	BgWhiteMeta   = SGRMeta{Flag: BgWhite, Name: BgWhiteName, Code: BgWhiteCode}

	FgBrightBlackMeta   = SGRMeta{Flag: FgBrightBlack, Name: FgBrightBlackName, Code: FgBrightBlackCode}
	FgBrightRedMeta     = SGRMeta{Flag: FgBrightRed, Name: FgBrightRedName, Code: FgBrightRedCode}
	FgBrightGreenMeta   = SGRMeta{Flag: FgBrightGreen, Name: FgBrightGreenName, Code: FgBrightGreenCode}
	FgBrightYellowMeta  = SGRMeta{Flag: FgBrightYellow, Name: FgBrightYellowName, Code: FgBrightYellowCode}
	FgBrightBlueMeta    = SGRMeta{Flag: FgBrightBlue, Name: FgBrightBlueName, Code: FgBrightBlueCode}
	FgBrightMagentaMeta = SGRMeta{Flag: FgBrightMagenta, Name: FgBrightMagentaName, Code: FgBrightMagentaCode}
	FgBrightCyanMeta    = SGRMeta{Flag: FgBrightCyan, Name: FgBrightCyanName, Code: FgBrightCyanCode}
	FgBrightWhiteMeta   = SGRMeta{Flag: FgBrightWhite, Name: FgBrightWhiteName, Code: FgBrightWhiteCode}

	BgBrightBlackMeta   = SGRMeta{Flag: BgBrightBlack, Name: BgBrightBlackName, Code: BgBrightBlackCode}
	BgBrightRedMeta     = SGRMeta{Flag: BgBrightRed, Name: BgBrightRedName, Code: BgBrightRedCode}
	BgBrightGreenMeta   = SGRMeta{Flag: BgBrightGreen, Name: BgBrightGreenName, Code: BgBrightGreenCode}
	BgBrightYellowMeta  = SGRMeta{Flag: BgBrightYellow, Name: BgBrightYellowName, Code: BgBrightYellowCode}
	BgBrightBlueMeta    = SGRMeta{Flag: BgBrightBlue, Name: BgBrightBlueName, Code: BgBrightBlueCode}
	BgBrightMagentaMeta = SGRMeta{Flag: BgBrightMagenta, Name: BgBrightMagentaName, Code: BgBrightMagentaCode}
	BgBrightCyanMeta    = SGRMeta{Flag: BgBrightCyan, Name: BgBrightCyanName, Code: BgBrightCyanCode}
	BgBrightWhiteMeta   = SGRMeta{Flag: BgBrightWhite, Name: BgBrightWhiteName, Code: BgBrightWhiteCode}

	CustomFgColorMeta         = SGRMeta{Flag: CustomFgColor, Name: CustomFgColorName, Code: CustomFgColorCode}
	DefaultFgColorMeta        = SGRMeta{Flag: DefaultFgColor, Name: DefaultFgColorName, Code: DefaultFgColorCode}
	CustomBgColorMeta         = SGRMeta{Flag: CustomBgColor, Name: CustomBgColorName, Code: CustomBgColorCode}
	DefaultBgColorMeta        = SGRMeta{Flag: DefaultBgColor, Name: DefaultBgColorName, Code: DefaultBgColorCode}
	CustomUnderlineColorMeta  = SGRMeta{Flag: CustomUnderlineColor, Name: CustomUnderlineColorName, Code: CustomUnderlineColorCode}
	DefaultUnderlineColorMeta = SGRMeta{Flag: DefaultUnderlineColor, Name: DefaultUnderlineColorName, Code: DefaultUnderlineColorCode}

	PrimaryFontMeta      = SGRMeta{Flag: PrimaryFont, Name: PrimaryFontName, Code: PrimaryFontCode}
	AlternativeFont1Meta = SGRMeta{Flag: AlternativeFont1, Name: AlternativeFont1Name, Code: AlternativeFont1Code}
	AlternativeFont2Meta = SGRMeta{Flag: AlternativeFont2, Name: AlternativeFont2Name, Code: AlternativeFont2Code}
	AlternativeFont3Meta = SGRMeta{Flag: AlternativeFont3, Name: AlternativeFont3Name, Code: AlternativeFont3Code}
	AlternativeFont4Meta = SGRMeta{Flag: AlternativeFont4, Name: AlternativeFont4Name, Code: AlternativeFont4Code}
	AlternativeFont5Meta = SGRMeta{Flag: AlternativeFont5, Name: AlternativeFont5Name, Code: AlternativeFont5Code}
	AlternativeFont6Meta = SGRMeta{Flag: AlternativeFont6, Name: AlternativeFont6Name, Code: AlternativeFont6Code}
	AlternativeFont7Meta = SGRMeta{Flag: AlternativeFont7, Name: AlternativeFont7Name, Code: AlternativeFont7Code}

	AlternativeFont8Meta = SGRMeta{Flag: AlternativeFont8, Name: AlternativeFont8Name, Code: AlternativeFont8Code}
	AlternativeFont9Meta = SGRMeta{Flag: AlternativeFont9, Name: AlternativeFont9Name, Code: AlternativeFont9Code}
)

var flagMetaMap = map[Flag]*SGRMeta{
	Reset: &ResetMeta,

	Bold:                 &BoldMeta,
	Faint:                &FaintMeta,
	Italic:               &ItalicMeta,
	Fraktur:              &FrakturMeta,
	Reversed:             &ReversedMeta,
	NormalIntensity:      &NormalIntensityMeta,
	NotItalicBlackLetter: &NotItalicBlackLetterMeta,
	NotReversed:          &NotReversedMeta,

	Underline:        &UnderlineMeta,
	DoublyUnderlined: &DoublyUnderlinedMeta,
	SlowBlink:        &SlowBlinkMeta,
	RapidBlink:       &RapidBlinkMeta,
	Overlined:        &OverlinedMeta,
	NotUnderlined:    &NotUnderlinedMeta,
	NotBlinking:      &NotBlinkingMeta,
	NotOverlined:     &NotOverlinedMeta,

	Conceal:                    &ConcealMeta,
	CrossedOut:                 &CrossedOutMeta,
	DisableProportionalSpacing: &DisableProportionalSpacingMeta,
	Reveal:                     &RevealMeta,
	NotCrossedOut:              &NotCrossedOutMeta,
	ProportionalSpacing:        &ProportionalSpacingMeta,

	Framed:                  &FramedMeta,
	Encircled:               &EncircledMeta,
	Superscript:             &SuperscriptMeta,
	Subscript:               &SubscriptMeta,
	NotFramedEncircled:      &NotFramedEncircledMeta,
	NotSuperscriptSubscript: &NotSuperscriptSubscriptMeta,

	IdeogramUnderline:       &IdeogramUnderlineMeta,
	IdeogramDoubleUnderline: &IdeogramDoubleUnderlineMeta,
	IdeogramOverline:        &IdeogramOverlineMeta,
	IdeogramDoubleOverline:  &IdeogramDoubleOverlineMeta,
	IdeogramStressMarking:   &IdeogramStressMarkingMeta,
	NoIdeogramAttributes:    &NoIdeogramAttributesMeta,

	FgBlack:   &FgBlackMeta,
	FgRed:     &FgRedMeta,
	FgGreen:   &FgGreenMeta,
	FgYellow:  &FgYellowMeta,
	FgBlue:    &FgBlueMeta,
	FgMagenta: &FgMagentaMeta,
	FgCyan:    &FgCyanMeta,
	FgWhite:   &FgWhiteMeta,

	BgBlack:   &BgBlackMeta,
	BgRed:     &BgRedMeta,
	BgGreen:   &BgGreenMeta,
	BgYellow:  &BgYellowMeta,
	BgBlue:    &BgBlueMeta,
	BgMagenta: &BgMagentaMeta,
	BgCyan:    &BgCyanMeta,
	BgWhite:   &BgWhiteMeta,

	FgBrightBlack:   &FgBrightBlackMeta,
	FgBrightRed:     &FgBrightRedMeta,
	FgBrightGreen:   &FgBrightGreenMeta,
	FgBrightYellow:  &FgBrightYellowMeta,
	FgBrightBlue:    &FgBrightBlueMeta,
	FgBrightMagenta: &FgBrightMagentaMeta,
	FgBrightCyan:    &FgBrightCyanMeta,
	FgBrightWhite:   &FgBrightWhiteMeta,

	BgBrightBlack:   &BgBrightBlackMeta,
	BgBrightRed:     &BgBrightRedMeta,
	BgBrightGreen:   &BgBrightGreenMeta,
	BgBrightYellow:  &BgBrightYellowMeta,
	BgBrightBlue:    &BgBrightBlueMeta,
	BgBrightMagenta: &BgBrightMagentaMeta,
	BgBrightCyan:    &BgBrightCyanMeta,
	BgBrightWhite:   &BgBrightWhiteMeta,

	CustomFgColor:         &CustomFgColorMeta,
	DefaultFgColor:        &DefaultFgColorMeta,
	CustomBgColor:         &CustomBgColorMeta,
	DefaultBgColor:        &DefaultBgColorMeta,
	CustomUnderlineColor:  &CustomUnderlineColorMeta,
	DefaultUnderlineColor: &DefaultUnderlineColorMeta,

	PrimaryFont:      &PrimaryFontMeta,
	AlternativeFont1: &AlternativeFont1Meta,
	AlternativeFont2: &AlternativeFont2Meta,
	AlternativeFont3: &AlternativeFont3Meta,
	AlternativeFont4: &AlternativeFont4Meta,
	AlternativeFont5: &AlternativeFont5Meta,
	AlternativeFont6: &AlternativeFont6Meta,
	AlternativeFont7: &AlternativeFont7Meta,

	AlternativeFont8: &AlternativeFont8Meta,
	AlternativeFont9: &AlternativeFont9Meta,
}

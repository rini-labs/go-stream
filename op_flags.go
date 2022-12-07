package stream

type Type int

const (
	TypeSplitIterator Type = iota
	TypeStream
	TypeOp
	TypeTerminalOp
	TypeUpstreamTerminalOp
)

var (
	Types = []Type{TypeSplitIterator, TypeStream, TypeOp, TypeTerminalOp, TypeUpstreamTerminalOp}
)

const (
	bitsMissing   = 0b00
	BitsSet       = 0b01
	BitsClear     = 0b10
	BitsPreserver = 0b11
)

type maskBuilder struct {
	typeMap map[Type]int
}

func (mb *maskBuilder) Mask(t Type, i int) *maskBuilder {
	mb.typeMap[t] = i
	return mb
}

func (mb *maskBuilder) set(t Type) *maskBuilder {
	return mb.Mask(t, BitsSet)
}

func (mb *maskBuilder) clear(t Type) *maskBuilder {
	return mb.Mask(t, BitsClear)
}

func (mb *maskBuilder) setAndClear(t Type) *maskBuilder {
	return mb.Mask(t, BitsPreserver)
}

func (mb *maskBuilder) Build() map[Type]int {
	for _, t := range Types {
		if _, ok := mb.typeMap[t]; !ok {
			mb.typeMap[t] = bitsMissing
		}
	}
	return mb.typeMap
}

func set(t Type) *maskBuilder {
	return (&maskBuilder{typeMap: make(map[Type]int, len(Types))}).set(t)
}

func newOpFlag(position int, mb *maskBuilder) opFlag {
	position = position * 2
	return opFlag{
		maskTable:   mb.Build(),
		bitPosition: position,
		set:         BitsSet << position,
		clear:       BitsClear << position,
		preserve:    BitsPreserver << position,
	}
}

type opFlag struct {
	maskTable   map[Type]int
	bitPosition int
	set         int
	clear       int
	preserve    int
}

func (of opFlag) Set() int {
	return of.set
}

func (of opFlag) Clear() int {
	return of.clear
}

func (of opFlag) IsStreamFlag() bool {
	return of.maskTable[TypeStream] > 0
}

func (of opFlag) IsKnown(flags int) bool {
	return (flags & of.preserve) == of.set
}
func (of opFlag) IsCleared(flags int) bool {
	return (flags & of.preserve) == of.clear
}

func (of opFlag) IsPreserved(flags int) bool {
	return (flags & of.preserve) == of.preserve
}

func (of opFlag) CanSet(t Type) bool {
	return (of.maskTable[t] & BitsSet) > 0
}

var (
	OpFlagDistinct = newOpFlag(0, set(TypeSplitIterator).set(TypeStream).setAndClear(TypeOp))
	OpFlagSorted   = newOpFlag(1, set(TypeSplitIterator).set(TypeStream).setAndClear(TypeOp))
	OpFlagOrdered  = newOpFlag(2, set(TypeSplitIterator).set(TypeStream).setAndClear(TypeOp).clear(TypeTerminalOp).clear(TypeUpstreamTerminalOp))
	OpFlagSized    = newOpFlag(3, set(TypeSplitIterator).set(TypeStream).clear(TypeOp))
	OpShortCircuit = newOpFlag(12, set(TypeOp).set(TypeTerminalOp))

	OpFlags = []opFlag{OpFlagDistinct, OpFlagSorted, OpFlagOrdered, OpFlagSized, OpShortCircuit}
)

func createMask(t Type) int {
	mask := 0
	for _, of := range OpFlags {
		mask |= of.maskTable[t] << of.bitPosition
	}
	return mask
}

func createFlagMask() int {
	mask := 0
	for _, of := range OpFlags {
		mask |= of.preserve
	}
	return mask
}

var (
	MaskSplitIterator      = createMask(TypeSplitIterator)
	MaskStream             = createMask(TypeStream)
	MaskOp                 = createMask(TypeOp)
	MaskTerminalOp         = createMask(TypeTerminalOp)
	MaskUpstreamTerminalOp = createMask(TypeUpstreamTerminalOp)

	FlagMask = createFlagMask()

	FlagMaskIs      = MaskStream
	FlagMaskNot     = MaskStream << 1
	InitialOpsValue = FlagMaskIs | FlagMaskNot

	IsDistinct  = OpFlagDistinct.set
	NotDistinct = OpFlagDistinct.clear

	IsSorted  = OpFlagSorted.set
	NotSorted = OpFlagSorted.clear

	IsOrdered  = OpFlagOrdered.set
	NotOrdered = OpFlagOrdered.clear

	IsSized  = OpFlagSized.set
	NotSized = OpFlagSized.clear

	IsShortCircuit = OpShortCircuit.set
)

func getMask(flags int) int {
	if flags == 0 {
		return FlagMask
	}
	return ^(flags | ((FlagMaskIs & flags) << 1) | ((FlagMaskNot & flags) >> 1))
}

func CombineOpFlags(newStreamOrOpFlags int, prevCombOpFlags int) int {
	return (prevCombOpFlags & getMask(newStreamOrOpFlags)) | newStreamOrOpFlags
}

func ToStreamFlags(combOpFlags int) int {
	return ((^combOpFlags) >> 1) & FlagMaskIs & combOpFlags
}

func ToCharacteristics(streamFlags int) int {
	return streamFlags & MaskSplitIterator
}

func FromCharacteristics(characteristics int) int {
	return characteristics & MaskSplitIterator
}

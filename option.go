package walk

// Optional is a type with which Option arguments and methods can be applied.
type Optional interface{ *Model | *Field }

type Option[O Optional] func(O) O

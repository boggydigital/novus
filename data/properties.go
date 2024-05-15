package data

const (
	GetContentErrorsProperty   = "get-content-errors"
	DecodeErrorsProperty       = "decode-errors"
	MatchContentErrorsProperty = "match-content-errors"
	ReduceErrorsProperty       = "reduce-errors"

	CurrentElementsProperty  = "current-elements"
	AddedElementsProperty    = "added-elements"
	RemovedElementsProperty  = "removed-elements"
	PreviousElementsProperty = "previous-elements"
	SourceURLProperty        = "source-url"
)

func AllProperties() []string {
	return []string{
		GetContentErrorsProperty,
		DecodeErrorsProperty,
		MatchContentErrorsProperty,
		ReduceErrorsProperty,
		CurrentElementsProperty,
		AddedElementsProperty,
		RemovedElementsProperty,
		PreviousElementsProperty,
		SourceURLProperty,
	}
}

package analysis

type AnalysisInfo interface{}

type Analysis interface {
	GetInfo() AnalysisInfo
}

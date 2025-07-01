package mock

//go:generate mockgen -source=../usecase/analyzer.go -destination=./analyzer_mock.go -package=mock -typed -write_generate_directive

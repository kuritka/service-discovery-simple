package runner

// IServiceRunner is tunning commandline commands
type IServiceRunner interface {
	Run() error
	String() string
}

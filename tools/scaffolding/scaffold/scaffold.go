package scaffold

type ScaffoldWorkflow interface {
	PromptValues()
	GetDestinationDirectories() []string
	GetTemplateDirectory() string
	GetTemplateValues() interface{}
	PostProcess(flags *Args, tmpFolders string)
}

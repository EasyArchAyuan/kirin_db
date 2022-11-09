module record_manager

go 1.19

replace file_manager => ../file_manager

replace tx => ../tx

replace log_manager => ../log_manager

replace buffer_manager => ../buffer_manager

require (
	github.com/stretchr/testify v1.8.1
	tx v0.0.0-00010101000000-000000000000
)

require (
	buffer_manager v0.0.0-00010101000000-000000000000 // indirect
	file_manager v0.0.0-00010101000000-000000000000 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	log_manager v0.0.0-00010101000000-000000000000 // indirect
)

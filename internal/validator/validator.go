package validator

import (
	"encoding/json"
	"fmt"

	"github.com/adipras/api-gladiatore/internal/models"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) ValidateJSON(jsonData []byte) (*models.ServiceConfig, error) {
	var config models.ServiceConfig
	
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %v", err)
	}
	
	if err := v.validator.Struct(&config); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}
	
	if err := v.validateBusinessRules(&config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

func (v *Validator) validateBusinessRules(config *models.ServiceConfig) error {
	tableMap := make(map[string]bool)
	for _, table := range config.Tables {
		if tableMap[table.Name] {
			return fmt.Errorf("duplicate table name: %s", table.Name)
		}
		tableMap[table.Name] = true
		
		fieldMap := make(map[string]bool)
		hasPrimaryKey := false
		for _, field := range table.Fields {
			if fieldMap[field.Name] {
				return fmt.Errorf("duplicate field name '%s' in table '%s'", field.Name, table.Name)
			}
			fieldMap[field.Name] = true
			
			if field.PrimaryKey {
				if hasPrimaryKey {
					return fmt.Errorf("table '%s' has multiple primary keys", table.Name)
				}
				hasPrimaryKey = true
			}
		}
		
		if !hasPrimaryKey {
			return fmt.Errorf("table '%s' must have a primary key", table.Name)
		}
	}
	
	for _, endpoint := range config.Endpoints {
		if endpoint.Operation != "custom" && endpoint.Table != "" {
			if !tableMap[endpoint.Table] {
				return fmt.Errorf("endpoint '%s' references non-existent table '%s'", endpoint.Path, endpoint.Table)
			}
		}
	}
	
	return nil
}
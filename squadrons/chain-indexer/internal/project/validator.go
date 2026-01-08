package project

import (
	"errors"
	"regexp"
)

var monthRegex = regexp.MustCompile(`^\d{4}-\d{2}$`)

func ValidateProjectCreate(p *Project) error {
	if p.ProjectID == "" {
		return errors.New("project_id is required")
	}
	if !monthRegex.MatchString(p.StartMonth) {
		return errors.New("start_month must be YYYY-MM")
	}
	if p.Step <= 0 {
		return errors.New("step must be > 0")
	}
	if p.IntervalSec <= 0 {
		return errors.New("interval_sec must be > 0")
	}
	if p.BasePath == "" {
		return errors.New("base_path is required")
	}
	return nil
}

func ValidateProjectUpdate(old, new *Project) error {
	// 不可变字段校验
	if old.ProjectID != new.ProjectID {
		return errors.New("project_id is immutable")
	}
	if old.StartMonth != new.StartMonth {
		return errors.New("start_month is immutable")
	}
	if old.EndMonth != new.EndMonth {
		return errors.New("end_month is immutable")
	}
	return nil
}

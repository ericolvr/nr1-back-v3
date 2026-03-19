package domain

// CompanyDashboard representa o dashboard para gestor de uma empresa
type CompanyDashboard struct {
	CompanyID                int64                 `json:"company_id"`
	CompanyName              string                `json:"company_name"`
	InProgressQuestionnaires []*TemplateInProgress `json:"in_progress_questionnaires"`
	UnreadNotifications      int64                 `json:"unread_notifications"`
	NotificationsPreview     []*Notification       `json:"notifications_preview"`
	PendingActionPlans       int                   `json:"pending_action_plans"`
	OverdueActionPlans       int                   `json:"overdue_action_plans"`
	PendingInvitations       int                   `json:"pending_invitations"`
	TotalEmployees           int                   `json:"total_employees"`
	TotalDepartments         int                   `json:"total_departments"`
	OverallRiskLevel         string                `json:"overall_risk_level"`
	Alerts                   []string              `json:"alerts"`
}

// TemplateInProgress representa um template em andamento
type TemplateInProgress struct {
	ID                        int64               `json:"id"`
	Name                      string              `json:"name"`
	Description               string              `json:"description"`
	Status                    string              `json:"status"`
	CreatedAt                 string              `json:"created_at"`
	TotalDepartments          int                 `json:"total_departments"`
	DepartmentsCompleted      int                 `json:"departments_completed"`
	DepartmentsInProgress     int                 `json:"departments_in_progress"`
	DepartmentsNotStarted     int                 `json:"departments_not_started"`
	TotalEmployees            int                 `json:"total_employees"`
	TotalInvitations          int                 `json:"total_invitations"`
	CompletedResponses        int                 `json:"completed_responses"`
	DepartmentsWithHighRisk   int                 `json:"departments_with_high_risk"`
	DepartmentsWithMediumRisk int                 `json:"departments_with_medium_risk"`
	DepartmentsWithLowRisk    int                 `json:"departments_with_low_risk"`
	OverallRiskLevel          string              `json:"overall_risk_level"`
	AverageScore              float64             `json:"average_score"`
	Departments               []*DepartmentStatus `json:"departments"`
	LastUpdated               string              `json:"last_updated"`
}

// DepartmentStatus representa o status de um departamento em um template
type DepartmentStatus struct {
	DepartmentID       int64   `json:"department_id"`
	DepartmentName     string  `json:"department_name"`
	TotalEmployees     int64   `json:"total_employees"`
	CompletedResponses int64   `json:"completed_responses"`
	PendingResponses   int64   `json:"pending_responses"`
	ResponseRate       float64 `json:"response_rate"`
	CanCalculateRisk   bool    `json:"can_calculate_risk"`
	Reliability        string  `json:"reliability"`
	AverageScore       float64 `json:"average_score"`
	RiskLevel          string  `json:"risk_level"`
	Status             string  `json:"status"` // completed, in_progress, not_started
	IsActive           bool    `json:"is_active"`
	CanClose           bool    `json:"can_close"`
	CanCloseReason     string  `json:"can_close_reason,omitempty"`

	// Thresholds de confiabilidade (da fórmula ativa)
	ReliabilityThresholds *ReliabilityThresholds `json:"reliability_thresholds,omitempty"`
}

// PartnerDashboard representa o dashboard para consultoria (Partner)
type PartnerDashboard struct {
	PartnerID            int64            `json:"partner_id"`
	PartnerName          string           `json:"partner_name"`
	CompaniesSummary     []*CompanyStatus `json:"companies_summary"`
	TotalCompanies       int              `json:"total_companies"`
	TotalActiveTemplates int              `json:"total_active_questionnaires"`
	CompaniesAtRisk      int              `json:"companies_at_risk"`
	OverallResponseRate  float64          `json:"overall_response_rate"`
	Alerts               []string         `json:"alerts"`
}

// CompanyStatus representa o status de uma empresa (visão do Partner)
type CompanyStatus struct {
	CompanyID         int64   `json:"company_id"`
	CompanyName       string  `json:"company_name"`
	ActiveTemplates   int     `json:"active_questionnaires"`
	ResponseRate      float64 `json:"response_rate"`
	RiskLevel         string  `json:"risk_level"`
	DepartmentsAtRisk int     `json:"departments_at_risk"`
	TotalDepartments  int     `json:"total_departments"`
}

// DepartmentDashboard representa o dashboard para supervisor de departamento
type DepartmentDashboard struct {
	DepartmentID     int64                      `json:"department_id"`
	DepartmentName   string                     `json:"department_name"`
	CompanyID        int64                      `json:"company_id"`
	CompanyName      string                     `json:"company_name"`
	ActiveTemplates  []*DepartmentQuestionnaire `json:"active_questionnaires"`
	EmployeesSummary *EmployeesSummary          `json:"employees_summary"`
	ActionPlans      []*ActionPlanSummary       `json:"action_plans"`
	RiskCategories   []*RiskCategorySummary     `json:"risk_categories"`
	Alerts           []string                   `json:"alerts"`
}

// DepartmentQuestionnaire representa um template do departamento
type DepartmentQuestionnaire struct {
	TemplateID     int64   `json:"template_id"`
	TemplateName   string  `json:"template_name"`
	TotalEmployees int64   `json:"total_employees"`
	Responses      int64   `json:"responses"`
	ResponseRate   float64 `json:"response_rate"`
	RiskLevel      string  `json:"risk_level"`
	CanClose       bool    `json:"can_close"`
}

// EmployeesSummary representa resumo de funcionários do departamento
type EmployeesSummary struct {
	Total        int      `json:"total"`
	Responded    int      `json:"responded"`
	Pending      int      `json:"pending"`
	PendingNames []string `json:"pending_names,omitempty"`
}

// ActionPlanSummary representa resumo de um action plan
type ActionPlanSummary struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Priority  string `json:"priority"`
	DueDate   string `json:"due_date,omitempty"`
	Status    string `json:"status"`
	IsOverdue bool   `json:"is_overdue"`
}

// GlobalDashboard representa o dashboard global/overview para gestor de RH
type GlobalDashboard struct {
	CompanyID   int64  `json:"company_id"`
	CompanyName string `json:"company_name"`

	// KPIs principais
	Metrics *GlobalMetrics `json:"metrics"`

	// Departamentos em avaliação
	DepartmentsOverview []*DepartmentOverview `json:"departments_overview"`

	// Alertas e ações necessárias
	Alerts *DashboardAlerts `json:"alerts"`

	// Resumo rápido
	QuickSummary *QuickSummary `json:"quick_summary"`
}

// GlobalMetrics representa as métricas principais do dashboard
type GlobalMetrics struct {
	ActiveAssessments      int     `json:"active_assessments"`       // Total de avaliações ativas
	ActiveAssessmentsDelta int     `json:"active_assessments_delta"` // Variação vs período anterior
	OverallResponseRate    float64 `json:"overall_response_rate"`    // Taxa de resposta global (%)
	ResponseRateDelta      float64 `json:"response_rate_delta"`      // Variação vs período anterior
	OverallRiskLevel       string  `json:"overall_risk_level"`       // low/medium/high
	DepartmentsAtRisk      int     `json:"departments_at_risk"`      // Departamentos com risco médio/alto
	DepartmentsAtRiskDelta int     `json:"departments_at_risk_delta"`
}

// DepartmentOverview representa visão geral de um departamento
type DepartmentOverview struct {
	DepartmentID   int64   `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	ResponseRate   float64 `json:"response_rate"`   // % de respostas
	TotalEmployees int64   `json:"total_employees"` // Total de funcionários
	Responded      int64   `json:"responded"`       // Já responderam
	RiskLevel      string  `json:"risk_level"`      // low/medium/high
	Status         string  `json:"status"`          // in_progress/can_close/closed
	CanClose       bool    `json:"can_close"`       // Pode encerrar?
	TemplateID     int64   `json:"template_id"`     // ID do template em avaliação
	TemplateName   string  `json:"template_name"`   // Nome do template
	AverageScore   float64 `json:"average_score"`   // Score médio (0-5)
	Reliability    string  `json:"reliability"`     // poor/acceptable/good/excellent
}

// DashboardAlerts representa alertas e ações necessárias
type DashboardAlerts struct {
	CanCloseCount      int      `json:"can_close_count"`      // Departamentos que podem encerrar
	CanCloseList       []string `json:"can_close_list"`       // Lista de nomes
	MediumRiskCount    int      `json:"medium_risk_count"`    // Departamentos com risco médio
	MediumRiskList     []string `json:"medium_risk_list"`     // Lista de nomes
	HighRiskCount      int      `json:"high_risk_count"`      // Departamentos com risco alto
	HighRiskList       []string `json:"high_risk_list"`       // Lista de nomes
	LowResponseCount   int      `json:"low_response_count"`   // Departamentos com baixa resposta (<30%)
	LowResponseList    []string `json:"low_response_list"`    // Lista de nomes
	PendingActionPlans int      `json:"pending_action_plans"` // Planos de ação pendentes
	OverdueActionPlans int      `json:"overdue_action_plans"` // Planos de ação atrasados
}

// QuickSummary representa resumo rápido de informações
type QuickSummary struct {
	TotalDepartments      int `json:"total_departments"`       // Total de departamentos
	DepartmentsInProgress int `json:"departments_in_progress"` // Em avaliação
	DepartmentsCompleted  int `json:"departments_completed"`   // Avaliações concluídas
	TotalEmployees        int `json:"total_employees"`         // Total de funcionários
	EmployeesResponded    int `json:"employees_responded"`     // Já responderam
	EmployeesPending      int `json:"employees_pending"`       // Pendentes
	ActiveActionPlans     int `json:"active_action_plans"`     // Planos de ação ativos
	CompletedActionPlans  int `json:"completed_action_plans"`  // Planos de ação concluídos
}

package prompts

const EnrichmentPromptTemplate = `
# Infrastructure as Code Analysis Enhancement

## Current Analysis

### Resources
{{range .Resources}}
- {{.Type}}.{{.Name}} ({{.File}}:{{.Line}})
{{end}}

### Existing Suggestions (Rule-Based)
{{range .BaseSuggestions}}
- [{{.Severity}}] {{.Message}}
  Recommendation: {{.Recommendation}}
{{end}}

### Knowledge Base Context

#### Relevant Best Practices
{{range .BestPractices}}
- {{.}}
{{end}}

#### Recommended Modules
{{range .Modules}}
- {{.Name}} ({{.Source}}) - {{.Description}}
{{end}}

## Your Task

As an Infrastructure as Code expert, analyze the above and provide:

1. **Enhanced Explanations**: For each existing suggestion, provide:
   - Why it matters (business impact)
   - How to implement (specific code example)
   - What could go wrong if ignored

2. **New Insights**: Identify additional improvements not caught by rules:
   - Architectural patterns that could be improved
   - Security considerations
   - Cost optimization opportunities
   - Module suggestions from the recommended list

3. **Prioritization**: Order suggestions by:
   - Impact (critical/high/medium/low)
   - Effort (easy/medium/hard)
   - ROI (quick wins first)

## Response Format

Respond ONLY with valid JSON:

{
  "enriched_suggestions": [
    {
      "original_id": "suggestion-uuid or null if new",
      "type": "security|cost|best_practice|architecture",
      "severity": "critical|high|medium|low",
      "title": "Brief title",
      "message": "Detailed explanation with business impact",
      "code_example": "# HCL code showing fix",
      "implementation_effort": "easy|medium|hard",
      "estimated_impact": "Description of impact",
      "why_it_matters": "Business justification",
      "references": ["https://...", "https://..."]
    }
  ],
  "architectural_insights": {
    "pattern_detected": "3-tier web app | microservices | ...",
    "strengths": ["..."],
    "areas_for_improvement": ["..."]
  },
  "priority_actions": [
    "Most important action to take first",
    "Second priority",
    "..."
  ]
}
`
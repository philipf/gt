---
type: kanban
created: {{ .CreatedAt.Format "2006-01-02 15:04" }}
updated: {{ .UpdatedAt.Format "2006-01-02 15:04" }}
status: {{.Status}}
channel: {{.Channel}}
externalId: {{.ExternalID}}
---

# {{.Title}}
{{.Description}}

## Notes
- [ ] Step 1

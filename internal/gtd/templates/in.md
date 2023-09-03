---
type: kanban
created: {{ .CreatedAt.Format "2006-01-02 15:04" }}
updated: {{ .UpdatedAt.Format "2006-01-02 15:04" }}
status: {{.Status}}
origin: {{.Channel}}
---

# {{.Title}}

An update about the Machine Learning course

## Notes

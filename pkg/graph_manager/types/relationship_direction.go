package types

// RelationshipDirection represents the temporal direction of a relationship
type RelationshipDirection string

const (
	// DirectionBackward indicates the relationship points to nodes that come before in execution order
	DirectionBackward RelationshipDirection = "backward"

	// DirectionForward indicates the relationship points to nodes that come after in execution order
	DirectionForward RelationshipDirection = "forward"

	// DirectionNone indicates the relationship has no temporal ordering
	DirectionNone RelationshipDirection = "none"
)

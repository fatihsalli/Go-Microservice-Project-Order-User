package graphQL

import "fmt"

func GenerateGraphQLQuery(userId, status string) string {
	// Create GraphQL query
	query := fmt.Sprintf(`
        query {
            orders(userId: "%s", status: "%s") {
                id
                userId
                status
            }
        }
    `, userId, status)

	return query
}

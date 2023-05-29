package graphQL

import "fmt"

func GenerateGraphQLQuery(status string) string {
	// Create GraphQL query
	query := fmt.Sprintf(`
		query {
			orders(status: "%s") {
				id
				userId
				status
				address {
					id
					address
					city
					district
					type
					default {
						isDefaultInvoiceAddress
						isDefaultRegularAddress
					}
				}
				invoiceAddress {
					id
					address
					city
					district
					type
					default {
						isDefaultInvoiceAddress
						isDefaultRegularAddress
					}
				}
				product {
					name
					quantity
					price
				}
				total
				createdAt
				updatedAt
			}
		}
	`, status)

	return query
}

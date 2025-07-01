package servers

import (
	"fmt"
)

var graphqlSandBoxHtmlTemplate = `
    <!DOCTYPE html>
    <html lang="en">
    <body style="margin: 0; overflow-x: hidden; overflow-y: hidden">
    <div id="sandbox" style="height:100vh; width:100vw;"></div>
    <script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
    <script>
    new window.EmbeddedSandbox({
    target: "#sandbox",
    initialEndpoint: "%s",
    });
    </script>
    </body>
    </html>
`

func getMainURL(baseUrl string) string {
	if baseUrl == "" {
		baseUrl = "http://localhost:3000/api/v1"
	}

	return baseUrl
}

func getGraphqlSandboxHtml(baseUrl string) []byte {
	url := getMainURL(baseUrl)
	return []byte(fmt.Sprintf(graphqlSandBoxHtmlTemplate, url+"/graphql"))
}

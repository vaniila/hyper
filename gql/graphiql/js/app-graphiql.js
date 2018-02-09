$(function (global) {

  global.renderGraphiql = function (elem) {

    var config = new graphiqlWorkspace.AppConfig("graphiql", {
      defaultUrl: 'http://localhost:4000/graphql',
      defaultWebsocketUrl: 'ws://localhost:4000/_s',
      defaultQuery:
        "# Welcome to GraphiQL\n" +
        "#\n" +
        "# GraphiQL is an in-browser IDE for writing, validating, and\n" +
        "# testing GraphQL queries.\n" +
        "#\n" +
        "# Type queries into this side of the screen, and you will\n" +
        "# see intelligent typeaheads aware of the current GraphQL type schema and\n" +
        "# live syntax and validation errors highlighted within the text.\n" +
        "#\n" +
        "# To bring up the auto-complete at any point, just press Ctrl-Space.\n" +
        "#\n" +
        "# Press the run button above, or Cmd-Enter to execute the query, and the result\n" +
        "# will appear in the pane to the right.\n\n" +
        "query RebelsShipsQuery {\n  rebels {\n    name\n    ships(first: 1) {\n      edges {\n" +
        "        node {\n          name \n        }\n      }\n    }\n  }\n}"
    });

    ReactDOM.render(
      React.createElement(graphiqlWorkspace.GraphiQLWorkspace, {config: config}),
      document.getElementById('workspace')
    );

  }
}(window))
